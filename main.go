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

package main

import (
	"arduino.cc/builder"
	"arduino.cc/builder/constants"
	"arduino.cc/builder/gohasissues"
	"arduino.cc/builder/i18n"
	"arduino.cc/builder/utils"
	"flag"
	"fmt"
	"github.com/go-errors/errors"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

const VERSION = "1.0.0-beta5"

type slice []string

func (h *slice) String() string {
	return fmt.Sprint(*h)
}

func (h *slice) Set(csv string) error {
	var values []string
	if strings.Contains(csv, string(os.PathListSeparator)) {
		values = strings.Split(csv, string(os.PathListSeparator))
	} else {
		values = strings.Split(csv, ",")
	}

	for _, value := range values {
		value = strings.TrimSpace(value)
		*h = append(*h, value)
	}

	return nil
}

var compileFlag *bool
var dumpPrefsFlag *bool
var hardwareFoldersFlag slice
var toolsFoldersFlag slice
var librariesFoldersFlag slice
var customBuildPropertiesFlag slice
var fqbnFlag *string
var ideVersionFlag *string
var buildPathFlag *string
var verboseFlag *bool
var debugLevelFlag *int
var warningsLevelFlag *string
var loggerFlag *string
var versionFlag *bool

func init() {
	compileFlag = flag.Bool("compile", false, "compiles the given sketch")
	dumpPrefsFlag = flag.Bool("dump-prefs", false, "dumps build properties used when compiling")
	flag.Var(&hardwareFoldersFlag, "hardware", "Specify a 'hardware' folder. Can be added multiple times for specifying multiple 'hardware' folders")
	flag.Var(&toolsFoldersFlag, "tools", "Specify a 'tools' folder. Can be added multiple times for specifying multiple 'tools' folders")
	flag.Var(&librariesFoldersFlag, "libraries", "Specify a 'libraries' folder. Can be added multiple times for specifying multiple 'libraries' folders")
	flag.Var(&customBuildPropertiesFlag, "prefs", "Specify a custom preference. Can be added multiple times for specifying multiple custom preferences")
	fqbnFlag = flag.String("fqbn", "", "fully qualified board name")
	ideVersionFlag = flag.String("ide-version", "10600", "fake IDE version")
	buildPathFlag = flag.String("build-path", "", "build path")
	verboseFlag = flag.Bool("verbose", false, "if 'true' prints lots of stuff")
	debugLevelFlag = flag.Int("debug-level", builder.DEFAULT_DEBUG_LEVEL, "Turns on debugging messages. The higher, the chattier")
	warningsLevelFlag = flag.String("warnings", "", "Sets warnings level. Available values are 'none', 'default', 'more' and 'all'")
	loggerFlag = flag.String("logger", "human", "Sets type of logger. Available values are 'human', 'machine'")
	versionFlag = flag.Bool("version", false, "prints version and exits")
}

func main() {
	flag.Parse()

	if *versionFlag {
		fmt.Println("Arduino Builder " + VERSION)
		fmt.Println("Copyright (C) 2015 Arduino LLC and contributors")
		fmt.Println("See https://www.arduino.cc/ and https://github.com/arduino/arduino-builder/graphs/contributors")
		fmt.Println("This is free software; see the source for copying conditions.  There is NO")
		fmt.Println("warranty; not even for MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.")
		defer os.Exit(0)
		return
	}

	compile := *compileFlag
	dumpPrefs := *dumpPrefsFlag

	if compile && dumpPrefs {
		fmt.Println("You can either specify --compile or --dump-prefs, not both")
		defer os.Exit(1)
		return
	}

	if !compile && !dumpPrefs {
		compile = true
	}

	context := make(map[string]interface{})

	hardware, err := toSliceOfUnquoted(hardwareFoldersFlag)
	if err != nil {
		printCompleteError(err)
		defer os.Exit(1)
		return
	}

	if len(hardware) == 0 {
		fmt.Println("Parameter 'hardware' is mandatory")
		flag.Usage()
		defer os.Exit(1)
		return
	}
	context[constants.CTX_HARDWARE_FOLDERS] = hardware

	tools, err := toSliceOfUnquoted(toolsFoldersFlag)
	if err != nil {
		printCompleteError(err)
		defer os.Exit(1)
		return
	}

	if len(tools) == 0 {
		fmt.Println("Parameter 'tools' is mandatory")
		flag.Usage()
		defer os.Exit(1)
		return
	}
	context[constants.CTX_TOOLS_FOLDERS] = tools

	libraries, err := toSliceOfUnquoted(librariesFoldersFlag)
	if err != nil {
		printCompleteError(err)
		defer os.Exit(1)
		return
	}
	context[constants.CTX_LIBRARIES_FOLDERS] = libraries

	customBuildProperties, err := toSliceOfUnquoted(customBuildPropertiesFlag)
	if err != nil {
		printCompleteError(err)
		defer os.Exit(1)
		return
	}
	context[constants.CTX_CUSTOM_BUILD_PROPERTIES] = customBuildProperties

	fqbn, err := gohasissues.Unquote(*fqbnFlag)
	if err != nil {
		printCompleteError(err)
		defer os.Exit(1)
		return
	}

	if fqbn == "" {
		fmt.Println("Parameter 'fqbn' is mandatory")
		flag.Usage()
		defer os.Exit(1)
		return
	}
	context[constants.CTX_FQBN] = fqbn

	buildPath, err := gohasissues.Unquote(*buildPathFlag)
	if err != nil {
		printCompleteError(err)
		defer os.Exit(1)
		return
	}

	if buildPath != "" {
		_, err := os.Stat(buildPath)
		if err != nil {
			fmt.Println(err)
			defer os.Exit(1)
			return
		}

		err = os.MkdirAll(buildPath, os.FileMode(0755))
		if err != nil {
			printCompleteError(err)
			defer os.Exit(1)
			return
		}
	}
	context[constants.CTX_BUILD_PATH] = buildPath

	if compile && flag.NArg() == 0 {
		fmt.Println("Last parameter must be the sketch to compile")
		flag.Usage()
		defer os.Exit(1)
		return
	}

	if flag.NArg() > 0 {
		sketchLocation := flag.Arg(0)
		sketchLocation, err := gohasissues.Unquote(sketchLocation)
		if err != nil {
			printCompleteError(err)
			defer os.Exit(1)
			return
		}
		context[constants.CTX_SKETCH_LOCATION] = sketchLocation
	}

	context[constants.CTX_VERBOSE] = *verboseFlag
	context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION] = *ideVersionFlag

	if *warningsLevelFlag != "" {
		context[constants.CTX_WARNINGS_LEVEL] = *warningsLevelFlag
	}

	if *debugLevelFlag > -1 {
		context[constants.CTX_DEBUG_LEVEL] = *debugLevelFlag
	}

	if *loggerFlag == "machine" {
		context[constants.CTX_LOGGER] = i18n.MachineLogger{}
	} else {
		context[constants.CTX_LOGGER] = i18n.HumanLogger{}
	}

	if compile {
		err = builder.RunBuilder(context)
	} else if dumpPrefs {
		err = builder.RunParseHardwareAndDumpBuildProperties(context)
	}

	exitCode := 0
	if err != nil {
		err = utils.WrapError(err)

		fmt.Println(err)

		if utils.DebugLevel(context) >= 10 {
			fmt.Println(err.(*errors.Error).ErrorStack())
		}

		exitCode = toExitCode(err)
	}

	defer os.Exit(exitCode)
}

func toExitCode(err error) int {
	if exiterr, ok := err.(*exec.ExitError); ok {
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus()
		}
	}
	return 1
}

func toSliceOfUnquoted(value slice) ([]string, error) {
	var values []string
	for _, v := range value {
		v, err := gohasissues.Unquote(v)
		if err != nil {
			return nil, err
		}
		values = append(values, v)
	}
	return values, nil
}

func printCompleteError(err error) {
	err = utils.WrapError(err)
	fmt.Println(err.(*errors.Error).ErrorStack())
}
