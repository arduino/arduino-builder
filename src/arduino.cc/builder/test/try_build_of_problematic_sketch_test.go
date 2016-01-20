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

package test

import (
	"arduino.cc/builder"
	"arduino.cc/builder/constants"
	"os"
	"path/filepath"
	"testing"
)

func TestTryBuild001(t *testing.T) {
	tryBuild(t, "sketch_with_inline_function", "sketch.ino")
}

func TestTryBuild002(t *testing.T) {
	tryBuild(t, "sketch_with_function_signature_inside_ifdef", "sketch.ino")
}

func TestTryBuild003(t *testing.T) {
	tryPreprocess(t, "sketch_no_functions", "main.ino")
}

func TestTryBuild004(t *testing.T) {
	tryBuild(t, "sketch_with_const", "sketch.ino")
}

func TestTryBuild005(t *testing.T) {
	tryBuild(t, "sketch_with_old_lib", "sketch.ino")
}

func TestTryBuild006(t *testing.T) {
	tryBuild(t, "sketch_with_macosx_garbage", "sketch.ino")
}

func TestTryBuild007(t *testing.T) {
	tryBuild(t, "sketch_with_config", "sketch_with_config.ino")
}

// XXX: Failing sketch, typename not supported
//func TestTryBuild008(t *testing.T) {
//	tryBuild(t, "sketch_with_typename", "sketch.ino")
//}

func TestTryBuild009(t *testing.T) {
	tryBuild(t, "sketch_with_usbcon", "sketch.ino")
}

func TestTryBuild010(t *testing.T) {
	tryBuild(t, "sketch_with_namespace", "sketch.ino")
}

func TestTryBuild011(t *testing.T) {
	tryBuild(t, "sketch_with_inline_function", "sketch.ino")
}

func TestTryBuild012(t *testing.T) {
	tryBuild(t, "sketch_with_default_args", "sketch.ino")
}

func TestTryBuild013(t *testing.T) {
	tryBuild(t, "sketch_with_class", "sketch.ino")
}

func TestTryBuild014(t *testing.T) {
	tryBuild(t, "sketch_with_backup_files", "sketch.ino")
}

func TestTryBuild015(t *testing.T) {
	tryBuild(t, "sketch_with_subfolders")
}

// This is a sketch that fails to build on purpose
//func TestTryBuild016(t *testing.T) {
//	tryBuild(t, "sketch_that_checks_if_SPI_has_transactions_and_includes_missing_Ethernet", "sketch.ino")
//}

func TestTryBuild017(t *testing.T) {
	tryPreprocess(t, "sketch_no_functions_two_files", "main.ino")
}

func TestTryBuild018(t *testing.T) {
	tryBuild(t, "sketch_that_checks_if_SPI_has_transactions", "sketch.ino")
}

func TestTryBuild019(t *testing.T) {
	tryBuild(t, "sketch_with_ifdef", "sketch.ino")
}

func TestTryBuild020(t *testing.T) {
	context := makeDefaultContext(t)
	context[constants.CTX_OTHER_LIBRARIES_FOLDERS] = []string{"dependent_libraries", "libraries"}
	tryPreprocessWithContext(t, context, "sketch_with_dependend_libraries", "sketch.ino")
}

func TestTryBuild021(t *testing.T) {
	tryBuild(t, "sketch_with_function_pointer", "sketch.ino")
}

func TestTryBuild022(t *testing.T) {
	context := makeDefaultContext(t)
	context[constants.CTX_FQBN] = "arduino:samd:arduino_zero_native"
	tryBuildWithContext(t, context, "sketch_usbhost", "sketch_usbhost.ino")
}

func TestTryBuild023(t *testing.T) {
	tryBuild(t, "sketch1", "sketch.ino")
}

func TestTryBuild024(t *testing.T) {
	tryBuild(t, "sketch2", "SketchWithIfDef.ino")
}

// The library for this sketch is missing
//func TestTryBuild025(t *testing.T) {
//	tryBuild(t, "sketch3", "Baladuino.ino")
//}

func TestTryBuild026(t *testing.T) {
	tryBuild(t, "sketch4", "CharWithEscapedDoubleQuote.ino")
}

func TestTryBuild027(t *testing.T) {
	tryBuild(t, "sketch5", "IncludeBetweenMultilineComment.ino")
}

func TestTryBuild028(t *testing.T) {
	tryBuild(t, "sketch6", "LineContinuations.ino")
}

func TestTryBuild029(t *testing.T) {
	tryBuild(t, "sketch7", "StringWithComment.ino")
}

func TestTryBuild030(t *testing.T) {
	tryBuild(t, "sketch8", "SketchWithStruct.ino")
}

func TestTryBuild031(t *testing.T) {
	tryBuild(t, "sketch9", "sketch.ino")
}

func TestTryBuild032(t *testing.T) {
	tryBuild(t, "sketch10", "sketch.ino")
}

func TestTryBuild033(t *testing.T) {
	tryBuild(t, "sketch_that_includes_arduino_h", "sketch_that_includes_arduino_h.ino")
}

func TestTryBuild034(t *testing.T) {
	tryBuild(t, "sketch_with_static_asserts", "sketch_with_static_asserts.ino")
}

func makeDefaultContext(t *testing.T) map[string]interface{} {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})
	buildPath := SetupBuildPath(t, context)
	defer os.RemoveAll(buildPath)

	context[constants.CTX_HARDWARE_FOLDERS] = []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware", "downloaded_board_manager_stuff"}
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools", "downloaded_board_manager_stuff"}
	context[constants.CTX_BUILT_IN_LIBRARIES_FOLDERS] = []string{"downloaded_libraries"}
	context[constants.CTX_FQBN] = "arduino:avr:leonardo"
	context[constants.CTX_BUILD_PROPERTIES_RUNTIME_IDE_VERSION] = "10607"
	context[constants.CTX_OTHER_LIBRARIES_FOLDERS] = []string{"libraries"}
	context[constants.CTX_VERBOSE] = true
	context[constants.CTX_DEBUG_PREPROCESSOR] = true

	return context
}

func tryBuild(t *testing.T, sketchPath ...string) {
	context := makeDefaultContext(t)
	tryBuildWithContext(t, context, sketchPath...)
}

func tryBuildWithContext(t *testing.T, context map[string]interface{}, sketchPath ...string) {
	sketchLocation := filepath.Join(sketchPath...)
	context[constants.CTX_SKETCH_LOCATION] = sketchLocation

	err := builder.RunBuilder(context)
	NoError(t, err, "Build error for "+sketchLocation)
}

func tryPreprocess(t *testing.T, sketchPath ...string) {
	context := makeDefaultContext(t)
	tryPreprocessWithContext(t, context, sketchPath...)
}

func tryPreprocessWithContext(t *testing.T, context map[string]interface{}, sketchPath ...string) {
	sketchLocation := filepath.Join(sketchPath...)
	context[constants.CTX_SKETCH_LOCATION] = sketchLocation

	err := builder.RunPreprocess(context)
	NoError(t, err, "Build error for "+sketchLocation)
}
