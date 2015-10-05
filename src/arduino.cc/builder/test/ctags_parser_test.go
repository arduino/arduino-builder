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
	"arduino.cc/builder/types"
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

	ctagsParser := builder.CTagsParser{PrototypesField: constants.CTX_PROTOTYPES}
	ctagsParser.Run(context)

	prototypes := context[constants.CTX_PROTOTYPES].([]*types.Prototype)

	require.Equal(t, 5, len(prototypes))
	require.Equal(t, "void setup();", prototypes[0].Prototype)
	require.Equal(t, "void loop();", prototypes[1].Prototype)
	require.Equal(t, "void digitalCommand(YunClient client);", prototypes[2].Prototype)
	require.Equal(t, "void analogCommand(YunClient client);", prototypes[3].Prototype)
	require.Equal(t, "void modeCommand(YunClient client);", prototypes[4].Prototype)
}

func TestCTagsParserShouldListTemplates(t *testing.T) {
	context := make(map[string]interface{})

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserShouldListTemplates.txt"))
	NoError(t, err)

	context[constants.CTX_CTAGS_OUTPUT] = string(bytes)

	ctagsParser := builder.CTagsParser{PrototypesField: constants.CTX_PROTOTYPES}
	ctagsParser.Run(context)

	prototypes := context[constants.CTX_PROTOTYPES].([]*types.Prototype)

	require.Equal(t, 3, len(prototypes))
	require.Equal(t, "template <typename T> T minimum (T a, T b);", prototypes[0].Prototype)
	require.Equal(t, "void setup();", prototypes[1].Prototype)
	require.Equal(t, "void loop();", prototypes[2].Prototype)
}

func TestCTagsParserShouldListTemplates2(t *testing.T) {
	context := make(map[string]interface{})

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserShouldListTemplates2.txt"))
	NoError(t, err)

	context[constants.CTX_CTAGS_OUTPUT] = string(bytes)

	ctagsParser := builder.CTagsParser{PrototypesField: constants.CTX_PROTOTYPES}
	ctagsParser.Run(context)

	prototypes := context[constants.CTX_PROTOTYPES].([]*types.Prototype)

	require.Equal(t, 4, len(prototypes))
	require.Equal(t, "void setup();", prototypes[0].Prototype)
	require.Equal(t, "void loop();", prototypes[1].Prototype)
	require.Equal(t, "template <class T> int SRAM_writeAnything(int ee, const T& value);", prototypes[2].Prototype)
	require.Equal(t, "template <class T> int SRAM_readAnything(int ee, T& value);", prototypes[3].Prototype)
}

func TestCTagsParserShouldDealWithClasses(t *testing.T) {
	context := make(map[string]interface{})

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserShouldDealWithClasses.txt"))
	NoError(t, err)

	context[constants.CTX_CTAGS_OUTPUT] = string(bytes)

	ctagsParser := builder.CTagsParser{PrototypesField: constants.CTX_PROTOTYPES}
	ctagsParser.Run(context)

	prototypes := context[constants.CTX_PROTOTYPES].([]*types.Prototype)

	require.Equal(t, 0, len(prototypes))
}

func TestCTagsParserShouldDealWithStructs(t *testing.T) {
	context := make(map[string]interface{})

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserShouldDealWithStructs.txt"))
	NoError(t, err)

	context[constants.CTX_CTAGS_OUTPUT] = string(bytes)

	ctagsParser := builder.CTagsParser{PrototypesField: constants.CTX_PROTOTYPES}
	ctagsParser.Run(context)

	prototypes := context[constants.CTX_PROTOTYPES].([]*types.Prototype)

	require.Equal(t, 3, len(prototypes))
	require.Equal(t, "void setup();", prototypes[0].Prototype)
	require.Equal(t, "void loop();", prototypes[1].Prototype)
	require.Equal(t, "void dostuff(A_NEW_TYPE * bar);", prototypes[2].Prototype)
}

func TestCTagsParserShouldDealWithMacros(t *testing.T) {
	context := make(map[string]interface{})

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserShouldDealWithMacros.txt"))
	NoError(t, err)

	context[constants.CTX_CTAGS_OUTPUT] = string(bytes)

	ctagsParser := builder.CTagsParser{PrototypesField: constants.CTX_PROTOTYPES}
	ctagsParser.Run(context)

	prototypes := context[constants.CTX_PROTOTYPES].([]*types.Prototype)

	require.Equal(t, 5, len(prototypes))
	require.Equal(t, "void setup();", prototypes[0].Prototype)
	require.Equal(t, "void loop();", prototypes[1].Prototype)
	require.Equal(t, "void debug();", prototypes[2].Prototype)
	require.Equal(t, "void disabledIsDefined();", prototypes[3].Prototype)
	require.Equal(t, "int useMyType(MyType type);", prototypes[4].Prototype)
}

func TestCTagsParserShouldDealFunctionWithDifferentSignatures(t *testing.T) {
	context := make(map[string]interface{})

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserShouldDealFunctionWithDifferentSignatures.txt"))
	NoError(t, err)

	context[constants.CTX_CTAGS_OUTPUT] = string(bytes)

	ctagsParser := builder.CTagsParser{PrototypesField: constants.CTX_PROTOTYPES}
	ctagsParser.Run(context)

	prototypes := context[constants.CTX_PROTOTYPES].([]*types.Prototype)

	require.Equal(t, 1, len(prototypes))
	require.Equal(t, "boolean getBytes( byte addr, int amount );", prototypes[0].Prototype)
}

func TestCTagsParserClassMembersAreFilteredOut(t *testing.T) {
	context := make(map[string]interface{})

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserClassMembersAreFilteredOut.txt"))
	NoError(t, err)

	context[constants.CTX_CTAGS_OUTPUT] = string(bytes)

	ctagsParser := builder.CTagsParser{PrototypesField: constants.CTX_PROTOTYPES}
	ctagsParser.Run(context)

	prototypes := context[constants.CTX_PROTOTYPES].([]*types.Prototype)

	require.Equal(t, 2, len(prototypes))
	require.Equal(t, "void setup();", prototypes[0].Prototype)
	require.Equal(t, "void loop();", prototypes[1].Prototype)
}

func TestCTagsParserStructWithFunctions(t *testing.T) {
	context := make(map[string]interface{})

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserStructWithFunctions.txt"))
	NoError(t, err)

	context[constants.CTX_CTAGS_OUTPUT] = string(bytes)

	ctagsParser := builder.CTagsParser{PrototypesField: constants.CTX_PROTOTYPES}
	ctagsParser.Run(context)

	prototypes := context[constants.CTX_PROTOTYPES].([]*types.Prototype)

	require.Equal(t, 2, len(prototypes))
	require.Equal(t, "void setup();", prototypes[0].Prototype)
	require.Equal(t, "void loop();", prototypes[1].Prototype)
}

func TestCTagsParserDefaultArguments(t *testing.T) {
	context := make(map[string]interface{})

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserDefaultArguments.txt"))
	NoError(t, err)

	context[constants.CTX_CTAGS_OUTPUT] = string(bytes)

	ctagsParser := builder.CTagsParser{PrototypesField: constants.CTX_PROTOTYPES}
	ctagsParser.Run(context)

	prototypes := context[constants.CTX_PROTOTYPES].([]*types.Prototype)

	require.Equal(t, 2, len(prototypes))
	require.Equal(t, "void setup();", prototypes[0].Prototype)
	require.Equal(t, "void loop();", prototypes[1].Prototype)
}
