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
	"strings"

	"github.com/arduino/arduino-builder/constants"
	"github.com/arduino/arduino-builder/i18n"
	"github.com/arduino/arduino-builder/types"
	"github.com/arduino/arduino-builder/utils"
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

	sketchInBuildPath := filepath.Join(buildPath, sketchFileName)
	sketchInSubfolder := filepath.Join(buildPath, constants.FOLDER_SKETCH, sketchFileName)

	availableExtensions := []string{".hex", ".bin"}
	builtSketchPath := constants.EMPTY_STRING

	extension := ""

	for _, extension = range availableExtensions {
		if _, err := os.Stat(sketchInBuildPath + extension); err == nil {
			builtSketchPath = sketchInBuildPath + extension
			break
		} else if _, err := os.Stat(sketchInSubfolder + extension); err == nil {
			builtSketchPath = sketchInSubfolder + extension
			break
		}
	}

	if builtSketchPath == constants.EMPTY_STRING {
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

	mergedSketchPath := filepath.Join(filepath.Dir(builtSketchPath), sketchFileName+".with_bootloader"+extension)

	var err error
	if extension == ".hex" {
		err = mergeHex(builtSketchPath, bootloaderPath, mergedSketchPath)
	} else {
		ldscript := buildProperties[constants.BUILD_PROPERTIES_BUILD_LDSCRIPT]
		variantFolder := buildProperties[constants.BUILD_PROPERTIES_BUILD_VARIANT_PATH]
		ldscriptPath := filepath.Join(variantFolder, ldscript)
		err = mergeBin(builtSketchPath, ldscriptPath, bootloaderPath, mergedSketchPath)
	}

	return err
}

func hexLineOnlyContainsFF(line string) bool {
	//:206FE000FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFB1
	if len(line) <= 11 {
		return false
	}
	byteArray := []byte(line)
	for _, char := range byteArray[9:(len(byteArray) - 2)] {
		if char != 'F' {
			return false
		}
	}
	return true
}

func extractActualBootloader(bootloader []string) []string {

	var realBootloader []string

	// skip until we find a line full of FFFFFF (except address and checksum)
	for i, row := range bootloader {
		if hexLineOnlyContainsFF(row) {
			realBootloader = bootloader[i:len(bootloader)]
			break
		}
	}

	// drop all "empty" lines
	for i, row := range realBootloader {
		if !hexLineOnlyContainsFF(row) {
			realBootloader = realBootloader[i:len(realBootloader)]
			break
		}
	}

	if len(realBootloader) == 0 {
		// we didn't find any line full of FFFF, thus it's a standalone bootloader
		realBootloader = bootloader
	}

	return realBootloader
}

func mergeHex(builtSketchPath, bootloaderPath, mergedSketchPath string) error {
	sketch, err := utils.ReadFileToRows(builtSketchPath)
	if err != nil {
		return i18n.WrapError(err)
	}
	sketch = sketch[:len(sketch)-2]

	bootloader, err := utils.ReadFileToRows(bootloaderPath)
	if err != nil {
		return i18n.WrapError(err)
	}

	realBootloader := extractActualBootloader(bootloader)

	for _, row := range realBootloader {
		sketch = append(sketch, row)
	}

	return utils.WriteFile(mergedSketchPath, strings.Join(sketch, "\n"))
}

func mergeBin(builtSketchPath, ldscriptPath, bootloaderPath, mergedSketchPath string) error {
	// 0xFF means empty
	// only works if the bootloader is at the beginning of the flash
	// only works if the flash address space is mapped at 0x00

	// METHOD 1: (non appliable to most architectures)
	// remove all comments from linkerscript
	// find NAMESPACE of .text section -> FLASH
	// find ORIGIN of FLASH section

	// METHOD 2:
	// Round the bootloader to the next "power of 2" bytes boundary

	// generate a byte[FLASH] full of 0xFF and bitwise OR with bootloader BIN
	// merge this slice with sketch BIN

	bootloader, _ := ioutil.ReadFile(bootloaderPath)
	sketch, _ := ioutil.ReadFile(builtSketchPath)

	paddedBootloaderLen := nextPowerOf2(len(bootloader))

	padding := make([]byte, paddedBootloaderLen-len(bootloader))
	for i, _ := range padding {
		padding[i] = 0xFF
	}

	bootloader = append(bootloader, padding...)
	sketch = append(bootloader, sketch...)

	return ioutil.WriteFile(mergedSketchPath, sketch, 0644)
}

func nextPowerOf2(v int) int {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return v
}
