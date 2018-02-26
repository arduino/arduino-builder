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

	"github.com/arduino/arduino-builder/constants"
	"github.com/arduino/arduino-builder/i18n"
	"github.com/arduino/arduino-builder/types"
	"github.com/arduino/arduino-builder/utils"
	"github.com/arduino/go-properties-map"
	"github.com/bcmi-labs/arduino-cli/cores"
)

type HardwareLoader struct{}

func (s *HardwareLoader) Run(ctx *types.Context) error {
	logger := ctx.GetLogger()

	packages := &cores.Packages{
		Packages:   map[string]*cores.Package{},
		Properties: properties.Map{},
	}

	folders := ctx.HardwareFolders
	folders, err := utils.AbsolutizePaths(folders)
	if err != nil {
		return i18n.WrapError(err)
	}

	for _, folder := range folders {
		stat, err := os.Stat(folder)
		if err != nil {
			return i18n.WrapError(err)
		}
		if !stat.IsDir() {
			return i18n.ErrorfWithLogger(logger, constants.MSG_MUST_BE_A_FOLDER, folder)
		}

		hardwarePlatformTxt, err := properties.SafeLoad(filepath.Join(folder, constants.FILE_PLATFORM_TXT))
		if err != nil {
			return i18n.WrapError(err)
		}
		packages.Properties.Merge(hardwarePlatformTxt)

		subfolders, err := utils.ReadDirFiltered(folder, utils.FilterDirs)
		if err != nil {
			return i18n.WrapError(err)
		}
		subfolders = utils.FilterOutFoldersByNames(subfolders, constants.FOLDER_TOOLS)

		for _, subfolder := range subfolders {
			subfolderPath := filepath.Join(folder, subfolder.Name())
			packageId := subfolder.Name()

			if _, err := os.Stat(filepath.Join(subfolderPath, constants.FOLDER_HARDWARE)); err == nil {
				subfolderPath = filepath.Join(subfolderPath, constants.FOLDER_HARDWARE)
			}

			targetPackage := getOrCreatePackage(packages, packageId)
			err = loadPackage(targetPackage, subfolderPath)
			if err != nil {
				return i18n.WrapError(err)
			}
			packages.Packages[packageId] = targetPackage
		}
	}

	ctx.Hardware = packages

	return nil
}

func getOrCreatePackage(packages *cores.Packages, packageId string) *cores.Package {
	if _, ok := packages.Packages[packageId]; ok {
		return packages.Packages[packageId]
	}

	targetPackage := cores.Package{}
	targetPackage.Name = packageId
	targetPackage.Platforms = map[string]*cores.Platform{}
	targetPackage.Packages = packages
	// targetPackage.Properties = properties.Map{}

	return &targetPackage
}

func loadPackage(targetPackage *cores.Package, folder string) error {
	// packagePlatformTxt, err := properties.SafeLoad(filepath.Join(folder, constants.FILE_PLATFORM_TXT))
	// if err != nil {
	// 	return i18n.WrapError(err)
	// }
	// targetPackage.Properties.Merge(packagePlatformTxt)

	subfolders, err := utils.ReadDirFiltered(folder, utils.FilterDirs)
	if err != nil {
		return i18n.WrapError(err)
	}

	subfolders = utils.FilterOutFoldersByNames(subfolders, constants.FOLDER_TOOLS)

	platforms := targetPackage.Platforms
	for _, subfolder := range subfolders {
		subfolderPath := filepath.Join(folder, subfolder.Name())
		platformId := subfolder.Name()

		_, err := os.Stat(filepath.Join(subfolderPath, constants.FILE_BOARDS_TXT))
		if err != nil && os.IsNotExist(err) {
			theOnlySubfolder, err := utils.TheOnlySubfolderOf(subfolderPath)
			if err != nil {
				return i18n.WrapError(err)
			}

			if theOnlySubfolder != constants.EMPTY_STRING {
				subfolderPath = filepath.Join(subfolderPath, theOnlySubfolder)
			}
		}

		platform := getOrCreatePlatform(platforms, platformId)
		err = loadPlatform(platform.Releases[""], subfolderPath)
		if err != nil {
			return i18n.WrapError(err)
		}
		platforms[platformId] = platform
	}

	return nil
}

func getOrCreatePlatform(platforms map[string]*cores.Platform, platformId string) *cores.Platform {
	if _, ok := platforms[platformId]; ok {
		return platforms[platformId]
	}

	targetPlatform := &cores.Platform{
		Architecture: platformId,
		Releases:     map[string]*cores.PlatformRelease{},
	}
	release := &cores.PlatformRelease{
		Boards:      map[string]*cores.Board{},
		Properties:  properties.Map{},
		Programmers: map[string]properties.Map{},
		Platform:    targetPlatform,
	}
	targetPlatform.Releases[""] = release

	return targetPlatform
}

func loadPlatform(targetPlatform *cores.PlatformRelease, folder string) error {
	_, err := os.Stat(filepath.Join(folder, constants.FILE_BOARDS_TXT))
	if err != nil && !os.IsNotExist(err) {
		return i18n.WrapError(err)
	}

	if os.IsNotExist(err) {
		return nil
	}

	targetPlatform.Folder = folder

	err = loadBoards(targetPlatform.Boards, folder)
	if err != nil {
		return i18n.WrapError(err)
	}

	platformTxt, err := properties.SafeLoad(filepath.Join(folder, constants.FILE_PLATFORM_TXT))
	if err != nil {
		return i18n.WrapError(err)
	}

	localPlatformProperties, err := properties.SafeLoad(filepath.Join(folder, constants.FILE_PLATFORM_LOCAL_TXT))
	if err != nil {
		return i18n.WrapError(err)
	}

	targetPlatform.Properties = targetPlatform.Properties.Clone()
	targetPlatform.Properties.Merge(platformTxt)
	targetPlatform.Properties.Merge(localPlatformProperties)

	programmersProperties, err := properties.SafeLoad(filepath.Join(folder, constants.FILE_PROGRAMMERS_TXT))
	if err != nil {
		return i18n.WrapError(err)
	}
	targetPlatform.Programmers = properties.MergeMapsOfProperties(
		map[string]properties.Map{},
		targetPlatform.Programmers,
		programmersProperties.FirstLevelOf())

	return nil
}

func loadBoards(boards map[string]*cores.Board, folder string) error {
	boardsProperties, err := properties.Load(filepath.Join(folder, constants.FILE_BOARDS_TXT))
	if err != nil {
		return i18n.WrapError(err)
	}

	localProperties, err := properties.SafeLoad(filepath.Join(folder, constants.FILE_BOARDS_LOCAL_TXT))
	if err != nil {
		return i18n.WrapError(err)
	}

	boardsProperties = boardsProperties.Merge(localProperties)

	propertiesByBoardId := boardsProperties.FirstLevelOf()
	delete(propertiesByBoardId, constants.BOARD_PROPERTIES_MENU)

	for boardID, boardProperties := range propertiesByBoardId {
		boardProperties[constants.ID] = boardID
		board := getOrCreateBoard(boards, boardID)
		board.Properties.Merge(boardProperties)
		boards[boardID] = board
	}

	return nil
}

func getOrCreateBoard(boards map[string]*cores.Board, boardId string) *cores.Board {
	if _, ok := boards[boardId]; ok {
		return boards[boardId]
	}

	board := cores.Board{}
	board.BoardId = boardId
	board.Properties = properties.Map{}

	return &board
}
