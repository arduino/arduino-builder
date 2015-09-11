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

package builder_utils

import (
	"arduino.cc/builder/constants"
	"arduino.cc/builder/i18n"
	"arduino.cc/builder/props"
	"arduino.cc/builder/utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func CompileFilesRecursive(objectFiles []string, sourcePath string, buildPath string, buildProperties map[string]string, includes []string, verbose bool, warningsLevel string, logger i18n.Logger) ([]string, error) {
	objectFiles, err := CompileFiles(objectFiles, sourcePath, false, buildPath, buildProperties, includes, verbose, warningsLevel, logger)
	if err != nil {
		return nil, utils.WrapError(err)
	}

	folders, err := utils.ReadDirFiltered(sourcePath, utils.FilterDirs)
	if err != nil {
		return nil, utils.WrapError(err)
	}

	for _, folder := range folders {
		objectFiles, err = CompileFilesRecursive(objectFiles, filepath.Join(sourcePath, folder.Name()), filepath.Join(buildPath, folder.Name()), buildProperties, includes, verbose, warningsLevel, logger)
		if err != nil {
			return nil, utils.WrapError(err)
		}
	}

	return objectFiles, nil
}

func CompileFiles(objectFiles []string, sourcePath string, recurse bool, buildPath string, buildProperties map[string]string, includes []string, verbose bool, warningsLevel string, logger i18n.Logger) ([]string, error) {
	objectFiles, err := compileFilesWithExtensionWithRecipe(objectFiles, sourcePath, recurse, buildPath, buildProperties, includes, ".S", constants.RECIPE_S_PATTERN, verbose, warningsLevel, logger)
	if err != nil {
		return nil, utils.WrapError(err)
	}
	objectFiles, err = compileFilesWithExtensionWithRecipe(objectFiles, sourcePath, recurse, buildPath, buildProperties, includes, ".c", constants.RECIPE_C_PATTERN, verbose, warningsLevel, logger)
	if err != nil {
		return nil, utils.WrapError(err)
	}
	objectFiles, err = compileFilesWithExtensionWithRecipe(objectFiles, sourcePath, recurse, buildPath, buildProperties, includes, ".cpp", constants.RECIPE_CPP_PATTERN, verbose, warningsLevel, logger)
	if err != nil {
		return nil, utils.WrapError(err)
	}
	return objectFiles, nil
}

func compileFilesWithExtensionWithRecipe(objectFiles []string, sourcePath string, recurse bool, buildPath string, buildProperties map[string]string, includes []string, extension string, recipe string, verbose bool, warningsLevel string, logger i18n.Logger) ([]string, error) {
	sources, err := findFilesInFolder(sourcePath, extension, recurse)
	if err != nil {
		return nil, utils.WrapError(err)
	}
	return compileWithRecipe(objectFiles, sourcePath, sources, buildPath, buildProperties, includes, recipe, verbose, warningsLevel, logger)
}

func findFilesInFolder(sourcePath string, extension string, recurse bool) ([]string, error) {
	files, err := utils.ReadDirFiltered(sourcePath, utils.FilterFilesWithExtension(extension))
	if err != nil {
		return nil, utils.WrapError(err)
	}
	var sources []string
	for _, file := range files {
		sources = append(sources, filepath.Join(sourcePath, file.Name()))
	}

	if recurse {
		folders, err := utils.ReadDirFiltered(sourcePath, utils.FilterDirs)
		if err != nil {
			return nil, utils.WrapError(err)
		}

		for _, folder := range folders {
			otherSources, err := findFilesInFolder(filepath.Join(sourcePath, folder.Name()), extension, recurse)
			if err != nil {
				return nil, utils.WrapError(err)
			}
			sources = append(sources, otherSources...)
		}
	}

	return sources, nil
}

func compileWithRecipe(objectFiles []string, sourcePath string, sources []string, buildPath string, buildProperties map[string]string, includes []string, recipe string, verbose bool, warningsLevel string, logger i18n.Logger) ([]string, error) {
	for _, source := range sources {
		properties := utils.MergeMapsOfStrings(make(map[string]string), buildProperties)
		properties[constants.BUILD_PROPERTIES_COMPILER_WARNING_FLAGS] = properties[constants.BUILD_PROPERTIES_COMPILER_WARNING_FLAGS+"."+warningsLevel]
		properties[constants.BUILD_PROPERTIES_INCLUDES] = strings.Join(includes, constants.SPACE)
		properties[constants.BUILD_PROPERTIES_SOURCE_FILE] = source
		relativeSource, err := filepath.Rel(sourcePath, source)
		if err != nil {
			return nil, utils.WrapError(err)
		}
		properties[constants.BUILD_PROPERTIES_OBJECT_FILE] = filepath.Join(buildPath, relativeSource+".o")

		err = os.MkdirAll(filepath.Dir(properties[constants.BUILD_PROPERTIES_OBJECT_FILE]), os.FileMode(0755))
		if err != nil {
			return nil, utils.WrapError(err)
		}

		sourceFileStat, err := os.Stat(properties[constants.BUILD_PROPERTIES_SOURCE_FILE])
		if err != nil {
			return nil, utils.WrapError(err)
		}

		objectFileStat, err := os.Stat(properties[constants.BUILD_PROPERTIES_OBJECT_FILE])
		if err != nil && !os.IsNotExist(err) {
			return nil, utils.WrapError(err)
		}

		if !objFileIsUpToDateWithSourceFile(sourceFileStat, objectFileStat) {
			_, err = ExecRecipe(properties, recipe, false, verbose, verbose, logger)
			if err != nil {
				return nil, utils.WrapError(err)
			}
		} else if verbose {
			logger.Println(constants.MSG_USING_PREVIOUS_COMPILED_FILE, properties[constants.BUILD_PROPERTIES_OBJECT_FILE])
		}

		objectFiles = append(objectFiles, properties[constants.BUILD_PROPERTIES_OBJECT_FILE])
	}
	return objectFiles, nil
}

func objFileIsUpToDateWithSourceFile(sourceFileStat, objectFileStat os.FileInfo) bool {
	return objectFileStat != nil && sourceFileStat.ModTime().Before(objectFileStat.ModTime())
}

func ExecRecipe(properties map[string]string, recipe string, removeUnsetProperties bool, echoCommandLine bool, echoOutput bool, logger i18n.Logger) ([]byte, error) {
	pattern := properties[recipe]
	if pattern == constants.EMPTY_STRING {
		return nil, utils.ErrorfWithLogger(logger, constants.MSG_PATTERN_MISSING, recipe)
	}

	var err error
	commandLine := props.ExpandPropsInString(properties, pattern)
	if removeUnsetProperties {
		commandLine, err = props.DeleteUnexpandedPropsFromString(commandLine)
		if err != nil {
			return nil, utils.WrapError(err)
		}
	}

	command, err := utils.PrepareCommand(commandLine, logger)
	if err != nil {
		return nil, utils.WrapError(err)
	}

	if echoCommandLine {
		fmt.Println(commandLine)
	}

	if echoOutput {
		command.Stdout = os.Stdout
	}

	command.Stderr = os.Stderr

	if echoOutput {
		err := command.Run()
		return nil, utils.WrapError(err)
	}

	bytes, err := command.Output()
	return bytes, utils.WrapError(err)
}
