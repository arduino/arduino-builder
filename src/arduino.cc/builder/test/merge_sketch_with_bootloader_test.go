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
	"arduino.cc/builder"
	"arduino.cc/builder/constants"
	"arduino.cc/builder/types"
	"arduino.cc/builder/utils"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMergeSketchWithBootloader(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	ctx := &types.Context{
		HardwareFolders:         []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"},
		ToolsFolders:            []string{"downloaded_tools"},
		BuiltInLibrariesFolders: []string{"downloaded_libraries"},
		OtherLibrariesFolders:   []string{"libraries"},
		SketchLocation:          filepath.Join("sketch1", "sketch.ino"),
		FQBN:                    "arduino:avr:uno",
		ArduinoAPIVersion:       "10600",
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	err := utils.EnsureFolderExists(filepath.Join(buildPath, "sketch"))
	NoError(t, err)

	fakeSketchHex := "row 1\n" +
		"row 2\n"
	err = utils.WriteFile(filepath.Join(buildPath, "sketch", "sketch.ino.hex"), fakeSketchHex)
	NoError(t, err)

	commands := []types.Command{
		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},
		&builder.MergeSketchWithBootloader{},
	}

	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}

	bytes, err := ioutil.ReadFile(filepath.Join(buildPath, "sketch", "sketch.ino.with_bootloader.hex"))
	NoError(t, err)
	mergedSketchHex := string(bytes)

	require.True(t, strings.HasPrefix(mergedSketchHex, "row 1\n:107E0000112484B714BE81FFF0D085E080938100F7\n"))
	require.True(t, strings.HasSuffix(mergedSketchHex, ":0400000300007E007B\n:00000001FF\n"))
}

func TestMergeSketchWithBootloaderSketchInBuildPath(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	ctx := &types.Context{
		HardwareFolders:         []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"},
		ToolsFolders:            []string{"downloaded_tools"},
		BuiltInLibrariesFolders: []string{"downloaded_libraries"},
		OtherLibrariesFolders:   []string{"libraries"},
		SketchLocation:          filepath.Join("sketch1", "sketch.ino"),
		FQBN:                    "arduino:avr:uno",
		ArduinoAPIVersion:       "10600",
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	err := utils.EnsureFolderExists(filepath.Join(buildPath, "sketch"))
	NoError(t, err)

	fakeSketchHex := "row 1\n" +
		"row 2\n"
	err = utils.WriteFile(filepath.Join(buildPath, "sketch.ino.hex"), fakeSketchHex)
	NoError(t, err)

	commands := []types.Command{
		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},
		&builder.MergeSketchWithBootloader{},
	}

	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}

	bytes, err := ioutil.ReadFile(filepath.Join(buildPath, "sketch.ino.with_bootloader.hex"))
	NoError(t, err)
	mergedSketchHex := string(bytes)

	require.True(t, strings.HasPrefix(mergedSketchHex, "row 1\n:107E0000112484B714BE81FFF0D085E080938100F7\n"))
	require.True(t, strings.HasSuffix(mergedSketchHex, ":0400000300007E007B\n:00000001FF\n"))
}

func TestMergeSketchWithBootloaderWhenNoBootloaderAvailable(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	ctx := &types.Context{
		HardwareFolders:         []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"},
		ToolsFolders:            []string{"downloaded_tools"},
		BuiltInLibrariesFolders: []string{"downloaded_libraries"},
		OtherLibrariesFolders:   []string{"libraries"},
		SketchLocation:          filepath.Join("sketch1", "sketch.ino"),
		FQBN:                    "arduino:avr:uno",
		ArduinoAPIVersion:       "10600",
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	commands := []types.Command{
		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},
	}

	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}

	buildProperties := ctx.BuildProperties
	delete(buildProperties, constants.BUILD_PROPERTIES_BOOTLOADER_NOBLINK)
	delete(buildProperties, constants.BUILD_PROPERTIES_BOOTLOADER_FILE)

	command := &builder.MergeSketchWithBootloader{}
	err := command.Run(ctx)
	NoError(t, err)

	_, err = os.Stat(filepath.Join(buildPath, "sketch.ino.with_bootloader.hex"))
	require.Error(t, err)
	require.True(t, os.IsNotExist(err))
}

func TestMergeSketchWithBootloaderPathIsParameterized(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	ctx := &types.Context{
		HardwareFolders:         []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware", "user_hardware"},
		ToolsFolders:            []string{"downloaded_tools"},
		BuiltInLibrariesFolders: []string{"downloaded_libraries"},
		OtherLibrariesFolders:   []string{"libraries"},
		SketchLocation:          filepath.Join("sketch1", "sketch.ino"),
		FQBN:                    "my_avr_platform:avr:mymega:cpu=atmega2560",
		ArduinoAPIVersion:       "10600",
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	err := utils.EnsureFolderExists(filepath.Join(buildPath, "sketch"))
	NoError(t, err)

	fakeSketchHex := "row 1\n" +
		"row 2\n"
	err = utils.WriteFile(filepath.Join(buildPath, "sketch", "sketch.ino.hex"), fakeSketchHex)
	NoError(t, err)

	commands := []types.Command{
		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},
		&builder.MergeSketchWithBootloader{},
	}

	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}

	bytes, err := ioutil.ReadFile(filepath.Join(buildPath, "sketch", "sketch.ino.with_bootloader.hex"))
	NoError(t, err)
	mergedSketchHex := string(bytes)

	require.True(t, strings.HasPrefix(mergedSketchHex, "row 1\n:020000023000CC"))
	require.True(t, strings.HasSuffix(mergedSketchHex, ":040000033000E000E9\n:00000001FF\n"))
}
