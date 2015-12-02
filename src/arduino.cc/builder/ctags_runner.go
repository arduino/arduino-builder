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
	"os"
)

type CTagsRunner struct{}

func (s *CTagsRunner) Run(context map[string]interface{}) error {
	buildProperties := context[constants.CTX_BUILD_PROPERTIES].(map[string]string)
	ctagsTargetFilePath := context[constants.CTX_CTAGS_TEMP_FILE_PATH].(string)
	logger := context[constants.CTX_LOGGER].(i18n.Logger)

	properties := utils.MergeMapsOfStrings(make(map[string]string), buildProperties, props.SubTree(props.SubTree(buildProperties, constants.BUILD_PROPERTIES_TOOLS_KEY), constants.CTAGS))
	properties[constants.BUILD_PROPERTIES_SOURCE_FILE] = ctagsTargetFilePath

	pattern := properties[constants.BUILD_PROPERTIES_PATTERN]
	if pattern == constants.EMPTY_STRING {
		return utils.Errorf(context, constants.MSG_PATTERN_MISSING, constants.CTAGS)
	}

	commandLine := props.ExpandPropsInString(properties, pattern)
	command, err := utils.PrepareCommand(commandLine, logger)
	if err != nil {
		return utils.WrapError(err)
	}

	verbose := context[constants.CTX_VERBOSE].(bool)
	if verbose {
		logger.UnformattedFprintln(os.Stdout, commandLine)
	}

	sourceBytes, err := command.Output()
	if err != nil {
		return utils.WrapError(err)
	}

	context[constants.CTX_CTAGS_OUTPUT] = string(sourceBytes)

	return nil
}
