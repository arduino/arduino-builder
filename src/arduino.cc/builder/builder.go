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
	"reflect"
	"strconv"
	"time"

	"arduino.cc/builder/constants"
	"arduino.cc/builder/i18n"
	"arduino.cc/builder/phases"
	"arduino.cc/builder/types"
	"arduino.cc/builder/utils"
)

var MAIN_FILE_VALID_EXTENSIONS = map[string]bool{".ino": true, ".pde": true}
var ADDITIONAL_FILE_VALID_EXTENSIONS = map[string]bool{".h": true, ".c": true, ".hpp": true, ".hh": true, ".cpp": true, ".s": true}
var ADDITIONAL_FILE_VALID_EXTENSIONS_NO_HEADERS = map[string]bool{".c": true, ".cpp": true, ".s": true}

var LIBRARY_MANDATORY_PROPERTIES = []string{constants.LIBRARY_NAME, constants.LIBRARY_VERSION, constants.LIBRARY_AUTHOR, constants.LIBRARY_MAINTAINER}
var LIBRARY_NOT_SO_MANDATORY_PROPERTIES = []string{constants.LIBRARY_SENTENCE, constants.LIBRARY_PARAGRAPH, constants.LIBRARY_URL}
var LIBRARY_CATEGORIES = map[string]bool{
	"Display":             true,
	"Communication":       true,
	"Signal Input/Output": true,
	"Sensors":             true,
	"Device Control":      true,
	"Timing":              true,
	"Data Storage":        true,
	"Data Processing":     true,
	"Other":               true,
	"Uncategorized":       true,
}

const DEFAULT_DEBUG_LEVEL = 5
const DEFAULT_WARNINGS_LEVEL = "none"
const DEFAULT_SOFTWARE = "ARDUINO"
const DEFAULT_BUILD_CORE = "arduino"

type Builder struct{}

func (s *Builder) Run(ctx *types.Context) error {
	commands := []types.Command{
		&GenerateBuildPathIfMissing{},
		&EnsureBuildPathExists{},

		&ContainerSetupHardwareToolsLibsSketchAndProps{},

		&ContainerBuildOptions{},

		&WarnAboutPlatformRewrites{},

		&RecipeByPrefixSuffixRunner{Prefix: constants.HOOKS_PREBUILD, Suffix: constants.HOOKS_PATTERN_SUFFIX},

		&ContainerMergeCopySketchFiles{},

		utils.LogIfVerbose(constants.LOG_LEVEL_INFO, "Detecting libraries used..."),
		&ContainerFindIncludes{},

		&WarnAboutArchIncompatibleLibraries{},

		utils.LogIfVerbose(constants.LOG_LEVEL_INFO, "Generating function prototypes..."),
		&ContainerAddPrototypes{},

		utils.LogIfVerbose(constants.LOG_LEVEL_INFO, "Compiling sketch..."),
		&RecipeByPrefixSuffixRunner{Prefix: constants.HOOKS_SKETCH_PREBUILD, Suffix: constants.HOOKS_PATTERN_SUFFIX},
		&phases.SketchBuilder{},
		&RecipeByPrefixSuffixRunner{Prefix: constants.HOOKS_SKETCH_POSTBUILD, Suffix: constants.HOOKS_PATTERN_SUFFIX},

		utils.LogIfVerbose(constants.LOG_LEVEL_INFO, "Compiling libraries..."),
		&RecipeByPrefixSuffixRunner{Prefix: constants.HOOKS_LIBRARIES_PREBUILD, Suffix: constants.HOOKS_PATTERN_SUFFIX},
		&UnusedCompiledLibrariesRemover{},
		&phases.LibrariesBuilder{},
		&RecipeByPrefixSuffixRunner{Prefix: constants.HOOKS_LIBRARIES_POSTBUILD, Suffix: constants.HOOKS_PATTERN_SUFFIX},

		utils.LogIfVerbose(constants.LOG_LEVEL_INFO, "Compiling core..."),
		&RecipeByPrefixSuffixRunner{Prefix: constants.HOOKS_CORE_PREBUILD, Suffix: constants.HOOKS_PATTERN_SUFFIX},
		&phases.CoreBuilder{},
		&RecipeByPrefixSuffixRunner{Prefix: constants.HOOKS_CORE_POSTBUILD, Suffix: constants.HOOKS_PATTERN_SUFFIX},

		utils.LogIfVerbose(constants.LOG_LEVEL_INFO, "Linking everything together..."),
		&RecipeByPrefixSuffixRunner{Prefix: constants.HOOKS_LINKING_PRELINK, Suffix: constants.HOOKS_PATTERN_SUFFIX},
		&phases.Linker{},
		&RecipeByPrefixSuffixRunner{Prefix: constants.HOOKS_LINKING_POSTLINK, Suffix: constants.HOOKS_PATTERN_SUFFIX},

		&RecipeByPrefixSuffixRunner{Prefix: constants.HOOKS_OBJCOPY_PREOBJCOPY, Suffix: constants.HOOKS_PATTERN_SUFFIX},
		&RecipeByPrefixSuffixRunner{Prefix: "recipe.objcopy.", Suffix: constants.HOOKS_PATTERN_SUFFIX},
		&RecipeByPrefixSuffixRunner{Prefix: constants.HOOKS_OBJCOPY_POSTOBJCOPY, Suffix: constants.HOOKS_PATTERN_SUFFIX},

		&MergeSketchWithBootloader{},

		&RecipeByPrefixSuffixRunner{Prefix: constants.HOOKS_POSTBUILD, Suffix: constants.HOOKS_PATTERN_SUFFIX},
	}

	mainErr := runCommands(ctx, commands, true)

	commands = []types.Command{
		&PrintUsedAndNotUsedLibraries{SketchError: mainErr != nil},

		&PrintUsedLibrariesIfVerbose{},

		&phases.Sizer{SketchError: mainErr != nil},
	}
	otherErr := runCommands(ctx, commands, false)

	if mainErr != nil {
		return mainErr
	}

	return otherErr
}

type Preprocess struct{}

func (s *Preprocess) Run(ctx *types.Context) error {
	commands := []types.Command{
		&GenerateBuildPathIfMissing{},
		&EnsureBuildPathExists{},

		&ContainerSetupHardwareToolsLibsSketchAndProps{},

		&ContainerBuildOptions{},

		&RecipeByPrefixSuffixRunner{Prefix: constants.HOOKS_PREBUILD, Suffix: constants.HOOKS_PATTERN_SUFFIX},

		&ContainerMergeCopySketchFiles{},

		&ContainerFindIncludes{},

		&WarnAboutArchIncompatibleLibraries{},

		&ContainerAddPrototypes{},

		&PrintPreprocessedSource{},
	}

	return runCommands(ctx, commands, true)
}

type ParseHardwareAndDumpBuildProperties struct{}

func (s *ParseHardwareAndDumpBuildProperties) Run(ctx *types.Context) error {
	commands := []types.Command{
		&GenerateBuildPathIfMissing{},

		&ContainerSetupHardwareToolsLibsSketchAndProps{},

		&DumpBuildProperties{},
	}

	return runCommands(ctx, commands, true)
}

func runCommands(ctx *types.Context, commands []types.Command, progressEnabled bool) error {
	commandsLength := len(commands)
	progressForEachCommand := float32(100) / float32(commandsLength)

	progress := float32(0)
	for _, command := range commands {
		PrintRingNameIfDebug(ctx, command)
		printProgressIfProgressEnabledAndMachineLogger(progressEnabled, ctx, progress)
		err := command.Run(ctx)
		if err != nil {
			return i18n.WrapError(err)
		}
		progress += progressForEachCommand
	}

	printProgressIfProgressEnabledAndMachineLogger(progressEnabled, ctx, 100)

	return nil
}

func printProgressIfProgressEnabledAndMachineLogger(progressEnabled bool, ctx *types.Context, progress float32) {
	if !progressEnabled {
		return
	}

	log := ctx.GetLogger()
	if log.Name() == "machine" {
		log.Println(constants.LOG_LEVEL_INFO, constants.MSG_PROGRESS, strconv.FormatFloat(float64(progress), 'f', 2, 32))
	}
}

func PrintRingNameIfDebug(ctx *types.Context, command types.Command) {
	if ctx.DebugLevel >= 10 {
		ctx.GetLogger().Fprintln(os.Stdout, constants.LOG_LEVEL_DEBUG, constants.MSG_RUNNING_COMMAND, strconv.FormatInt(time.Now().Unix(), 10), reflect.Indirect(reflect.ValueOf(command)).Type().Name())
	}
}

func RunBuilder(ctx *types.Context) error {
	command := Builder{}
	return command.Run(ctx)
}

func RunParseHardwareAndDumpBuildProperties(ctx *types.Context) error {
	command := ParseHardwareAndDumpBuildProperties{}
	return command.Run(ctx)
}

func RunPreprocess(ctx *types.Context) error {
	command := Preprocess{}
	return command.Run(ctx)
}
