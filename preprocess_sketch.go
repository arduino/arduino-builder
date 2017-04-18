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
	"fmt"
	"path/filepath"

	"io/ioutil"

	"github.com/arduino/arduino-builder/constants"
	"github.com/arduino/arduino-builder/i18n"
	"github.com/arduino/arduino-builder/types"
	"github.com/arduino/arduino-builder/utils"
)

type PreprocessSketch struct{}

func (s *PreprocessSketch) Run(ctx *types.Context) error {
	sourceFile := filepath.Join(ctx.SketchBuildPath, filepath.Base(ctx.Sketch.MainFile.Name)+".cpp")
	commands := []types.Command{
		&GCCPreprocRunner{SourceFilePath: sourceFile, TargetFileName: constants.FILE_CTAGS_TARGET_FOR_GCC_MINUS_E, Includes: ctx.IncludeFolders},
		&ArduinoPreprocessorRunner{},
		&SketchSaver{},
	}

	for _, command := range commands {
		PrintRingNameIfDebug(ctx, command)
		err := command.Run(ctx)
		if err != nil {
			return i18n.WrapError(err)
		}
	}

	return nil
}

type ArduinoPreprocessorRunner struct{}

func (s *ArduinoPreprocessorRunner) Run(ctx *types.Context) error {
	buildProperties := ctx.BuildProperties
	targetFilePath := ctx.FileToRead
	logger := ctx.GetLogger()

	properties := buildProperties.Clone()
	toolProps := buildProperties.SubTree("tools").SubTree("arduino-preprocessor")
	properties.Merge(toolProps)
	properties[constants.BUILD_PROPERTIES_SOURCE_FILE] = targetFilePath

	pattern := properties[constants.BUILD_PROPERTIES_PATTERN]
	if pattern == constants.EMPTY_STRING {
		return i18n.ErrorfWithLogger(logger, constants.MSG_PATTERN_MISSING, "arduino-preprocessor")
	}

	commandLine := properties.ExpandPropsInString(pattern)
	command, err := utils.PrepareCommand(commandLine, logger)
	if err != nil {
		return i18n.WrapError(err)
	}

	verbose := ctx.Verbose
	if verbose {
		fmt.Println(commandLine)
	}

	sourceBytes, err := command.Output()
	if err != nil {
		return i18n.WrapError(err)
	}

	ctx.Source = string(sourceBytes)
	if ctx.Source == "" {
		// If the output is empty the original sketch may pass through
		buf, err := ioutil.ReadFile(targetFilePath)
		if err != nil {
			return i18n.WrapError(err)
		}
		ctx.Source = string(buf)
	}
	//fmt.Printf("PREPROCESSOR OUTPUT:\n%s\n", ctx.Source)
	return nil
}
