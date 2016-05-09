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
	"arduino.cc/builder/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGCCMinusMOutputParser(t *testing.T) {
	ctx := &types.Context{}

	output := "sketch_with_config.o: sketch_with_config.ino config.h de\\ bug.h Bridge.h\n" +
		"\n" +
		"config.h:\n" +
		"\n" +
		"de\\ bug.h:\n" +
		"\n" +
		"Bridge.h:\n"

	ctx.OutputGccMinusM = output

	parser := builder.GCCMinusMOutputParser{}
	err := parser.Run(ctx)
	NoError(t, err)

	includes := ctx.Includes
	require.Equal(t, 3, len(includes))
	require.Equal(t, "config.h", includes[0])
	require.Equal(t, "de bug.h", includes[1])
	require.Equal(t, "Bridge.h", includes[2])
}

func TestGCCMinusMOutputParserEmptyOutput(t *testing.T) {
	ctx := &types.Context{}

	output := "sketch.ino.o: /tmp/test699709208/sketch/sketch.ino.cpp"

	ctx.OutputGccMinusM = output

	parser := builder.GCCMinusMOutputParser{}
	err := parser.Run(ctx)
	NoError(t, err)

	includes := ctx.Includes
	require.Equal(t, 0, len(includes))
}

func TestGCCMinusMOutputParserFirstLineOnMultipleLines(t *testing.T) {
	ctx := &types.Context{}

	output := "sketch_with_config.ino.o: \\\n" +
		" /tmp/test097286304/sketch/sketch_with_config.ino.cpp \\\n" +
		" /tmp/test097286304/sketch/config.h \\\n" +
		" /tmp/test097286304/sketch/includes/de\\ bug.h Bridge.h\n" +
		"\n" +
		"/tmp/test097286304/sketch/config.h:\n" +
		"\n" +
		"/tmp/test097286304/sketch/includes/de\\ bug.h:\n" +
		"\n" +
		"Bridge.h:\n"

	ctx.OutputGccMinusM = output

	parser := builder.GCCMinusMOutputParser{}
	err := parser.Run(ctx)
	NoError(t, err)

	includes := ctx.Includes
	require.Equal(t, 3, len(includes))
	require.Equal(t, "/tmp/test097286304/sketch/config.h", includes[0])
	require.Equal(t, "/tmp/test097286304/sketch/includes/de bug.h", includes[1])
	require.Equal(t, "Bridge.h", includes[2])
}
