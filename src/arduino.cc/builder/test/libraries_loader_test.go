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
	"path/filepath"
	"sort"
	"testing"
)

func TestLoadLibrariesAVR(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})

	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_FQBN] = "arduino:avr:leonardo"
	context[constants.CTX_LIBRARIES_FOLDERS] = []string{"libraries", "downloaded_libraries"}

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},
		&builder.AddAdditionalEntriesToContext{},
		&builder.HardwareLoader{},
		&builder.TargetBoardResolver{},
		&builder.LibrariesLoader{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	libraries := context[constants.CTX_LIBRARIES].([]*types.Library)
	require.Equal(t, 14, len(libraries))

	sort.Sort(ByLibraryName(libraries))

	idx := 0

	require.Equal(t, "ANewLibrary-master", libraries[idx].Name)

	idx++
	require.Equal(t, "Adafruit_PN532", libraries[idx].Name)
	require.Equal(t, Abs(t, "downloaded_libraries/Adafruit_PN532"), libraries[idx].Folder)
	require.Equal(t, Abs(t, "downloaded_libraries/Adafruit_PN532"), libraries[idx].SrcFolder)
	require.Equal(t, 1, len(libraries[idx].Archs))
	require.Equal(t, constants.LIBRARY_ALL_ARCHS, libraries[idx].Archs[0])
	require.False(t, libraries[idx].IsLegacy)

	idx++
	require.Equal(t, "Audio", libraries[idx].Name)

	idx++
	require.Equal(t, "Balanduino", libraries[idx].Name)
	require.True(t, libraries[idx].IsLegacy)

	idx++
	bridgeLib := libraries[idx]
	require.Equal(t, "Bridge", bridgeLib.Name)
	require.Equal(t, Abs(t, "downloaded_libraries/Bridge"), bridgeLib.Folder)
	require.Equal(t, Abs(t, "downloaded_libraries/Bridge/src"), bridgeLib.SrcFolder)
	require.Equal(t, 1, len(bridgeLib.Archs))
	require.Equal(t, constants.LIBRARY_ALL_ARCHS, bridgeLib.Archs[0])
	require.Equal(t, "Arduino", bridgeLib.Author)
	require.Equal(t, "Arduino <info@arduino.cc>", bridgeLib.Maintainer)

	idx++
	require.Equal(t, "CapacitiveSensor", libraries[idx].Name)
	idx++
	require.Equal(t, "EEPROM", libraries[idx].Name)
	idx++
	require.Equal(t, "FakeAudio", libraries[idx].Name)
	idx++
	require.Equal(t, "IRremote", libraries[idx].Name)
	idx++
	require.Equal(t, "Robot_IR_Remote", libraries[idx].Name)
	idx++
	require.Equal(t, "SPI", libraries[idx].Name)
	idx++
	require.Equal(t, "SPI", libraries[idx].Name)
	idx++
	require.Equal(t, "SoftwareSerial", libraries[idx].Name)
	idx++
	require.Equal(t, "Wire", libraries[idx].Name)

	headerToLibraries := context[constants.CTX_HEADER_TO_LIBRARIES].(map[string][]*types.Library)
	require.Equal(t, 2, len(headerToLibraries["Audio.h"]))
	require.Equal(t, "FakeAudio", headerToLibraries["Audio.h"][0].Name)
	require.Equal(t, "Audio", headerToLibraries["Audio.h"][1].Name)
	require.Equal(t, 1, len(headerToLibraries["FakeAudio.h"]))
	require.Equal(t, "FakeAudio", headerToLibraries["FakeAudio.h"][0].Name)
	require.Equal(t, 1, len(headerToLibraries["Adafruit_PN532.h"]))
	require.Equal(t, "Adafruit_PN532", headerToLibraries["Adafruit_PN532.h"][0].Name)

	require.Equal(t, 2, len(headerToLibraries["IRremote.h"]))

	libraries = headerToLibraries["IRremote.h"]
	sort.Sort(ByLibraryName(libraries))

	require.Equal(t, "IRremote", libraries[0].Name)
	require.Equal(t, "Robot_IR_Remote", libraries[1].Name)
}

func TestLoadLibrariesSAM(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})

	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_FQBN] = "arduino:sam:arduino_due_x_dbg"
	context[constants.CTX_LIBRARIES_FOLDERS] = []string{"libraries", "downloaded_libraries"}

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},
		&builder.AddAdditionalEntriesToContext{},
		&builder.HardwareLoader{},
		&builder.TargetBoardResolver{},
		&builder.LibrariesLoader{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	libraries := context[constants.CTX_LIBRARIES].([]*types.Library)
	require.Equal(t, 12, len(libraries))

	sort.Sort(ByLibraryName(libraries))

	idx := 0
	require.Equal(t, "ANewLibrary-master", libraries[idx].Name)
	idx++
	require.Equal(t, "Adafruit_PN532", libraries[idx].Name)
	idx++
	require.Equal(t, "Audio", libraries[idx].Name)
	idx++
	require.Equal(t, "Balanduino", libraries[idx].Name)
	idx++
	require.Equal(t, "Bridge", libraries[idx].Name)
	idx++
	require.Equal(t, "CapacitiveSensor", libraries[idx].Name)
	idx++
	require.Equal(t, "FakeAudio", libraries[idx].Name)
	idx++
	require.Equal(t, "IRremote", libraries[idx].Name)
	idx++
	require.Equal(t, "Robot_IR_Remote", libraries[idx].Name)
	idx++
	require.Equal(t, "SPI", libraries[idx].Name)
	idx++
	require.Equal(t, "SPI", libraries[idx].Name)
	idx++
	require.Equal(t, "Wire", libraries[idx].Name)

	headerToLibraries := context[constants.CTX_HEADER_TO_LIBRARIES].(map[string][]*types.Library)

	require.Equal(t, 2, len(headerToLibraries["Audio.h"]))
	libraries = headerToLibraries["Audio.h"]
	sort.Sort(ByLibraryName(libraries))
	require.Equal(t, "Audio", libraries[0].Name)
	require.Equal(t, "FakeAudio", libraries[1].Name)

	require.Equal(t, 1, len(headerToLibraries["FakeAudio.h"]))
	require.Equal(t, "FakeAudio", headerToLibraries["FakeAudio.h"][0].Name)
	require.Equal(t, 2, len(headerToLibraries["IRremote.h"]))
	require.Equal(t, "IRremote", headerToLibraries["IRremote.h"][0].Name)
	require.Equal(t, "Robot_IR_Remote", headerToLibraries["IRremote.h"][1].Name)
}

func TestLoadLibrariesAVRNoDuplicateLibrariesFolders(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})

	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_FQBN] = "arduino:avr:leonardo"
	context[constants.CTX_LIBRARIES_FOLDERS] = []string{"libraries", "downloaded_libraries", filepath.Join("downloaded_hardware", "arduino", "avr", "libraries")}

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},
		&builder.AddAdditionalEntriesToContext{},
		&builder.HardwareLoader{},
		&builder.TargetBoardResolver{},
		&builder.LibrariesLoader{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	require.Equal(t, 3, len(context[constants.CTX_LIBRARIES_FOLDERS].([]string)))
}
