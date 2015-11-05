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
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadSketchWithFolder(t *testing.T) {
	context := make(map[string]interface{})
	context[constants.CTX_SKETCH_LOCATION] = "sketch1"

	loggerCommand := builder.SetupHumanLoggerIfMissing{}
	err := loggerCommand.Run(context)
	NoError(t, err)

	loader := builder.SketchLoader{}
	err = loader.Run(context)

	require.Error(t, err)

	sketch := context[constants.CTX_SKETCH]
	require.Nil(t, sketch)
}

func TestLoadSketchNonExistentPath(t *testing.T) {
	context := make(map[string]interface{})
	context[constants.CTX_SKETCH_LOCATION] = "asdasd78128123981723981273asdasd"

	loggerCommand := builder.SetupHumanLoggerIfMissing{}
	err := loggerCommand.Run(context)
	NoError(t, err)

	loader := builder.SketchLoader{}
	err = loader.Run(context)

	require.Error(t, err)

	sketch := context[constants.CTX_SKETCH]
	require.Nil(t, sketch)
}

func TestLoadSketch(t *testing.T) {
	context := make(map[string]interface{})
	context[constants.CTX_SKETCH_LOCATION] = filepath.Join("sketch1", "sketch.ino")

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},
		&builder.SketchLoader{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	sketch := context[constants.CTX_SKETCH].(*types.Sketch)
	require.NotNil(t, sketch)

	require.True(t, strings.Index(sketch.MainFile.Name, "sketch.ino") != -1)

	require.Equal(t, 2, len(sketch.OtherSketchFiles))
	require.True(t, strings.Index(sketch.OtherSketchFiles[0].Name, "old.pde") != -1)
	require.True(t, strings.Index(sketch.OtherSketchFiles[1].Name, "other.ino") != -1)

	require.Equal(t, 2, len(sketch.AdditionalFiles))
	require.True(t, strings.Index(sketch.AdditionalFiles[0].Name, "header.h") != -1)
	require.True(t, strings.Index(sketch.AdditionalFiles[1].Name, "helper.h") != -1)
}

func TestFailToLoadSketchFromFolder(t *testing.T) {
	context := make(map[string]interface{})
	context[constants.CTX_SKETCH_LOCATION] = "./sketch1"

	loggerCommand := builder.SetupHumanLoggerIfMissing{}
	err := loggerCommand.Run(context)
	NoError(t, err)

	loader := builder.SketchLoader{}
	err = loader.Run(context)
	require.Error(t, err)

	sketch := context[constants.CTX_SKETCH]
	require.Nil(t, sketch)
}

func TestLoadSketchFromFolder(t *testing.T) {
	context := make(map[string]interface{})
	context[constants.CTX_SKETCH_LOCATION] = "sketch_with_subfolders"

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},
		&builder.SketchLoader{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	sketch := context[constants.CTX_SKETCH].(*types.Sketch)
	require.NotNil(t, sketch)

	require.True(t, strings.Index(sketch.MainFile.Name, "sketch_with_subfolders.ino") != -1)

	require.Equal(t, 2, len(sketch.AdditionalFiles))
	require.True(t, strings.Index(sketch.AdditionalFiles[0].Name, "other.cpp") != -1)
	require.True(t, strings.Index(sketch.AdditionalFiles[1].Name, "other.h") != -1)
}

func TestLoadSketchWithBackup(t *testing.T) {
	context := make(map[string]interface{})
	context[constants.CTX_SKETCH_LOCATION] = filepath.Join("sketch_with_backup_files", "sketch.ino")

	commands := []types.Command{
		&builder.SetupHumanLoggerIfMissing{},
		&builder.SketchLoader{},
	}

	for _, command := range commands {
		err := command.Run(context)
		NoError(t, err)
	}

	sketch := context[constants.CTX_SKETCH].(*types.Sketch)
	require.NotNil(t, sketch)

	require.True(t, strings.Index(sketch.MainFile.Name, "sketch.ino") != -1)

	require.Equal(t, 0, len(sketch.AdditionalFiles))
	require.Equal(t, 0, len(sketch.OtherSketchFiles))
}
