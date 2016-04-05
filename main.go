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
 * Copyright 2015 Matthijs Kooijman
 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"arduino.cc/builder"
	"arduino.cc/builder/constants"
	"arduino.cc/builder/gohasissues"
	"arduino.cc/builder/i18n"
	"arduino.cc/builder/utils"
	"github.com/go-errors/errors"
)

const VERSION = "1.3.16"

const FLAG_ACTION_COMPILE = "compile"
const FLAG_ACTION_PREPROCESS = "preprocess"
const FLAG_ACTION_DUMP_PREFS = "dump-prefs"
const FLAG_BUILD_OPTIONS_FILE = "build-options-file"
const FLAG_HARDWARE = "hardware"
const FLAG_TOOLS = "tools"
const FLAG_BUILT_IN_LIBRARIES = "built-in-libraries"
const FLAG_LIBRARIES = "libraries"
const FLAG_PREFS = "prefs"
const FLAG_FQBN = "fqbn"
const FLAG_IDE_VERSION = "ide-version"
const FLAG_CORE_API_VERSION = "core-api-version"
const FLAG_BUILD_PATH = "build-path"
const FLAG_VERBOSE = "verbose"
const FLAG_QUIET = "quiet"
const FLAG_DEBUG_LEVEL = "debug-level"
const FLAG_WARNINGS = "warnings"
const FLAG_WARNINGS_NONE = "none"
const FLAG_WARNINGS_DEFAULT = "default"
const FLAG_WARNINGS_MORE = "more"
const FLAG_WARNINGS_ALL = "all"
const FLAG_LOGGER = "logger"
const FLAG_LOGGER_HUMAN = "human"
const FLAG_LOGGER_MACHINE = "machine"
const FLAG_VERSION = "version"
const FLAG_VID_PID = "vid-pid"

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
var preprocessFlag *bool
var dumpPrefsFlag *bool
var buildOptionsFileFlag *string
var hardwareFoldersFlag slice
var toolsFoldersFlag slice
var librariesBuiltInFoldersFlag slice
var librariesFoldersFlag slice
var customBuildPropertiesFlag slice
var fqbnFlag *string
var coreAPIVersionFlag *string
var ideVersionFlag *string
var buildPathFlag *string
var verboseFlag *bool
var quietFlag *bool
var debugLevelFlag *int
var warningsLevelFlag *string
var loggerFlag *string
var versionFlag *bool
var vidPidFlag *string

func init() {
	compileFlag = flag.Bool(FLAG_ACTION_COMPILE, false, "compiles the given sketch")
	preprocessFlag = flag.Bool(FLAG_ACTION_PREPROCESS, false, "preprocess the given sketch")
	dumpPrefsFlag = flag.Bool(FLAG_ACTION_DUMP_PREFS, false, "dumps build properties used when compiling")
	buildOptionsFileFlag = flag.String(FLAG_BUILD_OPTIONS_FILE, "", "Instead of specifying --"+FLAG_HARDWARE+", --"+FLAG_TOOLS+" etc every time, you can load all such options from a file")
	flag.Var(&hardwareFoldersFlag, FLAG_HARDWARE, "Specify a 'hardware' folder. Can be added multiple times for specifying multiple 'hardware' folders")
	flag.Var(&toolsFoldersFlag, FLAG_TOOLS, "Specify a 'tools' folder. Can be added multiple times for specifying multiple 'tools' folders")
	flag.Var(&librariesBuiltInFoldersFlag, FLAG_BUILT_IN_LIBRARIES, "Specify a built-in 'libraries' folder. These are low priority libraries. Can be added multiple times for specifying multiple built-in 'libraries' folders")
	flag.Var(&librariesFoldersFlag, FLAG_LIBRARIES, "Specify a 'libraries' folder. Can be added multiple times for specifying multiple 'libraries' folders")
	flag.Var(&customBuildPropertiesFlag, FLAG_PREFS, "Specify a custom preference. Can be added multiple times for specifying multiple custom preferences")
	fqbnFlag = flag.String(FLAG_FQBN, "", "fully qualified board name")
	coreAPIVersionFlag = flag.String(FLAG_CORE_API_VERSION, "10600", "version of core APIs (used to populate ARDUINO #define)")
	ideVersionFlag = flag.String(FLAG_IDE_VERSION, "10600", "[deprecated] use '"+FLAG_CORE_API_VERSION+"' instead")
	buildPathFlag = flag.String(FLAG_BUILD_PATH, "", "build path")
	verboseFlag = flag.Bool(FLAG_VERBOSE, false, "if 'true' prints lots of stuff")
	quietFlag = flag.Bool(FLAG_QUIET, false, "if 'true' doesn't print any warnings or progress or whatever")
	debugLevelFlag = flag.Int(FLAG_DEBUG_LEVEL, builder.DEFAULT_DEBUG_LEVEL, "Turns on debugging messages. The higher, the chattier")
	warningsLevelFlag = flag.String(FLAG_WARNINGS, "", "Sets warnings level. Available values are '"+FLAG_WARNINGS_NONE+"', '"+FLAG_WARNINGS_DEFAULT+"', '"+FLAG_WARNINGS_MORE+"' and '"+FLAG_WARNINGS_ALL+"'")
	loggerFlag = flag.String(FLAG_LOGGER, FLAG_LOGGER_HUMAN, "Sets type of logger. Available values are '"+FLAG_LOGGER_HUMAN+"', '"+FLAG_LOGGER_MACHINE+"'")
	versionFlag = flag.Bool(FLAG_VERSION, false, "prints version and exits")
	vidPidFlag = flag.String(FLAG_VID_PID, "", "specify to use vid/pid specific build properties, as defined in boards.txt")
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

	context := make(map[string]interface{})

	buildOptions := make(map[string]string)
	if *buildOptionsFileFlag != "" {
		if _, err := os.Stat(*buildOptionsFileFlag); err == nil {
			data, err := ioutil.ReadFile(*buildOptionsFileFlag)
			if err != nil {
				printCompleteError(err)
				defer os.Exit(1)
				return
			}
			err = json.Unmarshal(data, &buildOptions)
			if err != nil {
				printCompleteError(err)
				defer os.Exit(1)
				return
			}
		}
	}

	var err error
	printStackTrace := false
	err, printStackTrace = setContextSliceKeyOrLoadItFromOptions(context, hardwareFoldersFlag, buildOptions, constants.CTX_HARDWARE_FOLDERS, FLAG_HARDWARE, true)
	if err != nil {
		printError(err, printStackTrace)
		defer os.Exit(1)
		return
	}

	err, printStackTrace = setContextSliceKeyOrLoadItFromOptions(context, toolsFoldersFlag, buildOptions, constants.CTX_TOOLS_FOLDERS, FLAG_TOOLS, true)
	if err != nil {
		printError(err, printStackTrace)
		defer os.Exit(1)
		return
	}

	err, printStackTrace = setContextSliceKeyOrLoadItFromOptions(context, librariesFoldersFlag, buildOptions, constants.CTX_OTHER_LIBRARIES_FOLDERS, FLAG_LIBRARIES, false)
	if err != nil {
		printError(err, printStackTrace)
		defer os.Exit(1)
		return
	}

	err, printStackTrace = setContextSliceKeyOrLoadItFromOptions(context, librariesBuiltInFoldersFlag, buildOptions, constants.CTX_BUILT_IN_LIBRARIES_FOLDERS, FLAG_BUILT_IN_LIBRARIES, false)
	if err != nil {
		printError(err, printStackTrace)
		defer os.Exit(1)
		return
	}

	err, printStackTrace = setContextSliceKeyOrLoadItFromOptions(context, customBuildPropertiesFlag, buildOptions, constants.CTX_CUSTOM_BUILD_PROPERTIES, FLAG_PREFS, false)
	if err != nil {
		printError(err, printStackTrace)
		defer os.Exit(1)
		return
	}

	fqbn, err := gohasissues.Unquote(*fqbnFlag)
	if err != nil {
		printCompleteError(err)
		defer os.Exit(1)
		return
	}

	if fqbn == "" {
		fqbn = buildOptions[constants.CTX_FQBN]
	}

	if fqbn == "" {
		printErrorMessageAndFlagUsage(errors.New("Parameter '" + FLAG_FQBN + "' is mandatory"))
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
			fmt.Fprintln(os.Stderr, err)
			defer os.Exit(1)
			return
		}

		err = utils.EnsureFolderExists(buildPath)
		if err != nil {
			printCompleteError(err)
			defer os.Exit(1)
			return
		}
	}
	context[constants.CTX_BUILD_PATH] = buildPath

	if *vidPidFlag != "" {
		context[constants.CTX_VIDPID] = *vidPidFlag
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

	if *verboseFlag && *quietFlag {
		*verboseFlag = false
		*quietFlag = false
	}

	context[constants.CTX_VERBOSE] = *verboseFlag

	coreAPIVersion := ""
	if utils.MapStringStringHas(buildOptions, constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION) {
		coreAPIVersion = buildOptions[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION]
	} else {
		// if deprecated 'ideVersionFlag' has been used...
		if *coreAPIVersionFlag == "10600" && *ideVersionFlag != "10600" {
			coreAPIVersion = *ideVersionFlag
		} else {
			coreAPIVersion = *coreAPIVersionFlag
		}
	}
	context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION] = coreAPIVersion

	if *warningsLevelFlag != "" {
		context[constants.CTX_WARNINGS_LEVEL] = *warningsLevelFlag
	}

	if *debugLevelFlag > -1 {
		context[constants.CTX_DEBUG_LEVEL] = *debugLevelFlag
	}

	if *quietFlag {
		context[constants.CTX_LOGGER] = i18n.NoopLogger{}
	} else if *loggerFlag == FLAG_LOGGER_MACHINE {
		context[constants.CTX_LOGGER] = i18n.MachineLogger{}
	} else {
		context[constants.CTX_LOGGER] = i18n.HumanLogger{}
	}

	if *dumpPrefsFlag {
		err = builder.RunParseHardwareAndDumpBuildProperties(context)
	} else if *preprocessFlag {
		err = builder.RunPreprocess(context)
	} else {
		if flag.NArg() == 0 {
			fmt.Fprintln(os.Stderr, "Last parameter must be the sketch to compile")
			flag.Usage()
			defer os.Exit(1)
			return
		}
		err = builder.RunBuilder(context)
	}

	exitCode := 0
	if err != nil {
		err = utils.WrapError(err)

		fmt.Fprintln(os.Stderr, err)

		if utils.DebugLevel(context) >= 10 {
			fmt.Fprintln(os.Stderr, err.(*errors.Error).ErrorStack())
		}

		exitCode = toExitCode(err)
	}

	defer os.Exit(exitCode)
}

func setContextSliceKeyOrLoadItFromOptions(context map[string]interface{}, cliFlag slice, buildOptions map[string]string, contextKey string, paramName string, mandatory bool) (error, bool) {
	values, err := toSliceOfUnquoted(cliFlag)
	if err != nil {
		return err, true
	}

	if len(values) == 0 && len(buildOptions[contextKey]) > 0 {
		values = strings.Split(buildOptions[contextKey], ",")
	}

	if mandatory && len(values) == 0 {
		return errors.New("Parameter '" + paramName + "' is mandatory"), false
	}

	context[contextKey] = values

	return nil, false
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

func printError(err error, printStackTrace bool) {
	if printStackTrace {
		printCompleteError(err)
	} else {
		printErrorMessageAndFlagUsage(err)
	}
}

func printCompleteError(err error) {
	err = utils.WrapError(err)
	fmt.Fprintln(os.Stderr, err.(*errors.Error).ErrorStack())
}

func printErrorMessageAndFlagUsage(err error) {
	fmt.Fprintln(os.Stderr, err)
	flag.Usage()
}
