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
	"path/filepath"

	"github.com/arduino/arduino-builder/constants"
	"github.com/arduino/arduino-builder/i18n"
	"github.com/arduino/arduino-builder/types"
	"github.com/arduino/arduino-builder/utils"
)

type ContainerAddPrototypes struct{}

func (s *ContainerAddPrototypes) Run(ctx *types.Context) error {
	sourceFile := filepath.Join(ctx.SketchBuildPath, filepath.Base(ctx.Sketch.MainFile.Name)+".cpp")

	// Generate the full pathname for the preproc output file
	err := utils.EnsureFolderExists(ctx.PreprocPath)
	if err != nil {
		return i18n.WrapError(err)
	}
	targetFilePath := filepath.Join(ctx.PreprocPath, constants.FILE_CTAGS_TARGET_FOR_GCC_MINUS_E)

	// Run preprocessor
	err = GCCPreprocRunner(ctx, sourceFile, targetFilePath, ctx.IncludeFolders)
	if err != nil {
		return i18n.WrapError(err)
	}
	commands := []types.Command{
		&ReadFileAndStoreInContext{FileToRead: targetFilePath, Target: &ctx.SourceGccMinusE},
		&FilterSketchSource{Source: &ctx.SourceGccMinusE},
		&CTagsTargetFileSaver{Source: &ctx.SourceGccMinusE, TargetFileName: constants.FILE_CTAGS_TARGET_FOR_GCC_MINUS_E},
		&CTagsRunner{},
		&PrototypesAdder{},
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
