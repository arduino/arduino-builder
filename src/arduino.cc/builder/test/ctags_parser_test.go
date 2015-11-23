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
	"arduino.cc/builder"
	"arduino.cc/builder/constants"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestCTagsParserShouldListPrototypes(t *testing.T) {
	context := make(map[string]interface{})

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserShouldListPrototypes.txt"))
	NoError(t, err)

	context[constants.CTX_CTAGS_OUTPUT] = string(bytes)

	ctagsParser := builder.CTagsParser{}
	ctagsParser.Run(context)

	ctags := context[constants.CTX_CTAGS_OF_PREPROC_SOURCE].([]*builder.CTag)

	require.Equal(t, 8, len(ctags))
	idx := 0
	require.Equal(t, "server", ctags[idx].FunctionName)
	require.Equal(t, "variable", ctags[idx].Kind)
	require.Equal(t, "/tmp/sketch7210316334309249705.cpp", ctags[idx].Filename)
	idx++
	require.Equal(t, "setup", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
	require.Equal(t, "/tmp/sketch7210316334309249705.cpp", ctags[idx].Filename)
	idx++
	require.Equal(t, "loop", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
	require.Equal(t, "/tmp/sketch7210316334309249705.cpp", ctags[idx].Filename)
	idx++
	require.Equal(t, "process", ctags[idx].FunctionName)
	require.Equal(t, "prototype", ctags[idx].Kind)
	require.Equal(t, "/tmp/sketch7210316334309249705.cpp", ctags[idx].Filename)
	idx++
	require.Equal(t, "process", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
	require.Equal(t, "/tmp/sketch7210316334309249705.cpp", ctags[idx].Filename)
	idx++
	require.Equal(t, "digitalCommand", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
	require.Equal(t, "/tmp/sketch7210316334309249705.cpp", ctags[idx].Filename)
	idx++
	require.Equal(t, "analogCommand", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
	require.Equal(t, "/tmp/sketch7210316334309249705.cpp", ctags[idx].Filename)
	idx++
	require.Equal(t, "modeCommand", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
	require.Equal(t, "/tmp/sketch7210316334309249705.cpp", ctags[idx].Filename)
}

func TestCTagsParserShouldListTemplates(t *testing.T) {
	context := make(map[string]interface{})

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserShouldListTemplates.txt"))
	NoError(t, err)

	context[constants.CTX_CTAGS_OUTPUT] = string(bytes)

	ctagsParser := builder.CTagsParser{}
	ctagsParser.Run(context)

	ctags := context[constants.CTX_CTAGS_OF_PREPROC_SOURCE].([]*builder.CTag)

	require.Equal(t, 3, len(ctags))
	idx := 0
	require.Equal(t, "minimum", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
	require.Equal(t, "(T a, T b)", ctags[idx].Signature)
	idx++
	require.Equal(t, "setup", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
	idx++
	require.Equal(t, "loop", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
}

func TestCTagsParserShouldListTemplates2(t *testing.T) {
	context := make(map[string]interface{})

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserShouldListTemplates2.txt"))
	NoError(t, err)

	context[constants.CTX_CTAGS_OUTPUT] = string(bytes)

	ctagsParser := builder.CTagsParser{}
	ctagsParser.Run(context)

	ctags := context[constants.CTX_CTAGS_OF_PREPROC_SOURCE].([]*builder.CTag)

	require.Equal(t, 4, len(ctags))
	idx := 0
	require.Equal(t, "setup", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
	idx++
	require.Equal(t, "loop", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
	idx++
	require.Equal(t, "SRAM_writeAnything", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
	require.Equal(t, "(int ee, const T& value)", ctags[idx].Signature)
	idx++
	require.Equal(t, "SRAM_readAnything", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
	require.Equal(t, "(int ee, T& value)", ctags[idx].Signature)
}

func TestCTagsParserShouldDealWithClasses(t *testing.T) {
	context := make(map[string]interface{})

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserShouldDealWithClasses.txt"))
	NoError(t, err)

	context[constants.CTX_CTAGS_OUTPUT] = string(bytes)

	ctagsParser := builder.CTagsParser{}
	ctagsParser.Run(context)

	ctags := context[constants.CTX_CTAGS_OF_PREPROC_SOURCE].([]*builder.CTag)

	require.Equal(t, 2, len(ctags))
	idx := 0
	require.Equal(t, "SleepCycle", ctags[idx].FunctionName)
	require.Equal(t, "prototype", ctags[idx].Kind)
	idx++
	require.Equal(t, "SleepCycle", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
}

func TestCTagsParserShouldDealWithStructs(t *testing.T) {
	context := make(map[string]interface{})

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserShouldDealWithStructs.txt"))
	NoError(t, err)

	context[constants.CTX_CTAGS_OUTPUT] = string(bytes)

	ctagsParser := builder.CTagsParser{}
	ctagsParser.Run(context)

	ctags := context[constants.CTX_CTAGS_OF_PREPROC_SOURCE].([]*builder.CTag)

	require.Equal(t, 5, len(ctags))
	idx := 0
	require.Equal(t, "A_NEW_TYPE", ctags[idx].FunctionName)
	require.Equal(t, "struct", ctags[idx].Kind)
	idx++
	require.Equal(t, "foo", ctags[idx].FunctionName)
	require.Equal(t, "variable", ctags[idx].Kind)
	require.Equal(t, "struct:A_NEW_TYPE", ctags[idx].Typeref)
	idx++
	require.Equal(t, "setup", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
	idx++
	require.Equal(t, "loop", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
	idx++
	require.Equal(t, "dostuff", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
}

func TestCTagsParserShouldDealWithMacros(t *testing.T) {
	context := make(map[string]interface{})

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserShouldDealWithMacros.txt"))
	NoError(t, err)

	context[constants.CTX_CTAGS_OUTPUT] = string(bytes)

	ctagsParser := builder.CTagsParser{}
	ctagsParser.Run(context)

	ctags := context[constants.CTX_CTAGS_OF_PREPROC_SOURCE].([]*builder.CTag)

	require.Equal(t, 8, len(ctags))
	idx := 0
	require.Equal(t, "DEBUG", ctags[idx].FunctionName)
	require.Equal(t, "macro", ctags[idx].Kind)
	idx++
	require.Equal(t, "DISABLED", ctags[idx].FunctionName)
	require.Equal(t, "macro", ctags[idx].Kind)
	idx++
	require.Equal(t, "hello", ctags[idx].FunctionName)
	require.Equal(t, "variable", ctags[idx].Kind)
	idx++
	require.Equal(t, "setup", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
	idx++
	require.Equal(t, "loop", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
	idx++
	require.Equal(t, "debug", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
	idx++
	require.Equal(t, "disabledIsDefined", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
	idx++
	require.Equal(t, "useMyType", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
}

func TestCTagsParserShouldDealFunctionWithDifferentSignatures(t *testing.T) {
	context := make(map[string]interface{})

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserShouldDealFunctionWithDifferentSignatures.txt"))
	NoError(t, err)

	context[constants.CTX_CTAGS_OUTPUT] = string(bytes)

	ctagsParser := builder.CTagsParser{}
	ctagsParser.Run(context)

	ctags := context[constants.CTX_CTAGS_OF_PREPROC_SOURCE].([]*builder.CTag)

	require.Equal(t, 3, len(ctags))
	idx := 0
	require.Equal(t, "getBytes", ctags[idx].FunctionName)
	require.Equal(t, "prototype", ctags[idx].Kind)
	idx++
	require.Equal(t, "getBytes", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
	idx++
	require.Equal(t, "getBytes", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
}

func TestCTagsParserClassMembersAreFilteredOut(t *testing.T) {
	context := make(map[string]interface{})

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserClassMembersAreFilteredOut.txt"))
	NoError(t, err)

	context[constants.CTX_CTAGS_OUTPUT] = string(bytes)

	ctagsParser := builder.CTagsParser{}
	ctagsParser.Run(context)

	ctags := context[constants.CTX_CTAGS_OF_PREPROC_SOURCE].([]*builder.CTag)

	require.Equal(t, 5, len(ctags))
	idx := 0
	require.Equal(t, "set_values", ctags[idx].FunctionName)
	require.Equal(t, "prototype", ctags[idx].Kind)
	require.Equal(t, "Rectangle", ctags[idx].Class)
	idx++
	require.Equal(t, "area", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
	require.Equal(t, "Rectangle", ctags[idx].Class)
	idx++
	require.Equal(t, "set_values", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
	require.Equal(t, "Rectangle", ctags[idx].Class)
	idx++
	require.Equal(t, "setup", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
	idx++
	require.Equal(t, "loop", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
}

func TestCTagsParserStructWithFunctions(t *testing.T) {
	context := make(map[string]interface{})

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserStructWithFunctions.txt"))
	NoError(t, err)

	context[constants.CTX_CTAGS_OUTPUT] = string(bytes)

	ctagsParser := builder.CTagsParser{}
	ctagsParser.Run(context)

	ctags := context[constants.CTX_CTAGS_OF_PREPROC_SOURCE].([]*builder.CTag)

	require.Equal(t, 8, len(ctags))
	idx := 0
	require.Equal(t, "sensorData", ctags[idx].FunctionName)
	require.Equal(t, "struct", ctags[idx].Kind)
	idx++
	require.Equal(t, "sensorData", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
	require.Equal(t, "sensorData", ctags[idx].Struct)
	idx++
	require.Equal(t, "sensorData", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
	require.Equal(t, "sensorData", ctags[idx].Struct)
	idx++
	require.Equal(t, "sensors", ctags[idx].FunctionName)
	require.Equal(t, "variable", ctags[idx].Kind)
	idx++
	require.Equal(t, "sensor1", ctags[idx].FunctionName)
	require.Equal(t, "variable", ctags[idx].Kind)
	idx++
	require.Equal(t, "sensor2", ctags[idx].FunctionName)
	require.Equal(t, "variable", ctags[idx].Kind)
	idx++
	require.Equal(t, "setup", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
	idx++
	require.Equal(t, "loop", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
}

func TestCTagsParserDefaultArguments(t *testing.T) {
	context := make(map[string]interface{})

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserDefaultArguments.txt"))
	NoError(t, err)

	context[constants.CTX_CTAGS_OUTPUT] = string(bytes)

	ctagsParser := builder.CTagsParser{}
	ctagsParser.Run(context)

	ctags := context[constants.CTX_CTAGS_OF_PREPROC_SOURCE].([]*builder.CTag)

	require.Equal(t, 3, len(ctags))
	idx := 0
	require.Equal(t, "test", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
	require.Equal(t, "(int x = 1)", ctags[idx].Signature)
	idx++
	require.Equal(t, "setup", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
	idx++
	require.Equal(t, "loop", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
}

func TestCTagsParserNamespace(t *testing.T) {
	context := make(map[string]interface{})

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserNamespace.txt"))
	NoError(t, err)

	context[constants.CTX_CTAGS_OUTPUT] = string(bytes)

	ctagsParser := builder.CTagsParser{}
	ctagsParser.Run(context)

	ctags := context[constants.CTX_CTAGS_OF_PREPROC_SOURCE].([]*builder.CTag)

	require.Equal(t, 3, len(ctags))
	idx := 0
	require.Equal(t, "value", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
	require.Equal(t, "Test", ctags[idx].Namespace)
	idx++
	require.Equal(t, "setup", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
	idx++
	require.Equal(t, "loop", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
}

func TestCTagsParserStatic(t *testing.T) {
	context := make(map[string]interface{})

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserStatic.txt"))
	NoError(t, err)

	context[constants.CTX_CTAGS_OUTPUT] = string(bytes)

	ctagsParser := builder.CTagsParser{}
	ctagsParser.Run(context)

	ctags := context[constants.CTX_CTAGS_OF_PREPROC_SOURCE].([]*builder.CTag)

	require.Equal(t, 3, len(ctags))
	idx := 0
	require.Equal(t, "setup", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
	idx++
	require.Equal(t, "loop", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
	idx++
	require.Equal(t, "doStuff", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
}

func TestCTagsParserFunctionPointer(t *testing.T) {
	context := make(map[string]interface{})

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserFunctionPointer.txt"))
	NoError(t, err)

	context[constants.CTX_CTAGS_OUTPUT] = string(bytes)

	ctagsParser := builder.CTagsParser{}
	ctagsParser.Run(context)

	ctags := context[constants.CTX_CTAGS_OF_PREPROC_SOURCE].([]*builder.CTag)

	require.Equal(t, 4, len(ctags))
	idx := 0
	require.Equal(t, "t1Callback", ctags[idx].FunctionName)
	require.Equal(t, "variable", ctags[idx].Kind)
	idx++
	require.Equal(t, "t1Callback", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
	idx++
	require.Equal(t, "setup", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
	idx++
	require.Equal(t, "loop", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
}

func TestCTagsParserFunctionPointers(t *testing.T) {
	context := make(map[string]interface{})

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserFunctionPointers.txt"))
	NoError(t, err)

	context[constants.CTX_CTAGS_OUTPUT] = string(bytes)

	ctagsParser := builder.CTagsParser{}
	ctagsParser.Run(context)

	ctags := context[constants.CTX_CTAGS_OF_PREPROC_SOURCE].([]*builder.CTag)

	require.Equal(t, 5, len(ctags))
	idx := 0
	require.Equal(t, "setup", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
	idx++
	require.Equal(t, "loop", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
	idx++
	require.Equal(t, "func", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
	idx++
	require.Equal(t, "funcArr", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
	require.Equal(t, "()", ctags[idx].Signature)
	idx++
	require.Equal(t, "funcCombo", ctags[idx].FunctionName)
	require.Equal(t, "function", ctags[idx].Kind)
	require.Equal(t, "(void (*(&in)[5])(int))", ctags[idx].Signature)

}
