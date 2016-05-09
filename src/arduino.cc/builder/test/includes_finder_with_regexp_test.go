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
	"sort"
	"testing"
)

func TestIncludesFinderWithRegExp(t *testing.T) {
	ctx := &types.Context{}

	output := "/some/path/sketch.ino:1:17: fatal error: SPI.h: No such file or directory\n" +
		"#include <SPI.h>\n" +
		"^\n" +
		"compilation terminated."
	ctx.Source = output

	parser := builder.IncludesFinderWithRegExp{Source: &ctx.Source}
	err := parser.Run(ctx)
	NoError(t, err)

	includes := ctx.Includes
	require.Equal(t, 1, len(includes))
	require.Equal(t, "SPI.h", includes[0])
}

func TestIncludesFinderWithRegExpEmptyOutput(t *testing.T) {
	ctx := &types.Context{}

	output := ""

	ctx.Source = output

	parser := builder.IncludesFinderWithRegExp{Source: &ctx.Source}
	err := parser.Run(ctx)
	NoError(t, err)

	includes := ctx.Includes
	require.Equal(t, 0, len(includes))
}

func TestIncludesFinderWithRegExpPreviousIncludes(t *testing.T) {
	ctx := &types.Context{
		Includes: []string{"test.h"},
	}

	output := "/some/path/sketch.ino:1:17: fatal error: SPI.h: No such file or directory\n" +
		"#include <SPI.h>\n" +
		"^\n" +
		"compilation terminated."

	ctx.Source = output

	parser := builder.IncludesFinderWithRegExp{Source: &ctx.Source}
	err := parser.Run(ctx)
	NoError(t, err)

	includes := ctx.Includes
	require.Equal(t, 2, len(includes))
	sort.Strings(includes)
	require.Equal(t, "SPI.h", includes[0])
	require.Equal(t, "test.h", includes[1])
}

func TestIncludesFinderWithRegExpPaddedIncludes(t *testing.T) {
	ctx := &types.Context{}

	output := "/some/path/sketch.ino:1:33: fatal error: Wire.h: No such file or directory\n" +
		" #               include <Wire.h>\n" +
		"                                 ^\n" +
		"compilation terminated.\n"
	ctx.Source = output

	parser := builder.IncludesFinderWithRegExp{Source: &ctx.Source}
	err := parser.Run(ctx)
	NoError(t, err)

	includes := ctx.Includes
	require.Equal(t, 1, len(includes))
	sort.Strings(includes)
	require.Equal(t, "Wire.h", includes[0])
}

func TestIncludesFinderWithRegExpPaddedIncludes2(t *testing.T) {
	ctx := &types.Context{}

	output := "/some/path/sketch.ino:1:33: fatal error: Wire.h: No such file or directory\n" +
		" #\t\t\tinclude <Wire.h>\n" +
		"                                 ^\n" +
		"compilation terminated.\n"
	ctx.Source = output

	parser := builder.IncludesFinderWithRegExp{Source: &ctx.Source}
	err := parser.Run(ctx)
	NoError(t, err)

	includes := ctx.Includes
	require.Equal(t, 1, len(includes))
	sort.Strings(includes)
	require.Equal(t, "Wire.h", includes[0])
}

func TestIncludesFinderWithRegExpPaddedIncludes3(t *testing.T) {
	ctx := &types.Context{}

	output := "/some/path/sketch.ino:1:33: fatal error: SPI.h: No such file or directory\n" +
		"compilation terminated.\n"

	ctx.Source = output

	parser := builder.IncludesFinderWithRegExp{Source: &ctx.Source}
	err := parser.Run(ctx)
	NoError(t, err)

	includes := ctx.Includes
	require.Equal(t, 1, len(includes))
	sort.Strings(includes)
	require.Equal(t, "SPI.h", includes[0])
}
