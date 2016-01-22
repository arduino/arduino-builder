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

package test

import (
	"arduino.cc/arduino-builder/Godeps/_workspace/src/github.com/stretchr/testify/require"
	"arduino.cc/arduino-builder/builder/constants"
	"arduino.cc/arduino-builder/builder/ctags"
	"arduino.cc/arduino-builder/builder/types"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestCTagsParserShouldListPrototypes(t *testing.T) {
	context := make(map[string]interface{})

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserShouldListPrototypes.txt"))
	NoError(t, err)

	context[constants.CTX_CTAGS_OUTPUT] = string(bytes)

	ctagsParser := ctags.CTagsParser{}
	ctagsParser.Run(context)

	tags := context[constants.CTX_CTAGS_OF_PREPROC_SOURCE].([]*types.CTag)

	require.Equal(t, 8, len(tags))
	idx := 0
	require.Equal(t, "server", tags[idx].FunctionName)
	require.Equal(t, "variable", tags[idx].Kind)
	require.Equal(t, "/tmp/sketch7210316334309249705.cpp", tags[idx].Filename)
	idx++
	require.Equal(t, "setup", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
	require.Equal(t, "/tmp/sketch7210316334309249705.cpp", tags[idx].Filename)
	idx++
	require.Equal(t, "loop", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
	require.Equal(t, "/tmp/sketch7210316334309249705.cpp", tags[idx].Filename)
	idx++
	require.Equal(t, "process", tags[idx].FunctionName)
	require.Equal(t, "prototype", tags[idx].Kind)
	require.Equal(t, "/tmp/sketch7210316334309249705.cpp", tags[idx].Filename)
	idx++
	require.Equal(t, "process", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
	require.Equal(t, "/tmp/sketch7210316334309249705.cpp", tags[idx].Filename)
	idx++
	require.Equal(t, "digitalCommand", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
	require.Equal(t, "/tmp/sketch7210316334309249705.cpp", tags[idx].Filename)
	idx++
	require.Equal(t, "analogCommand", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
	require.Equal(t, "/tmp/sketch7210316334309249705.cpp", tags[idx].Filename)
	idx++
	require.Equal(t, "modeCommand", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
	require.Equal(t, "/tmp/sketch7210316334309249705.cpp", tags[idx].Filename)
}

func TestCTagsParserShouldListTemplates(t *testing.T) {
	context := make(map[string]interface{})

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserShouldListTemplates.txt"))
	NoError(t, err)

	context[constants.CTX_CTAGS_OUTPUT] = string(bytes)

	ctagsParser := ctags.CTagsParser{}
	ctagsParser.Run(context)

	tags := context[constants.CTX_CTAGS_OF_PREPROC_SOURCE].([]*types.CTag)

	require.Equal(t, 3, len(tags))
	idx := 0
	require.Equal(t, "minimum", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
	require.Equal(t, "template <typename T> T minimum (T a, T b);", tags[idx].Prototype)
	idx++
	require.Equal(t, "setup", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
	require.Equal(t, "void setup();", tags[idx].Prototype)
	idx++
	require.Equal(t, "loop", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
	require.Equal(t, "void loop();", tags[idx].Prototype)
}

func TestCTagsParserShouldListTemplates2(t *testing.T) {
	context := make(map[string]interface{})

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserShouldListTemplates2.txt"))
	NoError(t, err)

	context[constants.CTX_CTAGS_OUTPUT] = string(bytes)

	ctagsParser := ctags.CTagsParser{}
	ctagsParser.Run(context)

	tags := context[constants.CTX_CTAGS_OF_PREPROC_SOURCE].([]*types.CTag)

	require.Equal(t, 4, len(tags))
	idx := 0
	require.Equal(t, "setup", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
	idx++
	require.Equal(t, "loop", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
	idx++
	require.Equal(t, "SRAM_writeAnything", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
	require.Equal(t, "template <class T> int SRAM_writeAnything(int ee, const T& value);", tags[idx].Prototype)
	idx++
	require.Equal(t, "SRAM_readAnything", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
	require.Equal(t, "template <class T> int SRAM_readAnything(int ee, T& value);", tags[idx].Prototype)
}

func TestCTagsParserShouldDealWithClasses(t *testing.T) {
	context := make(map[string]interface{})

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserShouldDealWithClasses.txt"))
	NoError(t, err)

	context[constants.CTX_CTAGS_OUTPUT] = string(bytes)

	ctagsParser := ctags.CTagsParser{}
	ctagsParser.Run(context)

	tags := context[constants.CTX_CTAGS_OF_PREPROC_SOURCE].([]*types.CTag)

	require.Equal(t, 2, len(tags))
	idx := 0
	require.Equal(t, "SleepCycle", tags[idx].FunctionName)
	require.Equal(t, "prototype", tags[idx].Kind)
	idx++
	require.Equal(t, "SleepCycle", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
}

func TestCTagsParserShouldDealWithStructs(t *testing.T) {
	context := make(map[string]interface{})

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserShouldDealWithStructs.txt"))
	NoError(t, err)

	context[constants.CTX_CTAGS_OUTPUT] = string(bytes)

	ctagsParser := ctags.CTagsParser{}
	ctagsParser.Run(context)

	tags := context[constants.CTX_CTAGS_OF_PREPROC_SOURCE].([]*types.CTag)

	require.Equal(t, 5, len(tags))
	idx := 0
	require.Equal(t, "A_NEW_TYPE", tags[idx].FunctionName)
	require.Equal(t, "struct", tags[idx].Kind)
	idx++
	require.Equal(t, "foo", tags[idx].FunctionName)
	require.Equal(t, "variable", tags[idx].Kind)
	require.Equal(t, "struct:A_NEW_TYPE", tags[idx].Typeref)
	idx++
	require.Equal(t, "setup", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
	idx++
	require.Equal(t, "loop", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
	idx++
	require.Equal(t, "dostuff", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
}

func TestCTagsParserShouldDealWithMacros(t *testing.T) {
	context := make(map[string]interface{})

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserShouldDealWithMacros.txt"))
	NoError(t, err)

	context[constants.CTX_CTAGS_OUTPUT] = string(bytes)

	ctagsParser := ctags.CTagsParser{}
	ctagsParser.Run(context)

	tags := context[constants.CTX_CTAGS_OF_PREPROC_SOURCE].([]*types.CTag)

	require.Equal(t, 8, len(tags))
	idx := 0
	require.Equal(t, "DEBUG", tags[idx].FunctionName)
	require.Equal(t, "macro", tags[idx].Kind)
	idx++
	require.Equal(t, "DISABLED", tags[idx].FunctionName)
	require.Equal(t, "macro", tags[idx].Kind)
	idx++
	require.Equal(t, "hello", tags[idx].FunctionName)
	require.Equal(t, "variable", tags[idx].Kind)
	idx++
	require.Equal(t, "setup", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
	idx++
	require.Equal(t, "loop", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
	idx++
	require.Equal(t, "debug", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
	idx++
	require.Equal(t, "disabledIsDefined", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
	idx++
	require.Equal(t, "useMyType", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
}

func TestCTagsParserShouldDealFunctionWithDifferentSignatures(t *testing.T) {
	context := make(map[string]interface{})

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserShouldDealFunctionWithDifferentSignatures.txt"))
	NoError(t, err)

	context[constants.CTX_CTAGS_OUTPUT] = string(bytes)

	ctagsParser := ctags.CTagsParser{}
	ctagsParser.Run(context)

	tags := context[constants.CTX_CTAGS_OF_PREPROC_SOURCE].([]*types.CTag)

	require.Equal(t, 3, len(tags))
	idx := 0
	require.Equal(t, "getBytes", tags[idx].FunctionName)
	require.Equal(t, "prototype", tags[idx].Kind)
	idx++
	require.Equal(t, "getBytes", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
	idx++
	require.Equal(t, "getBytes", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
}

func TestCTagsParserClassMembersAreFilteredOut(t *testing.T) {
	context := make(map[string]interface{})

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserClassMembersAreFilteredOut.txt"))
	NoError(t, err)

	context[constants.CTX_CTAGS_OUTPUT] = string(bytes)

	ctagsParser := ctags.CTagsParser{}
	ctagsParser.Run(context)

	tags := context[constants.CTX_CTAGS_OF_PREPROC_SOURCE].([]*types.CTag)

	require.Equal(t, 5, len(tags))
	idx := 0
	require.Equal(t, "set_values", tags[idx].FunctionName)
	require.Equal(t, "prototype", tags[idx].Kind)
	require.Equal(t, "Rectangle", tags[idx].Class)
	idx++
	require.Equal(t, "area", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
	require.Equal(t, "Rectangle", tags[idx].Class)
	idx++
	require.Equal(t, "set_values", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
	require.Equal(t, "Rectangle", tags[idx].Class)
	idx++
	require.Equal(t, "setup", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
	idx++
	require.Equal(t, "loop", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
}

func TestCTagsParserStructWithFunctions(t *testing.T) {
	context := make(map[string]interface{})

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserStructWithFunctions.txt"))
	NoError(t, err)

	context[constants.CTX_CTAGS_OUTPUT] = string(bytes)

	ctagsParser := ctags.CTagsParser{}
	ctagsParser.Run(context)

	tags := context[constants.CTX_CTAGS_OF_PREPROC_SOURCE].([]*types.CTag)

	require.Equal(t, 8, len(tags))
	idx := 0
	require.Equal(t, "sensorData", tags[idx].FunctionName)
	require.Equal(t, "struct", tags[idx].Kind)
	idx++
	require.Equal(t, "sensorData", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
	require.Equal(t, "sensorData", tags[idx].Struct)
	idx++
	require.Equal(t, "sensorData", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
	require.Equal(t, "sensorData", tags[idx].Struct)
	idx++
	require.Equal(t, "sensors", tags[idx].FunctionName)
	require.Equal(t, "variable", tags[idx].Kind)
	idx++
	require.Equal(t, "sensor1", tags[idx].FunctionName)
	require.Equal(t, "variable", tags[idx].Kind)
	idx++
	require.Equal(t, "sensor2", tags[idx].FunctionName)
	require.Equal(t, "variable", tags[idx].Kind)
	idx++
	require.Equal(t, "setup", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
	idx++
	require.Equal(t, "loop", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
}

func TestCTagsParserDefaultArguments(t *testing.T) {
	context := make(map[string]interface{})

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserDefaultArguments.txt"))
	NoError(t, err)

	context[constants.CTX_CTAGS_OUTPUT] = string(bytes)

	ctagsParser := ctags.CTagsParser{}
	ctagsParser.Run(context)

	tags := context[constants.CTX_CTAGS_OF_PREPROC_SOURCE].([]*types.CTag)

	require.Equal(t, 3, len(tags))
	idx := 0
	require.Equal(t, "test", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
	require.Equal(t, "void test(int x = 1);", tags[idx].Prototype)
	idx++
	require.Equal(t, "setup", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
	idx++
	require.Equal(t, "loop", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
}

func TestCTagsParserNamespace(t *testing.T) {
	context := make(map[string]interface{})

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserNamespace.txt"))
	NoError(t, err)

	context[constants.CTX_CTAGS_OUTPUT] = string(bytes)

	ctagsParser := ctags.CTagsParser{}
	ctagsParser.Run(context)

	tags := context[constants.CTX_CTAGS_OF_PREPROC_SOURCE].([]*types.CTag)

	require.Equal(t, 3, len(tags))
	idx := 0
	require.Equal(t, "value", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
	require.Equal(t, "Test", tags[idx].Namespace)
	idx++
	require.Equal(t, "setup", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
	idx++
	require.Equal(t, "loop", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
}

func TestCTagsParserStatic(t *testing.T) {
	context := make(map[string]interface{})

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserStatic.txt"))
	NoError(t, err)

	context[constants.CTX_CTAGS_OUTPUT] = string(bytes)

	ctagsParser := ctags.CTagsParser{}
	ctagsParser.Run(context)

	tags := context[constants.CTX_CTAGS_OF_PREPROC_SOURCE].([]*types.CTag)

	require.Equal(t, 3, len(tags))
	idx := 0
	require.Equal(t, "setup", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
	idx++
	require.Equal(t, "loop", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
	idx++
	require.Equal(t, "doStuff", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
}

func TestCTagsParserFunctionPointer(t *testing.T) {
	context := make(map[string]interface{})

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserFunctionPointer.txt"))
	NoError(t, err)

	context[constants.CTX_CTAGS_OUTPUT] = string(bytes)

	ctagsParser := ctags.CTagsParser{}
	ctagsParser.Run(context)

	tags := context[constants.CTX_CTAGS_OF_PREPROC_SOURCE].([]*types.CTag)

	require.Equal(t, 4, len(tags))
	idx := 0
	require.Equal(t, "t1Callback", tags[idx].FunctionName)
	require.Equal(t, "variable", tags[idx].Kind)
	idx++
	require.Equal(t, "t1Callback", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
	idx++
	require.Equal(t, "setup", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
	idx++
	require.Equal(t, "loop", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
}

func TestCTagsParserFunctionPointers(t *testing.T) {
	context := make(map[string]interface{})

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserFunctionPointers.txt"))
	NoError(t, err)

	context[constants.CTX_CTAGS_OUTPUT] = string(bytes)

	ctagsParser := ctags.CTagsParser{}
	ctagsParser.Run(context)

	tags := context[constants.CTX_CTAGS_OF_PREPROC_SOURCE].([]*types.CTag)

	require.Equal(t, 5, len(tags))
	idx := 0
	require.Equal(t, "setup", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
	idx++
	require.Equal(t, "loop", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
	idx++
	require.Equal(t, "func", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
	idx++
	require.Equal(t, "funcArr", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
	require.Equal(t, "int funcArr();", tags[idx].Prototype)
	idx++
	require.Equal(t, "funcCombo", tags[idx].FunctionName)
	require.Equal(t, "function", tags[idx].Kind)
	require.Equal(t, "void funcCombo(void (*(&in)[5])(int));", tags[idx].Prototype)
}
