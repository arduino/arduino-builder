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
	"arduino.cc/builder/types"
	"arduino.cc/builder/utils"
	"os"
	"path/filepath"
	"strings"
)

type IncludesToIncludeFolders struct{}

func (s *IncludesToIncludeFolders) Run(context map[string]interface{}) error {
	if !utils.MapHas(context, constants.CTX_LIBRARIES) {
		return nil
	}
	includes := context[constants.CTX_INCLUDES].([]string)
	headerToLibraries := context[constants.CTX_HEADER_TO_LIBRARIES].(map[string][]*types.Library)
	debugLevel := utils.DebugLevel(context)
	logger := context[constants.CTX_LOGGER].(i18n.Logger)
	platform := context[constants.CTX_TARGET_PLATFORM].(*types.Platform)
	actualPlatform := context[constants.CTX_ACTUAL_PLATFORM].(*types.Platform)

	var importedLibraries []*types.Library
	if utils.MapHas(context, constants.CTX_IMPORTED_LIBRARIES) {
		importedLibraries = context[constants.CTX_IMPORTED_LIBRARIES].([]*types.Library)
	}
	newlyImportedLibraries, err := resolveLibraries(includes, headerToLibraries, importedLibraries, []*types.Platform{platform, actualPlatform}, debugLevel, logger)
	if err != nil {
		return utils.WrapError(err)
	}

	foldersWithSources := context[constants.CTX_FOLDERS_WITH_SOURCES_QUEUE].(*types.UniqueStringQueue)

	for _, newlyImportedLibrary := range newlyImportedLibraries {
		if !sliceContainsLibrary(importedLibraries, newlyImportedLibrary) {
			importedLibraries = append(importedLibraries, newlyImportedLibrary)
			foldersWithSources.Push(newlyImportedLibrary.SrcFolder)
		}
	}

	context[constants.CTX_IMPORTED_LIBRARIES] = importedLibraries

	buildProperties := context[constants.CTX_BUILD_PROPERTIES].(map[string]string)
	verbose := context[constants.CTX_VERBOSE].(bool)
	includeFolders := resolveIncludeFolders(newlyImportedLibraries, buildProperties, verbose)
	context[constants.CTX_INCLUDE_FOLDERS] = includeFolders

	return nil
}

func resolveIncludeFolders(importedLibraries []*types.Library, buildProperties map[string]string, verbose bool) []string {
	var includeFolders []string
	includeFolders = append(includeFolders, buildProperties[constants.BUILD_PROPERTIES_BUILD_CORE_PATH])
	if buildProperties[constants.BUILD_PROPERTIES_BUILD_VARIANT_PATH] != constants.EMPTY_STRING {
		includeFolders = append(includeFolders, buildProperties[constants.BUILD_PROPERTIES_BUILD_VARIANT_PATH])
	}

	for _, library := range importedLibraries {
		includeFolders = append(includeFolders, library.SrcFolder)
	}

	return includeFolders
}

//FIXME it's also resolving previously resolved libraries
func resolveLibraries(includes []string, headerToLibraries map[string][]*types.Library, previousImportedLibraries []*types.Library, platforms []*types.Platform, debugLevel int, logger i18n.Logger) ([]*types.Library, error) {
	markImportedLibrary := make(map[*types.Library]bool)
	for _, header := range includes {
		resolveLibrary(header, headerToLibraries, markImportedLibrary, previousImportedLibraries, platforms, debugLevel, logger)
	}

	var importedLibraries []*types.Library
	for library, _ := range markImportedLibrary {
		importedLibraries = append(importedLibraries, library)
	}

	return importedLibraries, nil
}

func resolveLibrary(header string, headerToLibraries map[string][]*types.Library, markImportedLibrary map[*types.Library]bool, previousImportedLibraries []*types.Library, platforms []*types.Platform, debugLevel int, logger i18n.Logger) {
	libraries := headerToLibraries[header]

	if libraries == nil {
		return
	}

	if len(libraries) == 1 {
		markImportedLibrary[libraries[0]] = true
		return
	}

	var library *types.Library

	for _, platform := range platforms {
		if platform != nil && library == nil {
			librariesWithinSpecifiedPlatform := librariesWithinPlatform(libraries, platform)
			library = findBestLibraryWithHeader(header, librariesWithinSpecifiedPlatform)
		}
	}

	for _, platform := range platforms {
		if platform != nil && library == nil {
			library = findBestLibraryWithHeader(header, librariesCompatibleWithPlatform(libraries, platform))
		}
	}

	if library == nil {
		library = findBestLibraryWithHeader(header, libraries)
	}

	if library == nil {
		library = libraries[0]
	}

	if debugLevel > 0 && !sliceContainsLibrary(previousImportedLibraries, library) {
		logger.Fprintln(os.Stderr, constants.MSG_LIBRARIES_MULTIPLE_LIBS_FOUND_FOR, header)
		logger.Fprintln(os.Stderr, constants.MSG_LIBRARIES_USED, library.Folder)
		for _, notUsedLibrary := range libraries {
			if library != notUsedLibrary {
				logger.Fprintln(os.Stderr, constants.MSG_LIBRARIES_NOT_USED, notUsedLibrary.Folder)
			}
		}
	}
	markImportedLibrary[library] = true
}

func libraryCompatibleWithPlatform(library *types.Library, platform *types.Platform) bool {
	if len(library.Archs) == 0 {
		return true
	}
	if utils.SliceContains(library.Archs, constants.LIBRARY_ALL_ARCHS) {
		return true
	}
	return utils.SliceContains(library.Archs, platform.PlatformId)
}

func librariesCompatibleWithPlatform(libraries []*types.Library, platform *types.Platform) []*types.Library {
	var compatibleLibraries []*types.Library
	for _, library := range libraries {
		if libraryCompatibleWithPlatform(library, platform) {
			compatibleLibraries = append(compatibleLibraries, library)
		}
	}

	return compatibleLibraries
}

func librariesWithinPlatform(libraries []*types.Library, platform *types.Platform) []*types.Library {
	var librariesWithinSpecifiedPlatform []*types.Library
	for _, library := range libraries {
		cleanPlatformFolder := filepath.Clean(platform.Folder)
		cleanLibraryFolder := filepath.Clean(library.SrcFolder)
		if strings.Contains(cleanLibraryFolder, cleanPlatformFolder) {
			librariesWithinSpecifiedPlatform = append(librariesWithinSpecifiedPlatform, library)
		}
	}

	return librariesWithinSpecifiedPlatform

}

func findBestLibraryWithHeader(header string, libraries []*types.Library) *types.Library {
	headerName := strings.Replace(header, filepath.Ext(header), constants.EMPTY_STRING, -1)

	var library *types.Library
	for _, headerName := range []string{headerName, strings.ToLower(headerName)} {
		library = findLibWithName(headerName, libraries)
		if library != nil {
			return library
		}
		library = findLibWithName(headerName+"-master", libraries)
		if library != nil {
			return library
		}
		library = findLibWithNameStartingWith(headerName, libraries)
		if library != nil {
			return library
		}
		library = findLibWithNameEndingWith(headerName, libraries)
		if library != nil {
			return library
		}
		library = findLibWithNameContaining(headerName, libraries)
		if library != nil {
			return library
		}
	}

	return nil
}

func findLibWithName(name string, libraries []*types.Library) *types.Library {
	for _, library := range libraries {
		if library.Name == name {
			return library
		}
	}
	return nil
}

func findLibWithNameStartingWith(name string, libraries []*types.Library) *types.Library {
	for _, library := range libraries {
		if strings.HasPrefix(library.Name, name) {
			return library
		}
	}
	return nil
}

func findLibWithNameEndingWith(name string, libraries []*types.Library) *types.Library {
	for _, library := range libraries {
		if strings.HasSuffix(library.Name, name) {
			return library
		}
	}
	return nil
}

func findLibWithNameContaining(name string, libraries []*types.Library) *types.Library {
	for _, library := range libraries {
		if strings.Contains(library.Name, name) {
			return library
		}
	}
	return nil
}

// thank you golang: I can not use/recycle/adapt utils.SliceContains
func sliceContainsLibrary(slice []*types.Library, target *types.Library) bool {
	for _, value := range slice {
		if value.SrcFolder == target.SrcFolder {
			return true
		}
	}
	return false
}
