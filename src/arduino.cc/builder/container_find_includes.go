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

	"arduino.cc/builder/constants"
	"arduino.cc/builder/i18n"
	"arduino.cc/builder/types"
	"arduino.cc/builder/utils"
)

type ContainerFindIncludes struct{}

func (s *ContainerFindIncludes) Run(ctx *types.Context) error {
	appendIncludeFolder(ctx, ctx.BuildProperties[constants.BUILD_PROPERTIES_BUILD_CORE_PATH])
	if ctx.BuildProperties[constants.BUILD_PROPERTIES_BUILD_VARIANT_PATH] != constants.EMPTY_STRING {
		appendIncludeFolder(ctx, ctx.BuildProperties[constants.BUILD_PROPERTIES_BUILD_VARIANT_PATH])
	}

	sketchBuildPath := ctx.SketchBuildPath
	sketch := ctx.Sketch
	ctx.CollectedSourceFiles.Push(filepath.Join(sketchBuildPath, filepath.Base(sketch.MainFile.Name)+".cpp"))

	sourceFilePaths := ctx.CollectedSourceFiles
	queueSourceFilesFromFolder(sourceFilePaths, ctx.SketchBuildPath, /* recurse */ false)
	srcSubfolderPath := filepath.Join(ctx.SketchBuildPath, constants.SKETCH_FOLDER_SRC)
	if info, err := os.Stat(srcSubfolderPath); err == nil && info.IsDir() {
		queueSourceFilesFromFolder(sourceFilePaths, srcSubfolderPath, /* recurse */ true)
	}

	for !sourceFilePaths.Empty() {
		err := findIncludesUntilDone(ctx, sourceFilePaths.Pop().(string))
		if err != nil {
			return i18n.WrapError(err)
		}
	}

	err := runCommand(ctx, &FailIfImportedLibraryIsWrong{})
	if err != nil {
		return i18n.WrapError(err)
	}

	return nil
}

// Append the given folder to the include path.
func appendIncludeFolder(ctx *types.Context, folder string) {
	ctx.IncludeFolders = append(ctx.IncludeFolders, folder)
}

func runCommand(ctx *types.Context, command types.Command) error {
	PrintRingNameIfDebug(ctx, command)
	err := command.Run(ctx)
	if err != nil {
		return i18n.WrapError(err)
	}
	return nil
}

func findIncludesUntilDone(ctx *types.Context, sourceFilePath string) error {
	targetFilePath := utils.NULLFile()
	for {
		commands := []types.Command{
			&GCCPreprocRunnerForDiscoveringIncludes{SourceFilePath: sourceFilePath, TargetFilePath: targetFilePath},
			&IncludesFinderWithRegExp{Source: &ctx.SourceGccMinusE},
		}
		for _, command := range commands {
			err := runCommand(ctx, command)
			if err != nil {
				return i18n.WrapError(err)
			}
		}
		if ctx.IncludeJustFound == "" {
			// No missing includes found, we're done
			return nil
		}

		library := ResolveLibrary(ctx, ctx.IncludeJustFound)
		if library == nil {
			// Library could not be resolved, show error
			err := runCommand(ctx, &GCCPreprocRunner{TargetFileName: constants.FILE_CTAGS_TARGET_FOR_GCC_MINUS_E})
			return i18n.WrapError(err)
		}

		// Add this library to the list of libraries, the
		// include path and queue its source files for further
		// include scanning
		ctx.ImportedLibraries = append(ctx.ImportedLibraries, library)
		appendIncludeFolder(ctx, library.SrcFolder)
		sourceFolders := types.LibraryToSourceFolder(library)
		for _, sourceFolder := range sourceFolders {
			queueSourceFilesFromFolder(ctx.CollectedSourceFiles, sourceFolder.Folder, sourceFolder.Recurse)
		}
	}
}

func queueSourceFilesFromFolder(queue *types.UniqueStringQueue, folder string, recurse bool) error {
	extensions := func(ext string) bool { return ADDITIONAL_FILE_VALID_EXTENSIONS_NO_HEADERS[ext] }

	filePaths := []string{}
	err := utils.FindFilesInFolder(&filePaths, folder, extensions, recurse)
	if err != nil {
		return i18n.WrapError(err)
	}

	for _, filePath := range filePaths {
		queue.Push(filePath)
	}

	return nil
}
