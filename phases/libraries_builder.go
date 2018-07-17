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
	"github.com/arduino/go-properties-map"
)

var PRECOMPILED_LIBRARIES_VALID_EXTENSIONS_STATIC = map[string]bool{".a": true}
var PRECOMPILED_LIBRARIES_VALID_EXTENSIONS_DYNAMIC = map[string]bool{".so": true}

type LibrariesBuilder struct{}

func (s *LibrariesBuilder) Run(ctx *types.Context) error {
	librariesBuildPath := ctx.LibrariesBuildPath
	buildProperties := ctx.BuildProperties
	includes := ctx.IncludeFolders
	includes = utils.Map(includes, utils.WrapWithHyphenI)
	libraries := ctx.ImportedLibraries

	err := utils.EnsureFolderExists(librariesBuildPath)
	if err != nil {
		return i18n.WrapError(err)
	}

	objectFiles, err := compileLibraries(ctx, libraries, librariesBuildPath, buildProperties, includes)
	if err != nil {
		return i18n.WrapError(err)
	}

	ctx.LibrariesObjectFiles = objectFiles

	// Search for precompiled libraries
	fixLDFLAGforPrecompiledLibraries(ctx, libraries)

	return nil
}

func fixLDFLAGforPrecompiledLibraries(ctx *types.Context, libraries []*types.Library) error {

	for _, library := range libraries {
		if library.Precompiled {
			// add library src path to compiler.c.elf.extra_flags
			// use library.Name as lib name and srcPath/{mcpu} as location
			mcu := ctx.BuildProperties[constants.BUILD_PROPERTIES_BUILD_MCU]
			path := filepath.Join(library.SrcFolder, mcu)
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
			ctx.BuildProperties[constants.BUILD_PROPERTIES_COMPILER_LIBRARIES_LDFLAGS] += "\"-L" + path + "\" " + libs_cmd + " "
		}
	}
	return nil
}

func compileLibraries(ctx *types.Context, libraries []*types.Library, buildPath string, buildProperties properties.Map, includes []string) ([]string, error) {
	objectFiles := []string{}
	for _, library := range libraries {
		libraryObjectFiles, err := compileLibrary(ctx, library, buildPath, buildProperties, includes)
		if err != nil {
			return nil, i18n.WrapError(err)
		}
		objectFiles = append(objectFiles, libraryObjectFiles...)
	}

	return objectFiles, nil

}

func compileLibrary(ctx *types.Context, library *types.Library, buildPath string, buildProperties properties.Map, includes []string) ([]string, error) {
	logger := ctx.GetLogger()
	if ctx.Verbose {
		logger.Println(constants.LOG_LEVEL_INFO, "Compiling library \"{0}\"", library.Name)
	}
	libraryBuildPath := filepath.Join(buildPath, library.Name)

	err := utils.EnsureFolderExists(libraryBuildPath)
	if err != nil {
		return nil, i18n.WrapError(err)
	}

	objectFiles := []string{}

	if library.Precompiled {
		// search for files with PRECOMPILED_LIBRARIES_VALID_EXTENSIONS
		extensions := func(ext string) bool { return PRECOMPILED_LIBRARIES_VALID_EXTENSIONS_STATIC[ext] }

		filePaths := []string{}
		mcu := buildProperties[constants.BUILD_PROPERTIES_BUILD_MCU]
		err := utils.FindFilesInFolder(&filePaths, filepath.Join(library.SrcFolder, mcu), extensions, true)
		if err != nil {
			return nil, i18n.WrapError(err)
		}
		for _, path := range filePaths {
			if strings.Contains(filepath.Base(path), library.RealName) {
				objectFiles = append(objectFiles, path)
			}
		}
	}

	if library.Layout == types.LIBRARY_RECURSIVE {
		objectFiles, err = builder_utils.CompileFilesRecursive(ctx, objectFiles, library.SrcFolder, libraryBuildPath, buildProperties, includes)
		if err != nil {
			return nil, i18n.WrapError(err)
		}
		if library.DotALinkage {
			archiveFile, err := builder_utils.ArchiveCompiledFiles(ctx, libraryBuildPath, library.Name+".a", objectFiles, buildProperties)
			if err != nil {
				return nil, i18n.WrapError(err)
			}
			objectFiles = []string{archiveFile}
		}
	} else {
		if library.UtilityFolder != "" {
			includes = append(includes, utils.WrapWithHyphenI(library.UtilityFolder))
		}
		objectFiles, err = builder_utils.CompileFiles(ctx, objectFiles, library.SrcFolder, false, libraryBuildPath, buildProperties, includes)
		if err != nil {
			return nil, i18n.WrapError(err)
		}

		if library.UtilityFolder != "" {
			utilityBuildPath := filepath.Join(libraryBuildPath, constants.LIBRARY_FOLDER_UTILITY)
			objectFiles, err = builder_utils.CompileFiles(ctx, objectFiles, library.UtilityFolder, false, utilityBuildPath, buildProperties, includes)
			if err != nil {
				return nil, i18n.WrapError(err)
			}
		}
	}

	return objectFiles, nil
}
