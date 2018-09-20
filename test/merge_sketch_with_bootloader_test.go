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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/arduino/arduino-builder"
	"github.com/arduino/arduino-builder/constants"
	"github.com/arduino/arduino-builder/types"
	"github.com/arduino/arduino-builder/utils"
	"github.com/stretchr/testify/require"
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

	fakeSketchHex := ":100000000C9434000C9446000C9446000C9446006A\n" +
		":100010000C9446000C9446000C9446000C94460048\n" +
		":100020000C9446000C9446000C9446000C94460038\n" +
		":100030000C9446000C9446000C9446000C94460028\n" +
		":100040000C9448000C9446000C9446000C94460016\n" +
		":100050000C9446000C9446000C9446000C94460008\n" +
		":100060000C9446000C94460011241FBECFEFD8E03C\n" +
		":10007000DEBFCDBF21E0A0E0B1E001C01D92A930FC\n" +
		":10008000B207E1F70E9492000C94DC000C9400008F\n" +
		":100090001F920F920FB60F9211242F933F938F93BD\n" +
		":1000A0009F93AF93BF938091050190910601A0911A\n" +
		":1000B0000701B09108013091040123E0230F2D378F\n" +
		":1000C00020F40196A11DB11D05C026E8230F02965C\n" +
		":1000D000A11DB11D20930401809305019093060199\n" +
		":1000E000A0930701B0930801809100019091010154\n" +
		":1000F000A0910201B09103010196A11DB11D809351\n" +
		":10010000000190930101A0930201B0930301BF91FC\n" +
		":10011000AF919F918F913F912F910F900FBE0F90B4\n" +
		":100120001F901895789484B5826084BD84B58160F1\n" +
		":1001300084BD85B5826085BD85B5816085BD8091B2\n" +
		":100140006E00816080936E0010928100809181002A\n" +
		":100150008260809381008091810081608093810022\n" +
		":10016000809180008160809380008091B1008460E4\n" +
		":100170008093B1008091B00081608093B000809145\n" +
		":100180007A00846080937A0080917A008260809304\n" +
		":100190007A0080917A00816080937A0080917A0061\n" +
		":1001A000806880937A001092C100C0E0D0E0209770\n" +
		":0C01B000F1F30E940000FBCFF894FFCF99\n" +
		":00000001FF\n"
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

	require.True(t, strings.HasPrefix(mergedSketchHex, ":100000000C9434000C9446000C9446000C9446006A\n"))
	require.True(t, strings.HasSuffix(mergedSketchHex, ":00000001FF\n"))
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

	fakeSketchHex := ":100000000C9434000C9446000C9446000C9446006A\n" +
		":100010000C9446000C9446000C9446000C94460048\n" +
		":100020000C9446000C9446000C9446000C94460038\n" +
		":100030000C9446000C9446000C9446000C94460028\n" +
		":100040000C9448000C9446000C9446000C94460016\n" +
		":100050000C9446000C9446000C9446000C94460008\n" +
		":100060000C9446000C94460011241FBECFEFD8E03C\n" +
		":10007000DEBFCDBF21E0A0E0B1E001C01D92A930FC\n" +
		":10008000B207E1F70E9492000C94DC000C9400008F\n" +
		":100090001F920F920FB60F9211242F933F938F93BD\n" +
		":1000A0009F93AF93BF938091050190910601A0911A\n" +
		":1000B0000701B09108013091040123E0230F2D378F\n" +
		":1000C00020F40196A11DB11D05C026E8230F02965C\n" +
		":1000D000A11DB11D20930401809305019093060199\n" +
		":1000E000A0930701B0930801809100019091010154\n" +
		":1000F000A0910201B09103010196A11DB11D809351\n" +
		":10010000000190930101A0930201B0930301BF91FC\n" +
		":10011000AF919F918F913F912F910F900FBE0F90B4\n" +
		":100120001F901895789484B5826084BD84B58160F1\n" +
		":1001300084BD85B5826085BD85B5816085BD8091B2\n" +
		":100140006E00816080936E0010928100809181002A\n" +
		":100150008260809381008091810081608093810022\n" +
		":10016000809180008160809380008091B1008460E4\n" +
		":100170008093B1008091B00081608093B000809145\n" +
		":100180007A00846080937A0080917A008260809304\n" +
		":100190007A0080917A00816080937A0080917A0061\n" +
		":1001A000806880937A001092C100C0E0D0E0209770\n" +
		":0C01B000F1F30E940000FBCFF894FFCF99\n" +
		":00000001FF\n"
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

	require.True(t, strings.HasPrefix(mergedSketchHex, ":100000000C9434000C9446000C9446000C9446006A\n"))
	require.True(t, strings.HasSuffix(mergedSketchHex, ":00000001FF\n"))
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

	fakeSketchHex := ":100000000C9434000C9446000C9446000C9446006A\n" +
		":100010000C9446000C9446000C9446000C94460048\n" +
		":100020000C9446000C9446000C9446000C94460038\n" +
		":100030000C9446000C9446000C9446000C94460028\n" +
		":100040000C9448000C9446000C9446000C94460016\n" +
		":100050000C9446000C9446000C9446000C94460008\n" +
		":100060000C9446000C94460011241FBECFEFD8E03C\n" +
		":10007000DEBFCDBF21E0A0E0B1E001C01D92A930FC\n" +
		":10008000B207E1F70E9492000C94DC000C9400008F\n" +
		":100090001F920F920FB60F9211242F933F938F93BD\n" +
		":1000A0009F93AF93BF938091050190910601A0911A\n" +
		":1000B0000701B09108013091040123E0230F2D378F\n" +
		":1000C00020F40196A11DB11D05C026E8230F02965C\n" +
		":1000D000A11DB11D20930401809305019093060199\n" +
		":1000E000A0930701B0930801809100019091010154\n" +
		":1000F000A0910201B09103010196A11DB11D809351\n" +
		":10010000000190930101A0930201B0930301BF91FC\n" +
		":10011000AF919F918F913F912F910F900FBE0F90B4\n" +
		":100120001F901895789484B5826084BD84B58160F1\n" +
		":1001300084BD85B5826085BD85B5816085BD8091B2\n" +
		":100140006E00816080936E0010928100809181002A\n" +
		":100150008260809381008091810081608093810022\n" +
		":10016000809180008160809380008091B1008460E4\n" +
		":100170008093B1008091B00081608093B000809145\n" +
		":100180007A00846080937A0080917A008260809304\n" +
		":100190007A0080917A00816080937A0080917A0061\n" +
		":1001A000806880937A001092C100C0E0D0E0209770\n" +
		":0C01B000F1F30E940000FBCFF894FFCF99\n" +
		":00000001FF\n"
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

	require.True(t, strings.HasPrefix(mergedSketchHex, ":100000000C9434000C9446000C9446000C9446006A\n"))
	require.True(t, strings.HasSuffix(mergedSketchHex, ":00000001FF\n"))
}
