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
	"arduino.cc/builder/utils"
	"encoding/json"
	"reflect"
	"strings"
)

type CreateBuildOptionsMap struct{}

func (s *CreateBuildOptionsMap) Run(context map[string]interface{}) error {
	buildOptions := make(map[string]string)

	buildOptionsMapKeys := []string{
		constants.CTX_HARDWARE_FOLDERS,
		constants.CTX_TOOLS_FOLDERS,
		constants.CTX_LIBRARIES_FOLDERS,
		constants.CTX_FQBN,
		constants.CTX_SKETCH_LOCATION,
		constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION,
		constants.CTX_CUSTOM_BUILD_PROPERTIES,
	}

	for _, key := range buildOptionsMapKeys {
		if utils.MapHas(context, key) {
			originalValue := context[key]
			value := constants.EMPTY_STRING
			kindOfValue := reflect.TypeOf(originalValue).Kind()
			if kindOfValue == reflect.Slice {
				value = strings.Join(originalValue.([]string), ",")
			} else if kindOfValue == reflect.String {
				value = originalValue.(string)
			} else {
				return utils.Errorf(context, constants.MSG_UNHANDLED_TYPE_IN_CONTEXT, kindOfValue.String(), key)
			}

			buildOptions[key] = value
		}
	}

	context[constants.CTX_BUILD_OPTIONS] = buildOptions

	bytes, err := json.MarshalIndent(buildOptions, "", "  ")
	if err != nil {
		return utils.WrapError(err)
	}

	context[constants.CTX_BUILD_OPTIONS_JSON] = string(bytes)

	return nil
}
