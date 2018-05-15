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
	"os"
	"path/filepath"
	"strings"

	"github.com/arduino/arduino-builder/i18n"
	"github.com/arduino/arduino-builder/types"
	"github.com/arduino/arduino-builder/utils"
	"github.com/arduino/go-properties-map"
	"github.com/bcmi-labs/arduino-cli/arduino/libraries"
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
		sortedLibrariesFolders = appendPathToLibrariesFolders(sortedLibrariesFolders, filepath.Join(actualPlatform.Folder, "libraries"))
	}

	sortedLibrariesFolders = appendPathToLibrariesFolders(sortedLibrariesFolders, filepath.Join(platform.Folder, "libraries"))

	librariesFolders := ctx.OtherLibrariesFolders
	librariesFolders, err = utils.AbsolutizePaths(librariesFolders)
	if err != nil {
		return i18n.WrapError(err)
	}
	sortedLibrariesFolders = utils.AppendIfNotPresent(sortedLibrariesFolders, librariesFolders...)

	ctx.LibrariesFolders = sortedLibrariesFolders

	var libs []*libraries.Library
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
			libs = append(libs, library)
		}
	}

	ctx.Libraries = libs

	headerToLibraries := make(map[string][]*libraries.Library)
	for _, library := range libs {
		headers, err := utils.ReadDirFiltered(library.SrcFolder, utils.FilterFilesWithExtensions(".h", ".hpp", ".hh"))
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

func makeLibrary(libraryFolder string, debugLevel int, logger i18n.Logger) (*libraries.Library, error) {
	if _, err := os.Stat(filepath.Join(libraryFolder, "library.properties")); os.IsNotExist(err) {
		return makeLegacyLibrary(libraryFolder)
	}
	return makeNewLibrary(libraryFolder, debugLevel, logger)
}

func addUtilityFolder(library *libraries.Library) {
	utilitySourcePath := filepath.Join(library.Folder, "utility")
	stat, err := os.Stat(utilitySourcePath)
	if err == nil && stat.IsDir() {
		library.UtilityFolder = utilitySourcePath
	}
}

func makeNewLibrary(libraryFolder string, debugLevel int, logger i18n.Logger) (*libraries.Library, error) {
	libProperties, err := properties.Load(filepath.Join(libraryFolder, "library.properties"))
	if err != nil {
		return nil, i18n.WrapError(err)
	}

	if libProperties["maintainer"] == "" && libProperties["email"] != "" {
		libProperties["maintainer"] = libProperties["email"]
	}

	for _, propName := range libraries.MandatoryProperties {
		if libProperties[propName] == "" {
			libProperties[propName] = "-"
		}
	}

	library := &libraries.Library{}
	library.Folder = libraryFolder
	if stat, err := os.Stat(filepath.Join(libraryFolder, "src")); err == nil && stat.IsDir() {
		library.Layout = libraries.RecursiveLayout
		library.SrcFolder = filepath.Join(libraryFolder, "src")
	} else {
		library.Layout = libraries.FlatLayout
		library.SrcFolder = libraryFolder
		addUtilityFolder(library)
	}

	subFolders, err := utils.ReadDirFiltered(libraryFolder, utils.FilterDirs)
	if err != nil {
		return nil, i18n.WrapError(err)
	}

	if debugLevel >= 0 {
		for _, subFolder := range subFolders {
			if utils.IsSCCSOrHiddenFile(subFolder) {
				if !utils.IsSCCSFile(subFolder) && utils.IsHiddenFile(subFolder) {
					logger.Fprintln(os.Stdout, "warn",
						"WARNING: Spurious {0} folder in '{1}' library",
						filepath.Base(subFolder.Name()), libProperties["name"])
				}
			}
		}
	}

	if libProperties["architectures"] == "" {
		libProperties["architectures"] = "*"
	}
	library.Architectures = []string{}
	for _, arch := range strings.Split(libProperties["architectures"], ",") {
		library.Architectures = append(library.Architectures, strings.TrimSpace(arch))
	}

	libProperties["category"] = strings.TrimSpace(libProperties["category"])
	if !libraries.ValidCategories[libProperties["category"]] {
		logger.Fprintln(os.Stdout, "warn",
			"WARNING: Category '{0}' in library {1} is not valid. Setting to '{2}'",
			libProperties["category"], libProperties["name"], "Uncategorized")
		libProperties["category"] = "Uncategorized"
	}
	library.Category = libProperties["category"]

	if libProperties["license"] == "" {
		libProperties["license"] = "Unspecified"
	}
	library.License = libProperties["license"]

	library.Name = filepath.Base(libraryFolder)
	library.RealName = strings.TrimSpace(libProperties["name"])
	library.Version = strings.TrimSpace(libProperties["version"])
	library.Author = strings.TrimSpace(libProperties["author"])
	library.Maintainer = strings.TrimSpace(libProperties["maintainer"])
	library.Sentence = strings.TrimSpace(libProperties["sentence"])
	library.Paragraph = strings.TrimSpace(libProperties["paragraph"])
	library.Website = strings.TrimSpace(libProperties["url"])
	library.IsLegacy = false
	library.DotALinkage = strings.TrimSpace(libProperties["dot_a_linkage"]) == "true"
	library.Precompiled = strings.TrimSpace(libProperties["precompiled"]) == "true"
	library.LDflags = strings.TrimSpace(libProperties["ldflags"])
	library.Properties = libProperties

	return library, nil
}

func makeLegacyLibrary(libraryFolder string) (*libraries.Library, error) {
	library := &libraries.Library{
		Folder:        libraryFolder,
		SrcFolder:     libraryFolder,
		Layout:        libraries.FlatLayout,
		Name:          filepath.Base(libraryFolder),
		Architectures: []string{"*"},
		IsLegacy:      true,
	}
	addUtilityFolder(library)
	return library, nil
}

func appendPathToLibrariesFolders(librariesFolders []string, newLibrariesFolder string) []string {
	if stat, err := os.Stat(newLibrariesFolder); os.IsNotExist(err) || !stat.IsDir() {
		return librariesFolders
	}

	return utils.AppendIfNotPresent(librariesFolders, newLibrariesFolder)
}
