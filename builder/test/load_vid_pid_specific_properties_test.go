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
	"github.com/arduino/arduino-builder/builder/props"
	"github.com/arduino/arduino-builder/builder/types"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadVIDPIDSpecificPropertiesWhenNoVIDPIDAreProvided(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})

	buildPath := SetupBuildPath(t, context)
	defer os.RemoveAll(buildPath)

	context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION] = "10600"
	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools", "./tools_builtin"}
	context[constants.CTX_FQBN] = "arduino:avr:micro"
	context[constants.CTX_SKETCH_LOCATION] = filepath.Join("sketch1", "sketch.ino")

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},
		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	buildProperties := context[constants.CTX_BUILD_PROPERTIES].(props.PropertiesMap)

	require.Equal(t, "0x0037", buildProperties["pid.0"])
	require.Equal(t, "\"Genuino Micro\"", buildProperties["vid.4.build.usb_product"])
	require.Equal(t, "0x8037", buildProperties["build.pid"])
}

func TestLoadVIDPIDSpecificProperties(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})

	buildPath := SetupBuildPath(t, context)
	defer os.RemoveAll(buildPath)

	context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION] = "10600"
	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools", "./tools_builtin"}
	context[constants.CTX_FQBN] = "arduino:avr:micro"
	context[constants.CTX_SKETCH_LOCATION] = filepath.Join("sketch1", "sketch.ino")
	context[constants.CTX_VIDPID] = "0x2341_0x0237"

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},
		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	buildProperties := context[constants.CTX_BUILD_PROPERTIES].(props.PropertiesMap)

	require.Equal(t, "0x0037", buildProperties["pid.0"])
	require.Equal(t, "\"Genuino Micro\"", buildProperties["vid.4.build.usb_product"])
	require.Equal(t, "0x2341", buildProperties["build.vid"])
	require.Equal(t, "0x8237", buildProperties["build.pid"])
}
