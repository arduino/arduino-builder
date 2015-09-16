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
	"arduino.cc/builder/utils"
	"crypto/md5"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
)

type GenerateBuildPathIfMissing struct{}

func (s *GenerateBuildPathIfMissing) Run(context map[string]interface{}) error {
	if utils.MapHas(context, constants.CTX_BUILD_PATH) && context[constants.CTX_BUILD_PATH].(string) != constants.EMPTY_STRING {
		return nil
	}

	sketchLocation := context[constants.CTX_SKETCH_LOCATION].(string)
	md5sumBytes := md5.Sum([]byte(sketchLocation))
	md5sum := hex.EncodeToString(md5sumBytes[:])

	buildPath := filepath.Join(os.TempDir(), "arduino-sketch-"+strings.ToUpper(md5sum))
	_, err := os.Stat(buildPath)
	if err != nil && !os.IsNotExist(err) {
		return utils.WrapError(err)
	}

	if utils.DebugLevel(context) > 5 {
		logger := context[constants.CTX_LOGGER].(i18n.Logger)
		logger.Fprintln(os.Stderr, constants.MSG_SETTING_BUILD_PATH, buildPath)
	}

	context[constants.CTX_BUILD_PATH] = buildPath

	return nil
}
