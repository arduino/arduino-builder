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
	"github.com/arduino/arduino-builder/builder/types"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestStoreBuildOptionsMap(t *testing.T) {
	context := make(map[string]interface{})

	buildPath := SetupBuildPath(t, context)
	defer os.RemoveAll(buildPath)

	context[constants.CTX_HARDWARE_FOLDERS] = "hardware"
	context[constants.CTX_TOOLS_FOLDERS] = "tools"
	context[constants.CTX_BUILT_IN_LIBRARIES_FOLDERS] = "built-in libraries"
	context[constants.CTX_OTHER_LIBRARIES_FOLDERS] = "libraries"
	context[constants.CTX_FQBN] = "fqbn"
	context[constants.CTX_SKETCH_LOCATION] = "sketchLocation"
	context[constants.CTX_VERBOSE] = true
	context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION] = "ideVersion"
	context[constants.CTX_DEBUG_LEVEL] = 5

	commands := []types.Command{
		&builder.CreateBuildOptionsMap{},
		&builder.StoreBuildOptionsMap{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	_, err := os.Stat(filepath.Join(buildPath, constants.BUILD_OPTIONS_FILE))
	NoError(t, err)

	bytes, err := ioutil.ReadFile(filepath.Join(buildPath, constants.BUILD_OPTIONS_FILE))
	NoError(t, err)

	require.Equal(t, "{\n"+
		"  \"builtInLibrariesFolders\": \"built-in libraries\",\n"+
		"  \"fqbn\": \"fqbn\",\n"+
		"  \"hardwareFolders\": \"hardware\",\n"+
		"  \"otherLibrariesFolders\": \"libraries\",\n"+
		"  \"runtime.ide.version\": \"ideVersion\",\n"+
		"  \"sketchLocation\": \"sketchLocation\",\n"+
		"  \"toolsFolders\": \"tools\"\n"+
		"}", string(bytes))
}
