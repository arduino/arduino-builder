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
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"arduino.cc/builder/constants"
	"arduino.cc/builder/i18n"
	"arduino.cc/builder/types"
	"arduino.cc/builder/utils"
)

type SketchLoader struct{}

func (s *SketchLoader) Run(ctx *types.Context) error {
	if ctx.SketchLocation == "" {
		return nil
	}

	sketchLocation := ctx.SketchLocation

	sketchLocation, err := filepath.Abs(sketchLocation)
	if err != nil {
		return i18n.WrapError(err)
	}
	mainSketchStat, err := os.Stat(sketchLocation)
	if err != nil {
		return i18n.WrapError(err)
	}
	if mainSketchStat.IsDir() {
		sketchLocation = filepath.Join(sketchLocation, mainSketchStat.Name()+".ino")
	}

	ctx.SketchLocation = sketchLocation

	allSketchFilePaths, err := collectAllSketchFiles(filepath.Dir(sketchLocation))
	if err != nil {
		return i18n.WrapError(err)
	}

	logger := ctx.GetLogger()

	if !utils.SliceContains(allSketchFilePaths, sketchLocation) {
		return i18n.ErrorfWithLogger(logger, constants.MSG_CANT_FIND_SKETCH_IN_PATH, sketchLocation, filepath.Dir(sketchLocation))
	}

	sketch, err := makeSketch(sketchLocation, allSketchFilePaths, logger)
	if err != nil {
		return i18n.WrapError(err)
	}

	ctx.SketchLocation = sketchLocation
	ctx.Sketch = sketch

	return nil
}

func collectAllSketchFiles(from string) ([]string, error) {
	filePaths := []string{}
	// Source files in the root are compiled, non-recursively. This
	// is the only place where .ino files can be present.
	rootExtensions := func(ext string) bool { return MAIN_FILE_VALID_EXTENSIONS[ext] || ADDITIONAL_FILE_VALID_EXTENSIONS[ext] }
	err := utils.FindFilesInFolder(&filePaths, from, rootExtensions, true /* recurse */)
	if err != nil {
		return nil, i18n.WrapError(err)
	}

	return filePaths, i18n.WrapError(err)
}

func makeSketch(sketchLocation string, allSketchFilePaths []string, logger i18n.Logger) (*types.Sketch, error) {
	sketchFilesMap := make(map[string]types.SketchFile)
	for _, sketchFilePath := range allSketchFilePaths {
		source, err := ioutil.ReadFile(sketchFilePath)
		if err != nil {
			return nil, i18n.WrapError(err)
		}
		sketchFilesMap[sketchFilePath] = types.SketchFile{Name: sketchFilePath, Source: string(source)}
	}

	mainFile := sketchFilesMap[sketchLocation]
	delete(sketchFilesMap, sketchLocation)

	additionalFiles := []types.SketchFile{}
	otherSketchFiles := []types.SketchFile{}
	mainFileDir := filepath.Dir(mainFile.Name)
	for _, sketchFile := range sketchFilesMap {
		ext := strings.ToLower(filepath.Ext(sketchFile.Name))
		if MAIN_FILE_VALID_EXTENSIONS[ext] {
			if filepath.Dir(sketchFile.Name) == mainFileDir {
				otherSketchFiles = append(otherSketchFiles, sketchFile)
			}
		} else if ADDITIONAL_FILE_VALID_EXTENSIONS[ext] {
			additionalFiles = append(additionalFiles, sketchFile)
		} else {
			return nil, i18n.ErrorfWithLogger(logger, constants.MSG_UNKNOWN_SKETCH_EXT, sketchFile.Name)
		}
	}

	sort.Sort(types.SketchFileSortByName(additionalFiles))
	sort.Sort(types.SketchFileSortByName(otherSketchFiles))

	return &types.Sketch{MainFile: mainFile, OtherSketchFiles: otherSketchFiles, AdditionalFiles: additionalFiles}, nil
}
