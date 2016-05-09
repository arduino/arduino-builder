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

package props

import (
	"arduino.cc/builder/constants"
	"arduino.cc/builder/i18n"
	"github.com/go-errors/errors"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strings"
)

type PropertiesMap map[string]string

var OSNAME string

func init() {
	switch value := runtime.GOOS; value {
	case "linux":
		OSNAME = runtime.GOOS
	case "windows":
		OSNAME = runtime.GOOS
	case "darwin":
		OSNAME = "macosx"
	default:
		panic("Unsupported OS")
	}
}

func Load(filepath string, logger i18n.Logger) (PropertiesMap, error) {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	text := string(bytes)
	text = strings.Replace(text, "\r\n", "\n", -1)
	text = strings.Replace(text, "\r", "\n", -1)

	properties := make(PropertiesMap)

	for _, line := range strings.Split(text, "\n") {
		err := properties.loadSingleLine(line)
		if err != nil {
			return nil, i18n.ErrorfWithLogger(logger, constants.MSG_WRONG_PROPERTIES_FILE, line, filepath)
		}
	}

	return properties, nil
}

func LoadFromSlice(lines []string, logger i18n.Logger) (PropertiesMap, error) {
	properties := make(PropertiesMap)

	for _, line := range lines {
		err := properties.loadSingleLine(line)
		if err != nil {
			return nil, i18n.ErrorfWithLogger(logger, constants.MSG_WRONG_PROPERTIES, line)
		}
	}

	return properties, nil
}

func (properties PropertiesMap) loadSingleLine(line string) error {
	line = strings.TrimSpace(line)

	if len(line) > 0 && line[0] != '#' {
		lineParts := strings.SplitN(line, "=", 2)
		if len(lineParts) != 2 {
			return errors.New("")
		}
		key := strings.TrimSpace(lineParts[0])
		value := strings.TrimSpace(lineParts[1])

		key = strings.Replace(key, "."+OSNAME, constants.EMPTY_STRING, 1)
		properties[key] = value
	}

	return nil
}

func SafeLoad(filepath string, logger i18n.Logger) (PropertiesMap, error) {
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return make(PropertiesMap), nil
	}

	properties, err := Load(filepath, logger)
	if err != nil {
		return nil, i18n.WrapError(err)
	}
	return properties, nil
}

func (aMap PropertiesMap) FirstLevelOf() map[string]PropertiesMap {
	newMap := make(map[string]PropertiesMap)
	for key, value := range aMap {
		if strings.Index(key, ".") == -1 {
			continue
		}
		keyParts := strings.SplitN(key, ".", 2)
		if newMap[keyParts[0]] == nil {
			newMap[keyParts[0]] = make(PropertiesMap)
		}
		newMap[keyParts[0]][keyParts[1]] = value
	}
	return newMap
}

func (aMap PropertiesMap) SubTree(key string) PropertiesMap {
	return aMap.FirstLevelOf()[key]
}

func (aMap PropertiesMap) ExpandPropsInString(str string) string {
	replaced := true
	for i := 0; i < 10 && replaced; i++ {
		replaced = false
		for key, value := range aMap {
			newStr := strings.Replace(str, "{"+key+"}", value, -1)
			replaced = replaced || str != newStr
			str = newStr
		}
	}
	return str
}

func (target PropertiesMap) Merge(sources ...PropertiesMap) PropertiesMap {
	for _, source := range sources {
		for key, value := range source {
			target[key] = value
		}
	}
	return target
}

func (aMap PropertiesMap) Clone() PropertiesMap {
	newMap := make(PropertiesMap)
	newMap.Merge(aMap)
	return newMap
}

func (aMap PropertiesMap) Equals(anotherMap PropertiesMap) bool {
	return reflect.DeepEqual(aMap, anotherMap)
}

func MergeMapsOfProperties(target map[string]PropertiesMap, sources ...map[string]PropertiesMap) map[string]PropertiesMap {
	for _, source := range sources {
		for key, value := range source {
			target[key] = value
		}
	}
	return target
}

func DeleteUnexpandedPropsFromString(str string) (string, error) {
	rxp, err := regexp.Compile("\\{.+?\\}")
	if err != nil {
		return constants.EMPTY_STRING, i18n.WrapError(err)
	}

	return rxp.ReplaceAllString(str, constants.EMPTY_STRING), nil
}
