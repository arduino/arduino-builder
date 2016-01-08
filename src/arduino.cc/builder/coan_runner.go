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
	"arduino.cc/builder/i18n"
	"arduino.cc/builder/props"
	"arduino.cc/builder/utils"
	"fmt"
	"path/filepath"
	"regexp"
)

var ALLOWED_ARG = regexp.MustCompile("^source$|\\-E$|\\-m$|\\-P$|\\-kb$|\\-D.*")

type CoanRunner struct{}

func (s *CoanRunner) Run(context map[string]interface{}) error {
	source := context[constants.CTX_SOURCE].(string)
	source += "\n"
	verbose := context[constants.CTX_VERBOSE].(bool)

	preprocPath := context[constants.CTX_PREPROC_PATH].(string)
	err := utils.EnsureFolderExists(preprocPath)
	if err != nil {
		return utils.WrapError(err)
	}

	coanTargetFileName := filepath.Join(preprocPath, constants.FILE_COAN_TARGET)
	err = utils.WriteFile(coanTargetFileName, source)
	if err != nil {
		return utils.WrapError(err)
	}

	buildProperties := context[constants.CTX_BUILD_PROPERTIES].(props.PropertiesMap)
	properties := buildProperties.Clone()
	properties.Merge(buildProperties.SubTree(constants.BUILD_PROPERTIES_TOOLS_KEY).SubTree(constants.COAN))
	properties[constants.BUILD_PROPERTIES_SOURCE_FILE] = coanTargetFileName

	pattern := properties[constants.BUILD_PROPERTIES_PATTERN]
	if pattern == constants.EMPTY_STRING {
		return utils.Errorf(context, constants.MSG_PATTERN_MISSING, constants.COAN)
	}

	logger := context[constants.CTX_LOGGER].(i18n.Logger)
	commandLine := properties.ExpandPropsInString(pattern)
	command, err := utils.PrepareCommandFilteredArgs(commandLine, filterAllowedArg, logger)

	if verbose {
		fmt.Println(commandLine)
	}

	sourceBytes, _ := command.Output()

	context[constants.CTX_SOURCE] = string(sourceBytes)

	return nil
}

func filterAllowedArg(idx int, part string, parts []string) bool {
	return idx == len(parts)-1 || ALLOWED_ARG.MatchString(part)
}
