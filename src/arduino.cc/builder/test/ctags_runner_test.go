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
	"path/filepath"
	"strings"
	"testing"
)

func TestCTagsRunner(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})

	buildPath := SetupBuildPath(t, context)
	defer os.RemoveAll(buildPath)

	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools"}
	context[constants.CTX_FQBN] = "arduino:avr:leonardo"
	context[constants.CTX_SKETCH_LOCATION] = filepath.Join("downloaded_libraries", "Bridge", "examples", "Bridge", "Bridge.ino")
	context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION] = "10600"
	context[constants.CTX_LIBRARIES_FOLDERS] = []string{"libraries", "downloaded_libraries"}
	context[constants.CTX_VERBOSE] = false

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},

		&builder.PrintUsedLibrariesIfVerbose{},
		&builder.WarnAboutArchIncompatibleLibraries{},
		&builder.CTagsTargetFileSaver{SourceField: constants.CTX_SOURCE},
		&builder.CTagsRunner{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	ctagsTempFileName := context[constants.CTX_CTAGS_TEMP_FILE_NAME].(string)
	expectedOutput := "server	" + ctagsTempFileName + "	/^BridgeServer server;$/;\"	kind:variable	line:32\n" +
		"setup	" + ctagsTempFileName + "	/^void setup() {$/;\"	kind:function	line:34	signature:()	returntype:void\n" +
		"loop	" + ctagsTempFileName + "	/^void loop() {$/;\"	kind:function	line:47	signature:()	returntype:void\n" +
		"process	" + ctagsTempFileName + "	/^void process(BridgeClient client) {$/;\"	kind:function	line:63	signature:(BridgeClient client)	returntype:void\n" +
		"digitalCommand	" + ctagsTempFileName + "	/^void digitalCommand(BridgeClient client) {$/;\"	kind:function	line:83	signature:(BridgeClient client)	returntype:void\n" +
		"analogCommand	" + ctagsTempFileName + "	/^void analogCommand(BridgeClient client) {$/;\"	kind:function	line:110	signature:(BridgeClient client)	returntype:void\n" +
		"modeCommand	" + ctagsTempFileName + "	/^void modeCommand(BridgeClient client) {$/;\"	kind:function	line:150	signature:(BridgeClient client)	returntype:void\n"

	require.Equal(t, expectedOutput, strings.Replace(context[constants.CTX_CTAGS_OUTPUT].(string), "\r\n", "\n", -1))
}

func TestCTagsRunnerSketchWithClass(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})

	buildPath := SetupBuildPath(t, context)
	defer os.RemoveAll(buildPath)

	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools"}
	context[constants.CTX_FQBN] = "arduino:avr:leonardo"
	context[constants.CTX_SKETCH_LOCATION] = filepath.Join("sketch_with_class", "sketch.ino")
	context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION] = "10600"
	context[constants.CTX_LIBRARIES_FOLDERS] = []string{"libraries", "downloaded_libraries"}
	context[constants.CTX_VERBOSE] = false

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},

		&builder.PrintUsedLibrariesIfVerbose{},
		&builder.WarnAboutArchIncompatibleLibraries{},
		&builder.CTagsTargetFileSaver{SourceField: constants.CTX_SOURCE},
		&builder.CTagsRunner{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	ctagsTempFileName := context[constants.CTX_CTAGS_TEMP_FILE_NAME].(string)
	expectedOutput := "set_values\t" + ctagsTempFileName + "\t/^    void set_values (int,int);$/;\"\tkind:prototype\tline:5\tclass:Rectangle\tsignature:(int,int)\treturntype:void\n" +
		"area\t" + ctagsTempFileName + "\t/^    int area() {return width*height;}$/;\"\tkind:function\tline:6\tclass:Rectangle\tsignature:()\treturntype:int\n" +
		"set_values\t" + ctagsTempFileName + "\t/^void Rectangle::set_values (int x, int y) {$/;\"\tkind:function\tline:9\tclass:Rectangle\tsignature:(int x, int y)\treturntype:void\n" +
		"setup\t" + ctagsTempFileName + "\t/^void setup() {$/;\"\tkind:function\tline:14\tsignature:()\treturntype:void\n" +
		"loop\t" + ctagsTempFileName + "\t/^void loop() {$/;\"\tkind:function\tline:18\tsignature:()\treturntype:void\n"

	require.Equal(t, expectedOutput, strings.Replace(context[constants.CTX_CTAGS_OUTPUT].(string), "\r\n", "\n", -1))
}

func TestCTagsRunnerSketchWithTypename(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})

	buildPath := SetupBuildPath(t, context)
	defer os.RemoveAll(buildPath)

	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools"}
	context[constants.CTX_FQBN] = "arduino:avr:leonardo"
	context[constants.CTX_SKETCH_LOCATION] = filepath.Join("sketch_with_typename", "sketch.ino")
	context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION] = "10600"
	context[constants.CTX_LIBRARIES_FOLDERS] = []string{"libraries", "downloaded_libraries"}
	context[constants.CTX_VERBOSE] = true

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},

		&builder.PrintUsedLibrariesIfVerbose{},
		&builder.WarnAboutArchIncompatibleLibraries{},
		&builder.CTagsTargetFileSaver{SourceField: constants.CTX_SOURCE},
		&builder.CTagsRunner{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	ctagsTempFileName := context[constants.CTX_CTAGS_TEMP_FILE_NAME].(string)
	expectedOutput := "Foo\t" + ctagsTempFileName + "\t/^  struct Foo{$/;\"\tkind:struct\tline:3\n" +
		"setup\t" + ctagsTempFileName + "\t/^void setup() {$/;\"\tkind:function\tline:7\tsignature:()\treturntype:void\n" +
		"loop\t" + ctagsTempFileName + "\t/^void loop() {}$/;\"\tkind:function\tline:11\tsignature:()\treturntype:void\n" +
		"func\t" + ctagsTempFileName + "\t/^typename Foo<char>::Bar func(){$/;\"\tkind:function\tline:13\tsignature:()\treturntype:Foo::Bar\n"

	require.Equal(t, expectedOutput, strings.Replace(context[constants.CTX_CTAGS_OUTPUT].(string), "\r\n", "\n", -1))
}
