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
	"arduino.cc/builder/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type SketchLoader struct{}

func (s *SketchLoader) Run(context map[string]interface{}) error {
	if !utils.MapHas(context, constants.CTX_SKETCH_LOCATION) {
		return nil
	}

	sketchLocation := context[constants.CTX_SKETCH_LOCATION].(string)

	sketchLocation, err := filepath.Abs(sketchLocation)
	if err != nil {
		return utils.WrapError(err)
	}
	mainSketchStat, err := os.Stat(sketchLocation)
	if err != nil {
		return utils.WrapError(err)
	}
	if mainSketchStat.IsDir() {
		sketchLocation = filepath.Join(sketchLocation, mainSketchStat.Name()+".ino")
	}
	context[constants.CTX_SKETCH_LOCATION] = sketchLocation

	allSketchFilePaths, err := collectAllSketchFiles(filepath.Dir(sketchLocation))
	if err != nil {
		return utils.WrapError(err)
	}

	if !utils.SliceContains(allSketchFilePaths, sketchLocation) {
		return utils.Errorf(context, constants.MSG_CANT_FIND_SKETCH_IN_PATH, sketchLocation, filepath.Dir(sketchLocation))
	}

	logger := context[constants.CTX_LOGGER].(i18n.Logger)
	sketch, err := makeSketch(sketchLocation, allSketchFilePaths, logger)

	context[constants.CTX_SKETCH_LOCATION] = sketchLocation
	context[constants.CTX_SKETCH] = sketch

	return nil
}

func collectAllSketchFiles(from string) ([]string, error) {
	filePaths := []string{}
	walkFunc := func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(currentPath))
		if !MAIN_FILE_VALID_EXTENSIONS[ext] && !ADDITIONAL_FILE_VALID_EXTENSIONS[ext] {
			return nil
		}
		currentFile, err := os.Open(currentPath)
		if err != nil {
			return nil
		}
		currentFile.Close()

		filePaths = append(filePaths, currentPath)
		return nil
	}
	err := gohasissues.Walk(from, walkFunc)
	return filePaths, utils.WrapError(err)
}

func makeSketch(sketchLocation string, allSketchFilePaths []string, logger i18n.Logger) (*types.Sketch, error) {
	sketchFilesMap := make(map[string]types.SketchFile)
	for _, sketchFilePath := range allSketchFilePaths {
		source, err := ioutil.ReadFile(sketchFilePath)
		if err != nil {
			return nil, utils.WrapError(err)
		}
		sketchFilesMap[sketchFilePath] = types.SketchFile{Name: sketchFilePath, Source: string(source)}
	}

	mainFile := sketchFilesMap[sketchLocation]
	delete(sketchFilesMap, sketchLocation)

	additionalFiles := []types.SketchFile{}
	otherSketchFiles := []types.SketchFile{}
	for _, sketchFile := range sketchFilesMap {
		ext := strings.ToLower(filepath.Ext(sketchFile.Name))
		if MAIN_FILE_VALID_EXTENSIONS[ext] {
			otherSketchFiles = append(otherSketchFiles, sketchFile)
		} else if ADDITIONAL_FILE_VALID_EXTENSIONS[ext] {
			additionalFiles = append(additionalFiles, sketchFile)
		} else {
			return nil, utils.ErrorfWithLogger(logger, constants.MSG_UNKNOWN_SKETCH_EXT, sketchFile.Name)
		}
	}

	sort.Sort(types.SketchFileSortByName(additionalFiles))
	sort.Sort(types.SketchFileSortByName(otherSketchFiles))

	return &types.Sketch{MainFile: mainFile, OtherSketchFiles: otherSketchFiles, AdditionalFiles: additionalFiles}, nil
}
