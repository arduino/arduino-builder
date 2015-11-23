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

package builder

import (
	"arduino.cc/builder/constants"
	"arduino.cc/builder/types"
	"arduino.cc/builder/utils"
	"path/filepath"
)

type AddAdditionalEntriesToContext struct{}

func (s *AddAdditionalEntriesToContext) Run(context map[string]interface{}) error {
	if utils.MapHas(context, constants.CTX_BUILD_PATH) {
		buildPath := context[constants.CTX_BUILD_PATH].(string)
		preprocPath, err := filepath.Abs(filepath.Join(buildPath, constants.FOLDER_PREPROC))
		if err != nil {
			return utils.WrapError(err)
		}
		sketchBuildPath, err := filepath.Abs(filepath.Join(buildPath, constants.FOLDER_SKETCH))
		if err != nil {
			return utils.WrapError(err)
		}
		librariesBuildPath, err := filepath.Abs(filepath.Join(buildPath, constants.FOLDER_LIBRARIES))
		if err != nil {
			return utils.WrapError(err)
		}
		coreBuildPath, err := filepath.Abs(filepath.Join(buildPath, constants.FOLDER_CORE))
		if err != nil {
			return utils.WrapError(err)
		}

		context[constants.CTX_PREPROC_PATH] = preprocPath
		context[constants.CTX_SKETCH_BUILD_PATH] = sketchBuildPath
		context[constants.CTX_LIBRARIES_BUILD_PATH] = librariesBuildPath
		context[constants.CTX_CORE_BUILD_PATH] = coreBuildPath
	}

	if !utils.MapHas(context, constants.CTX_WARNINGS_LEVEL) {
		context[constants.CTX_WARNINGS_LEVEL] = DEFAULT_WARNINGS_LEVEL
	}

	if !utils.MapHas(context, constants.CTX_VERBOSE) {
		context[constants.CTX_VERBOSE] = false
	}

	if !utils.MapHas(context, constants.CTX_DEBUG_LEVEL) {
		context[constants.CTX_DEBUG_LEVEL] = DEFAULT_DEBUG_LEVEL
	}

	sourceFiles := &types.UniqueStringQueue{}
	context[constants.CTX_COLLECTED_SOURCE_FILES_QUEUE] = sourceFiles
	foldersWithSources := &types.UniqueSourceFolderQueue{}
	context[constants.CTX_FOLDERS_WITH_SOURCES_QUEUE] = foldersWithSources

	context[constants.CTX_LIBRARY_RESOLUTION_RESULTS] = make(map[string]types.LibraryResolutionResult)
	context[constants.CTX_HARDWARE_REWRITE_RESULTS] = make(map[*types.Platform][]types.PlatforKeyRewrite)

	return nil
}
