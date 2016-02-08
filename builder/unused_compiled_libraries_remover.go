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
	"github.com/arduino/arduino-builder/builder/constants"
	"github.com/arduino/arduino-builder/builder/types"
	"github.com/arduino/arduino-builder/builder/utils"
	"io/ioutil"
	"os"
	"path/filepath"
)

type UnusedCompiledLibrariesRemover struct{}

func (s *UnusedCompiledLibrariesRemover) Run(context map[string]interface{}) error {
	librariesBuildPath := context[constants.CTX_LIBRARIES_BUILD_PATH].(string)
	libraries := context[constants.CTX_IMPORTED_LIBRARIES].([]*types.Library)

	_, err := os.Stat(librariesBuildPath)
	if err != nil && os.IsNotExist(err) {
		return nil
	}

	libraryNames := toLibraryNames(libraries)

	files, err := ioutil.ReadDir(librariesBuildPath)
	if err != nil {
		return utils.WrapError(err)
	}
	for _, file := range files {
		if file.IsDir() {
			if !utils.SliceContains(libraryNames, file.Name()) {
				err := os.RemoveAll(filepath.Join(librariesBuildPath, file.Name()))
				if err != nil {
					return utils.WrapError(err)
				}
			}
		}
	}

	return nil
}

func toLibraryNames(libraries []*types.Library) []string {
	libraryNames := []string{}
	for _, library := range libraries {
		libraryNames = append(libraryNames, library.Name)
	}
	return libraryNames
}
