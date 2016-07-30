/*
 * This file is part of Arduino Builder.
 *
 * Copyright 2016 Arduino LLC (http://www.arduino.cc/)
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
 */

package timeutils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTime(t *testing.T) {
	cet, err := time.LoadLocation("CET")
	require.NoError(t, err)
	ast, err := time.LoadLocation("Australia/Sydney")
	require.NoError(t, err)

	firstJanuary2015CET := time.Date(2015, 1, 1, 0, 0, 0, 0, cet)
	require.Equal(t, int64(1420066800), firstJanuary2015CET.Unix())
	require.Equal(t, int64(1420066800+3600), LocalUnix(firstJanuary2015CET))
	require.Equal(t, 3600, TimezoneOffsetNoDST(firstJanuary2015CET))
	require.Equal(t, 0, DaylightSavingsOffset(firstJanuary2015CET))

	fall2015CET := time.Date(2015, 9, 23, 0, 0, 0, 0, cet)
	require.Equal(t, int64(1442959200), fall2015CET.Unix())
	require.Equal(t, int64(1442959200+3600+3600), LocalUnix(fall2015CET))
	require.Equal(t, 3600, TimezoneOffsetNoDST(fall2015CET))
	require.Equal(t, 3600, DaylightSavingsOffset(fall2015CET))

	firstJan2015AST := time.Date(2015, 1, 1, 0, 0, 0, 0, ast)
	require.Equal(t, int64(1420030800), firstJan2015AST.Unix())
	require.Equal(t, int64(1420030800+36000+3600), LocalUnix(firstJan2015AST))
	require.Equal(t, 36000, TimezoneOffsetNoDST(firstJan2015AST))
	require.Equal(t, 3600, DaylightSavingsOffset(firstJan2015AST))

	fall2015AST := time.Date(2015, 9, 23, 0, 0, 0, 0, ast)
	require.Equal(t, int64(1442930400), fall2015AST.Unix())
	require.Equal(t, int64(1442930400+36000), LocalUnix(fall2015AST))
	require.Equal(t, 36000, TimezoneOffsetNoDST(fall2015AST))
	require.Equal(t, 0, DaylightSavingsOffset(fall2015AST))
}
