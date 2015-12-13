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
	"arduino.cc/builder/types"
	"arduino.cc/builder/utils"
	"strings"
)

type TargetBoardResolver struct{}

func (s *TargetBoardResolver) Run(context map[string]interface{}) error {
	fqbn := context[constants.CTX_FQBN].(string)

	fqbnParts := strings.Split(fqbn, ":")
	targetPackageName := fqbnParts[0]
	targetPlatformName := fqbnParts[1]
	targetBoardName := fqbnParts[2]

	packages := context[constants.CTX_HARDWARE].(*types.Packages)

	targetPackage := packages.Packages[targetPackageName]
	if targetPackage == nil {
		return utils.Errorf(context, constants.MSG_PACKAGE_UNKNOWN, targetPackageName)
	}

	targetPlatform := targetPackage.Platforms[targetPlatformName]
	if targetPlatform == nil {
		return utils.Errorf(context, constants.MSG_PLATFORM_UNKNOWN, targetPlatformName, targetPackageName)
	}

	targetBoard := targetPlatform.Boards[targetBoardName]
	if targetBoard == nil {
		return utils.Errorf(context, constants.MSG_BOARD_UNKNOWN, targetBoardName, targetPlatformName, targetPackageName)
	}

	context[constants.CTX_TARGET_PACKAGE] = targetPackage
	context[constants.CTX_TARGET_PLATFORM] = targetPlatform
	context[constants.CTX_TARGET_BOARD] = targetBoard

	if len(fqbnParts) > 3 {
		addAdditionalPropertiesToTargetBoard(targetBoard, fqbnParts[3])
	}

	core := targetBoard.Properties[constants.BUILD_PROPERTIES_BUILD_CORE]
	if core == constants.EMPTY_STRING {
		core = DEFAULT_BUILD_CORE
	}

	var corePlatform *types.Platform
	coreParts := strings.Split(core, ":")
	if len(coreParts) > 1 {
		core = coreParts[1]
		if packages.Packages[coreParts[0]] == nil {
			return utils.Errorf(context, constants.MSG_MISSING_CORE_FOR_BOARD, coreParts[0])

		}
		corePlatform = packages.Packages[coreParts[0]].Platforms[targetPlatform.PlatformId]
	}

	var actualPlatform *types.Platform
	if corePlatform != nil {
		actualPlatform = corePlatform
	} else {
		actualPlatform = targetPlatform
	}

	context[constants.CTX_BUILD_CORE] = core
	context[constants.CTX_ACTUAL_PLATFORM] = actualPlatform

	return nil
}

func addAdditionalPropertiesToTargetBoard(board *types.Board, options string) {
	optionsParts := strings.Split(options, ",")
	optionsParts = utils.Map(optionsParts, utils.TrimSpace)

	for _, optionPart := range optionsParts {
		parts := strings.Split(optionPart, "=")
		parts = utils.Map(parts, utils.TrimSpace)

		key := parts[0]
		value := parts[1]
		if key != constants.EMPTY_STRING && value != constants.EMPTY_STRING {
			menu := board.Properties.SubTree(constants.BOARD_PROPERTIES_MENU)
			if len(menu) == 0 {
				return
			}
			propertiesOfKey := menu.SubTree(key)
			if len(propertiesOfKey) == 0 {
				return
			}
			propertiesOfKeyValue := propertiesOfKey.SubTree(value)
			board.Properties.Merge(propertiesOfKeyValue)
		}
	}
}
