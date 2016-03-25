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
	"arduino.cc/builder/gohasissues"
	"arduino.cc/builder/i18n"
	"arduino.cc/builder/utils"
	"os"
	"path/filepath"
	"regexp"
)

type WipeoutBuildPathIfBuildOptionsChanged struct{}

func (s *WipeoutBuildPathIfBuildOptionsChanged) Run(context map[string]interface{}) error {
	if !utils.MapHas(context, constants.CTX_BUILD_OPTIONS_PREVIOUS_JSON) {
		return nil
	}
	buildOptionsJson := context[constants.CTX_BUILD_OPTIONS_JSON].(string)
	previousBuildOptionsJson := context[constants.CTX_BUILD_OPTIONS_PREVIOUS_JSON].(string)
	logger := context[constants.CTX_LOGGER].(i18n.Logger)

	if buildOptionsJson == previousBuildOptionsJson {
		return nil
	}

	re := regexp.MustCompile("(?m)^.*" + constants.CTX_SKETCH_LOCATION + ".*$[\r\n]+")
	buildOptionsJson = re.ReplaceAllString(buildOptionsJson, "")
	previousBuildOptionsJson = re.ReplaceAllString(previousBuildOptionsJson, "")

	// if the only difference is the sketch path skip deleting everything
	if buildOptionsJson == previousBuildOptionsJson {
		return nil
	}

	logger.Println(constants.LOG_LEVEL_INFO, constants.MSG_BUILD_OPTIONS_CHANGED)

	buildPath := context[constants.CTX_BUILD_PATH].(string)
	files, err := gohasissues.ReadDir(buildPath)
	if err != nil {
		return utils.WrapError(err)
	}
	for _, file := range files {
		os.RemoveAll(filepath.Join(buildPath, file.Name()))
	}

	return nil
}
