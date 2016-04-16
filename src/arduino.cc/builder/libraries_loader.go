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

package builder

import (
	"arduino.cc/builder/constants"
	"arduino.cc/builder/i18n"
	"arduino.cc/builder/props"
	"arduino.cc/builder/types"
	"arduino.cc/builder/utils"
	"os"
	"path/filepath"
	"strings"
)

type LibrariesLoader struct{}

func (s *LibrariesLoader) Run(ctx *types.Context) error {
	builtInLibrariesFolders := ctx.BuiltInLibrariesFolders
	builtInLibrariesFolders, err := utils.AbsolutizePaths(builtInLibrariesFolders)
	if err != nil {
		return i18n.WrapError(err)
	}
	sortedLibrariesFolders := []string{}
	sortedLibrariesFolders = utils.AppendIfNotPresent(sortedLibrariesFolders, builtInLibrariesFolders...)

	platform := ctx.TargetPlatform
	debugLevel := ctx.DebugLevel
	logger := ctx.GetLogger()

	actualPlatform := ctx.ActualPlatform
	if actualPlatform != platform {
		sortedLibrariesFolders = appendPathToLibrariesFolders(sortedLibrariesFolders, filepath.Join(actualPlatform.Folder, constants.FOLDER_LIBRARIES))
	}

	sortedLibrariesFolders = appendPathToLibrariesFolders(sortedLibrariesFolders, filepath.Join(platform.Folder, constants.FOLDER_LIBRARIES))

	librariesFolders := ctx.OtherLibrariesFolders
	librariesFolders, err = utils.AbsolutizePaths(librariesFolders)
	if err != nil {
		return i18n.WrapError(err)
	}
	sortedLibrariesFolders = utils.AppendIfNotPresent(sortedLibrariesFolders, librariesFolders...)

	ctx.LibrariesFolders = sortedLibrariesFolders

	var libraries []*types.Library
	for _, libraryFolder := range sortedLibrariesFolders {
		subFolders, err := utils.ReadDirFiltered(libraryFolder, utils.FilterDirs)
		if err != nil {
			return i18n.WrapError(err)
		}
		for _, subFolder := range subFolders {
			library, err := makeLibrary(filepath.Join(libraryFolder, subFolder.Name()), debugLevel, logger)
			if err != nil {
				return i18n.WrapError(err)
			}
			libraries = append(libraries, library)
		}
	}

	ctx.Libraries = libraries

	headerToLibraries := make(map[string][]*types.Library)
	for _, library := range libraries {
		headers, err := utils.ReadDirFiltered(library.SrcFolder, utils.FilterFilesWithExtension(".h"))
		if err != nil {
			return i18n.WrapError(err)
		}
		for _, header := range headers {
			headerFileName := header.Name()
			headerToLibraries[headerFileName] = append(headerToLibraries[headerFileName], library)
		}
	}

	ctx.HeaderToLibraries = headerToLibraries

	return nil
}

func makeLibrary(libraryFolder string, debugLevel int, logger i18n.Logger) (*types.Library, error) {
	if _, err := os.Stat(filepath.Join(libraryFolder, constants.LIBRARY_PROPERTIES)); os.IsNotExist(err) {
		return makeLegacyLibrary(libraryFolder)
	}
	return makeNewLibrary(libraryFolder, debugLevel, logger)
}

func makeNewLibrary(libraryFolder string, debugLevel int, logger i18n.Logger) (*types.Library, error) {
	properties, err := props.Load(filepath.Join(libraryFolder, constants.LIBRARY_PROPERTIES), logger)
	if err != nil {
		return nil, i18n.WrapError(err)
	}

	if properties[constants.LIBRARY_MAINTAINER] == constants.EMPTY_STRING && properties[constants.LIBRARY_EMAIL] != constants.EMPTY_STRING {
		properties[constants.LIBRARY_MAINTAINER] = properties[constants.LIBRARY_EMAIL]
	}

	for _, propName := range LIBRARY_NOT_SO_MANDATORY_PROPERTIES {
		if properties[propName] == constants.EMPTY_STRING {
			properties[propName] = "-"
		}
	}

	library := &types.Library{}
	if stat, err := os.Stat(filepath.Join(libraryFolder, constants.LIBRARY_FOLDER_SRC)); err == nil && stat.IsDir() {
		library.Layout = types.LIBRARY_RECURSIVE
		library.SrcFolder = filepath.Join(libraryFolder, constants.LIBRARY_FOLDER_SRC)
	} else {
		library.Layout = types.LIBRARY_FLAT
		library.SrcFolder = libraryFolder
	}

	subFolders, err := utils.ReadDirFiltered(libraryFolder, utils.FilterDirs)
	if err != nil {
		return nil, i18n.WrapError(err)
	}

	if debugLevel >= 0 {
		for _, subFolder := range subFolders {
			if utils.IsSCCSOrHiddenFile(subFolder) {
				if !utils.IsSCCSFile(subFolder) && utils.IsHiddenFile(subFolder) {
					logger.Fprintln(os.Stdout, constants.LOG_LEVEL_WARN, constants.MSG_WARNING_SPURIOUS_FILE_IN_LIB, filepath.Base(subFolder.Name()), properties[constants.LIBRARY_NAME])
				}
			}
		}
	}

	if properties[constants.LIBRARY_ARCHITECTURES] == constants.EMPTY_STRING {
		properties[constants.LIBRARY_ARCHITECTURES] = constants.LIBRARY_ALL_ARCHS
	}
	library.Archs = []string{}
	for _, arch := range strings.Split(properties[constants.LIBRARY_ARCHITECTURES], ",") {
		library.Archs = append(library.Archs, strings.TrimSpace(arch))
	}

	properties[constants.LIBRARY_CATEGORY] = strings.TrimSpace(properties[constants.LIBRARY_CATEGORY])
	if !LIBRARY_CATEGORIES[properties[constants.LIBRARY_CATEGORY]] {
		logger.Fprintln(os.Stdout, constants.LOG_LEVEL_WARN, constants.MSG_WARNING_LIB_INVALID_CATEGORY, properties[constants.LIBRARY_CATEGORY], properties[constants.LIBRARY_NAME], constants.LIB_CATEGORY_UNCATEGORIZED)
		properties[constants.LIBRARY_CATEGORY] = constants.LIB_CATEGORY_UNCATEGORIZED
	}
	library.Category = properties[constants.LIBRARY_CATEGORY]

	if properties[constants.LIBRARY_LICENSE] == constants.EMPTY_STRING {
		properties[constants.LIBRARY_LICENSE] = constants.LIB_LICENSE_UNSPECIFIED
	}
	library.License = properties[constants.LIBRARY_LICENSE]

	library.Folder = libraryFolder
	library.Name = filepath.Base(libraryFolder)
	library.Version = strings.TrimSpace(properties[constants.LIBRARY_VERSION])
	library.Author = strings.TrimSpace(properties[constants.LIBRARY_AUTHOR])
	library.Maintainer = strings.TrimSpace(properties[constants.LIBRARY_MAINTAINER])
	library.Sentence = strings.TrimSpace(properties[constants.LIBRARY_SENTENCE])
	library.Paragraph = strings.TrimSpace(properties[constants.LIBRARY_PARAGRAPH])
	library.URL = strings.TrimSpace(properties[constants.LIBRARY_URL])
	library.IsLegacy = false
	library.DotALinkage = strings.TrimSpace(properties[constants.LIBRARY_DOT_A_LINKAGE]) == "true"
	library.Properties = properties

	return library, nil
}

func makeLegacyLibrary(libraryFolder string) (*types.Library, error) {
	library := &types.Library{
		Folder:    libraryFolder,
		SrcFolder: libraryFolder,
		Layout:    types.LIBRARY_FLAT,
		Name:      filepath.Base(libraryFolder),
		Archs:     []string{constants.LIBRARY_ALL_ARCHS},
		IsLegacy:  true,
	}
	return library, nil
}

func appendPathToLibrariesFolders(librariesFolders []string, newLibrariesFolder string) []string {
	if stat, err := os.Stat(newLibrariesFolder); os.IsNotExist(err) || !stat.IsDir() {
		return librariesFolders
	}

	return utils.AppendIfNotPresent(librariesFolders, newLibrariesFolder)
}
