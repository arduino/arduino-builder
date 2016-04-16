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
	"github.com/stretchr/testify/require"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestBuilderEmptySketch(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	ctx := &types.Context{
		HardwareFolders:         []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"},
		ToolsFolders:            []string{"downloaded_tools"},
		BuiltInLibrariesFolders: []string{"downloaded_libraries"},
		OtherLibrariesFolders:   []string{"libraries"},
		SketchLocation:          filepath.Join("sketch1", "sketch.ino"),
		FQBN:                    "arduino:avr:uno",
		ArduinoAPIVersion:       "10600",
		Verbose:                 true,
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	ctx.DebugLevel = 10

	command := builder.Builder{}
	err := command.Run(ctx)
	NoError(t, err)

	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_CORE, "HardwareSerial.cpp.o"))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_PREPROC, constants.FILE_CTAGS_TARGET_FOR_GCC_MINUS_E))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_SKETCH, "sketch.ino.cpp.o"))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, "sketch.ino.elf"))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, "sketch.ino.hex"))
	NoError(t, err)
}

func TestBuilderBridge(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	ctx := &types.Context{
		HardwareFolders:         []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"},
		ToolsFolders:            []string{"downloaded_tools"},
		BuiltInLibrariesFolders: []string{"downloaded_libraries"},
		OtherLibrariesFolders:   []string{"libraries"},
		SketchLocation:          filepath.Join("downloaded_libraries", "Bridge", "examples", "Bridge", "Bridge.ino"),
		FQBN:                    "arduino:avr:leonardo",
		ArduinoAPIVersion:       "10600",
		Verbose:                 true,
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	command := builder.Builder{}
	err := command.Run(ctx)
	NoError(t, err)

	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_CORE, "HardwareSerial.cpp.o"))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_PREPROC, constants.FILE_CTAGS_TARGET_FOR_GCC_MINUS_E))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_SKETCH, "Bridge.ino.cpp.o"))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, "Bridge.ino.elf"))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, "Bridge.ino.hex"))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_LIBRARIES, "Bridge", "Mailbox.cpp.o"))
	NoError(t, err)
}

func TestBuilderSketchWithConfig(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	ctx := &types.Context{
		HardwareFolders:         []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"},
		ToolsFolders:            []string{"downloaded_tools"},
		BuiltInLibrariesFolders: []string{"downloaded_libraries"},
		OtherLibrariesFolders:   []string{"libraries"},
		SketchLocation:          filepath.Join("sketch_with_config", "sketch_with_config.ino"),
		FQBN:                    "arduino:avr:leonardo",
		ArduinoAPIVersion:       "10600",
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	command := builder.Builder{}
	err := command.Run(ctx)
	NoError(t, err)

	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_CORE, "HardwareSerial.cpp.o"))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_PREPROC, constants.FILE_CTAGS_TARGET_FOR_GCC_MINUS_E))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_SKETCH, "sketch_with_config.ino.cpp.o"))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, "sketch_with_config.ino.elf"))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, "sketch_with_config.ino.hex"))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_LIBRARIES, "Bridge", "Mailbox.cpp.o"))
	NoError(t, err)
}

func TestBuilderBridgeTwice(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	ctx := &types.Context{
		HardwareFolders:         []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"},
		ToolsFolders:            []string{"downloaded_tools"},
		BuiltInLibrariesFolders: []string{"downloaded_libraries"},
		OtherLibrariesFolders:   []string{"libraries"},
		SketchLocation:          filepath.Join("downloaded_libraries", "Bridge", "examples", "Bridge", "Bridge.ino"),
		FQBN:                    "arduino:avr:leonardo",
		ArduinoAPIVersion:       "10600",
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	command := builder.Builder{}
	err := command.Run(ctx)
	NoError(t, err)

	command = builder.Builder{}
	err = command.Run(ctx)
	NoError(t, err)

	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_CORE, "HardwareSerial.cpp.o"))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_PREPROC, constants.FILE_CTAGS_TARGET_FOR_GCC_MINUS_E))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_SKETCH, "Bridge.ino.cpp.o"))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, "Bridge.ino.elf"))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, "Bridge.ino.hex"))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_LIBRARIES, "Bridge", "Mailbox.cpp.o"))
	NoError(t, err)
}

func TestBuilderBridgeSAM(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	ctx := &types.Context{
		HardwareFolders:         []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"},
		ToolsFolders:            []string{"downloaded_tools"},
		BuiltInLibrariesFolders: []string{"downloaded_libraries"},
		OtherLibrariesFolders:   []string{"libraries"},
		SketchLocation:          filepath.Join("downloaded_libraries", "Bridge", "examples", "Bridge", "Bridge.ino"),
		FQBN:                    "arduino:sam:arduino_due_x_dbg",
		ArduinoAPIVersion:       "10600",
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	ctx.WarningsLevel = "all"

	command := builder.Builder{}
	err := command.Run(ctx)
	NoError(t, err)

	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_CORE, "syscalls_sam3.c.o"))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_CORE, "USB", "PluggableUSB.cpp.o"))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_CORE, "avr", "dtostrf.c.d"))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_PREPROC, constants.FILE_CTAGS_TARGET_FOR_GCC_MINUS_E))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_SKETCH, "Bridge.ino.cpp.o"))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, "Bridge.ino.elf"))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, "Bridge.ino.bin"))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_LIBRARIES, "Bridge", "Mailbox.cpp.o"))
	NoError(t, err)

	cmd := exec.Command(filepath.Join("downloaded_tools", "arm-none-eabi-gcc", "4.8.3-2014q1", "bin", "arm-none-eabi-objdump"), "-f", filepath.Join(buildPath, constants.FOLDER_CORE, "core.a"))
	bytes, err := cmd.CombinedOutput()
	NoError(t, err)
	require.NotContains(t, string(bytes), "variant.cpp.o")
}

func TestBuilderBridgeRedBearLab(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	ctx := &types.Context{
		HardwareFolders:         []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware", "downloaded_board_manager_stuff"},
		ToolsFolders:            []string{"downloaded_tools", "downloaded_board_manager_stuff"},
		BuiltInLibrariesFolders: []string{"downloaded_libraries"},
		OtherLibrariesFolders:   []string{"libraries"},
		SketchLocation:          filepath.Join("downloaded_libraries", "Bridge", "examples", "Bridge", "Bridge.ino"),
		FQBN:                    "RedBearLab:avr:blend",
		ArduinoAPIVersion:       "10600",
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	command := builder.Builder{}
	err := command.Run(ctx)
	NoError(t, err)

	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_CORE, "HardwareSerial.cpp.o"))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_PREPROC, constants.FILE_CTAGS_TARGET_FOR_GCC_MINUS_E))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_SKETCH, "Bridge.ino.cpp.o"))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, "Bridge.ino.elf"))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, "Bridge.ino.hex"))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_LIBRARIES, "Bridge", "Mailbox.cpp.o"))
	NoError(t, err)
}

func TestBuilderSketchNoFunctions(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	ctx := &types.Context{
		HardwareFolders:         []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware", "downloaded_board_manager_stuff"},
		ToolsFolders:            []string{"downloaded_tools", "downloaded_board_manager_stuff"},
		BuiltInLibrariesFolders: []string{"downloaded_libraries"},
		OtherLibrariesFolders:   []string{"libraries"},
		SketchLocation:          filepath.Join("sketch_no_functions", "main.ino"),
		FQBN:                    "RedBearLab:avr:blend",
		ArduinoAPIVersion:       "10600",
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	command := builder.Builder{}
	err := command.Run(ctx)
	require.Error(t, err)
}

func TestBuilderSketchWithBackup(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	ctx := &types.Context{
		HardwareFolders:         []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware", "downloaded_board_manager_stuff"},
		ToolsFolders:            []string{"downloaded_tools", "downloaded_board_manager_stuff"},
		BuiltInLibrariesFolders: []string{"downloaded_libraries"},
		OtherLibrariesFolders:   []string{"libraries"},
		SketchLocation:          filepath.Join("sketch_with_backup_files", "sketch.ino"),
		FQBN:                    "arduino:avr:uno",
		ArduinoAPIVersion:       "10600",
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	command := builder.Builder{}
	err := command.Run(ctx)
	NoError(t, err)
}

func TestBuilderSketchWithOldLib(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	ctx := &types.Context{
		HardwareFolders:         []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"},
		ToolsFolders:            []string{"downloaded_tools"},
		BuiltInLibrariesFolders: []string{"downloaded_libraries"},
		OtherLibrariesFolders:   []string{"libraries"},
		SketchLocation:          filepath.Join("sketch_with_old_lib", "sketch.ino"),
		FQBN:                    "arduino:avr:uno",
		ArduinoAPIVersion:       "10600",
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	command := builder.Builder{}
	err := command.Run(ctx)
	NoError(t, err)
}

func TestBuilderSketchWithSubfolders(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	ctx := &types.Context{
		HardwareFolders:         []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"},
		ToolsFolders:            []string{"downloaded_tools"},
		BuiltInLibrariesFolders: []string{"downloaded_libraries"},
		OtherLibrariesFolders:   []string{"libraries"},
		SketchLocation:          filepath.Join("sketch_with_subfolders", "sketch_with_subfolders.ino"),
		FQBN:                    "arduino:avr:uno",
		ArduinoAPIVersion:       "10600",
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	command := builder.Builder{}
	err := command.Run(ctx)
	NoError(t, err)
}

func TestBuilderSketchBuildPathContainsUnusedPreviouslyCompiledLibrary(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	ctx := &types.Context{
		HardwareFolders:         []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"},
		ToolsFolders:            []string{"downloaded_tools"},
		BuiltInLibrariesFolders: []string{"downloaded_libraries"},
		OtherLibrariesFolders:   []string{"libraries"},
		SketchLocation:          filepath.Join("downloaded_libraries", "Bridge", "examples", "Bridge", "Bridge.ino"),
		FQBN:                    "arduino:avr:leonardo",
		ArduinoAPIVersion:       "10600",
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	NoError(t, os.MkdirAll(filepath.Join(buildPath, constants.FOLDER_LIBRARIES, "SPI"), os.FileMode(0755)))

	command := builder.Builder{}
	err := command.Run(ctx)
	NoError(t, err)

	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_LIBRARIES, "SPI"))
	require.Error(t, err)
	require.True(t, os.IsNotExist(err))
	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_LIBRARIES, "Bridge"))
	NoError(t, err)
}
