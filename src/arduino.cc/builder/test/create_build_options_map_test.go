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
	"arduino.cc/builder/props"
	"arduino.cc/builder/types"
	"arduino.cc/builder/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateBuildOptionsMap(t *testing.T) {
	context := make(map[string]interface{})
	ctx := &types.Context{
		HardwareFolders:       []string{"hardware", "hardware2"},
		ToolsFolders:          []string{"tools"},
		OtherLibrariesFolders: []string{"libraries"},
		SketchLocation:        "sketchLocation",
		FQBN:                  "fqbn",
		ArduinoAPIVersion:     "ideVersion",
	}

	context[constants.CTX_BUILD_PATH] = "buildPath"
	context[constants.CTX_VERBOSE] = true
	ctx.DebugLevel = 5

	create := builder.CreateBuildOptionsMap{}
	err := create.Run(context, ctx)
	NoError(t, err)

	buildOptions := context[constants.CTX_BUILD_OPTIONS].(props.PropertiesMap)
	require.Equal(t, 8, len(utils.KeysOfMapOfString(buildOptions)))
	require.Equal(t, "hardware,hardware2", buildOptions["hardwareFolders"])
	require.Equal(t, "tools", buildOptions["toolsFolders"])
	require.Equal(t, "", buildOptions["builtInLibrariesFolders"])
	require.Equal(t, "", buildOptions["customBuildProperties"])
	require.Equal(t, "libraries", buildOptions["otherLibrariesFolders"])
	require.Equal(t, "fqbn", buildOptions["fqbn"])
	require.Equal(t, "sketchLocation", buildOptions["sketchLocation"])
	require.Equal(t, "ideVersion", buildOptions["runtime.ide.version"])

	require.Equal(t, "{\n"+
		"  \"builtInLibrariesFolders\": \"\",\n"+
		"  \"customBuildProperties\": \"\",\n"+
		"  \"fqbn\": \"fqbn\",\n"+
		"  \"hardwareFolders\": \"hardware,hardware2\",\n"+
		"  \"otherLibrariesFolders\": \"libraries\",\n"+
		"  \"runtime.ide.version\": \"ideVersion\",\n"+
		"  \"sketchLocation\": \"sketchLocation\",\n"+
		"  \"toolsFolders\": \"tools\"\n"+
		"}", context[constants.CTX_BUILD_OPTIONS_JSON].(string))
}
