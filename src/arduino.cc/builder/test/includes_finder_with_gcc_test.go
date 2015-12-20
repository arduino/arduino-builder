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
	"sort"
	"testing"
)

func TestIncludesFinderWithGCC(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})

	buildPath := SetupBuildPath(t, context)
	defer os.RemoveAll(buildPath)

	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools"}
	context[constants.CTX_FQBN] = "arduino:avr:leonardo"
	context[constants.CTX_SKETCH_LOCATION] = filepath.Join("sketch2", "SketchWithIfDef.ino")
	context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION] = "10600"
	context[constants.CTX_VERBOSE] = true

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	require.Nil(t, context[constants.CTX_INCLUDES])
}

func TestIncludesFinderWithGCCSketchWithConfig(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})

	buildPath := SetupBuildPath(t, context)
	defer os.RemoveAll(buildPath)

	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools"}
	context[constants.CTX_BUILT_IN_LIBRARIES_FOLDERS] = []string{"downloaded_libraries"}
	context[constants.CTX_OTHER_LIBRARIES_FOLDERS] = []string{"dependent_libraries", "libraries"}
	context[constants.CTX_FQBN] = "arduino:avr:leonardo"
	context[constants.CTX_SKETCH_LOCATION] = filepath.Join("sketch_with_config", "sketch_with_config.ino")
	context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION] = "10600"
	context[constants.CTX_VERBOSE] = true

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	require.NotNil(t, context[constants.CTX_INCLUDES])
	includes := context[constants.CTX_INCLUDES].([]string)
	require.Equal(t, 1, len(includes))
	require.True(t, utils.SliceContains(includes, "Bridge.h"))

	importedLibraries := context[constants.CTX_IMPORTED_LIBRARIES].([]*types.Library)
	require.Equal(t, 1, len(importedLibraries))
	require.Equal(t, "Bridge", importedLibraries[0].Name)
}

func TestIncludesFinderWithGCCSketchWithDependendLibraries(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})

	buildPath := SetupBuildPath(t, context)
	defer os.RemoveAll(buildPath)

	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools"}
	context[constants.CTX_OTHER_LIBRARIES_FOLDERS] = []string{"dependent_libraries"}
	context[constants.CTX_FQBN] = "arduino:avr:leonardo"
	context[constants.CTX_SKETCH_LOCATION] = filepath.Join("sketch_with_dependend_libraries", "sketch.ino")
	context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION] = "10600"
	context[constants.CTX_VERBOSE] = true

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	require.NotNil(t, context[constants.CTX_INCLUDES])
	includes := context[constants.CTX_INCLUDES].([]string)
	require.Equal(t, 4, len(includes))

	sort.Strings(includes)
	require.Equal(t, "library1.h", includes[0])
	require.Equal(t, "library2.h", includes[1])
	require.Equal(t, "library3.h", includes[2])
	require.Equal(t, "library4.h", includes[3])

	importedLibraries := context[constants.CTX_IMPORTED_LIBRARIES].([]*types.Library)
	require.Equal(t, 4, len(importedLibraries))

	sort.Sort(ByLibraryName(importedLibraries))
	require.Equal(t, "library1", importedLibraries[0].Name)
	require.Equal(t, "library2", importedLibraries[1].Name)
	require.Equal(t, "library3", importedLibraries[2].Name)
	require.Equal(t, "library4", importedLibraries[3].Name)
}

func TestIncludesFinderWithGCCSketchWithThatChecksIfSPIHasTransactions(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})

	buildPath := SetupBuildPath(t, context)
	defer os.RemoveAll(buildPath)

	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools"}
	context[constants.CTX_BUILT_IN_LIBRARIES_FOLDERS] = []string{"downloaded_libraries"}
	context[constants.CTX_OTHER_LIBRARIES_FOLDERS] = []string{"dependent_libraries", "libraries"}
	context[constants.CTX_FQBN] = "arduino:avr:leonardo"
	context[constants.CTX_SKETCH_LOCATION] = filepath.Join("sketch_that_checks_if_SPI_has_transactions", "sketch.ino")
	context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION] = "10600"
	context[constants.CTX_VERBOSE] = true

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	require.NotNil(t, context[constants.CTX_INCLUDES])
	includes := context[constants.CTX_INCLUDES].([]string)
	require.Equal(t, 1, len(includes))
	require.Equal(t, "SPI.h", includes[0])

	importedLibraries := context[constants.CTX_IMPORTED_LIBRARIES].([]*types.Library)
	require.Equal(t, 1, len(importedLibraries))
	require.Equal(t, "SPI", importedLibraries[0].Name)
}

func TestIncludesFinderWithGCCSketchWithThatChecksIfSPIHasTransactionsAndIncludesMissingLib(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})

	buildPath := SetupBuildPath(t, context)
	defer os.RemoveAll(buildPath)

	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools"}
	context[constants.CTX_BUILT_IN_LIBRARIES_FOLDERS] = []string{"downloaded_libraries"}
	context[constants.CTX_OTHER_LIBRARIES_FOLDERS] = []string{"dependent_libraries", "libraries"}
	context[constants.CTX_FQBN] = "arduino:avr:leonardo"
	context[constants.CTX_SKETCH_LOCATION] = filepath.Join("sketch_that_checks_if_SPI_has_transactions_and_includes_missing_Ethernet", "sketch.ino")
	context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION] = "10600"
	context[constants.CTX_VERBOSE] = true

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	command := &builder.ContainerFindIncludes{}
	err := command.Run(context)
	require.Error(t, err)

	require.NotNil(t, context[constants.CTX_INCLUDES])
	includes := context[constants.CTX_INCLUDES].([]string)
	require.Equal(t, 2, len(includes))
	sort.Strings(includes)
	require.Equal(t, "Inexistent.h", includes[0])
	require.Equal(t, "SPI.h", includes[1])

	importedLibraries := context[constants.CTX_IMPORTED_LIBRARIES].([]*types.Library)
	require.Equal(t, 1, len(importedLibraries))
	require.Equal(t, "SPI", importedLibraries[0].Name)
}
