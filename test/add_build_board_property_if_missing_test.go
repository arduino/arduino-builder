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
	"path/filepath"
	"testing"

	"github.com/arduino/arduino-builder"
	"github.com/arduino/arduino-builder/constants"
	"github.com/arduino/arduino-builder/json_package_index"
	"github.com/arduino/arduino-builder/types"
	"github.com/stretchr/testify/require"
)

func TestAddBuildBoardPropertyIfMissing(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	ctx := &types.Context{
		HardwareFolders: []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware", "user_hardware"},
		FQBN:            "my_avr_platform:avr:mymega",
	}

	commands := []types.Command{
		&builder.HardwareLoader{},
		&builder.TargetBoardResolver{},
		&builder.AddBuildBoardPropertyIfMissing{},
	}

	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}

	targetPackage := ctx.TargetPackage
	require.Equal(t, "my_avr_platform", targetPackage.PackageId)
	targetPlatform := ctx.TargetPlatform
	require.Equal(t, "avr", targetPlatform.PlatformId)
	targetBoard := ctx.TargetBoard
	require.Equal(t, "mymega", targetBoard.BoardId)
	require.Equal(t, constants.EMPTY_STRING, targetBoard.Properties[constants.BUILD_PROPERTIES_BUILD_MCU])
	require.Equal(t, "AVR_MYMEGA", targetBoard.Properties[constants.BUILD_PROPERTIES_BUILD_BOARD])
}

func TestAddBuildBoardPropertyIfMissingNotMissing(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	ctx := &types.Context{
		HardwareFolders: []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware", "user_hardware"},
		FQBN:            "my_avr_platform:avr:mymega:cpu=atmega2560",
	}

	commands := []types.Command{
		&builder.HardwareLoader{},
		&builder.TargetBoardResolver{},
		&builder.AddBuildBoardPropertyIfMissing{},
	}

	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}

	targetPackage := ctx.TargetPackage
	require.Equal(t, "my_avr_platform", targetPackage.PackageId)
	targetPlatform := ctx.TargetPlatform
	require.Equal(t, "avr", targetPlatform.PlatformId)
	targetBoard := ctx.TargetBoard
	require.Equal(t, "mymega", targetBoard.BoardId)
	require.Equal(t, "atmega2560", targetBoard.Properties[constants.BUILD_PROPERTIES_BUILD_MCU])
	require.Equal(t, "AVR_MEGA2560", targetBoard.Properties[constants.BUILD_PROPERTIES_BUILD_BOARD])
}

func TestCoreDependencyProperty(t *testing.T) {
	var paths []string
	path, _ := filepath.Abs(filepath.Join("..", "..", "json_package_index", "testdata"))
	paths = append(paths, path)

	DownloadCoresAndToolsAndLibraries(t)

	ctx := &types.Context{
		HardwareFolders: []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"},
	}

	commands := []types.Command{
		&builder.HardwareLoader{},
	}

	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}

	p, err := json_package_index.PackageIndexFoldersToPropertiesMap(ctx.Hardware, paths)
	require.NoError(t, err)

	version := ctx.Hardware.Packages["arduino"].Platforms["avr"].Properties["version"]

	if json_package_index.CompareVersions(version, "1.6.12") >= 0 {
		require.Equal(t, "{runtime.tools.avr-gcc-4.9.2-atmel3.5.3-arduino2.path}", p["attiny:avr:1.0.2"]["runtime.tools.avr-gcc.path"])
	} else {
		require.Equal(t, "{runtime.tools.avr-gcc-4.8.1-arduino5.path}", p["attiny:avr:1.0.2"]["runtime.tools.avr-gcc.path"])
	}
}
