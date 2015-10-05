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
 * Copyright 2015 Steve Marple
 */

package builder

import (
	"arduino.cc/builder/constants"
	"arduino.cc/builder/props"
	"github.com/go-errors/errors"
	"os"
        "path/filepath"
)

func GetSketchBuildPropertiesFilename(context map[string]interface{}) (string, error) {
	sketchLocation, ok := context[constants.CTX_SKETCH_LOCATION].(string)
	if !ok {
		return "", errors.New("Unknown sketch location")
	}

	filename := filepath.Join(filepath.Dir(sketchLocation), constants.SKETCH_BUILD_OPTIONS_TXT)
	return filename, nil
}


func GetSketchBuildProperties(context map[string]interface{}) (map[string]string, error) {
	filename, err := GetSketchBuildPropertiesFilename(context)
	if err != nil {
		return nil, err
	}
	_, err = os.Stat(filename)
	if err == nil {
		return props.SafeLoad(filename)
	}
	return make(map[string]string), nil
}

