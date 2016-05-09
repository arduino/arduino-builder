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
	"os"
	"path/filepath"
	"testing"
)

func TestSetupBuildProperties(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	ctx := &types.Context{
		HardwareFolders:   []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware", "user_hardware"},
		ToolsFolders:      []string{"downloaded_tools", "./tools_builtin"},
		SketchLocation:    filepath.Join("sketch1", "sketch.ino"),
		FQBN:              "arduino:avr:uno",
		ArduinoAPIVersion: "10600",
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	commands := []types.Command{
		&builder.AddAdditionalEntriesToContext{},
		&builder.HardwareLoader{},
		&builder.ToolsLoader{},
		&builder.TargetBoardResolver{},
		&builder.SketchLoader{},
		&builder.SetupBuildProperties{},
	}

	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}

	buildProperties := ctx.BuildProperties

	require.Equal(t, "ARDUINO", buildProperties[constants.BUILD_PROPERTIES_SOFTWARE])

	require.Equal(t, "uno", buildProperties[constants.ID])
	require.Equal(t, "Arduino/Genuino Uno", buildProperties["name"])
	require.Equal(t, "0x2341", buildProperties["vid.0"])
	require.Equal(t, "\"{compiler.path}{compiler.c.cmd}\" {compiler.c.flags} -mmcu={build.mcu} -DF_CPU={build.f_cpu} -DARDUINO={runtime.ide.version} -DARDUINO_{build.board} -DARDUINO_ARCH_{build.arch} {compiler.c.extra_flags} {build.extra_flags} {includes} \"{source_file}\" -o \"{object_file}\"", buildProperties["recipe.c.o.pattern"])
	require.Equal(t, "{path}/etc/avrdude.conf", buildProperties["tools.avrdude.config.path"])

	require.Equal(t, Abs(t, "downloaded_hardware/arduino/avr"), buildProperties[constants.BUILD_PROPERTIES_RUNTIME_PLATFORM_PATH])
	require.Equal(t, Abs(t, "downloaded_hardware/arduino"), buildProperties[constants.BUILD_PROPERTIES_RUNTIME_HARDWARE_PATH])
	require.Equal(t, "10600", buildProperties[constants.BUILD_PROPERTIES_RUNTIME_IDE_VERSION])
	require.NotEmpty(t, buildProperties[constants.BUILD_PROPERTIES_RUNTIME_OS])

	require.Equal(t, Abs(t, "./downloaded_tools/arm-none-eabi-gcc/4.8.3-2014q1"), buildProperties["runtime.tools.arm-none-eabi-gcc.path"])
	require.Equal(t, Abs(t, "./downloaded_tools/arm-none-eabi-gcc/4.8.3-2014q1"), buildProperties["runtime.tools.arm-none-eabi-gcc-4.8.3-2014q1.path"])
	require.Equal(t, Abs(t, "./downloaded_tools/bossac/1.6.1-arduino"), buildProperties["runtime.tools.bossac-1.6.1-arduino.path"])
	require.Equal(t, Abs(t, "./downloaded_tools/bossac/1.5-arduino"), buildProperties["runtime.tools.bossac-1.5-arduino.path"])
	require.True(t, buildProperties["runtime.tools.bossac.path"] == Abs(t, "./downloaded_tools/bossac/1.6.1-arduino") || buildProperties["runtime.tools.bossac.path"] == Abs(t, "./downloaded_tools/bossac/1.5-arduino"))
	require.Equal(t, Abs(t, "./downloaded_tools/avrdude/6.0.1-arduino5"), buildProperties["runtime.tools.avrdude.path"])
	require.Equal(t, Abs(t, "./downloaded_tools/avrdude/6.0.1-arduino5"), buildProperties["runtime.tools.avrdude-6.0.1-arduino5.path"])
	require.Equal(t, Abs(t, "./downloaded_tools/avr-gcc/4.8.1-arduino5"), buildProperties["runtime.tools.avr-gcc.path"])
	require.Equal(t, Abs(t, "./downloaded_tools/avr-gcc/4.8.1-arduino5"), buildProperties["runtime.tools.avr-gcc-4.8.1-arduino5.path"])

	require.Equal(t, Abs(t, "sketch1"), buildProperties[constants.BUILD_PROPERTIES_SOURCE_PATH])

	require.True(t, utils.MapStringStringHas(buildProperties, constants.BUILD_PROPERTIES_EXTRA_TIME_UTC))
	require.True(t, utils.MapStringStringHas(buildProperties, constants.BUILD_PROPERTIES_EXTRA_TIME_LOCAL))
	require.True(t, utils.MapStringStringHas(buildProperties, constants.BUILD_PROPERTIES_EXTRA_TIME_ZONE))
	require.True(t, utils.MapStringStringHas(buildProperties, constants.BUILD_PROPERTIES_EXTRA_TIME_DST))
}

func TestSetupBuildPropertiesWithSomeCustomOverrides(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	ctx := &types.Context{
		HardwareFolders:   []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"},
		ToolsFolders:      []string{"downloaded_tools", "./tools_builtin"},
		SketchLocation:    filepath.Join("sketch1", "sketch.ino"),
		FQBN:              "arduino:avr:uno",
		ArduinoAPIVersion: "10600",

		CustomBuildProperties: []string{"name=fake name", "tools.avrdude.config.path=non existent path with space and a ="},
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	commands := []types.Command{
		&builder.AddAdditionalEntriesToContext{},
		&builder.HardwareLoader{},
		&builder.ToolsLoader{},
		&builder.TargetBoardResolver{},
		&builder.SketchLoader{},
		&builder.SetupBuildProperties{},
		&builder.SetCustomBuildProperties{},
	}

	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}

	buildProperties := ctx.BuildProperties

	require.Equal(t, "ARDUINO", buildProperties[constants.BUILD_PROPERTIES_SOFTWARE])

	require.Equal(t, "uno", buildProperties[constants.ID])
	require.Equal(t, "fake name", buildProperties["name"])
	require.Equal(t, "\"{compiler.path}{compiler.c.cmd}\" {compiler.c.flags} -mmcu={build.mcu} -DF_CPU={build.f_cpu} -DARDUINO={runtime.ide.version} -DARDUINO_{build.board} -DARDUINO_ARCH_{build.arch} {compiler.c.extra_flags} {build.extra_flags} {includes} \"{source_file}\" -o \"{object_file}\"", buildProperties["recipe.c.o.pattern"])
	require.Equal(t, "non existent path with space and a =", buildProperties["tools.avrdude.config.path"])
}

func TestSetupBuildPropertiesUserHardware(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	ctx := &types.Context{
		HardwareFolders:   []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware", "user_hardware"},
		ToolsFolders:      []string{"downloaded_tools", "./tools_builtin"},
		SketchLocation:    filepath.Join("sketch1", "sketch.ino"),
		FQBN:              "my_avr_platform:avr:custom_yun",
		ArduinoAPIVersion: "10600",
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	commands := []types.Command{
		&builder.AddAdditionalEntriesToContext{},
		&builder.HardwareLoader{},
		&builder.ToolsLoader{},
		&builder.TargetBoardResolver{},
		&builder.SketchLoader{},
		&builder.SetupBuildProperties{},
	}

	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}

	buildProperties := ctx.BuildProperties

	require.Equal(t, "ARDUINO", buildProperties[constants.BUILD_PROPERTIES_SOFTWARE])

	require.Equal(t, "custom_yun", buildProperties[constants.ID])
	require.Equal(t, "caterina/Caterina-custom_yun.hex", buildProperties[constants.BUILD_PROPERTIES_BOOTLOADER_FILE])
	require.Equal(t, Abs(t, filepath.Join("user_hardware", "my_avr_platform", "avr")), buildProperties[constants.BUILD_PROPERTIES_RUNTIME_PLATFORM_PATH])
	require.Equal(t, Abs(t, filepath.Join("user_hardware", "my_avr_platform")), buildProperties[constants.BUILD_PROPERTIES_RUNTIME_HARDWARE_PATH])
}

func TestSetupBuildPropertiesWithMissingPropsFromParentPlatformTxtFiles(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	ctx := &types.Context{
		HardwareFolders:   []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware", "user_hardware"},
		ToolsFolders:      []string{"downloaded_tools", "./tools_builtin"},
		SketchLocation:    filepath.Join("sketch1", "sketch.ino"),
		FQBN:              "my_avr_platform:avr:custom_yun",
		ArduinoAPIVersion: "10600",
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

	require.Equal(t, "ARDUINO", buildProperties[constants.BUILD_PROPERTIES_SOFTWARE])

	require.Equal(t, "custom_yun", buildProperties[constants.ID])
	require.Equal(t, "Arduino Yún", buildProperties["name"])
	require.Equal(t, "0x2341", buildProperties["vid.0"])
	require.Equal(t, "\"{compiler.path}{compiler.c.cmd}\" {compiler.c.flags} -mmcu={build.mcu} -DF_CPU={build.f_cpu} -DARDUINO={runtime.ide.version} -DARDUINO_{build.board} -DARDUINO_ARCH_{build.arch} {compiler.c.extra_flags} {build.extra_flags} {includes} \"{source_file}\" -o \"{object_file}\"", buildProperties["recipe.c.o.pattern"])
	require.Equal(t, "{path}/etc/avrdude.conf", buildProperties["tools.avrdude.config.path"])

	require.Equal(t, Abs(t, "user_hardware/my_avr_platform/avr"), buildProperties[constants.BUILD_PROPERTIES_RUNTIME_PLATFORM_PATH])
	require.Equal(t, Abs(t, "user_hardware/my_avr_platform"), buildProperties[constants.BUILD_PROPERTIES_RUNTIME_HARDWARE_PATH])
	require.Equal(t, "10600", buildProperties[constants.BUILD_PROPERTIES_RUNTIME_IDE_VERSION])
	require.NotEmpty(t, buildProperties[constants.BUILD_PROPERTIES_RUNTIME_OS])

	require.Equal(t, Abs(t, "./downloaded_tools/arm-none-eabi-gcc/4.8.3-2014q1"), buildProperties["runtime.tools.arm-none-eabi-gcc.path"])
	require.Equal(t, Abs(t, "./downloaded_tools/arm-none-eabi-gcc/4.8.3-2014q1"), buildProperties["runtime.tools.arm-none-eabi-gcc-4.8.3-2014q1.path"])
	require.Equal(t, Abs(t, "./downloaded_tools/bossac/1.6.1-arduino"), buildProperties["runtime.tools.bossac-1.6.1-arduino.path"])
	require.Equal(t, Abs(t, "./downloaded_tools/bossac/1.5-arduino"), buildProperties["runtime.tools.bossac-1.5-arduino.path"])
	require.True(t, buildProperties["runtime.tools.bossac.path"] == Abs(t, "./downloaded_tools/bossac/1.6.1-arduino") || buildProperties["runtime.tools.bossac.path"] == Abs(t, "./downloaded_tools/bossac/1.5-arduino"))
	require.Equal(t, Abs(t, "./downloaded_tools/avrdude/6.0.1-arduino5"), buildProperties["runtime.tools.avrdude.path"])
	require.Equal(t, Abs(t, "./downloaded_tools/avrdude/6.0.1-arduino5"), buildProperties["runtime.tools.avrdude-6.0.1-arduino5.path"])
	require.Equal(t, Abs(t, "./downloaded_tools/avr-gcc/4.8.1-arduino5"), buildProperties["runtime.tools.avr-gcc.path"])
	require.Equal(t, Abs(t, "./downloaded_tools/avr-gcc/4.8.1-arduino5"), buildProperties["runtime.tools.avr-gcc-4.8.1-arduino5.path"])

	require.Equal(t, Abs(t, "sketch1"), buildProperties[constants.BUILD_PROPERTIES_SOURCE_PATH])

	require.True(t, utils.MapStringStringHas(buildProperties, constants.BUILD_PROPERTIES_EXTRA_TIME_UTC))
	require.True(t, utils.MapStringStringHas(buildProperties, constants.BUILD_PROPERTIES_EXTRA_TIME_LOCAL))
	require.True(t, utils.MapStringStringHas(buildProperties, constants.BUILD_PROPERTIES_EXTRA_TIME_ZONE))
	require.True(t, utils.MapStringStringHas(buildProperties, constants.BUILD_PROPERTIES_EXTRA_TIME_DST))
}
