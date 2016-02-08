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

func TestTargetBoardResolverUno(t *testing.T) {
	context := make(map[string]interface{})
	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_FQBN] = "arduino:avr:uno"

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},
		&builder.HardwareLoader{},
		&builder.TargetBoardResolver{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	targetPackage := context[constants.CTX_TARGET_PACKAGE].(*types.Package)
	require.Equal(t, "arduino", targetPackage.PackageId)
	targetPlatform := context[constants.CTX_TARGET_PLATFORM].(*types.Platform)
	require.Equal(t, "avr", targetPlatform.PlatformId)
	targetBoard := context[constants.CTX_TARGET_BOARD].(*types.Board)
	require.Equal(t, "uno", targetBoard.BoardId)
	require.Equal(t, "atmega328p", targetBoard.Properties[constants.BUILD_PROPERTIES_BUILD_MCU])
}

func TestTargetBoardResolverDue(t *testing.T) {
	context := make(map[string]interface{})
	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_FQBN] = "arduino:sam:arduino_due_x"

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},
		&builder.HardwareLoader{},
		&builder.TargetBoardResolver{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	targetPackage := context[constants.CTX_TARGET_PACKAGE].(*types.Package)
	require.Equal(t, "arduino", targetPackage.PackageId)
	targetPlatform := context[constants.CTX_TARGET_PLATFORM].(*types.Platform)
	require.Equal(t, "sam", targetPlatform.PlatformId)
	targetBoard := context[constants.CTX_TARGET_BOARD].(*types.Board)
	require.Equal(t, "arduino_due_x", targetBoard.BoardId)
	require.Equal(t, "cortex-m3", targetBoard.Properties[constants.BUILD_PROPERTIES_BUILD_MCU])
}

func TestTargetBoardResolverMega1280(t *testing.T) {
	context := make(map[string]interface{})
	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_FQBN] = "arduino:avr:mega:cpu=atmega1280"

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},
		&builder.HardwareLoader{},
		&builder.TargetBoardResolver{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	targetPackage := context[constants.CTX_TARGET_PACKAGE].(*types.Package)
	require.Equal(t, "arduino", targetPackage.PackageId)
	targetPlatform := context[constants.CTX_TARGET_PLATFORM].(*types.Platform)
	require.Equal(t, "avr", targetPlatform.PlatformId)
	targetBoard := context[constants.CTX_TARGET_BOARD].(*types.Board)
	require.Equal(t, "mega", targetBoard.BoardId)
	require.Equal(t, "atmega1280", targetBoard.Properties[constants.BUILD_PROPERTIES_BUILD_MCU])
	require.Equal(t, "AVR_MEGA", targetBoard.Properties[constants.BUILD_PROPERTIES_BUILD_BOARD])
}

func TestTargetBoardResolverMega2560(t *testing.T) {
	context := make(map[string]interface{})
	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_FQBN] = "arduino:avr:mega:cpu=atmega2560"

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},
		&builder.HardwareLoader{},
		&builder.TargetBoardResolver{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	targetPackage := context[constants.CTX_TARGET_PACKAGE].(*types.Package)
	require.Equal(t, "arduino", targetPackage.PackageId)
	targetPlatform := context[constants.CTX_TARGET_PLATFORM].(*types.Platform)
	require.Equal(t, "avr", targetPlatform.PlatformId)
	targetBoard := context[constants.CTX_TARGET_BOARD].(*types.Board)
	require.Equal(t, "mega", targetBoard.BoardId)
	require.Equal(t, "atmega2560", targetBoard.Properties[constants.BUILD_PROPERTIES_BUILD_MCU])
	require.Equal(t, "AVR_MEGA2560", targetBoard.Properties[constants.BUILD_PROPERTIES_BUILD_BOARD])
}

func TestTargetBoardResolverCustomYun(t *testing.T) {
	context := make(map[string]interface{})
	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware", "user_hardware"}
	context[constants.CTX_FQBN] = "my_avr_platform:avr:custom_yun"

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},
		&builder.HardwareLoader{},
		&builder.TargetBoardResolver{},
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
	require.Equal(t, "custom_yun", targetBoard.BoardId)
	require.Equal(t, "atmega32u4", targetBoard.Properties[constants.BUILD_PROPERTIES_BUILD_MCU])
	require.Equal(t, "AVR_YUN", targetBoard.Properties[constants.BUILD_PROPERTIES_BUILD_BOARD])
}

func TestTargetBoardResolverCustomCore(t *testing.T) {
	context := make(map[string]interface{})
	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware", "user_hardware"}
	context[constants.CTX_FQBN] = "watterott:avr:attiny841:core=spencekonde,info=info"

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},
		&builder.HardwareLoader{},
		&builder.TargetBoardResolver{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	targetPackage := context[constants.CTX_TARGET_PACKAGE].(*types.Package)
	require.Equal(t, "watterott", targetPackage.PackageId)
	targetPlatform := context[constants.CTX_TARGET_PLATFORM].(*types.Platform)
	require.Equal(t, "avr", targetPlatform.PlatformId)
	targetBoard := context[constants.CTX_TARGET_BOARD].(*types.Board)
	require.Equal(t, "attiny841", targetBoard.BoardId)
	require.Equal(t, "tiny841", context[constants.CTX_BUILD_CORE].(string))
	require.Equal(t, "tiny14", targetBoard.Properties[constants.BUILD_PROPERTIES_BUILD_VARIANT])
}
