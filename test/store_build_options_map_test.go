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
	"github.com/arduino/arduino-builder"
	"github.com/arduino/arduino-builder/constants"
	"github.com/arduino/arduino-builder/types"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestStoreBuildOptionsMap(t *testing.T) {
	ctx := &types.Context{
		HardwareFolders:         []string{"hardware"},
		ToolsFolders:            []string{"tools"},
		BuiltInLibrariesFolders: []string{"built-in libraries"},
		OtherLibrariesFolders:   []string{"libraries"},
		SketchLocation:          "sketchLocation",
		FQBN:                    "fqbn",
		ArduinoAPIVersion:       "ideVersion",
		CustomBuildProperties:   []string{"custom=prop"},
		Verbose:                 true,
		DebugLevel:              5,
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	commands := []types.Command{
		&builder.CreateBuildOptionsMap{},
		&builder.StoreBuildOptionsMap{},
	}

	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}

	_, err := os.Stat(filepath.Join(buildPath, constants.BUILD_OPTIONS_FILE))
	NoError(t, err)

	bytes, err := ioutil.ReadFile(filepath.Join(buildPath, constants.BUILD_OPTIONS_FILE))
	NoError(t, err)

	require.Equal(t, "{\n"+
		"  \"additionalFiles\": \"\",\n"+
		"  \"builtInLibrariesFolders\": \"built-in libraries\",\n"+
		"  \"customBuildProperties\": \"custom=prop\",\n"+
		"  \"fqbn\": \"fqbn\",\n"+
		"  \"hardwareFolders\": \"hardware\",\n"+
		"  \"otherLibrariesFolders\": \"libraries\",\n"+
		"  \"runtime.ide.version\": \"ideVersion\",\n"+
		"  \"sketchLocation\": \"sketchLocation\",\n"+
		"  \"toolsFolders\": \"tools\"\n"+
		"}", string(bytes))
}
