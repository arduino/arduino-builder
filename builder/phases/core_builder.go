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

type CoreBuilder struct{}

func (s *CoreBuilder) Run(context map[string]interface{}) error {
	coreBuildPath := context[constants.CTX_CORE_BUILD_PATH].(string)
	buildProperties := context[constants.CTX_BUILD_PROPERTIES].(props.PropertiesMap)
	verbose := context[constants.CTX_VERBOSE].(bool)
	warningsLevel := context[constants.CTX_WARNINGS_LEVEL].(string)
	logger := context[constants.CTX_LOGGER].(i18n.Logger)

	err := utils.EnsureFolderExists(coreBuildPath)
	if err != nil {
		return utils.WrapError(err)
	}

	archiveFile, objectFiles, err := compileCore(coreBuildPath, buildProperties, verbose, warningsLevel, logger)
	if err != nil {
		return utils.WrapError(err)
	}

	context[constants.CTX_ARCHIVE_FILE_PATH_CORE] = archiveFile
	context[constants.CTX_OBJECT_FILES_CORE] = objectFiles

	return nil
}

func compileCore(buildPath string, buildProperties props.PropertiesMap, verbose bool, warningsLevel string, logger i18n.Logger) (string, []string, error) {
	coreFolder := buildProperties[constants.BUILD_PROPERTIES_BUILD_CORE_PATH]
	variantFolder := buildProperties[constants.BUILD_PROPERTIES_BUILD_VARIANT_PATH]

	includes := []string{}
	includes = append(includes, coreFolder)
	if variantFolder != constants.EMPTY_STRING {
		includes = append(includes, variantFolder)
	}
	includes = utils.Map(includes, utils.WrapWithHyphenI)

	var err error

	variantObjectFiles := []string{}
	if variantFolder != constants.EMPTY_STRING {
		variantObjectFiles, err = builder_utils.CompileFiles(variantObjectFiles, variantFolder, true, buildPath, buildProperties, includes, verbose, warningsLevel, logger)
		if err != nil {
			return "", nil, utils.WrapError(err)
		}
	}

	coreObjectFiles, err := builder_utils.CompileFiles([]string{}, coreFolder, true, buildPath, buildProperties, includes, verbose, warningsLevel, logger)
	if err != nil {
		return "", nil, utils.WrapError(err)
	}

	archiveFile, err := builder_utils.ArchiveCompiledFiles(buildPath, "core.a", coreObjectFiles, buildProperties, verbose, logger)
	if err != nil {
		return "", nil, utils.WrapError(err)
	}

	return archiveFile, variantObjectFiles, nil
}
