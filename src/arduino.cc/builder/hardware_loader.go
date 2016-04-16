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
)

type HardwareLoader struct{}

func (s *HardwareLoader) Run(ctx *types.Context) error {
	logger := ctx.GetLogger()

	packages := &types.Packages{}
	packages.Packages = make(map[string]*types.Package)
	packages.Properties = make(map[string]string)

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

		hardwarePlatformTxt, err := props.SafeLoad(filepath.Join(folder, constants.FILE_PLATFORM_TXT), logger)
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
			err = loadPackage(targetPackage, subfolderPath, logger)
			if err != nil {
				return i18n.WrapError(err)
			}
			packages.Packages[packageId] = targetPackage
		}
	}

	ctx.Hardware = packages

	return nil
}

func getOrCreatePackage(packages *types.Packages, packageId string) *types.Package {
	if _, ok := packages.Packages[packageId]; ok {
		return packages.Packages[packageId]
	}

	targetPackage := types.Package{}
	targetPackage.PackageId = packageId
	targetPackage.Platforms = make(map[string]*types.Platform)
	targetPackage.Properties = make(map[string]string)

	return &targetPackage
}

func loadPackage(targetPackage *types.Package, folder string, logger i18n.Logger) error {
	packagePlatformTxt, err := props.SafeLoad(filepath.Join(folder, constants.FILE_PLATFORM_TXT), logger)
	if err != nil {
		return i18n.WrapError(err)
	}
	targetPackage.Properties.Merge(packagePlatformTxt)

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
		err = loadPlatform(platform, targetPackage.PackageId, subfolderPath, logger)
		if err != nil {
			return i18n.WrapError(err)
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
	targetPlatform.Programmers = make(map[string]props.PropertiesMap)

	return &targetPlatform
}

func loadPlatform(targetPlatform *types.Platform, packageId string, folder string, logger i18n.Logger) error {
	_, err := os.Stat(filepath.Join(folder, constants.FILE_BOARDS_TXT))
	if err != nil && !os.IsNotExist(err) {
		return i18n.WrapError(err)
	}

	if os.IsNotExist(err) {
		return nil
	}

	targetPlatform.Folder = folder

	err = loadBoards(targetPlatform.Boards, packageId, targetPlatform.PlatformId, folder, logger)
	if err != nil {
		return i18n.WrapError(err)
	}

	assignDefaultBoardToPlatform(targetPlatform)

	platformTxt, err := props.SafeLoad(filepath.Join(folder, constants.FILE_PLATFORM_TXT), logger)
	if err != nil {
		return i18n.WrapError(err)
	}

	localPlatformProperties, err := props.SafeLoad(filepath.Join(folder, constants.FILE_PLATFORM_LOCAL_TXT), logger)
	if err != nil {
		return i18n.WrapError(err)
	}

	targetPlatform.Properties = targetPlatform.Properties.Clone()
	targetPlatform.Properties.Merge(platformTxt)
	targetPlatform.Properties.Merge(localPlatformProperties)

	programmersProperties, err := props.SafeLoad(filepath.Join(folder, constants.FILE_PROGRAMMERS_TXT), logger)
	if err != nil {
		return i18n.WrapError(err)
	}
	targetPlatform.Programmers = props.MergeMapsOfProperties(make(map[string]props.PropertiesMap), targetPlatform.Programmers, programmersProperties.FirstLevelOf())

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

func loadBoards(boards map[string]*types.Board, packageId string, platformId string, folder string, logger i18n.Logger) error {
	properties, err := props.Load(filepath.Join(folder, constants.FILE_BOARDS_TXT), logger)
	if err != nil {
		return i18n.WrapError(err)
	}

	localProperties, err := props.SafeLoad(filepath.Join(folder, constants.FILE_BOARDS_LOCAL_TXT), logger)
	if err != nil {
		return i18n.WrapError(err)
	}

	properties = properties.Merge(localProperties)

	propertiesByBoardId := properties.FirstLevelOf()
	delete(propertiesByBoardId, constants.BOARD_PROPERTIES_MENU)

	for boardId, properties := range propertiesByBoardId {
		properties[constants.ID] = boardId
		board := getOrCreateBoard(boards, boardId)
		board.Properties.Merge(properties)
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
	board.Properties = make(props.PropertiesMap)

	return &board
}
