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

	"arduino.cc/builder/builder_utils"
	"arduino.cc/builder/constants"
	"arduino.cc/builder/i18n"
	"arduino.cc/builder/types"
	"arduino.cc/builder/utils"
	"arduino.cc/properties"
)

type LibrariesBuilder struct{}

func (s *LibrariesBuilder) Run(ctx *types.Context) error {
	librariesBuildPath := ctx.LibrariesBuildPath
	buildProperties := ctx.BuildProperties
	includes := ctx.IncludeFolders
	includes = utils.Map(includes, utils.WrapWithHyphenI)
	libraries := ctx.ImportedLibraries
	verbose := ctx.Verbose
	warningsLevel := ctx.WarningsLevel
	logger := ctx.GetLogger()

	err := utils.EnsureFolderExists(librariesBuildPath)
	if err != nil {
		return i18n.WrapError(err)
	}

	objectFiles, err := compileLibraries(libraries, librariesBuildPath, buildProperties, includes, verbose, warningsLevel, logger)
	if err != nil {
		return i18n.WrapError(err)
	}

	ctx.LibrariesObjectFiles = objectFiles

	return nil
}

func compileLibraries(libraries []*types.Library, buildPath string, buildProperties properties.Map, includes []string, verbose bool, warningsLevel string, logger i18n.Logger) ([]string, error) {
	objectFiles := []string{}
	for _, library := range libraries {
		libraryObjectFiles, err := compileLibrary(library, buildPath, buildProperties, includes, verbose, warningsLevel, logger)
		if err != nil {
			return nil, i18n.WrapError(err)
		}
		objectFiles = append(objectFiles, libraryObjectFiles...)
	}

	return objectFiles, nil

}

func compileLibrary(library *types.Library, buildPath string, buildProperties properties.Map, includes []string, verbose bool, warningsLevel string, logger i18n.Logger) ([]string, error) {
	if verbose {
		logger.Println(constants.LOG_LEVEL_INFO, "Compiling library \"{0}\"", library.Name)
	}
	libraryBuildPath := filepath.Join(buildPath, library.Name)

	err := utils.EnsureFolderExists(libraryBuildPath)
	if err != nil {
		return nil, i18n.WrapError(err)
	}

	objectFiles := []string{}
	if library.Layout == types.LIBRARY_RECURSIVE {
		objectFiles, err = builder_utils.CompileFilesRecursive(objectFiles, library.SrcFolder, libraryBuildPath, buildProperties, includes, verbose, warningsLevel, logger)
		if err != nil {
			return nil, i18n.WrapError(err)
		}
		if library.DotALinkage {
			archiveFile, err := builder_utils.ArchiveCompiledFiles(libraryBuildPath, library.Name+".a", objectFiles, buildProperties, verbose, logger)
			if err != nil {
				return nil, i18n.WrapError(err)
			}
			objectFiles = []string{archiveFile}
		}
	} else {
		if library.UtilityFolder != "" {
			includes = append(includes, utils.WrapWithHyphenI(library.UtilityFolder))
		}
		objectFiles, err = builder_utils.CompileFiles(objectFiles, library.SrcFolder, false, libraryBuildPath, buildProperties, includes, verbose, warningsLevel, logger)
		if err != nil {
			return nil, i18n.WrapError(err)
		}

		if library.UtilityFolder != "" {
			utilityBuildPath := filepath.Join(libraryBuildPath, constants.LIBRARY_FOLDER_UTILITY)
			objectFiles, err = builder_utils.CompileFiles(objectFiles, library.UtilityFolder, false, utilityBuildPath, buildProperties, includes, verbose, warningsLevel, logger)
			if err != nil {
				return nil, i18n.WrapError(err)
			}
		}
	}

	return objectFiles, nil
}
