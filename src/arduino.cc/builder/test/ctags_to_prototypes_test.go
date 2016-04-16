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
	"arduino.cc/builder/ctags"
	"arduino.cc/builder/types"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"path/filepath"
	"testing"
)

type CollectCtagsFromPreprocSource struct{}

func (*CollectCtagsFromPreprocSource) Run(ctx *types.Context) error {
	ctx.CTagsCollected = ctx.CTagsOfPreprocessedSource
	return nil
}

func TestCTagsToPrototypesShouldListPrototypes(t *testing.T) {
	ctx := &types.Context{}

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserShouldListPrototypes.txt"))
	NoError(t, err)

	ctx.CTagsOutput = string(bytes)

	commands := []types.Command{
		&ctags.CTagsParser{},
		&CollectCtagsFromPreprocSource{},
		&ctags.CTagsToPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}

	prototypes := ctx.Prototypes

	require.Equal(t, 5, len(prototypes))
	require.Equal(t, "void setup();", prototypes[0].Prototype)
	require.Equal(t, "/tmp/sketch7210316334309249705.cpp", prototypes[0].File)
	require.Equal(t, "void loop();", prototypes[1].Prototype)
	require.Equal(t, "void digitalCommand(YunClient client);", prototypes[2].Prototype)
	require.Equal(t, "void analogCommand(YunClient client);", prototypes[3].Prototype)
	require.Equal(t, "void modeCommand(YunClient client);", prototypes[4].Prototype)

	require.Equal(t, 33, ctx.PrototypesLineWhereToInsert)
}

func TestCTagsToPrototypesShouldListTemplates(t *testing.T) {
	ctx := &types.Context{}

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserShouldListTemplates.txt"))
	NoError(t, err)

	ctx.CTagsOutput = string(bytes)

	commands := []types.Command{
		&ctags.CTagsParser{},
		&CollectCtagsFromPreprocSource{},
		&ctags.CTagsToPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}

	prototypes := ctx.Prototypes

	require.Equal(t, 3, len(prototypes))
	require.Equal(t, "template <typename T> T minimum (T a, T b);", prototypes[0].Prototype)
	require.Equal(t, "/tmp/sketch8398023134925534708.cpp", prototypes[0].File)
	require.Equal(t, "void setup();", prototypes[1].Prototype)
	require.Equal(t, "void loop();", prototypes[2].Prototype)

	require.Equal(t, 2, ctx.PrototypesLineWhereToInsert)
}

func TestCTagsToPrototypesShouldListTemplates2(t *testing.T) {
	ctx := &types.Context{}

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserShouldListTemplates2.txt"))
	NoError(t, err)

	ctx.CTagsOutput = string(bytes)

	commands := []types.Command{
		&ctags.CTagsParser{},
		&CollectCtagsFromPreprocSource{},
		&ctags.CTagsToPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}

	prototypes := ctx.Prototypes

	require.Equal(t, 4, len(prototypes))
	require.Equal(t, "void setup();", prototypes[0].Prototype)
	require.Equal(t, "/tmp/sketch463160524247569568.cpp", prototypes[0].File)
	require.Equal(t, "void loop();", prototypes[1].Prototype)
	require.Equal(t, "template <class T> int SRAM_writeAnything(int ee, const T& value);", prototypes[2].Prototype)
	require.Equal(t, "template <class T> int SRAM_readAnything(int ee, T& value);", prototypes[3].Prototype)

	require.Equal(t, 1, ctx.PrototypesLineWhereToInsert)
}

func TestCTagsToPrototypesShouldDealWithClasses(t *testing.T) {
	ctx := &types.Context{}

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserShouldDealWithClasses.txt"))
	NoError(t, err)

	ctx.CTagsOutput = string(bytes)

	commands := []types.Command{
		&ctags.CTagsParser{},
		&CollectCtagsFromPreprocSource{},
		&ctags.CTagsToPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}

	prototypes := ctx.Prototypes

	require.Equal(t, 0, len(prototypes))

	require.Equal(t, 8, ctx.PrototypesLineWhereToInsert)
}

func TestCTagsToPrototypesShouldDealWithStructs(t *testing.T) {
	ctx := &types.Context{}

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserShouldDealWithStructs.txt"))
	NoError(t, err)

	ctx.CTagsOutput = string(bytes)

	commands := []types.Command{
		&ctags.CTagsParser{},
		&CollectCtagsFromPreprocSource{},
		&ctags.CTagsToPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}

	prototypes := ctx.Prototypes

	require.Equal(t, 3, len(prototypes))
	require.Equal(t, "void setup();", prototypes[0].Prototype)
	require.Equal(t, "/tmp/sketch8930345717354294915.cpp", prototypes[0].File)
	require.Equal(t, "void loop();", prototypes[1].Prototype)
	require.Equal(t, "void dostuff(A_NEW_TYPE * bar);", prototypes[2].Prototype)

	require.Equal(t, 9, ctx.PrototypesLineWhereToInsert)
}

func TestCTagsToPrototypesShouldDealWithMacros(t *testing.T) {
	ctx := &types.Context{}

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserShouldDealWithMacros.txt"))
	NoError(t, err)

	ctx.CTagsOutput = string(bytes)

	commands := []types.Command{
		&ctags.CTagsParser{},
		&CollectCtagsFromPreprocSource{},
		&ctags.CTagsToPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}

	prototypes := ctx.Prototypes

	require.Equal(t, 5, len(prototypes))
	require.Equal(t, "void setup();", prototypes[0].Prototype)
	require.Equal(t, "/tmp/sketch5976699731718729500.cpp", prototypes[0].File)
	require.Equal(t, "void loop();", prototypes[1].Prototype)
	require.Equal(t, "void debug();", prototypes[2].Prototype)
	require.Equal(t, "void disabledIsDefined();", prototypes[3].Prototype)
	require.Equal(t, "int useMyType(MyType type);", prototypes[4].Prototype)

	require.Equal(t, 18, ctx.PrototypesLineWhereToInsert)
}

func TestCTagsToPrototypesShouldDealFunctionWithDifferentSignatures(t *testing.T) {
	ctx := &types.Context{}

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserShouldDealFunctionWithDifferentSignatures.txt"))
	NoError(t, err)

	ctx.CTagsOutput = string(bytes)

	commands := []types.Command{
		&ctags.CTagsParser{},
		&CollectCtagsFromPreprocSource{},
		&ctags.CTagsToPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}

	prototypes := ctx.Prototypes

	require.Equal(t, 1, len(prototypes))
	require.Equal(t, "boolean getBytes( byte addr, int amount );", prototypes[0].Prototype)
	require.Equal(t, "/tmp/test260613593/preproc/ctags_target.cpp", prototypes[0].File)

	require.Equal(t, 5031, ctx.PrototypesLineWhereToInsert)
}

func TestCTagsToPrototypesClassMembersAreFilteredOut(t *testing.T) {
	ctx := &types.Context{}

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserClassMembersAreFilteredOut.txt"))
	NoError(t, err)

	ctx.CTagsOutput = string(bytes)

	commands := []types.Command{
		&ctags.CTagsParser{},
		&CollectCtagsFromPreprocSource{},
		&ctags.CTagsToPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}

	prototypes := ctx.Prototypes

	require.Equal(t, 2, len(prototypes))
	require.Equal(t, "void setup();", prototypes[0].Prototype)
	require.Equal(t, "/tmp/test834438754/preproc/ctags_target.cpp", prototypes[0].File)
	require.Equal(t, "void loop();", prototypes[1].Prototype)

	require.Equal(t, 14, ctx.PrototypesLineWhereToInsert)
}

func TestCTagsToPrototypesStructWithFunctions(t *testing.T) {
	ctx := &types.Context{}

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserStructWithFunctions.txt"))
	NoError(t, err)

	ctx.CTagsOutput = string(bytes)

	commands := []types.Command{
		&ctags.CTagsParser{},
		&CollectCtagsFromPreprocSource{},
		&ctags.CTagsToPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}

	prototypes := ctx.Prototypes

	require.Equal(t, 2, len(prototypes))
	require.Equal(t, "void setup();", prototypes[0].Prototype)
	require.Equal(t, "/tmp/build7315640391316178285.tmp/preproc/ctags_target.cpp", prototypes[0].File)
	require.Equal(t, "void loop();", prototypes[1].Prototype)

	require.Equal(t, 16, ctx.PrototypesLineWhereToInsert)
}

func TestCTagsToPrototypesDefaultArguments(t *testing.T) {
	ctx := &types.Context{}

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserDefaultArguments.txt"))
	NoError(t, err)

	ctx.CTagsOutput = string(bytes)

	commands := []types.Command{
		&ctags.CTagsParser{},
		&CollectCtagsFromPreprocSource{},
		&ctags.CTagsToPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}

	prototypes := ctx.Prototypes

	require.Equal(t, 3, len(prototypes))
	require.Equal(t, "void test(int x = 1);", prototypes[0].Prototype)
	require.Equal(t, "void setup();", prototypes[1].Prototype)
	require.Equal(t, "/tmp/test179252494/preproc/ctags_target.cpp", prototypes[1].File)
	require.Equal(t, "void loop();", prototypes[2].Prototype)

	require.Equal(t, 2, ctx.PrototypesLineWhereToInsert)
}

func TestCTagsToPrototypesNamespace(t *testing.T) {
	ctx := &types.Context{}

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserNamespace.txt"))
	NoError(t, err)

	ctx.CTagsOutput = string(bytes)

	commands := []types.Command{
		&ctags.CTagsParser{},
		&CollectCtagsFromPreprocSource{},
		&ctags.CTagsToPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}

	prototypes := ctx.Prototypes

	require.Equal(t, 2, len(prototypes))
	require.Equal(t, "void setup();", prototypes[0].Prototype)
	require.Equal(t, "/tmp/test030883150/preproc/ctags_target.cpp", prototypes[0].File)
	require.Equal(t, "void loop();", prototypes[1].Prototype)

	require.Equal(t, 8, ctx.PrototypesLineWhereToInsert)
}

func TestCTagsToPrototypesStatic(t *testing.T) {
	ctx := &types.Context{}

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserStatic.txt"))
	NoError(t, err)

	ctx.CTagsOutput = string(bytes)

	commands := []types.Command{
		&ctags.CTagsParser{},
		&CollectCtagsFromPreprocSource{},
		&ctags.CTagsToPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}

	prototypes := ctx.Prototypes

	require.Equal(t, 3, len(prototypes))
	require.Equal(t, "void setup();", prototypes[0].Prototype)
	require.Equal(t, "/tmp/test542833488/preproc/ctags_target.cpp", prototypes[0].File)
	require.Equal(t, "void loop();", prototypes[1].Prototype)
	require.Equal(t, "void doStuff();", prototypes[2].Prototype)
	require.Equal(t, "static", prototypes[2].Modifiers)

	require.Equal(t, 2, ctx.PrototypesLineWhereToInsert)
}

func TestCTagsToPrototypesFunctionPointer(t *testing.T) {
	ctx := &types.Context{}

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserFunctionPointer.txt"))
	NoError(t, err)

	ctx.CTagsOutput = string(bytes)

	commands := []types.Command{
		&ctags.CTagsParser{},
		&CollectCtagsFromPreprocSource{},
		&ctags.CTagsToPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}

	prototypes := ctx.Prototypes

	require.Equal(t, 3, len(prototypes))
	require.Equal(t, "void t1Callback();", prototypes[0].Prototype)
	require.Equal(t, "/tmp/test547238273/preproc/ctags_target.cpp", prototypes[0].File)
	require.Equal(t, "void setup();", prototypes[1].Prototype)
	require.Equal(t, "void loop();", prototypes[2].Prototype)

	require.Equal(t, 2, ctx.PrototypesLineWhereToInsert)
}

func TestCTagsToPrototypesFunctionPointers(t *testing.T) {
	ctx := &types.Context{}

	bytes, err := ioutil.ReadFile(filepath.Join("ctags_output", "TestCTagsParserFunctionPointers.txt"))
	NoError(t, err)

	ctx.CTagsOutput = string(bytes)

	commands := []types.Command{
		&ctags.CTagsParser{},
		&CollectCtagsFromPreprocSource{},
		&ctags.CTagsToPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}

	prototypes := ctx.Prototypes
	require.Equal(t, 2, len(prototypes))
	require.Equal(t, "void setup();", prototypes[0].Prototype)
	require.Equal(t, "/tmp/test907446433/preproc/ctags_target.cpp", prototypes[0].File)
	require.Equal(t, "void loop();", prototypes[1].Prototype)

	require.Equal(t, 2, ctx.PrototypesLineWhereToInsert)
}
