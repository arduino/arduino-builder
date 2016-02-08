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
	"github.com/arduino/arduino-builder/builder"
	"github.com/arduino/arduino-builder/builder/constants"
	"github.com/arduino/arduino-builder/builder/types"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCoanRunner(t *testing.T) {
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
		&builder.SketchSourceMerger{},
		&builder.CoanRunner{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	expectedSource := LoadAndInterpolate(t, filepath.Join("sketch2", "SketchWithIfDef.resolved.directives.txt"), context)
	require.Equal(t, expectedSource, strings.Replace(context[constants.CTX_SOURCE].(string), "\r\n", "\n", -1))

}
