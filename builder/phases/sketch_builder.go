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
)

type SketchBuilder struct{}

func (s *SketchBuilder) Run(context map[string]interface{}) error {
	sketchBuildPath := context[constants.CTX_SKETCH_BUILD_PATH].(string)
	buildProperties := context[constants.CTX_BUILD_PROPERTIES].(props.PropertiesMap)
	includes := context[constants.CTX_INCLUDE_FOLDERS].([]string)
	includes = utils.Map(includes, utils.WrapWithHyphenI)
	verbose := context[constants.CTX_VERBOSE].(bool)
	warningsLevel := context[constants.CTX_WARNINGS_LEVEL].(string)
	logger := context[constants.CTX_LOGGER].(i18n.Logger)

	err := utils.EnsureFolderExists(sketchBuildPath)
	if err != nil {
		return utils.WrapError(err)
	}

	var objectFiles []string
	objectFiles, err = builder_utils.CompileFiles(objectFiles, sketchBuildPath, true, sketchBuildPath, buildProperties, includes, verbose, warningsLevel, logger)
	if err != nil {
		return utils.WrapError(err)
	}

	context[constants.CTX_OBJECT_FILES_SKETCH] = objectFiles

	return nil
}
