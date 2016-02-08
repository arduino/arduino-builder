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
	"github.com/arduino/arduino-builder/builder/constants"
	"github.com/arduino/arduino-builder/builder/props"
	"github.com/arduino/arduino-builder/builder/types"
	"github.com/arduino/arduino-builder/builder/utils"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type SetupBuildProperties struct{}

func (s *SetupBuildProperties) Run(context map[string]interface{}) error {
	packages := context[constants.CTX_HARDWARE].(*types.Packages)

	targetPlatform := context[constants.CTX_TARGET_PLATFORM].(*types.Platform)
	actualPlatform := context[constants.CTX_ACTUAL_PLATFORM].(*types.Platform)
	targetBoard := context[constants.CTX_TARGET_BOARD].(*types.Board)

	buildProperties := make(props.PropertiesMap)
	buildProperties.Merge(actualPlatform.Properties)
	buildProperties.Merge(targetPlatform.Properties)
	buildProperties.Merge(targetBoard.Properties)

	if utils.MapHas(context, constants.CTX_BUILD_PATH) {
		buildProperties[constants.BUILD_PROPERTIES_BUILD_PATH] = context[constants.CTX_BUILD_PATH].(string)
	}
	if utils.MapHas(context, constants.CTX_SKETCH) {
		buildProperties[constants.BUILD_PROPERTIES_BUILD_PROJECT_NAME] = filepath.Base(context[constants.CTX_SKETCH].(*types.Sketch).MainFile.Name)
	}
	buildProperties[constants.BUILD_PROPERTIES_BUILD_ARCH] = strings.ToUpper(targetPlatform.PlatformId)

	buildProperties[constants.BUILD_PROPERTIES_BUILD_CORE] = context[constants.CTX_BUILD_CORE].(string)
	buildProperties[constants.BUILD_PROPERTIES_BUILD_CORE_PATH] = filepath.Join(actualPlatform.Folder, constants.FOLDER_CORES, buildProperties[constants.BUILD_PROPERTIES_BUILD_CORE])
	buildProperties[constants.BUILD_PROPERTIES_BUILD_SYSTEM_PATH] = filepath.Join(actualPlatform.Folder, constants.FOLDER_SYSTEM)
	buildProperties[constants.BUILD_PROPERTIES_RUNTIME_PLATFORM_PATH] = targetPlatform.Folder
	buildProperties[constants.BUILD_PROPERTIES_RUNTIME_HARDWARE_PATH] = filepath.Join(targetPlatform.Folder, "..")
	buildProperties[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION] = context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION].(string)
	buildProperties[constants.IDE_VERSION] = context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION].(string)
	buildProperties[constants.BUILD_PROPERTIES_RUNTIME_OS] = utils.PrettyOSName()

	variant := buildProperties[constants.BUILD_PROPERTIES_BUILD_VARIANT]
	if variant == constants.EMPTY_STRING {
		buildProperties[constants.BUILD_PROPERTIES_BUILD_VARIANT_PATH] = constants.EMPTY_STRING
	} else {
		var variantPlatform *types.Platform
		variantParts := strings.Split(variant, ":")
		if len(variantParts) > 1 {
			variantPlatform = packages.Packages[variantParts[0]].Platforms[targetPlatform.PlatformId]
			variant = variantParts[1]
		} else {
			variantPlatform = targetPlatform
		}
		buildProperties[constants.BUILD_PROPERTIES_BUILD_VARIANT_PATH] = filepath.Join(variantPlatform.Folder, constants.FOLDER_VARIANTS, variant)
	}

	tools := context[constants.CTX_TOOLS].([]*types.Tool)
	for _, tool := range tools {
		buildProperties[constants.BUILD_PROPERTIES_RUNTIME_TOOLS_PREFIX+tool.Name+constants.BUILD_PROPERTIES_RUNTIME_TOOLS_SUFFIX] = tool.Folder
		buildProperties[constants.BUILD_PROPERTIES_RUNTIME_TOOLS_PREFIX+tool.Name+"-"+tool.Version+constants.BUILD_PROPERTIES_RUNTIME_TOOLS_SUFFIX] = tool.Folder
	}

	if !utils.MapStringStringHas(buildProperties, constants.BUILD_PROPERTIES_SOFTWARE) {
		buildProperties[constants.BUILD_PROPERTIES_SOFTWARE] = DEFAULT_SOFTWARE
	}

	if utils.MapHas(context, constants.CTX_SKETCH_LOCATION) {
		sourcePath, err := filepath.Abs(context[constants.CTX_SKETCH_LOCATION].(string))
		if err != nil {
			return err
		}
		sourcePath = filepath.Dir(sourcePath)
		buildProperties[constants.BUILD_PROPERTIES_SOURCE_PATH] = sourcePath
	}

	now := time.Now()
	buildProperties[constants.BUILD_PROPERTIES_EXTRA_TIME_UTC] = strconv.FormatInt(now.Unix(), 10)
	buildProperties[constants.BUILD_PROPERTIES_EXTRA_TIME_LOCAL] = strconv.FormatInt(utils.LocalUnix(now), 10)
	buildProperties[constants.BUILD_PROPERTIES_EXTRA_TIME_ZONE] = strconv.Itoa(utils.TimezoneOffset())
	buildProperties[constants.BUILD_PROPERTIES_EXTRA_TIME_DST] = strconv.Itoa(utils.DaylightSavingsOffset(now))

	context[constants.CTX_BUILD_PROPERTIES] = buildProperties

	return nil
}
