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
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"strconv"
	"strings"
	"syscall"

	"github.com/arduino/arduino-builder/grpc"
	"github.com/arduino/arduino-cli/arduino/cores"
	"github.com/arduino/arduino-cli/legacy/builder"
	"github.com/arduino/arduino-cli/legacy/builder/i18n"
	"github.com/arduino/arduino-cli/legacy/builder/types"
	paths "github.com/arduino/go-paths-helper"
	properties "github.com/arduino/go-properties-orderedmap"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const VERSION = "1.6.1"

type foldersFlag []string

func (h *foldersFlag) String() string {
	return fmt.Sprint(*h)
}

func (h *foldersFlag) Set(folder string) error {
	*h = append(*h, folder)
	return nil
}

type propertiesFlag []string

func (h *propertiesFlag) String() string {
	return fmt.Sprint(*h)
}

func (h *propertiesFlag) Set(value string) error {
	value = strings.TrimSpace(value)
	*h = append(*h, value)

	return nil
}

func main() {
	var hardwareFoldersFlag foldersFlag
	var toolsFoldersFlag foldersFlag
	var librariesBuiltInFoldersFlag foldersFlag
	var librariesFoldersFlag foldersFlag
	var customBuildPropertiesFlag propertiesFlag

	preprocessFlag := flag.Bool("preprocess", false, "preprocess the given sketch")
	dumpPrefsFlag := flag.Bool("dump-prefs", false, "dumps build properties used when compiling")
	codeCompleteAtFlag := flag.String("code-complete-at", "", "output code completions for sketch at a specific location. Location format is \"file:line:col\"")
	buildOptionsFileFlag := flag.String("build-options-file", "", "Instead of specifying --hardware, --tools etc every time, you can load all such options from a file")
	flag.Var(&hardwareFoldersFlag, "hardware", "Specify a 'hardware' folder. Can be added multiple times for specifying multiple 'hardware' folders")
	flag.Var(&toolsFoldersFlag, "tools", "Specify a 'tools' folder. Can be added multiple times for specifying multiple 'tools' folders")
	flag.Var(&librariesBuiltInFoldersFlag, "built-in-libraries", "Specify a built-in 'libraries' folder. These are low priority libraries. Can be added multiple times for specifying multiple built-in 'libraries' folders")
	flag.Var(&librariesFoldersFlag, "libraries", "Specify a 'libraries' folder. Can be added multiple times for specifying multiple 'libraries' folders")
	flag.Var(&customBuildPropertiesFlag, "prefs", "Specify a custom preference. Can be added multiple times for specifying multiple custom preferences")
	fqbnFlag := flag.String("fqbn", "", "fully qualified board name")
	coreAPIVersionFlag := flag.String("core-api-version", "10600", "version of core APIs (used to populate ARDUINO #define)")
	ideVersionFlag := flag.String("ide-version", "10600", "[deprecated] use 'core-api-version' instead")
	buildPathFlag := flag.String("build-path", "", "build path")
	buildCachePathFlag := flag.String("build-cache", "", "builds of 'core.a' are saved into this folder to be cached and reused")
	verboseFlag := flag.Bool("verbose", false, "if 'true' prints lots of stuff")
	quietFlag := flag.Bool("quiet", false, "if 'true' doesn't print any warnings or progress or whatever")
	debugLevelFlag := flag.Int("debug-level", builder.DEFAULT_DEBUG_LEVEL, "Turns on debugging messages. The higher, the chattier")
	warningsLevelFlag := flag.String("warnings", "", "Sets warnings level. Available values are 'none', 'default', 'more' and 'all'")
	loggerFlag := flag.String("logger", "human", "Sets type of logger. Available values are 'human', 'humantags', 'machine'")
	versionFlag := flag.Bool("version", false, "prints version and exits")
	daemonFlag := flag.Bool("daemon", false, "daemonizes and serves its functions via rpc")
	vidPidFlag := flag.String("vid-pid", "", "specify to use vid/pid specific build properties, as defined in boards.txt")
	jobsFlag := flag.Int("jobs", 0, "specify how many concurrent gcc processes should run at the same time. Defaults to the number of available cores on the running machine")
	traceFlag := flag.Bool("trace", false, "traces the whole process lifecycle")
	experimentalFeatures := flag.Bool("experimental", false, "enables experimental features")
	// Not used anymore, kept only because the Arduino IDE still provides this flag
	flag.Bool("compile", true, "[deprecated] this is now the default action")

	flag.Parse()

	if *traceFlag {
		f, err := os.Create("trace.out")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		f2, err := os.Create("profile.out")
		if err != nil {
			panic(err)
		}
		defer f2.Close()

		pprof.StartCPUProfile(f2)
		defer pprof.StopCPUProfile()

		err = trace.Start(f)
		if err != nil {
			panic(err)
		}
		defer trace.Stop()
	}

	if *versionFlag {
		fmt.Println("Arduino Builder " + VERSION)
		fmt.Println("Copyright (C) 2015 Arduino LLC and contributors")
		fmt.Println("See https://www.arduino.cc/ and https://github.com/arduino/arduino-builder/graphs/contributors")
		fmt.Println("This is free software; see the source for copying conditions.  There is NO")
		fmt.Println("warranty; not even for MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.")
		return
	}

	if *jobsFlag > 0 {
		runtime.GOMAXPROCS(*jobsFlag)
	} else {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	ctx := &types.Context{}
	ctx.IgnoreSketchFolderNameErrors = true

	// place here all experimental features that should live under this flag
	if *experimentalFeatures {
		ctx.UseArduinoPreprocessor = true
	}

	if *daemonFlag {
		ctx.SetLogger(i18n.NoopLogger{})
		grpc.RegisterAndServeJsonRPC(ctx)
	}

	if *buildOptionsFileFlag != "" {
		buildOptions := properties.NewMap()
		if _, err := os.Stat(*buildOptionsFileFlag); err == nil {
			data, err := ioutil.ReadFile(*buildOptionsFileFlag)
			if err != nil {
				printCompleteError(err)
			}
			err = json.Unmarshal(data, &buildOptions)
			if err != nil {
				printCompleteError(err)
			}
		}
		ctx.InjectBuildOptions(buildOptions)
	}

	// FLAG_HARDWARE
	if hardwareFolders, err := toSliceOfUnquoted(hardwareFoldersFlag); err != nil {
		printCompleteError(err)
	} else if len(hardwareFolders) > 0 {
		ctx.HardwareDirs = paths.NewPathList(hardwareFolders...)
	}
	if len(ctx.HardwareDirs) == 0 {
		printErrorMessageAndFlagUsage(errors.New("Parameter 'hardware' is mandatory"))
	}

	// FLAG_TOOLS
	if toolsFolders, err := toSliceOfUnquoted(toolsFoldersFlag); err != nil {
		printCompleteError(err)
	} else if len(toolsFolders) > 0 {
		ctx.BuiltInToolsDirs = paths.NewPathList(toolsFolders...)
	}
	if len(ctx.BuiltInToolsDirs) == 0 {
		printErrorMessageAndFlagUsage(errors.New("Parameter 'tools' is mandatory"))
	}

	// FLAG_LIBRARIES
	if librariesFolders, err := toSliceOfUnquoted(librariesFoldersFlag); err != nil {
		printCompleteError(err)
	} else if len(librariesFolders) > 0 {
		ctx.OtherLibrariesDirs = paths.NewPathList(librariesFolders...)
	}

	// FLAG_BUILT_IN_LIBRARIES
	if librariesBuiltInFolders, err := toSliceOfUnquoted(librariesBuiltInFoldersFlag); err != nil {
		printCompleteError(err)
	} else if len(librariesBuiltInFolders) > 0 {
		ctx.BuiltInLibrariesDirs = paths.NewPathList(librariesBuiltInFolders...)
	}

	// FLAG_PREFS
	if customBuildProperties, err := toSliceOfUnquoted(customBuildPropertiesFlag); err != nil {
		printCompleteError(err)
	} else if len(customBuildProperties) > 0 {
		ctx.CustomBuildProperties = customBuildProperties
	}

	// FLAG_FQBN
	if fqbnIn, err := unquote(*fqbnFlag); err != nil {
		printCompleteError(err)
	} else if fqbnIn != "" {
		if fqbn, err := cores.ParseFQBN(fqbnIn); err != nil {
			printCompleteError(err)
		} else {
			ctx.FQBN = fqbn
		}
	}
	if ctx.FQBN == nil {
		printErrorMessageAndFlagUsage(errors.New("Parameter 'fqbn' is mandatory"))
	}

	// FLAG_BUILD_PATH
	if *buildPathFlag != "" {
		buildPathUnquoted, err := unquote(*buildPathFlag)
		if err != nil {
			printCompleteError(err)
		}
		buildPath := paths.New(buildPathUnquoted)

		if _, err := buildPath.Stat(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err := buildPath.MkdirAll(); err != nil {
			printCompleteError(err)
		}
		ctx.BuildPath, _ = buildPath.Abs()
	}

	// FLAG_BUILD_CACHE
	if *buildCachePathFlag != "" {
		buildCachePathUnquoted, err := unquote(*buildCachePathFlag)
		if err != nil {
			printCompleteError(err)
		}
		buildCachePath := paths.New(buildCachePathUnquoted)
		if buildCachePath != nil {
			if err := buildCachePath.MkdirAll(); err != nil {
				printCompleteError(err)
			}
		}
		ctx.BuildCachePath = buildCachePath
	}

	// FLAG_VID_PID
	if *vidPidFlag != "" {
		ctx.USBVidPid = *vidPidFlag
	}

	if flag.NArg() > 0 {
		sketchLocationUnquoted, err := unquote(flag.Arg(0))
		if err != nil {
			printCompleteError(err)
		}
		ctx.SketchLocation = paths.New(sketchLocationUnquoted)
	}

	if *verboseFlag && *quietFlag {
		*verboseFlag = false
		*quietFlag = false
	}

	ctx.Verbose = *verboseFlag

	// FLAG_IDE_VERSION
	if ctx.ArduinoAPIVersion == "" {
		// if deprecated "--ideVersionFlag" has been used...
		if *coreAPIVersionFlag == "10600" && *ideVersionFlag != "10600" {
			ctx.ArduinoAPIVersion = *ideVersionFlag
		} else {
			ctx.ArduinoAPIVersion = *coreAPIVersionFlag
		}
	}

	if *warningsLevelFlag != "" {
		ctx.WarningsLevel = *warningsLevelFlag
	}

	if *debugLevelFlag > -1 {
		ctx.DebugLevel = *debugLevelFlag
	}
	if ctx.DebugLevel < 10 {
		logrus.SetOutput(ioutil.Discard)
	}

	if *quietFlag {
		ctx.SetLogger(i18n.NoopLogger{})
	} else if *loggerFlag == "machine" {
		ctx.SetLogger(i18n.MachineLogger{})
		ctx.Progress.PrintEnabled = true
	} else if *loggerFlag == "humantags" {
		ctx.SetLogger(i18n.HumanTagsLogger{})
	} else {
		ctx.SetLogger(i18n.HumanLogger{})
	}

	var err error
	if *dumpPrefsFlag {
		err = builder.RunParseHardwareAndDumpBuildProperties(ctx)
	} else if *preprocessFlag || *codeCompleteAtFlag != "" {
		ctx.CodeCompleteAt = *codeCompleteAtFlag
		err = builder.RunPreprocess(ctx)
	} else {
		if flag.NArg() == 0 {
			fmt.Fprintln(os.Stderr, "Last parameter must be the sketch to compile")
			flag.Usage()
			os.Exit(1)
		}
		err = builder.RunBuilder(ctx)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)

		if ctx.DebugLevel >= 10 {
			err = errors.WithStack(err)
			fmt.Fprintf(os.Stderr, "%+v\n", err)
		}
		os.Exit(toExitCode(err))
	}
}

func toExitCode(err error) int {
	if exiterr, ok := errors.Cause(err).(*exec.ExitError); ok {
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus()
		}
	}
	return 1
}

func toSliceOfUnquoted(value []string) ([]string, error) {
	var values []string
	for _, v := range value {
		v, err := unquote(v)
		if err != nil {
			return nil, err
		}
		values = append(values, v)
	}
	return values, nil
}

func unquote(s string) (string, error) {
	if stringStartsEndsWith(s, "'") {
		s = s[1 : len(s)-1]
	}

	if !stringStartsEndsWith(s, "\"") {
		return s, nil
	}

	return strconv.Unquote(s)
}

func stringStartsEndsWith(s string, c string) bool {
	return strings.HasPrefix(s, c) && strings.HasSuffix(s, c)
}

func printCompleteError(err error) {
	err = errors.WithStack(err)
	fmt.Fprintf(os.Stderr, "%+v\n", err)
	os.Exit(1)
}

func printErrorMessageAndFlagUsage(err error) {
	fmt.Fprintln(os.Stderr, err)
	flag.Usage()
	os.Exit(1)
}
