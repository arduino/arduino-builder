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
	"arduino.cc/builder/ctags"
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

	sketchLocation := Abs(t, filepath.Join("downloaded_libraries", "Bridge", "examples", "Bridge", "Bridge.ino"))
	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools"}
	context[constants.CTX_FQBN] = "arduino:avr:leonardo"
	context[constants.CTX_SKETCH_LOCATION] = sketchLocation
	context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION] = "10600"
	context[constants.CTX_BUILT_IN_LIBRARIES_FOLDERS] = []string{"downloaded_libraries"}
	context[constants.CTX_OTHER_LIBRARIES_FOLDERS] = []string{"libraries"}
	context[constants.CTX_VERBOSE] = true

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},

		&builder.PrintUsedLibrariesIfVerbose{},
		&builder.WarnAboutArchIncompatibleLibraries{},
		&builder.CTagsTargetFileSaver{SourceField: constants.CTX_SOURCE, TargetFileName: constants.FILE_CTAGS_TARGET},
		&ctags.CTagsRunner{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	sketchLocation = strings.Replace(sketchLocation, "\\", "\\\\", -1)
	expectedOutput := "server	" + sketchLocation + "	/^BridgeServer server;$/;\"	kind:variable	line:31\n" +
		"setup	" + sketchLocation + "	/^void setup() {$/;\"	kind:function	line:33	signature:()	returntype:void\n" +
		"loop	" + sketchLocation + "	/^void loop() {$/;\"	kind:function	line:46	signature:()	returntype:void\n" +
		"process	" + sketchLocation + "	/^void process(BridgeClient client) {$/;\"	kind:function	line:62	signature:(BridgeClient client)	returntype:void\n" +
		"digitalCommand	" + sketchLocation + "	/^void digitalCommand(BridgeClient client) {$/;\"	kind:function	line:82	signature:(BridgeClient client)	returntype:void\n" +
		"analogCommand	" + sketchLocation + "	/^void analogCommand(BridgeClient client) {$/;\"	kind:function	line:109	signature:(BridgeClient client)	returntype:void\n" +
		"modeCommand	" + sketchLocation + "	/^void modeCommand(BridgeClient client) {$/;\"	kind:function	line:149	signature:(BridgeClient client)	returntype:void\n"

	require.Equal(t, expectedOutput, strings.Replace(context[constants.CTX_CTAGS_OUTPUT].(string), "\r\n", "\n", -1))
}

func TestCTagsRunnerSketchWithClass(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})

	buildPath := SetupBuildPath(t, context)
	defer os.RemoveAll(buildPath)

	sketchLocation := Abs(t, filepath.Join("sketch_with_class", "sketch.ino"))
	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools"}
	context[constants.CTX_FQBN] = "arduino:avr:leonardo"
	context[constants.CTX_SKETCH_LOCATION] = sketchLocation
	context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION] = "10600"
	context[constants.CTX_BUILT_IN_LIBRARIES_FOLDERS] = []string{"downloaded_libraries"}
	context[constants.CTX_OTHER_LIBRARIES_FOLDERS] = []string{"libraries"}
	context[constants.CTX_VERBOSE] = true

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},

		&builder.PrintUsedLibrariesIfVerbose{},
		&builder.WarnAboutArchIncompatibleLibraries{},
		&builder.CTagsTargetFileSaver{SourceField: constants.CTX_SOURCE, TargetFileName: constants.FILE_CTAGS_TARGET},
		&ctags.CTagsRunner{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	sketchLocation = strings.Replace(sketchLocation, "\\", "\\\\", -1)
	expectedOutput := "set_values\t" + sketchLocation + "\t/^    void set_values (int,int);$/;\"\tkind:prototype\tline:4\tclass:Rectangle\tsignature:(int,int)\treturntype:void\n" +
		"area\t" + sketchLocation + "\t/^    int area() {return width*height;}$/;\"\tkind:function\tline:5\tclass:Rectangle\tsignature:()\treturntype:int\n" +
		"set_values\t" + sketchLocation + "\t/^void Rectangle::set_values (int x, int y) {$/;\"\tkind:function\tline:8\tclass:Rectangle\tsignature:(int x, int y)\treturntype:void\n" +
		"setup\t" + sketchLocation + "\t/^void setup() {$/;\"\tkind:function\tline:13\tsignature:()\treturntype:void\n" +
		"loop\t" + sketchLocation + "\t/^void loop() {$/;\"\tkind:function\tline:17\tsignature:()\treturntype:void\n"

	require.Equal(t, expectedOutput, strings.Replace(context[constants.CTX_CTAGS_OUTPUT].(string), "\r\n", "\n", -1))
}

func TestCTagsRunnerSketchWithTypename(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})

	buildPath := SetupBuildPath(t, context)
	defer os.RemoveAll(buildPath)

	sketchLocation := Abs(t, filepath.Join("sketch_with_typename", "sketch.ino"))
	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools"}
	context[constants.CTX_FQBN] = "arduino:avr:leonardo"
	context[constants.CTX_SKETCH_LOCATION] = sketchLocation
	context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION] = "10600"
	context[constants.CTX_BUILT_IN_LIBRARIES_FOLDERS] = []string{"downloaded_libraries"}
	context[constants.CTX_OTHER_LIBRARIES_FOLDERS] = []string{"libraries"}
	context[constants.CTX_VERBOSE] = true

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},

		&builder.PrintUsedLibrariesIfVerbose{},
		&builder.WarnAboutArchIncompatibleLibraries{},
		&builder.CTagsTargetFileSaver{SourceField: constants.CTX_SOURCE, TargetFileName: constants.FILE_CTAGS_TARGET},
		&ctags.CTagsRunner{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	sketchLocation = strings.Replace(sketchLocation, "\\", "\\\\", -1)
	expectedOutput := "Foo\t" + sketchLocation + "\t/^  struct Foo{$/;\"\tkind:struct\tline:2\n" +
		"setup\t" + sketchLocation + "\t/^void setup() {$/;\"\tkind:function\tline:6\tsignature:()\treturntype:void\n" +
		"loop\t" + sketchLocation + "\t/^void loop() {}$/;\"\tkind:function\tline:10\tsignature:()\treturntype:void\n" +
		"func\t" + sketchLocation + "\t/^typename Foo<char>::Bar func(){$/;\"\tkind:function\tline:12\tsignature:()\treturntype:Foo::Bar\n"

	require.Equal(t, expectedOutput, strings.Replace(context[constants.CTX_CTAGS_OUTPUT].(string), "\r\n", "\n", -1))
}

func TestCTagsRunnerSketchWithNamespace(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})

	buildPath := SetupBuildPath(t, context)
	defer os.RemoveAll(buildPath)

	sketchLocation := Abs(t, filepath.Join("sketch_with_namespace", "sketch.ino"))
	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools"}
	context[constants.CTX_FQBN] = "arduino:avr:leonardo"
	context[constants.CTX_SKETCH_LOCATION] = sketchLocation
	context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION] = "10600"
	context[constants.CTX_BUILT_IN_LIBRARIES_FOLDERS] = []string{"downloaded_libraries"}
	context[constants.CTX_OTHER_LIBRARIES_FOLDERS] = []string{"libraries"}
	context[constants.CTX_VERBOSE] = true

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},

		&builder.PrintUsedLibrariesIfVerbose{},
		&builder.WarnAboutArchIncompatibleLibraries{},
		&builder.CTagsTargetFileSaver{SourceField: constants.CTX_SOURCE, TargetFileName: constants.FILE_CTAGS_TARGET},
		&ctags.CTagsRunner{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	sketchLocation = strings.Replace(sketchLocation, "\\", "\\\\", -1)
	expectedOutput := "value\t" + sketchLocation + "\t/^\tint value() {$/;\"\tkind:function\tline:2\tnamespace:Test\tsignature:()\treturntype:int\n" +
		"setup\t" + sketchLocation + "\t/^void setup() {}$/;\"\tkind:function\tline:7\tsignature:()\treturntype:void\n" +
		"loop\t" + sketchLocation + "\t/^void loop() {}$/;\"\tkind:function\tline:8\tsignature:()\treturntype:void\n"

	require.Equal(t, expectedOutput, strings.Replace(context[constants.CTX_CTAGS_OUTPUT].(string), "\r\n", "\n", -1))
}

func TestCTagsRunnerSketchWithTemplates(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})

	buildPath := SetupBuildPath(t, context)
	defer os.RemoveAll(buildPath)

	sketchLocation := Abs(t, filepath.Join("sketch_with_templates_and_shift", "template_and_shift.cpp"))
	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools"}
	context[constants.CTX_FQBN] = "arduino:avr:leonardo"
	context[constants.CTX_SKETCH_LOCATION] = sketchLocation
	context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION] = "10600"
	context[constants.CTX_BUILT_IN_LIBRARIES_FOLDERS] = []string{"downloaded_libraries"}
	context[constants.CTX_OTHER_LIBRARIES_FOLDERS] = []string{"libraries"}
	context[constants.CTX_VERBOSE] = true

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},

		&builder.PrintUsedLibrariesIfVerbose{},
		&builder.WarnAboutArchIncompatibleLibraries{},
		&builder.CTagsTargetFileSaver{SourceField: constants.CTX_SOURCE, TargetFileName: constants.FILE_CTAGS_TARGET},
		&ctags.CTagsRunner{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	sketchLocation = strings.Replace(sketchLocation, "\\", "\\\\", -1)
	expectedOutput := "printGyro\t" + sketchLocation + "\t/^void printGyro()$/;\"\tkind:function\tline:10\tsignature:()\treturntype:void\n" +
		"bVar\t" + sketchLocation + "\t/^c< 8 > bVar;$/;\"\tkind:variable\tline:15\n" +
		"aVar\t" + sketchLocation + "\t/^c< 1<<8 > aVar;$/;\"\tkind:variable\tline:16\n" +
		"func\t" + sketchLocation + "\t/^template<int X> func( c< 1<<X> & aParam) {$/;\"\tkind:function\tline:18\tsignature:( c< 1<<X> & aParam)\treturntype:template\n"

	require.Equal(t, expectedOutput, strings.Replace(context[constants.CTX_CTAGS_OUTPUT].(string), "\r\n", "\n", -1))
}
