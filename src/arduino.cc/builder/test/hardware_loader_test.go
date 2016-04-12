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
	"path/filepath"
	"runtime"
	"testing"
)

func TestLoadHardware(t *testing.T) {
	context := make(map[string]interface{})
	context[constants.CTX_HARDWARE_FOLDERS] = []string{"downloaded_hardware", filepath.Join("..", "hardware"), "hardware"}

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},
		&builder.HardwareLoader{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	packages := context[constants.CTX_HARDWARE].(*types.Packages)
	require.Equal(t, 2, len(packages.Packages))
	require.NotNil(t, packages.Packages["arduino"])
	require.Equal(t, 2, len(packages.Packages["arduino"].Platforms))

	require.Equal(t, "uno", packages.Packages["arduino"].Platforms["avr"].Boards["uno"].BoardId)
	require.Equal(t, "uno", packages.Packages["arduino"].Platforms["avr"].Boards["uno"].Properties[constants.ID])

	require.Equal(t, "yun", packages.Packages["arduino"].Platforms["avr"].Boards["yun"].BoardId)
	require.Equal(t, "true", packages.Packages["arduino"].Platforms["avr"].Boards["yun"].Properties["upload.wait_for_upload_port"])

	require.Equal(t, "{build.usb_flags}", packages.Packages["arduino"].Platforms["avr"].Boards["robotMotor"].Properties["build.extra_flags"])

	require.Equal(t, "arduino_due_x", packages.Packages["arduino"].Platforms["sam"].Boards["arduino_due_x"].BoardId)

	require.Equal(t, "ATmega123", packages.Packages["arduino"].Platforms["avr"].Boards["diecimila"].Properties["menu.cpu.atmega123"])

	avrPlatform := packages.Packages["arduino"].Platforms["avr"]
	require.Equal(t, "Arduino AVR Boards", avrPlatform.Properties[constants.PLATFORM_NAME])
	require.Equal(t, "-v", avrPlatform.Properties["tools.avrdude.bootloader.params.verbose"])
	require.Equal(t, "/my/personal/avrdude", avrPlatform.Properties["tools.avrdude.cmd.path"])

	require.Equal(t, "AVRISP mkII", avrPlatform.Programmers["avrispmkii"][constants.PROGRAMMER_NAME])

	require.Equal(t, "{runtime.tools.ctags.path}", packages.Properties["tools.ctags.path"])
	require.Equal(t, "\"{cmd.path}\" -u --language-force=c++ -f - --c++-kinds=svpf --fields=KSTtzns --line-directives \"{source_file}\"", packages.Properties["tools.ctags.pattern"])
	require.Equal(t, "{runtime.tools.avrdude.path}", packages.Properties["tools.avrdude.path"])
	require.Equal(t, "-w -x c++ -E -CC", packages.Properties["preproc.macros.flags"])
}

func TestLoadHardwareMixingUserHardwareFolder(t *testing.T) {
	context := make(map[string]interface{})
	context[constants.CTX_HARDWARE_FOLDERS] = []string{"downloaded_hardware", filepath.Join("..", "hardware"), "hardware", "user_hardware"}

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},
		&builder.AddAdditionalEntriesToContext{},
		&builder.HardwareLoader{},
		&builder.PlatformKeysRewriteLoader{},
		&builder.RewriteHardwareKeys{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	packages := context[constants.CTX_HARDWARE].(*types.Packages)

	if runtime.GOOS == "windows" {
		//a package is a symlink, and windows does not support them
		require.Equal(t, 3, len(packages.Packages))
	} else {
		require.Equal(t, 4, len(packages.Packages))
	}

	require.NotNil(t, packages.Packages["arduino"])
	require.Equal(t, 2, len(packages.Packages["arduino"].Platforms))

	require.Equal(t, "uno", packages.Packages["arduino"].Platforms["avr"].Boards["uno"].BoardId)
	require.Equal(t, "uno", packages.Packages["arduino"].Platforms["avr"].Boards["uno"].Properties[constants.ID])

	require.Equal(t, "yun", packages.Packages["arduino"].Platforms["avr"].Boards["yun"].BoardId)
	require.Equal(t, "true", packages.Packages["arduino"].Platforms["avr"].Boards["yun"].Properties["upload.wait_for_upload_port"])

	require.Equal(t, "{build.usb_flags}", packages.Packages["arduino"].Platforms["avr"].Boards["robotMotor"].Properties["build.extra_flags"])

	require.Equal(t, "arduino_due_x", packages.Packages["arduino"].Platforms["sam"].Boards["arduino_due_x"].BoardId)

	avrPlatform := packages.Packages["arduino"].Platforms["avr"]
	require.Equal(t, "Arduino AVR Boards", avrPlatform.Properties[constants.PLATFORM_NAME])
	require.Equal(t, "-v", avrPlatform.Properties["tools.avrdude.bootloader.params.verbose"])
	require.Equal(t, "/my/personal/avrdude", avrPlatform.Properties["tools.avrdude.cmd.path"])

	require.Equal(t, "AVRISP mkII", avrPlatform.Programmers["avrispmkii"][constants.PROGRAMMER_NAME])

	require.Equal(t, "-w -x c++ -M -MG -MP", avrPlatform.Properties["preproc.includes.flags"])
	require.Equal(t, "-w -x c++ -E -CC", avrPlatform.Properties["preproc.macros.flags"])
	require.Equal(t, "\"{compiler.path}{compiler.cpp.cmd}\" {compiler.cpp.flags} {preproc.includes.flags} -mmcu={build.mcu} -DF_CPU={build.f_cpu} -DARDUINO={runtime.ide.version} -DARDUINO_{build.board} -DARDUINO_ARCH_{build.arch} {compiler.cpp.extra_flags} {build.extra_flags} {includes} \"{source_file}\"", avrPlatform.Properties[constants.RECIPE_PREPROC_INCLUDES])
	require.False(t, utils.MapStringStringHas(avrPlatform.Properties, "preproc.macros.compatibility_flags"))

	require.NotNil(t, packages.Packages["my_avr_platform"])
	myAVRPlatform := packages.Packages["my_avr_platform"]
	require.Equal(t, "hello world", myAVRPlatform.Properties["example"])
	myAVRPlatformAvrArch := myAVRPlatform.Platforms["avr"]
	require.Equal(t, "custom_yun", myAVRPlatformAvrArch.Boards["custom_yun"].BoardId)

	require.False(t, utils.MapStringStringHas(myAVRPlatformAvrArch.Properties, "preproc.includes.flags"))

	require.Equal(t, "{runtime.tools.ctags.path}", packages.Properties["tools.ctags.path"])
	require.Equal(t, "\"{cmd.path}\" -u --language-force=c++ -f - --c++-kinds=svpf --fields=KSTtzns --line-directives \"{source_file}\"", packages.Properties["tools.ctags.pattern"])
	require.Equal(t, "{runtime.tools.avrdude.path}", packages.Properties["tools.avrdude.path"])
	require.Equal(t, "-w -x c++ -E -CC", packages.Properties["preproc.macros.flags"])

	if runtime.GOOS != "windows" {
		require.NotNil(t, packages.Packages["my_symlinked_avr_platform"])
		require.NotNil(t, packages.Packages["my_symlinked_avr_platform"].Platforms["avr"])
	}
}

func TestLoadHardwareWithBoardManagerFolderStructure(t *testing.T) {
	context := make(map[string]interface{})
	context[constants.CTX_HARDWARE_FOLDERS] = []string{"downloaded_board_manager_stuff"}

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},
		&builder.HardwareLoader{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	packages := context[constants.CTX_HARDWARE].(*types.Packages)
	require.Equal(t, 3, len(packages.Packages))
	require.NotNil(t, packages.Packages["arduino"])
	require.Equal(t, 1, len(packages.Packages["arduino"].Platforms))
	require.NotNil(t, packages.Packages["RedBearLab"])
	require.Equal(t, 1, len(packages.Packages["RedBearLab"].Platforms))
	require.NotNil(t, packages.Packages["RFduino"])
	require.Equal(t, 0, len(packages.Packages["RFduino"].Platforms))

	samdPlatform := packages.Packages["arduino"].Platforms["samd"]
	require.Equal(t, 3, len(samdPlatform.Boards))

	require.Equal(t, "arduino_zero_edbg", samdPlatform.Boards["arduino_zero_edbg"].BoardId)
	require.Equal(t, "arduino_zero_edbg", samdPlatform.Boards["arduino_zero_edbg"].Properties[constants.ID])

	require.Equal(t, "arduino_zero", samdPlatform.Boards["arduino_zero_native"].Properties["build.variant"])
	require.Equal(t, "-D__SAMD21G18A__ {build.usb_flags}", samdPlatform.Boards["arduino_zero_native"].Properties["build.extra_flags"])

	require.Equal(t, "Arduino SAMD (32-bits ARM Cortex-M0+) Boards", samdPlatform.Properties[constants.PLATFORM_NAME])
	require.Equal(t, "-d3", samdPlatform.Properties["tools.openocd.erase.params.verbose"])

	require.Equal(t, 3, len(samdPlatform.Programmers))

	require.Equal(t, "Atmel EDBG", samdPlatform.Programmers["edbg"][constants.PROGRAMMER_NAME])
	require.Equal(t, "openocd", samdPlatform.Programmers["edbg"]["program.tool"])

	avrRedBearPlatform := packages.Packages["RedBearLab"].Platforms["avr"]
	require.Equal(t, 3, len(avrRedBearPlatform.Boards))

	require.Equal(t, "blend", avrRedBearPlatform.Boards["blend"].BoardId)
	require.Equal(t, "blend", avrRedBearPlatform.Boards["blend"].Properties[constants.ID])
	require.Equal(t, "arduino:arduino", avrRedBearPlatform.Boards["blend"].Properties[constants.BUILD_PROPERTIES_BUILD_CORE])
}

func TestLoadLotsOfHardware(t *testing.T) {
	context := make(map[string]interface{})

	context[constants.CTX_HARDWARE_FOLDERS] = []string{"downloaded_board_manager_stuff", "downloaded_hardware", filepath.Join("..", "hardware"), "hardware", "user_hardware"}

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},
		&builder.HardwareLoader{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	packages := context[constants.CTX_HARDWARE].(*types.Packages)

	if runtime.GOOS == "windows" {
		//a package is a symlink, and windows does not support them
		require.Equal(t, 5, len(packages.Packages))
	} else {
		require.Equal(t, 6, len(packages.Packages))
	}

	require.NotNil(t, packages.Packages["arduino"])
	require.NotNil(t, packages.Packages["my_avr_platform"])

	require.Equal(t, 3, len(packages.Packages["arduino"].Platforms))
	require.Equal(t, 20, len(packages.Packages["arduino"].Platforms["avr"].Boards))
	require.Equal(t, 2, len(packages.Packages["arduino"].Platforms["sam"].Boards))
	require.Equal(t, 3, len(packages.Packages["arduino"].Platforms["samd"].Boards))

	require.Equal(t, 1, len(packages.Packages["my_avr_platform"].Platforms))
	require.Equal(t, 2, len(packages.Packages["my_avr_platform"].Platforms["avr"].Boards))

	if runtime.GOOS != "windows" {
		require.Equal(t, 1, len(packages.Packages["my_symlinked_avr_platform"].Platforms))
		require.Equal(t, 2, len(packages.Packages["my_symlinked_avr_platform"].Platforms["avr"].Boards))
	}
}
