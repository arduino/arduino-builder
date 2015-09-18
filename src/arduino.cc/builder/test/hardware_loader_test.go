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
	"runtime"
	"testing"
)

func TestLoadHardware(t *testing.T) {
	context := make(map[string]interface{})
	context[constants.CTX_HARDWARE_FOLDERS] = []string{"downloaded_hardware", filepath.Join("..", "hardware"), "hardware"}

	loader := builder.HardwareLoader{}
	err := loader.Run(context)
	NoError(t, err)

	packages := context[constants.CTX_HARDWARE].(map[string]*types.Package)
	require.Equal(t, 1, len(packages))
	require.NotNil(t, packages["arduino"])
	require.Equal(t, 2, len(packages["arduino"].Platforms))

	require.Equal(t, "uno", packages["arduino"].Platforms["avr"].Boards["uno"].BoardId)
	require.Equal(t, "uno", packages["arduino"].Platforms["avr"].Boards["uno"].Properties[constants.ID])

	require.Equal(t, "yun", packages["arduino"].Platforms["avr"].Boards["yun"].BoardId)
	require.Equal(t, "true", packages["arduino"].Platforms["avr"].Boards["yun"].Properties["upload.wait_for_upload_port"])

	require.Equal(t, "{build.usb_flags}", packages["arduino"].Platforms["avr"].Boards["robotMotor"].Properties["build.extra_flags"])

	require.Equal(t, "arduino_due_x", packages["arduino"].Platforms["sam"].Boards["arduino_due_x"].BoardId)

	avrPlatform := packages["arduino"].Platforms["avr"]
	require.Equal(t, "Arduino AVR Boards", avrPlatform.Properties[constants.PLATFORM_NAME])
	require.Equal(t, "-v", avrPlatform.Properties["tools.avrdude.bootloader.params.verbose"])
	require.Equal(t, "/my/personal/avrdude", avrPlatform.Properties["tools.avrdude.cmd.path"])

	require.Equal(t, "AVRISP mkII", avrPlatform.Programmers["avrispmkii"][constants.PROGRAMMER_NAME])

	require.Equal(t, "{path}/ctags", avrPlatform.Properties["tools.ctags.cmd.path"])
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

	packages := context[constants.CTX_HARDWARE].(map[string]*types.Package)

	if runtime.GOOS == "windows" {
		//a package is a symlink, and windows does not support them
		require.Equal(t, 2, len(packages))
	} else {
		require.Equal(t, 3, len(packages))
	}

	require.NotNil(t, packages["arduino"])
	require.Equal(t, 2, len(packages["arduino"].Platforms))

	require.Equal(t, "uno", packages["arduino"].Platforms["avr"].Boards["uno"].BoardId)
	require.Equal(t, "uno", packages["arduino"].Platforms["avr"].Boards["uno"].Properties[constants.ID])

	require.Equal(t, "yun", packages["arduino"].Platforms["avr"].Boards["yun"].BoardId)
	require.Equal(t, "true", packages["arduino"].Platforms["avr"].Boards["yun"].Properties["upload.wait_for_upload_port"])

	require.Equal(t, "{build.usb_flags}", packages["arduino"].Platforms["avr"].Boards["robotMotor"].Properties["build.extra_flags"])

	require.Equal(t, "arduino_due_x", packages["arduino"].Platforms["sam"].Boards["arduino_due_x"].BoardId)

	avrPlatform := packages["arduino"].Platforms["avr"]
	require.Equal(t, "Arduino AVR Boards", avrPlatform.Properties[constants.PLATFORM_NAME])
	require.Equal(t, "-v", avrPlatform.Properties["tools.avrdude.bootloader.params.verbose"])
	require.Equal(t, "/my/personal/avrdude", avrPlatform.Properties["tools.avrdude.cmd.path"])

	require.Equal(t, "AVRISP mkII", avrPlatform.Programmers["avrispmkii"][constants.PROGRAMMER_NAME])

	require.Equal(t, "-w -x c++ -M -MG -MP", avrPlatform.Properties["preproc.includes.flags"])
	require.Equal(t, "-w -x c++ -E -CC", avrPlatform.Properties["preproc.macros.flags"])
	require.Equal(t, "{build.mbed_api_include} {build.nRF51822_api_include} {build.ble_api_include} {compiler.libsam.c.flags} {compiler.arm.cmsis.path} {build.variant_system_include}", avrPlatform.Properties["preproc.macros.compatibility_flags"])
	require.Equal(t, "\"{compiler.path}{compiler.cpp.cmd}\" {preproc.includes.flags} -DF_CPU={build.f_cpu} -DARDUINO={runtime.ide.version} -DARDUINO_{build.board} -DARDUINO_ARCH_{build.arch} {compiler.cpp.extra_flags} {build.extra_flags} {includes} \"{source_file}\"", avrPlatform.Properties[constants.RECIPE_PREPROC_INCLUDES])

	require.NotNil(t, packages["my_avr_platform"])
	myAVRPlatform := packages["my_avr_platform"].Platforms["avr"]
	require.Equal(t, "custom_yun", myAVRPlatform.Boards["custom_yun"].BoardId)
	require.Equal(t, "{path}/ctags", myAVRPlatform.Properties["tools.ctags.cmd.path"])
	require.Equal(t, "{runtime.tools.avr-gcc.path}/bin/", myAVRPlatform.Properties[constants.BUILD_PROPERTIES_COMPILER_PATH])
	require.Equal(t, "{runtime.tools.avrdude.path}", myAVRPlatform.Properties["tools.avrdude.path"])
	require.Equal(t, "{path}/bin/avrdude", myAVRPlatform.Properties["tools.avrdude.cmd.path"])
	require.Equal(t, "{path}/etc/avrdude.conf", myAVRPlatform.Properties["tools.avrdude.config.path"])

	require.Equal(t, "-w -x c++ -M -MG -MP", myAVRPlatform.Properties["preproc.includes.flags"])
	require.Equal(t, "-w -x c++ -E -CC", myAVRPlatform.Properties["preproc.macros.flags"])

	if runtime.GOOS != "windows" {
		require.NotNil(t, packages["my_symlinked_avr_platform"])
		require.NotNil(t, packages["my_symlinked_avr_platform"].Platforms["avr"])
	}
}

func TestLoadHardwareWithBoardManagerFolderStructure(t *testing.T) {
	context := make(map[string]interface{})
	context[constants.CTX_HARDWARE_FOLDERS] = []string{"downloaded_board_manager_stuff"}

	loader := builder.HardwareLoader{}
	err := loader.Run(context)
	NoError(t, err)

	packages := context[constants.CTX_HARDWARE].(map[string]*types.Package)
	require.Equal(t, 3, len(packages))
	require.NotNil(t, packages["arduino"])
	require.Equal(t, 1, len(packages["arduino"].Platforms))
	require.NotNil(t, packages["RedBearLab"])
	require.Equal(t, 1, len(packages["RedBearLab"].Platforms))
	require.NotNil(t, packages["RFduino"])
	require.Equal(t, 0, len(packages["RFduino"].Platforms))

	samdPlatform := packages["arduino"].Platforms["samd"]
	require.Equal(t, 2, len(samdPlatform.Boards))

	require.Equal(t, "arduino_zero_edbg", samdPlatform.Boards["arduino_zero_edbg"].BoardId)
	require.Equal(t, "arduino_zero_edbg", samdPlatform.Boards["arduino_zero_edbg"].Properties[constants.ID])

	require.Equal(t, "arduino_zero", samdPlatform.Boards["arduino_zero_native"].Properties["build.variant"])
	require.Equal(t, "-D__SAMD21G18A__ {build.usb_flags}", samdPlatform.Boards["arduino_zero_native"].Properties["build.extra_flags"])

	require.Equal(t, "Arduino SAMD (32-bits ARM Cortex-M0+) Boards", samdPlatform.Properties[constants.PLATFORM_NAME])
	require.Equal(t, "-d3", samdPlatform.Properties["tools.openocd.erase.params.verbose"])

	require.Equal(t, 3, len(samdPlatform.Programmers))

	require.Equal(t, "Atmel EDBG", samdPlatform.Programmers["edbg"][constants.PROGRAMMER_NAME])
	require.Equal(t, "openocd", samdPlatform.Programmers["edbg"]["program.tool"])

	avrRedBearPlatform := packages["RedBearLab"].Platforms["avr"]
	require.Equal(t, 3, len(avrRedBearPlatform.Boards))

	require.Equal(t, "blend", avrRedBearPlatform.Boards["blend"].BoardId)
	require.Equal(t, "blend", avrRedBearPlatform.Boards["blend"].Properties[constants.ID])
	require.Equal(t, "arduino:arduino", avrRedBearPlatform.Boards["blend"].Properties[constants.BUILD_PROPERTIES_BUILD_CORE])

}

func TestLoadLotsOfHardware(t *testing.T) {
	context := make(map[string]interface{})

	context[constants.CTX_HARDWARE_FOLDERS] = []string{"downloaded_board_manager_stuff", "downloaded_hardware", filepath.Join("..", "hardware"), "hardware", "user_hardware"}

	loader := builder.HardwareLoader{}
	err := loader.Run(context)
	NoError(t, err)

	packages := context[constants.CTX_HARDWARE].(map[string]*types.Package)

	if runtime.GOOS == "windows" {
		//a package is a symlink, and windows does not support them
		require.Equal(t, 4, len(packages))
	} else {
		require.Equal(t, 5, len(packages))
	}

	require.NotNil(t, packages["arduino"])
	require.NotNil(t, packages["my_avr_platform"])

	require.Equal(t, 3, len(packages["arduino"].Platforms))
	require.Equal(t, 20, len(packages["arduino"].Platforms["avr"].Boards))
	require.Equal(t, 2, len(packages["arduino"].Platforms["sam"].Boards))
	require.Equal(t, 2, len(packages["arduino"].Platforms["samd"].Boards))

	require.Equal(t, 1, len(packages["my_avr_platform"].Platforms))
	require.Equal(t, 2, len(packages["my_avr_platform"].Platforms["avr"].Boards))

	if runtime.GOOS != "windows" {
		require.Equal(t, 1, len(packages["my_symlinked_avr_platform"].Platforms))
		require.Equal(t, 2, len(packages["my_symlinked_avr_platform"].Platforms["avr"].Boards))
	}
}
