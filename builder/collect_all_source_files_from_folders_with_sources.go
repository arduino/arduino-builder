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
	"github.com/arduino/arduino-builder/builder/gohasissues"
	"github.com/arduino/arduino-builder/builder/types"
	"github.com/arduino/arduino-builder/builder/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type CollectAllSourceFilesFromFoldersWithSources struct{}

func (s *CollectAllSourceFilesFromFoldersWithSources) Run(context map[string]interface{}) error {
	foldersWithSources := context[constants.CTX_FOLDERS_WITH_SOURCES_QUEUE].(*types.UniqueSourceFolderQueue)
	sourceFiles := context[constants.CTX_COLLECTED_SOURCE_FILES_QUEUE].(*types.UniqueStringQueue)

	filePaths := []string{}
	for !foldersWithSources.Empty() {
		sourceFolder := foldersWithSources.Pop().(types.SourceFolder)
		var err error
		if sourceFolder.Recurse {
			err = collectByWalk(&filePaths, sourceFolder.Folder)
		} else {
			err = collectByReadDir(&filePaths, sourceFolder.Folder)
		}
		if err != nil {
			return utils.WrapError(err)
		}
	}

	for _, filePath := range filePaths {
		sourceFiles.Push(filePath)
	}

	return nil
}

func collectByWalk(filePaths *[]string, folder string) error {
	checkExtensionFunc := func(filePath string) bool {
		name := filepath.Base(filePath)
		ext := strings.ToLower(filepath.Ext(filePath))
		return !strings.HasPrefix(name, ".") && ADDITIONAL_FILE_VALID_EXTENSIONS_NO_HEADERS[ext]
	}
	walkFunc := utils.CollectAllReadableFiles(filePaths, checkExtensionFunc)
	err := gohasissues.Walk(folder, walkFunc)
	return utils.WrapError(err)
}

func collectByReadDir(filePaths *[]string, folder string) error {
	if _, err := os.Stat(folder); err != nil && os.IsNotExist(err) {
		return nil
	}

	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return utils.WrapError(err)
	}
	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file.Name()))
		if ADDITIONAL_FILE_VALID_EXTENSIONS_NO_HEADERS[ext] {
			*filePaths = append(*filePaths, filepath.Join(folder, file.Name()))
		}
	}
	return nil
}
