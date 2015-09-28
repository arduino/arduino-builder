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
	"os"
	"path/filepath"
)

type HardwareLoader struct{}

func (s *HardwareLoader) Run(context map[string]interface{}) error {
	mainHardwarePlatformTxt := make(map[string]string)

	packages := &types.Packages{}
	packages.Packages = make(map[string]*types.Package)
	packages.Properties = make(map[string]string)

	folders := context[constants.CTX_HARDWARE_FOLDERS].([]string)
	folders, err := utils.AbsolutizePaths(folders)
	if err != nil {
		return utils.WrapError(err)
	}

	for _, folder := range folders {
		stat, err := os.Stat(folder)
		if err != nil {
			return utils.WrapError(err)
		}
		if !stat.IsDir() {
			return utils.Errorf(context, constants.MSG_MUST_BE_A_FOLDER, folder)
		}

		if len(mainHardwarePlatformTxt) == 0 {
			mainHardwarePlatformTxt, err = props.SafeLoad(filepath.Join(folder, constants.FILE_PLATFORM_TXT))
			if err != nil {
				return utils.WrapError(err)
			}
		}
		hardwarePlatformTxt, err := props.SafeLoad(filepath.Join(folder, constants.FILE_PLATFORM_TXT))
		if err != nil {
			return utils.WrapError(err)
		}
		hardwarePlatformTxt = utils.MergeMapsOfStrings(make(map[string]string), mainHardwarePlatformTxt, hardwarePlatformTxt)

		subfolders, err := utils.ReadDirFiltered(folder, utils.FilterDirs)
		if err != nil {
			return utils.WrapError(err)
		}
		subfolders = utils.FilterOutFoldersByNames(subfolders, constants.FOLDER_TOOLS)

		for _, subfolder := range subfolders {
			subfolderPath := filepath.Join(folder, subfolder.Name())
			packageId := subfolder.Name()

			if _, err := os.Stat(filepath.Join(subfolderPath, constants.FOLDER_HARDWARE)); err == nil {
				subfolderPath = filepath.Join(subfolderPath, constants.FOLDER_HARDWARE)
			}

			targetPackage := getOrCreatePackage(packages, packageId)
			err = loadPackage(targetPackage, subfolderPath, hardwarePlatformTxt)
			if err != nil {
				return utils.WrapError(err)
			}
			packages.Packages[packageId] = targetPackage
		}
	}

	context[constants.CTX_HARDWARE] = packages

	return nil
}

func getOrCreatePackage(packages *types.Packages, packageId string) *types.Package {
	if _, ok := packages.Packages[packageId]; ok {
		return packages.Packages[packageId]
	}

	targetPackage := types.Package{}
	targetPackage.PackageId = packageId
	targetPackage.Platforms = make(map[string]*types.Platform)

	return &targetPackage
}

func loadPackage(targetPackage *types.Package, folder string, hardwarePlatformTxt map[string]string) error {
	packagePlatformTxt, err := props.SafeLoad(filepath.Join(folder, constants.FILE_PLATFORM_TXT))
	if err != nil {
		return utils.WrapError(err)
	}
	packagePlatformTxt = utils.MergeMapsOfStrings(make(map[string]string), hardwarePlatformTxt, packagePlatformTxt)

	subfolders, err := utils.ReadDirFiltered(folder, utils.FilterDirs)
	if err != nil {
		return utils.WrapError(err)
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
				return utils.WrapError(err)
			}

			if theOnlySubfolder != constants.EMPTY_STRING {
				subfolderPath = filepath.Join(subfolderPath, theOnlySubfolder)
			}
		}

		platform := getOrCreatePlatform(platforms, platformId)
		err = loadPlatform(platform, targetPackage.PackageId, subfolderPath, packagePlatformTxt)
		if err != nil {
			return utils.WrapError(err)
		}
		platforms[platformId] = platform
	}

	return nil
}

func getOrCreatePlatform(platforms map[string]*types.Platform, platformId string) *types.Platform {
	if _, ok := platforms[platformId]; ok {
		return platforms[platformId]
	}

	targetPlatform := types.Platform{}
	targetPlatform.PlatformId = platformId
	targetPlatform.Boards = make(map[string]*types.Board)
	targetPlatform.Properties = make(map[string]string)
	targetPlatform.Programmers = make(map[string]map[string]string)

	return &targetPlatform
}

func loadPlatform(targetPlatform *types.Platform, packageId string, folder string, packagePlatformTxt map[string]string) error {
	_, err := os.Stat(filepath.Join(folder, constants.FILE_BOARDS_TXT))
	if err != nil && !os.IsNotExist(err) {
		return utils.WrapError(err)
	}

	if os.IsNotExist(err) {
		return nil
	}

	targetPlatform.Folder = folder

	err = loadBoards(targetPlatform.Boards, packageId, targetPlatform.PlatformId, folder)
	if err != nil {
		return utils.WrapError(err)
	}

	assignDefaultBoardToPlatform(targetPlatform)

	platformProperties, err := props.SafeLoad(filepath.Join(folder, constants.FILE_PLATFORM_TXT))
	if err != nil {
		return utils.WrapError(err)
	}

	localPlatformProperties, err := props.SafeLoad(filepath.Join(folder, constants.FILE_PLATFORM_LOCAL_TXT))
	if err != nil {
		return utils.WrapError(err)
	}

	targetPlatform.Properties = utils.MergeMapsOfStrings(make(map[string]string), targetPlatform.Properties, packagePlatformTxt, platformProperties, localPlatformProperties)

	programmersProperties, err := props.SafeLoad(filepath.Join(folder, constants.FILE_PROGRAMMERS_TXT))
	if err != nil {
		return utils.WrapError(err)
	}
	targetPlatform.Programmers = utils.MergeMapsOfMapsOfStrings(make(map[string]map[string]string), targetPlatform.Programmers, props.FirstLevelOf(programmersProperties))

	return nil
}

func assignDefaultBoardToPlatform(targetPlatform *types.Platform) {
	if targetPlatform.DefaultBoard == nil {
		for _, board := range targetPlatform.Boards {
			if targetPlatform.DefaultBoard == nil {
				targetPlatform.DefaultBoard = board
			}
		}
	}
}

func loadBoards(boards map[string]*types.Board, packageId string, platformId string, folder string) error {
	properties, err := props.Load(filepath.Join(folder, constants.FILE_BOARDS_TXT))
	if err != nil {
		return utils.WrapError(err)
	}

	propertiesByBoardId := props.FirstLevelOf(properties)
	delete(propertiesByBoardId, constants.BOARD_PROPERTIES_MENU)

	for boardId, properties := range propertiesByBoardId {
		properties[constants.ID] = boardId
		board := getOrCreateBoard(boards, boardId)
		board.Properties = utils.MergeMapsOfStrings(board.Properties, properties)
		boards[boardId] = board
	}

	return nil
}

func getOrCreateBoard(boards map[string]*types.Board, boardId string) *types.Board {
	if _, ok := boards[boardId]; ok {
		return boards[boardId]
	}

	board := types.Board{}
	board.BoardId = boardId
	board.Properties = make(map[string]string)

	return &board
}
