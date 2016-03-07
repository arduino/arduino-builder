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
	"arduino.cc/builder/props"
	"arduino.cc/builder/types"
	"arduino.cc/builder/utils"
	"path/filepath"
	"strings"
)

type IncludesToIncludeFolders struct{}

func (s *IncludesToIncludeFolders) Run(context map[string]interface{}) error {
	includes := []string{}
	if utils.MapHas(context, constants.CTX_INCLUDES) {
		includes = context[constants.CTX_INCLUDES].([]string)
	}
	headerToLibraries := make(map[string][]*types.Library)
	if utils.MapHas(context, constants.CTX_HEADER_TO_LIBRARIES) {
		headerToLibraries = context[constants.CTX_HEADER_TO_LIBRARIES].(map[string][]*types.Library)
	}

	platform := context[constants.CTX_TARGET_PLATFORM].(*types.Platform)
	actualPlatform := context[constants.CTX_ACTUAL_PLATFORM].(*types.Platform)
	libraryResolutionResults := context[constants.CTX_LIBRARY_RESOLUTION_RESULTS].(map[string]types.LibraryResolutionResult)

	importedLibraries := []*types.Library{}
	if utils.MapHas(context, constants.CTX_IMPORTED_LIBRARIES) {
		importedLibraries = context[constants.CTX_IMPORTED_LIBRARIES].([]*types.Library)
	}
	newlyImportedLibraries, err := resolveLibraries(includes, headerToLibraries, importedLibraries, []*types.Platform{actualPlatform, platform}, libraryResolutionResults)
	if err != nil {
		return utils.WrapError(err)
	}

	foldersWithSources := context[constants.CTX_FOLDERS_WITH_SOURCES_QUEUE].(*types.UniqueSourceFolderQueue)

	for _, newlyImportedLibrary := range newlyImportedLibraries {
		if !sliceContainsLibrary(importedLibraries, newlyImportedLibrary) {
			importedLibraries = append(importedLibraries, newlyImportedLibrary)
			sourceFolders := types.LibraryToSourceFolder(newlyImportedLibrary)
			for _, sourceFolder := range sourceFolders {
				foldersWithSources.Push(sourceFolder)
			}
		}
	}

	context[constants.CTX_IMPORTED_LIBRARIES] = importedLibraries

	buildProperties := context[constants.CTX_BUILD_PROPERTIES].(props.PropertiesMap)
	verbose := context[constants.CTX_VERBOSE].(bool)
	includeFolders := resolveIncludeFolders(newlyImportedLibraries, buildProperties, verbose)
	context[constants.CTX_INCLUDE_FOLDERS] = includeFolders

	return nil
}

func resolveIncludeFolders(importedLibraries []*types.Library, buildProperties props.PropertiesMap, verbose bool) []string {
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
func resolveLibraries(includes []string, headerToLibraries map[string][]*types.Library, importedLibraries []*types.Library, platforms []*types.Platform, libraryResolutionResults map[string]types.LibraryResolutionResult) ([]*types.Library, error) {
	markImportedLibrary := make(map[*types.Library]bool)
	for _, library := range importedLibraries {
		markImportedLibrary[library] = true
	}
	for _, header := range includes {
		resolveLibrary(header, headerToLibraries, markImportedLibrary, platforms, libraryResolutionResults)
	}

	var newlyImportedLibraries []*types.Library
	for library, _ := range markImportedLibrary {
		newlyImportedLibraries = append(newlyImportedLibraries, library)
	}

	return newlyImportedLibraries, nil
}

func resolveLibrary(header string, headerToLibraries map[string][]*types.Library, markImportedLibrary map[*types.Library]bool, platforms []*types.Platform, libraryResolutionResults map[string]types.LibraryResolutionResult) {
	libraries := append([]*types.Library{}, headerToLibraries[header]...)

	if libraries == nil || len(libraries) == 0 {
		return
	}

	if len(libraries) == 1 {
		markImportedLibrary[libraries[0]] = true
		return
	}

	if markImportedLibraryContainsOneOfCandidates(markImportedLibrary, libraries) {
		return
	}

	reverse(libraries)

	librariesInPlatforms := librariesInSomePlatform(libraries, platforms)

	var library *types.Library

	for _, platform := range platforms {
		if platform != nil {
			library = findBestLibraryWithHeader(header, librariesCompatibleWithPlatform(libraries, platform, true))
		}
	}

	if library == nil {
		library = findBestLibraryWithHeader(header, libraries)
	}

	if library == nil {
		// reorder libraries to promote fully compatible ones
		for _, platform := range platforms {
			if platform != nil {
				libraries = append(librariesCompatibleWithPlatform(libraries, platform, false), libraries...)
			}
		}
		library = libraries[0]
	}

	library = useAlreadyImportedLibraryWithSameNameIfExists(library, markImportedLibrary)

	isLibraryFromPlatform := findLibraryIn(librariesInPlatforms, library) != nil
	libraryResolutionResults[header] = types.LibraryResolutionResult{Library: library, IsLibraryFromPlatform: isLibraryFromPlatform, NotUsedLibraries: filterOutLibraryFrom(libraries, library)}

	markImportedLibrary[library] = true
}

//facepalm. sort.Reverse needs an Interface that implements Len/Less/Swap. It's a slice! What else for reversing it?!?
func reverse(data []*types.Library) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

func librariesInSomePlatform(libraries []*types.Library, platforms []*types.Platform) []*types.Library {
	librariesInPlatforms := []*types.Library{}
	for _, platform := range platforms {
		if platform != nil {
			librariesWithinSpecifiedPlatform := librariesWithinPlatform(libraries, platform)
			librariesInPlatforms = append(librariesInPlatforms, librariesWithinSpecifiedPlatform...)
		}
	}
	return librariesInPlatforms
}

func markImportedLibraryContainsOneOfCandidates(markImportedLibrary map[*types.Library]bool, libraries []*types.Library) bool {
	for markedLibrary, _ := range markImportedLibrary {
		for _, library := range libraries {
			if markedLibrary == library {
				return true
			}
		}
	}
	return false
}

func useAlreadyImportedLibraryWithSameNameIfExists(library *types.Library, markImportedLibrary map[*types.Library]bool) *types.Library {
	for lib, _ := range markImportedLibrary {
		if lib.Name == library.Name {
			return lib
		}
	}
	return library
}

func filterOutLibraryFrom(libraries []*types.Library, libraryToRemove *types.Library) []*types.Library {
	filteredOutLibraries := []*types.Library{}
	for _, lib := range libraries {
		if lib != libraryToRemove {
			filteredOutLibraries = append(filteredOutLibraries, lib)
		}
	}
	return filteredOutLibraries
}

func findLibraryIn(libraries []*types.Library, library *types.Library) *types.Library {
	for _, lib := range libraries {
		if lib == library {
			return lib
		}
	}
	return nil
}

func libraryCompatibleWithPlatform(library *types.Library, platform *types.Platform) (bool, bool) {
	if len(library.Archs) == 0 {
		return true, true
	}
	if utils.SliceContains(library.Archs, constants.LIBRARY_ALL_ARCHS) {
		return true, true
	}
	return utils.SliceContains(library.Archs, platform.PlatformId), false
}

func libraryCompatibleWithAllPlatforms(library *types.Library) bool {
	if utils.SliceContains(library.Archs, constants.LIBRARY_ALL_ARCHS) {
		return true
	}
	return false
}

func librariesCompatibleWithPlatform(libraries []*types.Library, platform *types.Platform, reorder bool) []*types.Library {
	var compatibleLibraries []*types.Library
	for _, library := range libraries {
		compatible, generic := libraryCompatibleWithPlatform(library, platform)
		if compatible {
			if !generic && len(compatibleLibraries) != 0 && libraryCompatibleWithAllPlatforms(compatibleLibraries[0]) && reorder == true {
				//priority inversion
				compatibleLibraries = append([]*types.Library{library}, compatibleLibraries...)
			} else {
				compatibleLibraries = append(compatibleLibraries, library)
			}
		}
	}

	return compatibleLibraries
}

func librariesWithinPlatform(libraries []*types.Library, platform *types.Platform) []*types.Library {
	librariesWithinSpecifiedPlatform := []*types.Library{}
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
