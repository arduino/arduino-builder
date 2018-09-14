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
	"os"
	"path/filepath"

	"github.com/arduino/arduino-builder/constants"
	"github.com/arduino/arduino-builder/types"
	"github.com/arduino/arduino-builder/utils"
	"github.com/marcinbor85/gohex"
)

type MergeSketchWithBootloader struct{}

func (s *MergeSketchWithBootloader) Run(ctx *types.Context) error {
	buildProperties := ctx.BuildProperties
	if !utils.MapStringStringHas(buildProperties, constants.BUILD_PROPERTIES_BOOTLOADER_NOBLINK) && !utils.MapStringStringHas(buildProperties, constants.BUILD_PROPERTIES_BOOTLOADER_FILE) {
		return nil
	}

	buildPath := ctx.BuildPath
	sketch := ctx.Sketch
	sketchFileName := filepath.Base(sketch.MainFile.Name)
	logger := ctx.GetLogger()

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

	// Ignore merger errors for the first iteration
	merge(builtSketchPath, bootloaderPath, mergedSketchPath)

	return nil
}

func merge(builtSketchPath, bootloaderPath, mergedSketchPath string) error {
	bootFile, err := os.Open(bootloaderPath)
	if err != nil {
		return err
	}
	defer bootFile.Close()

	mem_boot := gohex.NewMemory()
	err = mem_boot.ParseIntelHex(bootFile)
	if err != nil {
		return err
	}

	buildFile, err := os.Open(builtSketchPath)
	if err != nil {
		return err
	}
	defer buildFile.Close()

	mem_sketch := gohex.NewMemory()
	err = mem_sketch.ParseIntelHex(buildFile)
	if err != nil {
		return err
	}

	mem_merge := gohex.NewMemory()

	for _, segment := range mem_boot.GetDataSegments() {
		err = mem_merge.AddBinary(segment.Address, segment.Data)
	}
	for _, segment := range mem_sketch.GetDataSegments() {
		err = mem_merge.AddBinary(segment.Address, segment.Data)
	}

	mergeFile, err := os.Create(mergedSketchPath)
	if err != nil {
		return err
	}
	defer mergeFile.Close()

	mem_merge.DumpIntelHex(mergeFile, 32)
	return nil
	//bytes := mem_merge.ToBinary(0, len, 0xFF)
	//return utils.WriteFile(mergedSketchPath, string(bytes))
}
