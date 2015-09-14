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
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestBuilderEmptySketch(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})

	buildPath := SetupBuildPath(t, context)
	defer os.RemoveAll(buildPath)

	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools"}
	context[constants.CTX_FQBN] = "arduino:avr:uno"
	context[constants.CTX_SKETCH_LOCATION] = filepath.Join("sketch1", "sketch.ino")
	context[constants.CTX_LIBRARIES_FOLDERS] = []string{"libraries", "downloaded_libraries"}
	context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION] = "10600"
	//	context[constants.CTX_VERBOSE] = true
	//	context[constants.CTX_DEBUG_LEVEL] = 10

	command := builder.Builder{}
	err := command.Run(context)
	NoError(t, err)

	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_CORE, "HardwareSerial.cpp.o"))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_PREPROC, constants.FILE_CTAGS_TARGET))
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

	context := make(map[string]interface{})

	buildPath := SetupBuildPath(t, context)
	defer os.RemoveAll(buildPath)

	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools"}
	context[constants.CTX_FQBN] = "arduino:avr:leonardo"
	context[constants.CTX_SKETCH_LOCATION] = filepath.Join("downloaded_libraries", "Bridge", "examples", "Bridge", "Bridge.ino")
	context[constants.CTX_LIBRARIES_FOLDERS] = []string{"libraries", "downloaded_libraries"}
	context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION] = "10600"

	command := builder.Builder{}
	err := command.Run(context)
	NoError(t, err)

	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_CORE, "HardwareSerial.cpp.o"))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_PREPROC, constants.FILE_CTAGS_TARGET))
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

	context := make(map[string]interface{})

	buildPath := SetupBuildPath(t, context)
	defer os.RemoveAll(buildPath)

	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools"}
	context[constants.CTX_FQBN] = "arduino:avr:leonardo"
	context[constants.CTX_SKETCH_LOCATION] = filepath.Join("sketch_with_config", "sketch_with_config.ino")
	context[constants.CTX_LIBRARIES_FOLDERS] = []string{"libraries", "downloaded_libraries"}
	context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION] = "10600"

	command := builder.Builder{}
	err := command.Run(context)
	NoError(t, err)

	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_CORE, "HardwareSerial.cpp.o"))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_PREPROC, constants.FILE_CTAGS_TARGET))
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

	context := make(map[string]interface{})

	buildPath := SetupBuildPath(t, context)
	defer os.RemoveAll(buildPath)

	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools"}
	context[constants.CTX_FQBN] = "arduino:avr:leonardo"
	context[constants.CTX_SKETCH_LOCATION] = filepath.Join("downloaded_libraries", "Bridge", "examples", "Bridge", "Bridge.ino")
	context[constants.CTX_LIBRARIES_FOLDERS] = []string{"libraries", "downloaded_libraries"}
	context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION] = "10600"

	command := builder.Builder{}
	err := command.Run(context)
	NoError(t, err)

	command = builder.Builder{}
	err = command.Run(context)
	NoError(t, err)

	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_CORE, "HardwareSerial.cpp.o"))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_PREPROC, constants.FILE_CTAGS_TARGET))
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

	context := make(map[string]interface{})

	buildPath := SetupBuildPath(t, context)
	defer os.RemoveAll(buildPath)

	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools"}
	context[constants.CTX_FQBN] = "arduino:sam:arduino_due_x_dbg"
	context[constants.CTX_SKETCH_LOCATION] = filepath.Join("downloaded_libraries", "Bridge", "examples", "Bridge", "Bridge.ino")
	context[constants.CTX_LIBRARIES_FOLDERS] = []string{"libraries", "downloaded_libraries"}
	context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION] = "10600"
	context[constants.CTX_WARNINGS_LEVEL] = "all"

	command := builder.Builder{}
	err := command.Run(context)
	NoError(t, err)

	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_CORE, "syscalls_sam3.c.o"))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_CORE, "USB", "HID.cpp.o"))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_CORE, "avr", "dtostrf.c.d"))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_PREPROC, constants.FILE_CTAGS_TARGET))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_SKETCH, "Bridge.ino.cpp.o"))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, "Bridge.ino.elf"))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, "Bridge.ino.bin"))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_LIBRARIES, "Bridge", "Mailbox.cpp.o"))
	NoError(t, err)
}

func TestBuilderBridgeRedBearLab(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})

	buildPath := SetupBuildPath(t, context)
	defer os.RemoveAll(buildPath)

	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware", "downloaded_board_manager_stuff"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools", "downloaded_board_manager_stuff"}
	context[constants.CTX_FQBN] = "RedBearLab:avr:blend"
	context[constants.CTX_SKETCH_LOCATION] = filepath.Join("downloaded_libraries", "Bridge", "examples", "Bridge", "Bridge.ino")
	context[constants.CTX_LIBRARIES_FOLDERS] = []string{"libraries", "downloaded_libraries"}
	context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION] = "10600"

	command := builder.Builder{}
	err := command.Run(context)
	NoError(t, err)

	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_CORE, "HardwareSerial.cpp.o"))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(buildPath, constants.FOLDER_PREPROC, constants.FILE_CTAGS_TARGET))
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

	context := make(map[string]interface{})

	buildPath := SetupBuildPath(t, context)
	defer os.RemoveAll(buildPath)

	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware", "downloaded_board_manager_stuff"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools", "downloaded_board_manager_stuff"}
	context[constants.CTX_FQBN] = "RedBearLab:avr:blend"
	context[constants.CTX_SKETCH_LOCATION] = filepath.Join("sketch_no_functions", "main.ino")
	context[constants.CTX_LIBRARIES_FOLDERS] = []string{"libraries", "downloaded_libraries"}
	context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION] = "10600"

	command := builder.Builder{}
	err := command.Run(context)
	require.Error(t, err)
}
