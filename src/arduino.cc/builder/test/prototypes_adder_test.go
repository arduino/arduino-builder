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
 * Copyright 2015 Matthijs Kooijman
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

func TestPrototypesAdderBridgeExample(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})

	buildPath := SetupBuildPath(t, context)
	defer os.RemoveAll(buildPath)

	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools"}
	context[constants.CTX_FQBN] = "arduino:avr:leonardo"
	context[constants.CTX_SKETCH_LOCATION] = filepath.Join("downloaded_libraries", "Bridge", "examples", "Bridge", "Bridge.ino")
	context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION] = "10600"
	context[constants.CTX_BUILT_IN_LIBRARIES_FOLDERS] = []string{"downloaded_libraries"}
	context[constants.CTX_OTHER_LIBRARIES_FOLDERS] = []string{"libraries"}
	context[constants.CTX_VERBOSE] = false

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},

		&builder.PrintUsedLibrariesIfVerbose{},
		&builder.WarnAboutArchIncompatibleLibraries{},

		&builder.ContainerAddPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	require.Equal(t, "#include <Arduino.h>\n#line 1\n", context[constants.CTX_INCLUDE_SECTION].(string))
	require.Equal(t, "void setup();\nvoid loop();\nvoid process(BridgeClient client);\nvoid digitalCommand(BridgeClient client);\nvoid analogCommand(BridgeClient client);\nvoid modeCommand(BridgeClient client);\n#line 33\n", context[constants.CTX_PROTOTYPE_SECTION].(string))
}

func TestPrototypesAdderSketchWithIfDef(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})

	buildPath := SetupBuildPath(t, context)
	defer os.RemoveAll(buildPath)

	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools"}
	context[constants.CTX_FQBN] = "arduino:avr:leonardo"
	context[constants.CTX_SKETCH_LOCATION] = filepath.Join("sketch2", "SketchWithIfDef.ino")
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

		&builder.ContainerAddPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	preprocessed := LoadAndInterpolate(t, filepath.Join("sketch2", "SketchWithIfDef.preprocessed.txt"), context)
	require.Equal(t, preprocessed, strings.Replace(context[constants.CTX_SOURCE].(string), "\r\n", "\n", -1))
}

func TestPrototypesAdderBaladuino(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})

	buildPath := SetupBuildPath(t, context)
	defer os.RemoveAll(buildPath)

	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools"}
	context[constants.CTX_FQBN] = "arduino:avr:leonardo"
	context[constants.CTX_SKETCH_LOCATION] = filepath.Join("sketch3", "Baladuino.ino")
	context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION] = "10600"
	context[constants.CTX_BUILT_IN_LIBRARIES_FOLDERS] = []string{"downloaded_libraries"}
	context[constants.CTX_OTHER_LIBRARIES_FOLDERS] = []string{"libraries"}
	context[constants.CTX_VERBOSE] = false

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},

		&builder.PrintUsedLibrariesIfVerbose{},
		&builder.WarnAboutArchIncompatibleLibraries{},

		&builder.ContainerAddPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	preprocessed := LoadAndInterpolate(t, filepath.Join("sketch3", "Baladuino.preprocessed.txt"), context)
	require.Equal(t, preprocessed, strings.Replace(context[constants.CTX_SOURCE].(string), "\r\n", "\n", -1))
}

func TestPrototypesAdderCharWithEscapedDoubleQuote(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})

	buildPath := SetupBuildPath(t, context)
	defer os.RemoveAll(buildPath)

	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools"}
	context[constants.CTX_FQBN] = "arduino:avr:leonardo"
	context[constants.CTX_SKETCH_LOCATION] = filepath.Join("sketch4", "CharWithEscapedDoubleQuote.ino")
	context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION] = "10600"
	context[constants.CTX_BUILT_IN_LIBRARIES_FOLDERS] = []string{"downloaded_libraries"}
	context[constants.CTX_OTHER_LIBRARIES_FOLDERS] = []string{"libraries"}
	context[constants.CTX_VERBOSE] = false

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},

		&builder.PrintUsedLibrariesIfVerbose{},
		&builder.WarnAboutArchIncompatibleLibraries{},

		&builder.ContainerAddPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	preprocessed := LoadAndInterpolate(t, filepath.Join("sketch4", "CharWithEscapedDoubleQuote.preprocessed.txt"), context)
	require.Equal(t, preprocessed, strings.Replace(context[constants.CTX_SOURCE].(string), "\r\n", "\n", -1))
}

func TestPrototypesAdderIncludeBetweenMultilineComment(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})

	buildPath := SetupBuildPath(t, context)
	defer os.RemoveAll(buildPath)

	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools"}
	context[constants.CTX_FQBN] = "arduino:sam:arduino_due_x_dbg"
	context[constants.CTX_SKETCH_LOCATION] = filepath.Join("sketch5", "IncludeBetweenMultilineComment.ino")
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

		&builder.ContainerAddPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	preprocessed := LoadAndInterpolate(t, filepath.Join("sketch5", "IncludeBetweenMultilineComment.preprocessed.txt"), context)
	require.Equal(t, preprocessed, strings.Replace(context[constants.CTX_SOURCE].(string), "\r\n", "\n", -1))
}

func TestPrototypesAdderLineContinuations(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})

	buildPath := SetupBuildPath(t, context)
	defer os.RemoveAll(buildPath)

	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools"}
	context[constants.CTX_FQBN] = "arduino:avr:leonardo"
	context[constants.CTX_SKETCH_LOCATION] = filepath.Join("sketch6", "/LineContinuations.ino")
	context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION] = "10600"
	context[constants.CTX_BUILT_IN_LIBRARIES_FOLDERS] = []string{"downloaded_libraries"}
	context[constants.CTX_OTHER_LIBRARIES_FOLDERS] = []string{"libraries"}
	context[constants.CTX_VERBOSE] = false

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},

		&builder.PrintUsedLibrariesIfVerbose{},
		&builder.WarnAboutArchIncompatibleLibraries{},

		&builder.ContainerAddPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	preprocessed := LoadAndInterpolate(t, filepath.Join("sketch6", "LineContinuations.preprocessed.txt"), context)
	require.Equal(t, preprocessed, strings.Replace(context[constants.CTX_SOURCE].(string), "\r\n", "\n", -1))
}

func TestPrototypesAdderStringWithComment(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})

	buildPath := SetupBuildPath(t, context)
	defer os.RemoveAll(buildPath)

	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools"}
	context[constants.CTX_FQBN] = "arduino:avr:leonardo"
	context[constants.CTX_SKETCH_LOCATION] = filepath.Join("sketch7", "StringWithComment.ino")
	context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION] = "10600"
	context[constants.CTX_BUILT_IN_LIBRARIES_FOLDERS] = []string{"downloaded_libraries"}
	context[constants.CTX_OTHER_LIBRARIES_FOLDERS] = []string{"libraries"}
	context[constants.CTX_VERBOSE] = false

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},

		&builder.PrintUsedLibrariesIfVerbose{},
		&builder.WarnAboutArchIncompatibleLibraries{},

		&builder.ContainerAddPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	preprocessed := LoadAndInterpolate(t, filepath.Join("sketch7", "StringWithComment.preprocessed.txt"), context)
	require.Equal(t, preprocessed, strings.Replace(context[constants.CTX_SOURCE].(string), "\r\n", "\n", -1))
}

func TestPrototypesAdderSketchWithStruct(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})

	buildPath := SetupBuildPath(t, context)
	defer os.RemoveAll(buildPath)

	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools"}
	context[constants.CTX_FQBN] = "arduino:avr:leonardo"
	context[constants.CTX_SKETCH_LOCATION] = filepath.Join("sketch8", "SketchWithStruct.ino")
	context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION] = "10600"
	context[constants.CTX_BUILT_IN_LIBRARIES_FOLDERS] = []string{"downloaded_libraries"}
	context[constants.CTX_OTHER_LIBRARIES_FOLDERS] = []string{"libraries"}
	context[constants.CTX_VERBOSE] = false

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},

		&builder.PrintUsedLibrariesIfVerbose{},
		&builder.WarnAboutArchIncompatibleLibraries{},

		&builder.ContainerAddPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	preprocessed := LoadAndInterpolate(t, filepath.Join("sketch8", "SketchWithStruct.preprocessed.txt"), context)
	require.Equal(t, preprocessed, strings.Replace(context[constants.CTX_SOURCE].(string), "\r\n", "\n", -1))
}

func TestPrototypesAdderSketchWithConfig(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})

	buildPath := SetupBuildPath(t, context)
	defer os.RemoveAll(buildPath)

	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools"}
	context[constants.CTX_FQBN] = "arduino:avr:leonardo"
	context[constants.CTX_SKETCH_LOCATION] = filepath.Join("sketch_with_config", "sketch_with_config.ino")
	context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION] = "10600"
	context[constants.CTX_BUILT_IN_LIBRARIES_FOLDERS] = []string{"downloaded_libraries"}
	context[constants.CTX_OTHER_LIBRARIES_FOLDERS] = []string{"libraries"}
	context[constants.CTX_VERBOSE] = false

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},

		&builder.PrintUsedLibrariesIfVerbose{},
		&builder.WarnAboutArchIncompatibleLibraries{},

		&builder.ContainerAddPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	require.Equal(t, "#include <Arduino.h>\n#line 1\n", context[constants.CTX_INCLUDE_SECTION].(string))
	require.Equal(t, "void setup();\nvoid loop();\n#line 13\n", context[constants.CTX_PROTOTYPE_SECTION].(string))

	preprocessed := LoadAndInterpolate(t, filepath.Join("sketch_with_config", "sketch_with_config.preprocessed.txt"), context)
	require.Equal(t, preprocessed, strings.Replace(context[constants.CTX_SOURCE].(string), "\r\n", "\n", -1))
}

func TestPrototypesAdderSketchNoFunctionsTwoFiles(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})

	buildPath := SetupBuildPath(t, context)
	defer os.RemoveAll(buildPath)

	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools"}
	context[constants.CTX_FQBN] = "arduino:avr:leonardo"
	context[constants.CTX_SKETCH_LOCATION] = filepath.Join("sketch_no_functions_two_files", "main.ino")
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

		&builder.ContainerAddPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	require.Nil(t, context[constants.CTX_INCLUDE_SECTION])
	require.Nil(t, context[constants.CTX_PROTOTYPE_SECTION])
}

func TestPrototypesAdderSketchNoFunctions(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})

	buildPath := SetupBuildPath(t, context)
	defer os.RemoveAll(buildPath)

	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools"}
	context[constants.CTX_FQBN] = "arduino:avr:leonardo"
	context[constants.CTX_SKETCH_LOCATION] = filepath.Join("sketch_no_functions", "main.ino")
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

		&builder.ContainerAddPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	require.Nil(t, context[constants.CTX_INCLUDE_SECTION])
	require.Nil(t, context[constants.CTX_PROTOTYPE_SECTION])
}

func TestPrototypesAdderSketchWithDefaultArgs(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})

	buildPath := SetupBuildPath(t, context)
	defer os.RemoveAll(buildPath)

	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools"}
	context[constants.CTX_FQBN] = "arduino:avr:leonardo"
	context[constants.CTX_SKETCH_LOCATION] = filepath.Join("sketch_with_default_args", "sketch.ino")
	context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION] = "10600"
	context[constants.CTX_BUILT_IN_LIBRARIES_FOLDERS] = []string{"downloaded_libraries"}
	context[constants.CTX_OTHER_LIBRARIES_FOLDERS] = []string{"libraries"}
	context[constants.CTX_VERBOSE] = false

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},

		&builder.PrintUsedLibrariesIfVerbose{},
		&builder.WarnAboutArchIncompatibleLibraries{},

		&builder.ContainerAddPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	require.Equal(t, "#include <Arduino.h>\n#line 1\n", context[constants.CTX_INCLUDE_SECTION].(string))
	require.Equal(t, "void setup();\nvoid loop();\n#line 1\n", context[constants.CTX_PROTOTYPE_SECTION].(string))
}

func TestPrototypesAdderSketchWithInlineFunction(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})

	buildPath := SetupBuildPath(t, context)
	defer os.RemoveAll(buildPath)

	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools"}
	context[constants.CTX_FQBN] = "arduino:avr:leonardo"
	context[constants.CTX_SKETCH_LOCATION] = filepath.Join("sketch_with_inline_function", "sketch.ino")
	context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION] = "10600"
	context[constants.CTX_BUILT_IN_LIBRARIES_FOLDERS] = []string{"downloaded_libraries"}
	context[constants.CTX_OTHER_LIBRARIES_FOLDERS] = []string{"libraries"}
	context[constants.CTX_VERBOSE] = false

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},

		&builder.PrintUsedLibrariesIfVerbose{},
		&builder.WarnAboutArchIncompatibleLibraries{},

		&builder.ContainerAddPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	require.Equal(t, "#include <Arduino.h>\n#line 1\n", context[constants.CTX_INCLUDE_SECTION].(string))
	require.Equal(t, "void setup();\nvoid loop();\nshort unsigned int testInt();\nint8_t testInline();\nuint8_t testAttribute();\n#line 1\n", context[constants.CTX_PROTOTYPE_SECTION].(string))
}

func TestPrototypesAdderSketchWithFunctionSignatureInsideIFDEF(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})

	buildPath := SetupBuildPath(t, context)
	defer os.RemoveAll(buildPath)

	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools"}
	context[constants.CTX_FQBN] = "arduino:avr:leonardo"
	context[constants.CTX_SKETCH_LOCATION] = filepath.Join("sketch_with_function_signature_inside_ifdef", "sketch.ino")
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

		&builder.ContainerAddPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	require.Equal(t, "#include <Arduino.h>\n#line 1\n", context[constants.CTX_INCLUDE_SECTION].(string))
	require.Equal(t, "void setup();\nvoid loop();\nint8_t adalight();\n#line 1\n", context[constants.CTX_PROTOTYPE_SECTION].(string))
}

func TestPrototypesAdderSketchWithUSBCON(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})

	buildPath := SetupBuildPath(t, context)
	defer os.RemoveAll(buildPath)

	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools"}
	context[constants.CTX_FQBN] = "arduino:avr:leonardo"
	context[constants.CTX_SKETCH_LOCATION] = filepath.Join("sketch_with_usbcon", "sketch.ino")
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

		&builder.ContainerAddPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	require.Equal(t, "#include <Arduino.h>\n#line 1\n", context[constants.CTX_INCLUDE_SECTION].(string))
	require.Equal(t, "void ciao();\nvoid setup();\nvoid loop();\n#line 3\n", context[constants.CTX_PROTOTYPE_SECTION].(string))
}

func TestPrototypesAdderSketchWithTypename(t *testing.T) {
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

		&builder.ContainerAddPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	require.Equal(t, "#include <Arduino.h>\n#line 1\n", context[constants.CTX_INCLUDE_SECTION].(string))
	require.Equal(t, "void setup();\nvoid loop();\n#line 6\n", context[constants.CTX_PROTOTYPE_SECTION].(string))
}
