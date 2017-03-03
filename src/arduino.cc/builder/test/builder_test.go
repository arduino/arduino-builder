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
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"arduino.cc/builder"
	"arduino.cc/builder/builder_utils"
	"arduino.cc/builder/constants"
	"arduino.cc/builder/types"
	"github.com/stretchr/testify/require"
)

func prepareBuilderTestContext(sketchPath, fqbn string) *types.Context {
	return &types.Context{
		SketchLocation:          sketchPath,
		FQBN:                    fqbn,
		HardwareFolders:         []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"},
		ToolsFolders:            []string{"downloaded_tools"},
		BuiltInLibrariesFolders: []string{"downloaded_libraries"},
		OtherLibrariesFolders:   []string{"libraries"},
		ArduinoAPIVersion:       "10600",
		Verbose:                 false,
	}
}

func TestBuilderEmptySketch(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	ctx := prepareBuilderTestContext(filepath.Join("sketch1", "sketch.ino"), "arduino:avr:uno")
	ctx.DebugLevel = 10

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	// Run builder
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

	ctx := prepareBuilderTestContext(filepath.Join("downloaded_libraries", "Bridge", "examples", "Bridge", "Bridge.ino"), "arduino:avr:leonardo")

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	// Run builder
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

	ctx := prepareBuilderTestContext(filepath.Join("sketch_with_config", "sketch_with_config.ino"), "arduino:avr:leonardo")

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	// Run builder
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

	ctx := prepareBuilderTestContext(filepath.Join("downloaded_libraries", "Bridge", "examples", "Bridge", "Bridge.ino"), "arduino:avr:leonardo")

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	// Run builder
	command := builder.Builder{}
	err := command.Run(ctx)
	NoError(t, err)

	// Run builder again
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

	ctx := prepareBuilderTestContext(filepath.Join("downloaded_libraries", "Bridge", "examples", "Bridge", "Bridge.ino"), "arduino:sam:arduino_due_x_dbg")
	ctx.WarningsLevel = "all"

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	// Run builder
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

	ctx := prepareBuilderTestContext(filepath.Join("downloaded_libraries", "Bridge", "examples", "Bridge", "Bridge.ino"), "RedBearLab:avr:blend")
	ctx.HardwareFolders = append(ctx.HardwareFolders, "downloaded_board_manager_stuff")
	ctx.ToolsFolders = append(ctx.ToolsFolders, "downloaded_board_manager_stuff")

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	// Run builder
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

	ctx := prepareBuilderTestContext(filepath.Join("sketch_no_functions", "main.ino"), "RedBearLab:avr:blend")
	ctx.HardwareFolders = append(ctx.HardwareFolders, "downloaded_board_manager_stuff")
	ctx.ToolsFolders = append(ctx.ToolsFolders, "downloaded_board_manager_stuff")

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	// Run builder
	command := builder.Builder{}
	err := command.Run(ctx)
	require.Error(t, err)
}

func TestBuilderSketchWithBackup(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	ctx := prepareBuilderTestContext(filepath.Join("sketch_with_backup_files", "sketch.ino"), "arduino:avr:uno")
	ctx.HardwareFolders = append(ctx.HardwareFolders, "downloaded_board_manager_stuff")
	ctx.ToolsFolders = append(ctx.ToolsFolders, "downloaded_board_manager_stuff")

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	// Run builder
	command := builder.Builder{}
	err := command.Run(ctx)
	NoError(t, err)
}

func TestBuilderSketchWithOldLib(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	ctx := prepareBuilderTestContext(filepath.Join("sketch_with_old_lib", "sketch.ino"), "arduino:avr:uno")

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	// Run builder
	command := builder.Builder{}
	err := command.Run(ctx)
	NoError(t, err)
}

func TestBuilderSketchWithSubfolders(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	ctx := prepareBuilderTestContext(filepath.Join("sketch_with_subfolders", "sketch_with_subfolders.ino"), "arduino:avr:uno")

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	// Run builder
	command := builder.Builder{}
	err := command.Run(ctx)
	NoError(t, err)
}

func TestBuilderSketchBuildPathContainsUnusedPreviouslyCompiledLibrary(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	ctx := prepareBuilderTestContext(filepath.Join("downloaded_libraries", "Bridge", "examples", "Bridge", "Bridge.ino"), "arduino:avr:leonardo")

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	NoError(t, os.MkdirAll(filepath.Join(buildPath, constants.FOLDER_LIBRARIES, "SPI"), os.FileMode(0755)))

	// Run builder
	command := builder.Builder{}
	err := command.Run(ctx)
	NoError(t, err)

	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_LIBRARIES, "SPI"))
	require.Error(t, err)
	require.True(t, os.IsNotExist(err))
	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_LIBRARIES, "Bridge"))
	NoError(t, err)
}

func TestBuilderWithBuildPathInSketchDir(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	ctx := prepareBuilderTestContext(filepath.Join("sketch1", "sketch.ino"), "arduino:avr:uno")

	var err error
	ctx.BuildPath, err = filepath.Abs(filepath.Join("sketch1", "build"))
	NoError(t, err)
	defer os.RemoveAll(ctx.BuildPath)

	// Run build
	command := builder.Builder{}
	err = command.Run(ctx)
	NoError(t, err)

	// Run build twice, to verify the build still works when the
	// build directory is present at the start
	err = command.Run(ctx)
	NoError(t, err)
}

func TestBuilderCacheCoreAFile(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	ctx := prepareBuilderTestContext(filepath.Join("sketch1", "sketch.ino"), "arduino:avr:uno")

	SetupBuildPath(t, ctx)
	defer os.RemoveAll(ctx.BuildPath)
	SetupBuildCachePath(t, ctx)
	defer os.RemoveAll(ctx.BuildCachePath)

	// Run build
	bldr := builder.Builder{}
	err := bldr.Run(ctx)
	NoError(t, err)

	// Pick timestamp of cached core
	coreFolder := filepath.Join("downloaded_hardware", "arduino", "avr")
	coreFileName := builder_utils.GetCachedCoreArchiveFileName(ctx.FQBN, coreFolder)
	cachedCoreFile := filepath.Join(ctx.CoreBuildCachePath, coreFileName)
	coreStatBefore, err := os.Stat(cachedCoreFile)
	require.NoError(t, err)

	// Run build again, to verify that the builder skips rebuilding core.a
	err = bldr.Run(ctx)
	NoError(t, err)

	coreStatAfterRebuild, err := os.Stat(cachedCoreFile)
	require.NoError(t, err)
	require.Equal(t, coreStatBefore.ModTime(), coreStatAfterRebuild.ModTime())

	// Touch a file of the core and check if the builder invalidate the cache
	time.Sleep(time.Second)
	now := time.Now().Local()
	err = os.Chtimes(filepath.Join(coreFolder, "cores", "arduino", "Arduino.h"), now, now)
	require.NoError(t, err)

	// Run build again, to verify that the builder rebuilds core.a
	err = bldr.Run(ctx)
	NoError(t, err)

	coreStatAfterTouch, err := os.Stat(cachedCoreFile)
	require.NoError(t, err)
	require.NotEqual(t, coreStatBefore.ModTime(), coreStatAfterTouch.ModTime())
}
