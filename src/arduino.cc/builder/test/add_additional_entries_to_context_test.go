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
	"testing"
)

func TestAddAdditionalEntriesToContextNoBuildPath(t *testing.T) {
	context := make(map[string]interface{})

	command := builder.AddAdditionalEntriesToContext{}
	NoError(t, command.Run(context))

	require.Nil(t, context[constants.CTX_PREPROC_PATH])
	require.Nil(t, context[constants.CTX_SKETCH_BUILD_PATH])
	require.Nil(t, context[constants.CTX_LIBRARIES_BUILD_PATH])
	require.Nil(t, context[constants.CTX_CORE_BUILD_PATH])

	require.NotNil(t, context[constants.CTX_WARNINGS_LEVEL])
	require.NotNil(t, context[constants.CTX_VERBOSE])
	require.NotNil(t, context[constants.CTX_DEBUG_LEVEL])
	require.NotNil(t, context[constants.CTX_LIBRARY_DISCOVERY_RECURSION_DEPTH])

	require.True(t, context[constants.CTX_COLLECTED_SOURCE_FILES_QUEUE].(*types.UniqueStringQueue).Empty())
	require.True(t, context[constants.CTX_FOLDERS_WITH_SOURCES_QUEUE].(*types.UniqueSourceFolderQueue).Empty())

	require.Equal(t, 0, len(context[constants.CTX_LIBRARY_RESOLUTION_RESULTS].(map[string]types.LibraryResolutionResult)))
}

func TestAddAdditionalEntriesToContextWithBuildPath(t *testing.T) {
	context := make(map[string]interface{})

	context[constants.CTX_BUILD_PATH] = "folder"

	command := builder.AddAdditionalEntriesToContext{}
	NoError(t, command.Run(context))

	require.Equal(t, Abs(t, filepath.Join("folder", constants.FOLDER_PREPROC)), context[constants.CTX_PREPROC_PATH])
	require.Equal(t, Abs(t, filepath.Join("folder", constants.FOLDER_SKETCH)), context[constants.CTX_SKETCH_BUILD_PATH])
	require.Equal(t, Abs(t, filepath.Join("folder", constants.FOLDER_LIBRARIES)), context[constants.CTX_LIBRARIES_BUILD_PATH])
	require.Equal(t, Abs(t, filepath.Join("folder", constants.FOLDER_CORE)), context[constants.CTX_CORE_BUILD_PATH])

	require.NotNil(t, context[constants.CTX_WARNINGS_LEVEL])
	require.NotNil(t, context[constants.CTX_VERBOSE])
	require.NotNil(t, context[constants.CTX_DEBUG_LEVEL])
	require.NotNil(t, context[constants.CTX_LIBRARY_DISCOVERY_RECURSION_DEPTH])

	require.True(t, context[constants.CTX_COLLECTED_SOURCE_FILES_QUEUE].(*types.UniqueStringQueue).Empty())
	require.True(t, context[constants.CTX_FOLDERS_WITH_SOURCES_QUEUE].(*types.UniqueSourceFolderQueue).Empty())

	require.Equal(t, 0, len(context[constants.CTX_LIBRARY_RESOLUTION_RESULTS].(map[string]types.LibraryResolutionResult)))
}
