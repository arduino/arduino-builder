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
	"github.com/arduino/arduino-builder/builder/i18n"
	"github.com/arduino/arduino-builder/builder/props"
	"github.com/arduino/arduino-builder/builder/types"
	"github.com/arduino/arduino-builder/builder/utils"
	"os"
	"path/filepath"
	"strings"
)

type MergeSketchWithBootloader struct{}

func (s *MergeSketchWithBootloader) Run(context map[string]interface{}) error {
	buildProperties := context[constants.CTX_BUILD_PROPERTIES].(props.PropertiesMap)
	if !utils.MapStringStringHas(buildProperties, constants.BUILD_PROPERTIES_BOOTLOADER_NOBLINK) && !utils.MapStringStringHas(buildProperties, constants.BUILD_PROPERTIES_BOOTLOADER_FILE) {
		return nil
	}

	buildPath := context[constants.CTX_BUILD_PATH].(string)
	sketch := context[constants.CTX_SKETCH].(*types.Sketch)
	sketchFileName := filepath.Base(sketch.MainFile.Name)
	logger := context[constants.CTX_LOGGER].(i18n.Logger)

	sketchInBuildPath := filepath.Join(buildPath, sketchFileName+".hex")
	sketchInSubfolder := filepath.Join(buildPath, constants.FOLDER_SKETCH, sketchFileName+".hex")

	builtSketchPath := constants.EMPTY_STRING
	if _, err := os.Stat(sketchInBuildPath); err == nil {
		builtSketchPath = sketchInBuildPath
	} else if _, err := os.Stat(sketchInSubfolder); err == nil {
		builtSketchPath = sketchInSubfolder
	} else {
		return nil
	}

	bootloader := constants.EMPTY_STRING
	if utils.MapStringStringHas(buildProperties, constants.BUILD_PROPERTIES_BOOTLOADER_NOBLINK) {
		bootloader = buildProperties[constants.BUILD_PROPERTIES_BOOTLOADER_NOBLINK]
	} else {
		bootloader = buildProperties[constants.BUILD_PROPERTIES_BOOTLOADER_FILE]
	}
	bootloader = buildProperties.ExpandPropsInString(bootloader)

	bootloaderPath := filepath.Join(buildProperties[constants.BUILD_PROPERTIES_RUNTIME_PLATFORM_PATH], constants.FOLDER_BOOTLOADERS, bootloader)
	if _, err := os.Stat(bootloaderPath); err != nil {
		logger.Fprintln(os.Stdout, constants.LOG_LEVEL_WARN, constants.MSG_BOOTLOADER_FILE_MISSING, bootloaderPath)
		return nil
	}

	mergedSketchPath := filepath.Join(filepath.Dir(builtSketchPath), sketchFileName+".with_bootloader.hex")

	err := merge(builtSketchPath, bootloaderPath, mergedSketchPath)

	return err
}

func merge(builtSketchPath, bootloaderPath, mergedSketchPath string) error {
	sketch, err := utils.ReadFileToRows(builtSketchPath)
	if err != nil {
		return utils.WrapError(err)
	}
	sketch = sketch[:len(sketch)-2]

	bootloader, err := utils.ReadFileToRows(bootloaderPath)
	if err != nil {
		return utils.WrapError(err)
	}

	for _, row := range bootloader {
		sketch = append(sketch, row)
	}

	return utils.WriteFile(mergedSketchPath, strings.Join(sketch, "\n"))
}
