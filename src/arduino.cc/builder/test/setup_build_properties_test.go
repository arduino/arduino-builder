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
	"arduino.cc/builder/props"
	"arduino.cc/builder/types"
	"arduino.cc/builder/utils"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestSetupBuildProperties(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})

	buildPath := SetupBuildPath(t, context)
	defer os.RemoveAll(buildPath)

	context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION] = "10600"
	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware", "user_hardware"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools", "./tools_builtin"}
	context[constants.CTX_FQBN] = "arduino:avr:uno"
	context[constants.CTX_SKETCH_LOCATION] = filepath.Join("sketch1", "sketch.ino")

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},
		&builder.AddAdditionalEntriesToContext{},
		&builder.HardwareLoader{},
		&builder.ToolsLoader{},
		&builder.TargetBoardResolver{},
		&builder.SketchLoader{},
		&builder.SetupBuildProperties{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	buildProperties := context[constants.CTX_BUILD_PROPERTIES].(map[string]string)

	require.Equal(t, "ARDUINO", buildProperties[constants.BUILD_PROPERTIES_SOFTWARE])

	require.Equal(t, "uno", buildProperties[constants.ID])
	require.Equal(t, "Arduino/Genuino Uno", buildProperties["name"])
	require.Equal(t, "0x2341", buildProperties["vid.0"])
	require.Equal(t, "\"{compiler.path}{compiler.c.cmd}\" {compiler.c.flags} -mmcu={build.mcu} -DF_CPU={build.f_cpu} -DARDUINO={runtime.ide.version} -DARDUINO_{build.board} -DARDUINO_ARCH_{build.arch} {compiler.c.extra_flags} {build.extra_flags} {includes} \"{source_file}\" -o \"{object_file}\"", buildProperties["recipe.c.o.pattern"])
	require.Equal(t, "{path}/etc/avrdude.conf", buildProperties["tools.avrdude.config.path"])

	coanProps := props.SubTree(props.SubTree(buildProperties, constants.BUILD_PROPERTIES_TOOLS_KEY), constants.COAN)
	require.Equal(t, "{path}/coan", coanProps["cmd.path"])
	require.Equal(t, "\"{cmd.path}\" source -m -E -P -kb {compiler.c.flags} -mmcu={build.mcu} -DF_CPU={build.f_cpu} -DARDUINO={runtime.ide.version} -DARDUINO_{build.board} -DARDUINO_ARCH_{build.arch} {compiler.c.extra_flags} {build.extra_flags} \"{source_file}\"", coanProps[constants.BUILD_PROPERTIES_PATTERN])

	require.Equal(t, Abs(t, "downloaded_hardware/arduino/avr"), buildProperties[constants.BUILD_PROPERTIES_RUNTIME_PLATFORM_PATH])
	require.Equal(t, Abs(t, "downloaded_hardware/arduino"), buildProperties[constants.BUILD_PROPERTIES_RUNTIME_HARDWARE_PATH])
	require.Equal(t, "10600", buildProperties[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION])
	require.NotEmpty(t, buildProperties[constants.BUILD_PROPERTIES_RUNTIME_OS])

	require.Equal(t, Abs(t, "./downloaded_tools/arm-none-eabi-gcc/4.8.3-2014q1"), buildProperties["runtime.tools.arm-none-eabi-gcc.path"])
	require.Equal(t, Abs(t, "./downloaded_tools/arm-none-eabi-gcc/4.8.3-2014q1"), buildProperties["runtime.tools.arm-none-eabi-gcc-4.8.3-2014q1.path"])
	require.Equal(t, Abs(t, "./downloaded_tools/bossac/1.3a-arduino"), buildProperties["runtime.tools.bossac-1.3a-arduino.path"])
	require.Equal(t, Abs(t, "./downloaded_tools/bossac/1.5-arduino"), buildProperties["runtime.tools.bossac-1.5-arduino.path"])
	require.True(t, buildProperties["runtime.tools.bossac.path"] == Abs(t, "./downloaded_tools/bossac/1.3a-arduino") || buildProperties["runtime.tools.bossac.path"] == Abs(t, "./downloaded_tools/bossac/1.5-arduino"))
	require.Equal(t, Abs(t, "./downloaded_tools/avrdude/6.0.1-arduino5"), buildProperties["runtime.tools.avrdude.path"])
	require.Equal(t, Abs(t, "./downloaded_tools/avrdude/6.0.1-arduino5"), buildProperties["runtime.tools.avrdude-6.0.1-arduino5.path"])
	require.Equal(t, Abs(t, "./downloaded_tools/avr-gcc/4.8.1-arduino5"), buildProperties["runtime.tools.avr-gcc.path"])
	require.Equal(t, Abs(t, "./downloaded_tools/avr-gcc/4.8.1-arduino5"), buildProperties["runtime.tools.avr-gcc-4.8.1-arduino5.path"])

	require.Equal(t, Abs(t, filepath.Join("sketch1", "sketch.ino")), buildProperties[constants.BUILD_PROPERTIES_SOURCE_PATH])

	require.True(t, utils.MapStringStringHas(buildProperties, constants.BUILD_PROPERTIES_EXTRA_TIME_UTC))
	require.True(t, utils.MapStringStringHas(buildProperties, constants.BUILD_PROPERTIES_EXTRA_TIME_LOCAL))
	require.True(t, utils.MapStringStringHas(buildProperties, constants.BUILD_PROPERTIES_EXTRA_TIME_ZONE))
	require.True(t, utils.MapStringStringHas(buildProperties, constants.BUILD_PROPERTIES_EXTRA_TIME_DST))
}

func TestSetupBuildPropertiesWithSomeCustomOverrides(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})

	buildPath := SetupBuildPath(t, context)
	defer os.RemoveAll(buildPath)

	context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION] = "10600"
	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools", "./tools_builtin"}
	context[constants.CTX_FQBN] = "arduino:avr:uno"
	context[constants.CTX_SKETCH_LOCATION] = filepath.Join("sketch1", "sketch.ino")
	context[constants.CTX_CUSTOM_BUILD_PROPERTIES] = []string{"name=fake name", "tools.avrdude.config.path=non existent path with space and a ="}

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},
		&builder.AddAdditionalEntriesToContext{},
		&builder.HardwareLoader{},
		&builder.ToolsLoader{},
		&builder.TargetBoardResolver{},
		&builder.SketchLoader{},
		&builder.SetupBuildProperties{},
		&builder.SetCustomBuildProperties{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	buildProperties := context[constants.CTX_BUILD_PROPERTIES].(map[string]string)

	require.Equal(t, "ARDUINO", buildProperties[constants.BUILD_PROPERTIES_SOFTWARE])

	require.Equal(t, "uno", buildProperties[constants.ID])
	require.Equal(t, "fake name", buildProperties["name"])
	require.Equal(t, "\"{compiler.path}{compiler.c.cmd}\" {compiler.c.flags} -mmcu={build.mcu} -DF_CPU={build.f_cpu} -DARDUINO={runtime.ide.version} -DARDUINO_{build.board} -DARDUINO_ARCH_{build.arch} {compiler.c.extra_flags} {build.extra_flags} {includes} \"{source_file}\" -o \"{object_file}\"", buildProperties["recipe.c.o.pattern"])
	require.Equal(t, "non existent path with space and a =", buildProperties["tools.avrdude.config.path"])
}
