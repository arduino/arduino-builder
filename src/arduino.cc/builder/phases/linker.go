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
	"arduino.cc/builder/builder_utils"
	"arduino.cc/builder/constants"
	"arduino.cc/builder/i18n"
	"arduino.cc/builder/utils"
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

	buildProperties := context[constants.CTX_BUILD_PROPERTIES].(map[string]string)
	verbose := context[constants.CTX_VERBOSE].(bool)
	warningsLevel := context[constants.CTX_WARNINGS_LEVEL].(string)
	logger := context[constants.CTX_LOGGER].(i18n.Logger)

	err := link(objectFiles, buildProperties, verbose, warningsLevel, logger)
	if err != nil {
		return utils.WrapError(err)
	}

	return nil
}

func link(objectFiles []string, buildProperties map[string]string, verbose bool, warningsLevel string, logger i18n.Logger) error {
	optRelax := constants.EMPTY_STRING
	if buildProperties[constants.BUILD_PROPERTIES_BUILD_MCU] == "atmega2560" {
		optRelax = ",--relax"
	}

	objectFiles = utils.Map(objectFiles, wrapWithDoubleQuotes)
	objectFileList := strings.Join(objectFiles, constants.SPACE)

	properties := utils.MergeMapsOfStrings(make(map[string]string), buildProperties)
	properties["compiler.c.elf.flags"] = properties["compiler.c.elf.flags"] + optRelax
	properties["compiler.warning_flags"] = properties["compiler.warning_flags."+warningsLevel]
	properties[constants.BUILD_PROPERTIES_ARCHIVE_FILE] = "core.a"
	properties[constants.BUILD_PROPERTIES_OBJECT_FILES] = objectFileList

	_, err := builder_utils.ExecRecipe(properties, "recipe.c.combine.pattern", false, verbose, verbose, logger)
	return err
}

func wrapWithDoubleQuotes(value string) string {
	return "\"" + value + "\""
}
