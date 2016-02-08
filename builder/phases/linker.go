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

package phases

import (
	"github.com/arduino/arduino-builder/builder/builder_utils"
	"github.com/arduino/arduino-builder/builder/constants"
	"github.com/arduino/arduino-builder/builder/i18n"
	"github.com/arduino/arduino-builder/builder/props"
	"github.com/arduino/arduino-builder/builder/utils"
	"path/filepath"
	"strings"
)

type Linker struct{}

func (s *Linker) Run(context map[string]interface{}) error {
	objectFilesSketch := context[constants.CTX_OBJECT_FILES_SKETCH].([]string)
	objectFilesLibraries := context[constants.CTX_OBJECT_FILES_LIBRARIES].([]string)
	objectFilesCore := context[constants.CTX_OBJECT_FILES_CORE].([]string)

	var objectFiles []string
	objectFiles = append(objectFiles, objectFilesSketch...)
	objectFiles = append(objectFiles, objectFilesLibraries...)
	objectFiles = append(objectFiles, objectFilesCore...)

	coreArchiveFilePath := context[constants.CTX_ARCHIVE_FILE_PATH_CORE].(string)
	buildPath := context[constants.CTX_BUILD_PATH].(string)
	coreDotARelPath, err := filepath.Rel(buildPath, coreArchiveFilePath)
	if err != nil {
		return utils.WrapError(err)
	}

	buildProperties := context[constants.CTX_BUILD_PROPERTIES].(props.PropertiesMap)
	verbose := context[constants.CTX_VERBOSE].(bool)
	warningsLevel := context[constants.CTX_WARNINGS_LEVEL].(string)
	logger := context[constants.CTX_LOGGER].(i18n.Logger)

	err = link(objectFiles, coreDotARelPath, coreArchiveFilePath, buildProperties, verbose, warningsLevel, logger)
	if err != nil {
		return utils.WrapError(err)
	}

	return nil
}

func link(objectFiles []string, coreDotARelPath string, coreArchiveFilePath string, buildProperties props.PropertiesMap, verbose bool, warningsLevel string, logger i18n.Logger) error {
	optRelax := addRelaxTrickIfATMEGA2560(buildProperties)

	objectFiles = utils.Map(objectFiles, wrapWithDoubleQuotes)
	objectFileList := strings.Join(objectFiles, constants.SPACE)

	properties := buildProperties.Clone()
	properties[constants.BUILD_PROPERTIES_COMPILER_C_ELF_FLAGS] = properties[constants.BUILD_PROPERTIES_COMPILER_C_ELF_FLAGS] + optRelax
	properties[constants.BUILD_PROPERTIES_COMPILER_WARNING_FLAGS] = properties[constants.BUILD_PROPERTIES_COMPILER_WARNING_FLAGS+"."+warningsLevel]
	properties[constants.BUILD_PROPERTIES_ARCHIVE_FILE] = coreDotARelPath
	properties[constants.BUILD_PROPERTIES_ARCHIVE_FILE_PATH] = coreArchiveFilePath
	properties[constants.BUILD_PROPERTIES_OBJECT_FILES] = objectFileList

	_, err := builder_utils.ExecRecipe(properties, constants.RECIPE_C_COMBINE_PATTERN, false, verbose, verbose, logger)
	return err
}

func wrapWithDoubleQuotes(value string) string {
	return "\"" + value + "\""
}

func addRelaxTrickIfATMEGA2560(buildProperties props.PropertiesMap) string {
	if buildProperties[constants.BUILD_PROPERTIES_BUILD_MCU] == "atmega2560" {
		return ",--relax"
	}
	return constants.EMPTY_STRING
}
