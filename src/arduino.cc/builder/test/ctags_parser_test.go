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
	"testing"
)

func TestCTagsParserShouldListPrototypes(t *testing.T) {
	context := make(map[string]interface{})

	context[constants.CTX_CTAGS_OUTPUT] = "server\t/tmp/sketch7210316334309249705.cpp\t/^YunServer server;$/;\"\tkind:variable\tline:31\n" +
		"setup\t/tmp/sketch7210316334309249705.cpp\t/^void setup() {$/;\"\tkind:function\tline:33\tsignature:()\treturntype:void\n" +
		"loop\t/tmp/sketch7210316334309249705.cpp\t/^void loop() {$/;\"\tkind:function\tline:46\tsignature:()\treturntype:void\n" +
		"process\t/tmp/sketch7210316334309249705.cpp\t/^void process(YunClient client);$/;\"\tkind:prototype\tline:61\tsignature:(YunClient client)\treturntype:void\n" +
		"process\t/tmp/sketch7210316334309249705.cpp\t/^void process(YunClient client) {$/;\"\tkind:function\tline:62\tsignature:(YunClient client)\treturntype:void\n" +
		"digitalCommand\t/tmp/sketch7210316334309249705.cpp\t/^void digitalCommand(YunClient client) {$/;\"\tkind:function\tline:82\tsignature:(YunClient client)\treturntype:void\n" +
		"analogCommand\t/tmp/sketch7210316334309249705.cpp\t/^void analogCommand(YunClient client) {$/;\"\tkind:function\tline:110\tsignature:(YunClient client)\treturntype:void\n" +
		"modeCommand\t/tmp/sketch7210316334309249705.cpp\t/^void modeCommand(YunClient client) {$/;\"\tkind:function\tline:151\tsignature:(YunClient client)\treturntype:void\n"

	ctagsParser := builder.CTagsParser{PrototypesField: constants.CTX_PROTOTYPES}
	ctagsParser.Run(context)

	prototypes := context[constants.CTX_PROTOTYPES].([]string)

	require.Equal(t, 5, len(prototypes))
	require.Equal(t, "void setup();", prototypes[0])
	require.Equal(t, "void loop();", prototypes[1])
	require.Equal(t, "void digitalCommand(YunClient client);", prototypes[2])
	require.Equal(t, "void analogCommand(YunClient client);", prototypes[3])
	require.Equal(t, "void modeCommand(YunClient client);", prototypes[4])
}

func TestCTagsParserShouldListTemplates(t *testing.T) {
	context := make(map[string]interface{})

	context[constants.CTX_CTAGS_OUTPUT] = "minimum\t/tmp/sketch8398023134925534708.cpp\t/^template <typename T> T minimum (T a, T b) $/;\"\tkind:function\tline:2\tsignature:(T a, T b)\treturntype:templateT\n" +
		"setup\t/tmp/sketch8398023134925534708.cpp\t/^void setup () $/;\"\tkind:function\tline:9\tsignature:()\treturntype:void\n" +
		"loop\t/tmp/sketch8398023134925534708.cpp\t/^void loop () { }$/;\"\tkind:function\tline:13\tsignature:()\treturntype:void\n"

	ctagsParser := builder.CTagsParser{PrototypesField: constants.CTX_PROTOTYPES}
	ctagsParser.Run(context)

	prototypes := context[constants.CTX_PROTOTYPES].([]string)

	require.Equal(t, 3, len(prototypes))
	require.Equal(t, "template <typename T> T minimum (T a, T b);", prototypes[0])
	require.Equal(t, "void setup();", prototypes[1])
	require.Equal(t, "void loop();", prototypes[2])
}

func TestCTagsParserShouldListTemplates2(t *testing.T) {
	context := make(map[string]interface{})

	context[constants.CTX_CTAGS_OUTPUT] = "setup\t/tmp/sketch463160524247569568.cpp\t/^void setup() {$/;\"\tkind:function\tline:1\tsignature:()\treturntype:void\n" +
		"loop\t/tmp/sketch463160524247569568.cpp\t/^void loop() {$/;\"\tkind:function\tline:6\tsignature:()\treturntype:void\n" +
		"SRAM_writeAnything\t/tmp/sketch463160524247569568.cpp\t/^template <class T> int SRAM_writeAnything(int ee, const T& value)$/;\"\tkind:function\tline:11\tsignature:(int ee, const T& value)\treturntype:template int\n" +
		"SRAM_readAnything\t/tmp/sketch463160524247569568.cpp\t/^template <class T> int SRAM_readAnything(int ee, T& value)$/;\"\tkind:function\tline:21\tsignature:(int ee, T& value)\treturntype:template int\n"

	ctagsParser := builder.CTagsParser{PrototypesField: constants.CTX_PROTOTYPES}
	ctagsParser.Run(context)

	prototypes := context[constants.CTX_PROTOTYPES].([]string)

	require.Equal(t, 4, len(prototypes))
	require.Equal(t, "void setup();", prototypes[0])
	require.Equal(t, "void loop();", prototypes[1])
	require.Equal(t, "template <class T> int SRAM_writeAnything(int ee, const T& value);", prototypes[2])
	require.Equal(t, "template <class T> int SRAM_readAnything(int ee, T& value);", prototypes[3])
}

func TestCTagsParserShouldDealWithClasses(t *testing.T) {
	context := make(map[string]interface{})

	context[constants.CTX_CTAGS_OUTPUT] = "SleepCycle\t/tmp/sketch9043227824785312266.cpp\t/^        SleepCycle( const char* name );$/;\"\tkind:prototype\tline:4\tsignature:( const char* name )\n" +
		"SleepCycle\t/tmp/sketch9043227824785312266.cpp\t/^    SleepCycle::SleepCycle( const char* name )$/;\"\tkind:function\tline:8\tsignature:( const char* name )\n"

	ctagsParser := builder.CTagsParser{PrototypesField: constants.CTX_PROTOTYPES}
	ctagsParser.Run(context)

	prototypes := context[constants.CTX_PROTOTYPES].([]string)

	require.Equal(t, 0, len(prototypes))
}

func TestCTagsParserShouldDealWithStructs(t *testing.T) {
	context := make(map[string]interface{})

	context[constants.CTX_CTAGS_OUTPUT] = "A_NEW_TYPE\t/tmp/sketch8930345717354294915.cpp\t/^struct A_NEW_TYPE {$/;\"\tkind:struct\tline:3\n" +
		"foo\t/tmp/sketch8930345717354294915.cpp\t/^} foo;$/;\"\tkind:variable\tline:7\ttyperef:struct:A_NEW_TYPE\n" +
		"setup\t/tmp/sketch8930345717354294915.cpp\t/^void setup() {$/;\"\tkind:function\tline:9\tsignature:()\treturntype:void\n" +
		"loop\t/tmp/sketch8930345717354294915.cpp\t/^void loop() {$/;\"\tkind:function\tline:13\tsignature:()\treturntype:void\n" +
		"dostuff\t/tmp/sketch8930345717354294915.cpp\t/^void dostuff (A_NEW_TYPE * bar)$/;\"\tkind:function\tline:17\tsignature:(A_NEW_TYPE * bar)\treturntype:void\n"

	ctagsParser := builder.CTagsParser{PrototypesField: constants.CTX_PROTOTYPES}
	ctagsParser.Run(context)

	prototypes := context[constants.CTX_PROTOTYPES].([]string)

	require.Equal(t, 3, len(prototypes))
	require.Equal(t, "void setup();", prototypes[0])
	require.Equal(t, "void loop();", prototypes[1])
	require.Equal(t, "void dostuff(A_NEW_TYPE * bar);", prototypes[2])
}

func TestCTagsParserShouldDealWithMacros(t *testing.T) {
	context := make(map[string]interface{})

	context[constants.CTX_CTAGS_OUTPUT] = "DEBUG\t/tmp/sketch5976699731718729500.cpp\t1;\"\tkind:macro\tline:1\n" +
		"DISABLED\t/tmp/sketch5976699731718729500.cpp\t2;\"\tkind:macro\tline:2\n" +
		"hello\t/tmp/sketch5976699731718729500.cpp\t/^String hello = \"world!\";$/;\"\tkind:variable\tline:16\n" +
		"setup\t/tmp/sketch5976699731718729500.cpp\t/^void setup() {$/;\"\tkind:function\tline:18\tsignature:()\treturntype:void\n" +
		"loop\t/tmp/sketch5976699731718729500.cpp\t/^void loop() {$/;\"\tkind:function\tline:23\tsignature:()\treturntype:void\n" +
		"debug\t/tmp/sketch5976699731718729500.cpp\t/^void debug() {$/;\"\tkind:function\tline:35\tsignature:()\treturntype:void\n" +
		"disabledIsDefined\t/tmp/sketch5976699731718729500.cpp\t/^void disabledIsDefined() {$/;\"\tkind:function\tline:46\tsignature:()\treturntype:void\n" +
		"useMyType\t/tmp/sketch5976699731718729500.cpp\t/^int useMyType(MyType type) {$/;\"\tkind:function\tline:50\tsignature:(MyType type)\treturntype:int\n"

	ctagsParser := builder.CTagsParser{PrototypesField: constants.CTX_PROTOTYPES}
	ctagsParser.Run(context)

	prototypes := context[constants.CTX_PROTOTYPES].([]string)

	require.Equal(t, 5, len(prototypes))
	require.Equal(t, "void setup();", prototypes[0])
	require.Equal(t, "void loop();", prototypes[1])
	require.Equal(t, "void debug();", prototypes[2])
	require.Equal(t, "void disabledIsDefined();", prototypes[3])
	require.Equal(t, "int useMyType(MyType type);", prototypes[4])
}

func TestCTagsParserShouldDealFunctionWithDifferentSignatures(t *testing.T) {
	context := make(map[string]interface{})

	context[constants.CTX_CTAGS_OUTPUT] = "getBytes	/tmp/test260613593/preproc/ctags_target.cpp	/^ void getBytes(unsigned char *buf, unsigned int bufsize, unsigned int index=0) const;$/;\"	kind:prototype	line:4330	signature:(unsigned char *buf, unsigned int bufsize, unsigned int index=0) const	returntype:void\n" +
		"getBytes	/tmp/test260613593/preproc/ctags_target.cpp	/^boolean getBytes( byte addr, int amount ) // updates the byte array \"received\" with the given amount of bytes, read from the given address$/;\"	kind:function	line:5031	signature:( byte addr, int amount )	returntype:boolean\n" +
		"getBytes	/tmp/test260613593/preproc/ctags_target.cpp	/^boolean getBytes( byte addr, int amount ) // updates the byte array \"received\" with the given amount of bytes, read from the given address$/;\"	kind:function	line:214	signature:( byte addr, int amount )	returntype:boolean"

	ctagsParser := builder.CTagsParser{PrototypesField: constants.CTX_PROTOTYPES}
	ctagsParser.Run(context)

	prototypes := context[constants.CTX_PROTOTYPES].([]string)

	require.Equal(t, 1, len(prototypes))
	require.Equal(t, "boolean getBytes( byte addr, int amount );", prototypes[0])
}

func TestCTagsParserClassMembersAreFilteredOut(t *testing.T) {
	context := make(map[string]interface{})

	context[constants.CTX_CTAGS_OUTPUT] = "set_values\t/tmp/test834438754/preproc/ctags_target.cpp\t/^    void set_values (int,int);$/;\"\tkind:prototype\tline:5\tclass:Rectangle\tsignature:(int,int)\treturntype:void\n" +
		"area\t/tmp/test834438754/preproc/ctags_target.cpp\t/^    int area() {return width*height;}$/;\"\tkind:function\tline:6\tclass:Rectangle\tsignature:()\treturntype:int\n" +
		"set_values\t/tmp/test834438754/preproc/ctags_target.cpp\t/^void Rectangle::set_values (int x, int y) {$/;\"\tkind:function\tline:9\tclass:Rectangle\tsignature:(int x, int y)\treturntype:void\n" +
		"setup\t/tmp/test834438754/preproc/ctags_target.cpp\t/^void setup() {$/;\"\tkind:function\tline:14\tsignature:()\treturntype:void\n" +
		"loop\t/tmp/test834438754/preproc/ctags_target.cpp\t/^void loop() {$/;\"\tkind:function\tline:18\tsignature:()\treturntype:void\n"

	ctagsParser := builder.CTagsParser{PrototypesField: constants.CTX_PROTOTYPES}
	ctagsParser.Run(context)

	prototypes := context[constants.CTX_PROTOTYPES].([]string)

	require.Equal(t, 2, len(prototypes))
	require.Equal(t, "void setup();", prototypes[0])
	require.Equal(t, "void loop();", prototypes[1])
}
