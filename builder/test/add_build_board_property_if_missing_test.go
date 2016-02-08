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

package test

import (
	"github.com/arduino/arduino-builder/builder"
	"github.com/arduino/arduino-builder/builder/constants"
	"github.com/arduino/arduino-builder/builder/types"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func TestAddBuildBoardPropertyIfMissing(t *testing.T) {
	context := make(map[string]interface{})
	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware", "user_hardware"}
	context[constants.CTX_FQBN] = "my_avr_platform:avr:mymega"

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},
		&builder.HardwareLoader{},
		&builder.TargetBoardResolver{},
		&builder.AddBuildBoardPropertyIfMissing{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	targetPackage := context[constants.CTX_TARGET_PACKAGE].(*types.Package)
	require.Equal(t, "my_avr_platform", targetPackage.PackageId)
	targetPlatform := context[constants.CTX_TARGET_PLATFORM].(*types.Platform)
	require.Equal(t, "avr", targetPlatform.PlatformId)
	targetBoard := context[constants.CTX_TARGET_BOARD].(*types.Board)
	require.Equal(t, "mymega", targetBoard.BoardId)
	require.Equal(t, constants.EMPTY_STRING, targetBoard.Properties[constants.BUILD_PROPERTIES_BUILD_MCU])
	require.Equal(t, "AVR_MYMEGA", targetBoard.Properties[constants.BUILD_PROPERTIES_BUILD_BOARD])
}

func TestAddBuildBoardPropertyIfMissingNotMissing(t *testing.T) {
	context := make(map[string]interface{})
	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware", "user_hardware"}
	context[constants.CTX_FQBN] = "my_avr_platform:avr:mymega:cpu=atmega2560"

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},
		&builder.HardwareLoader{},
		&builder.TargetBoardResolver{},
		&builder.AddBuildBoardPropertyIfMissing{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	targetPackage := context[constants.CTX_TARGET_PACKAGE].(*types.Package)
	require.Equal(t, "my_avr_platform", targetPackage.PackageId)
	targetPlatform := context[constants.CTX_TARGET_PLATFORM].(*types.Platform)
	require.Equal(t, "avr", targetPlatform.PlatformId)
	targetBoard := context[constants.CTX_TARGET_BOARD].(*types.Board)
	require.Equal(t, "mymega", targetBoard.BoardId)
	require.Equal(t, "atmega2560", targetBoard.Properties[constants.BUILD_PROPERTIES_BUILD_MCU])
	require.Equal(t, "AVR_MEGA2560", targetBoard.Properties[constants.BUILD_PROPERTIES_BUILD_BOARD])
}
