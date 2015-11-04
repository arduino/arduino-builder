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
 * Copyright 2015 Steve Marple
 */

package test

import (
	"arduino.cc/builder"
	"arduino.cc/builder/constants"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func TestSketchWithNoBuildProps(t *testing.T) {
	var err error
	context := make(map[string]interface{})

	context[constants.CTX_BUILD_PATH] = "buildPath"
	context[constants.CTX_HARDWARE_FOLDERS] = []string{"hardware"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"tools"}
	context[constants.CTX_LIBRARIES_FOLDERS] = []string{"libraries"}
	context[constants.CTX_FQBN] = "fqbn"
	context[constants.CTX_SKETCH_LOCATION] = filepath.Join("sketch_no_props", "sketch.ino")

	context[constants.CTX_VERBOSE] = true
	context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION] = "ideVersion"
	context[constants.CTX_SKETCH_BUILD_PROPERTIES], err = builder.GetSketchBuildProperties(context)
	NoError(t, err)

	context[constants.CTX_DEBUG_LEVEL] = 5

	create := builder.CreateBuildOptionsMap{}
	err = create.Run(context)
	NoError(t, err)

	buildOptions := context[constants.CTX_BUILD_OPTIONS].(map[string]interface{})
	sketchBuildProps := make(map[string]string)
	require.Equal(t, sketchBuildProps, buildOptions[constants.CTX_SKETCH_BUILD_PROPERTIES])
}

func TestSketchWithBuildProps(t *testing.T) {
	var err error
	context := make(map[string]interface{})

	context[constants.CTX_BUILD_PATH] = "buildPath"
	context[constants.CTX_HARDWARE_FOLDERS] = []string{"hardware"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"tools"}
	context[constants.CTX_LIBRARIES_FOLDERS] = []string{"libraries"}
	context[constants.CTX_FQBN] = "fqbn"
	context[constants.CTX_SKETCH_LOCATION] = filepath.Join("sketch_with_props", "sketch.ino")

	context[constants.CTX_VERBOSE] = true
	context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION] = "ideVersion"
	context[constants.CTX_SKETCH_BUILD_PROPERTIES], err = builder.GetSketchBuildProperties(context)
	NoError(t, err)

	context[constants.CTX_DEBUG_LEVEL] = 5

	create := builder.CreateBuildOptionsMap{}
	err = create.Run(context)
	NoError(t, err)

	buildOptions := context[constants.CTX_BUILD_OPTIONS].(map[string]interface{})
	sketchBuildProps := make(map[string]string)
	sketchBuildProps["compiler.c.extra_flags"] = "-D NDEBUG"
	sketchBuildProps["compiler.cpp.extra_flags"] = "-D NDEBUG -D TESTLIBRARY_BUFSIZE=100 -D ARDUINO_URL=\"https://www.arduino.cc/\""

	require.Equal(t, sketchBuildProps, buildOptions[constants.CTX_SKETCH_BUILD_PROPERTIES])

}
