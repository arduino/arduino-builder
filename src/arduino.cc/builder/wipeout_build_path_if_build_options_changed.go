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
	"arduino.cc/builder/types"
	"encoding/json"
	"os"
	"path/filepath"
)

type WipeoutBuildPathIfBuildOptionsChanged struct{}

func (s *WipeoutBuildPathIfBuildOptionsChanged) Run(ctx *types.Context) error {
	if ctx.BuildOptionsJsonPrevious == "" {
		return nil
	}
	buildOptionsJson := ctx.BuildOptionsJson
	previousBuildOptionsJson := ctx.BuildOptionsJsonPrevious
	logger := ctx.GetLogger()

	if buildOptionsJson == previousBuildOptionsJson {
		return nil
	}

	// unmarshal jsons and check if every field is equal except SketchLocation
	// if SketchLocation path is different but filename is the same, don't wipe
	var buildOptions map[string]string
	var previousBuildOptions map[string]string
	json.Unmarshal([]byte(buildOptionsJson), &buildOptions)
	json.Unmarshal([]byte(previousBuildOptionsJson), &previousBuildOptions)

	for key, _ := range buildOptions {
		if buildOptions[key] != previousBuildOptions[key] {
			if key == "sketchLocation" {
				if filepath.Base(buildOptions[key]) == filepath.Base(previousBuildOptions[key]) {
					return nil
				}
			} else {
				break
			}
		}
	}

	logger.Println(constants.LOG_LEVEL_INFO, constants.MSG_BUILD_OPTIONS_CHANGED)

	buildPath := ctx.BuildPath
	files, err := gohasissues.ReadDir(buildPath)
	if err != nil {
		return i18n.WrapError(err)
	}
	for _, file := range files {
		os.RemoveAll(filepath.Join(buildPath, file.Name()))
	}

	return nil
}
