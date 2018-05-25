/*
 * This file is part of Arduino Builder.
 *
 * Arduino Builder is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
 *
 * As a special exception, you may use this file as part of a free software
 * library without restriction.  Specifically, if other files instantiate
 * templates or use macros or inline functions from this file, or you compile
 * this file and link it with other files to produce an executable, this
 * file does not by itself cause the resulting executable to be covered by
 * the GNU General Public License.  This exception does not however
 * invalidate any other reasons why the executable file might be covered by
 * the GNU General Public License.
 *
 * Copyright 2015 Arduino LLC (http://www.arduino.cc/)
 */

package phases

import (
	"path/filepath"
	"strings"

	"github.com/arduino/arduino-builder/builder_utils"
	"github.com/arduino/arduino-builder/constants"
	"github.com/arduino/arduino-builder/i18n"
	"github.com/arduino/arduino-builder/types"
	"github.com/arduino/arduino-builder/utils"
	"github.com/arduino/go-paths-helper"
	"github.com/arduino/go-properties-map"
	"github.com/bcmi-labs/arduino-cli/arduino/libraries"
)

var PRECOMPILED_LIBRARIES_VALID_EXTENSIONS_STATIC = map[string]bool{".a": true}
var PRECOMPILED_LIBRARIES_VALID_EXTENSIONS_DYNAMIC = map[string]bool{".so": true}

type LibrariesBuilder struct{}

func (s *LibrariesBuilder) Run(ctx *types.Context) error {
	librariesBuildPath := ctx.LibrariesBuildPath
	buildProperties := ctx.BuildProperties
	includes := utils.Map(ctx.IncludeFolders.AsStrings(), utils.WrapWithHyphenI)
	libraries := ctx.ImportedLibraries
	verbose := ctx.Verbose
	warningsLevel := ctx.WarningsLevel
	logger := ctx.GetLogger()

	if err := librariesBuildPath.MkdirAll(); err != nil {
		return i18n.WrapError(err)
	}

	objectFiles, err := compileLibraries(libraries, librariesBuildPath, buildProperties, includes, verbose, warningsLevel, logger)
	if err != nil {
		return i18n.WrapError(err)
	}

	ctx.LibrariesObjectFiles = objectFiles

	// Search for precompiled libraries
	fixLDFLAGforPrecompiledLibraries(ctx, libraries)

	return nil
}

func fixLDFLAGforPrecompiledLibraries(ctx *types.Context, libraries []*libraries.Library) error {

	for _, library := range libraries {
		if library.Precompiled {
			// add library src path to compiler.c.elf.extra_flags
			// use library.Name as lib name and srcPath/{mcpu} as location
			mcu := ctx.BuildProperties[constants.BUILD_PROPERTIES_BUILD_MCU]
			path := library.SrcFolder.Join(mcu).String()
			// find all library names in the folder and prepend -l
			filePaths := []string{}
			libs_cmd := library.LDflags + " "
			extensions := func(ext string) bool { return PRECOMPILED_LIBRARIES_VALID_EXTENSIONS_DYNAMIC[ext] }
			utils.FindFilesInFolder(&filePaths, path, extensions, true)
			for _, lib := range filePaths {
				name := strings.TrimSuffix(filepath.Base(lib), filepath.Ext(lib))
				// strip "lib" first occurrence
				name = strings.Replace(name, "lib", "", 1)
				libs_cmd += "-l" + name + " "
			}
			ctx.BuildProperties[constants.BUILD_PROPERTIES_COMPILER_C_ELF_EXTRAFLAGS] += "\"-L" + path + "\" " + libs_cmd
		}
	}
	return nil
}

func compileLibraries(libraries []*libraries.Library, buildPath *paths.Path, buildProperties properties.Map, includes []string, verbose bool, warningsLevel string, logger i18n.Logger) (paths.PathList, error) {
	objectFiles := paths.NewPathList()
	for _, library := range libraries {
		libraryObjectFiles, err := compileLibrary(library, buildPath, buildProperties, includes, verbose, warningsLevel, logger)
		if err != nil {
			return nil, i18n.WrapError(err)
		}
		objectFiles = append(objectFiles, libraryObjectFiles...)
	}

	return objectFiles, nil
}

func compileLibrary(library *libraries.Library, buildPath *paths.Path, buildProperties properties.Map, includes []string, verbose bool, warningsLevel string, logger i18n.Logger) (paths.PathList, error) {
	if verbose {
		logger.Println(constants.LOG_LEVEL_INFO, "Compiling library \"{0}\"", library.Name)
	}
	libraryBuildPath := buildPath.Join(library.Name)

	if err := libraryBuildPath.MkdirAll(); err != nil {
		return nil, i18n.WrapError(err)
	}

	objectFiles := paths.NewPathList()

	if library.Precompiled {
		// search for files with PRECOMPILED_LIBRARIES_VALID_EXTENSIONS
		extensions := func(ext string) bool { return PRECOMPILED_LIBRARIES_VALID_EXTENSIONS_STATIC[ext] }

		filePaths := []string{}
		mcu := buildProperties[constants.BUILD_PROPERTIES_BUILD_MCU]
		err := utils.FindFilesInFolder(&filePaths, library.SrcFolder.Join(mcu).String(), extensions, true)
		if err != nil {
			return nil, i18n.WrapError(err)
		}
		for _, path := range filePaths {
			if strings.Contains(filepath.Base(path), library.RealName) {
				objectFiles.Add(paths.New(path))
			}
		}
	}

	if library.Layout == libraries.RecursiveLayout {
		libObjectFiles, err := builder_utils.CompileFilesRecursive(library.SrcFolder, libraryBuildPath, buildProperties, includes, verbose, warningsLevel, logger)
		if err != nil {
			return nil, i18n.WrapError(err)
		}
		if library.DotALinkage {
			archiveFile, err := builder_utils.ArchiveCompiledFiles(libraryBuildPath, paths.New(library.Name+".a"), libObjectFiles, buildProperties, verbose, logger)
			if err != nil {
				return nil, i18n.WrapError(err)
			}
			objectFiles.Add(archiveFile)
		} else {
			objectFiles.AddAll(libObjectFiles)
		}
	} else {
		if library.UtilityFolder != nil {
			includes = append(includes, utils.WrapWithHyphenI(library.UtilityFolder.String()))
		}
		libObjectFiles, err := builder_utils.CompileFiles(library.SrcFolder, false, libraryBuildPath, buildProperties, includes, verbose, warningsLevel, logger)
		if err != nil {
			return nil, i18n.WrapError(err)
		}
		objectFiles.AddAll(libObjectFiles)

		if library.UtilityFolder != nil {
			utilityBuildPath := libraryBuildPath.Join("utility")
			utilityObjectFiles, err := builder_utils.CompileFiles(library.UtilityFolder, false, utilityBuildPath, buildProperties, includes, verbose, warningsLevel, logger)
			if err != nil {
				return nil, i18n.WrapError(err)
			}
			objectFiles.AddAll(utilityObjectFiles)
		}
	}

	return objectFiles, nil
}
