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
	"arduino.cc/builder/constants"
	"arduino.cc/builder/i18n"
	"arduino.cc/builder/props"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"runtime"
	"testing"
)

func TestPropertiesBoardsTxt(t *testing.T) {
	properties, err := props.Load(filepath.Join("props", "boards.txt"), i18n.HumanLogger{})

	NoError(t, err)

	require.Equal(t, "Processor", properties["menu.cpu"])
	require.Equal(t, "32256", properties["ethernet.upload.maximum_size"])
	require.Equal(t, "{build.usb_flags}", properties["robotMotor.build.extra_flags"])

	ethernet := props.SubTree(properties, "ethernet")
	require.Equal(t, "Arduino Ethernet", ethernet[constants.LIBRARY_NAME])
}

func TestPropertiesTestTxt(t *testing.T) {
	properties, err := props.Load(filepath.Join("props", "test.txt"), i18n.HumanLogger{})

	NoError(t, err)

	require.Equal(t, 4, len(properties))
	require.Equal(t, "value = 1", properties["key"])

	switch value := runtime.GOOS; value {
	case "linux":
		require.Equal(t, "is linux", properties["which.os"])
	case "windows":
		require.Equal(t, "is windows", properties["which.os"])
	case "darwin":
		require.Equal(t, "is macosx", properties["which.os"])
	default:
		require.FailNow(t, "unsupported OS")
	}
}

func TestExpandPropsInString(t *testing.T) {
	aMap := make(map[string]string)
	aMap["key1"] = "42"
	aMap["key2"] = "{key1}"

	str := "{key1} == {key2} == true"

	str = props.ExpandPropsInString(aMap, str)
	require.Equal(t, "42 == 42 == true", str)
}

func TestExpandPropsInString2(t *testing.T) {
	aMap := make(map[string]string)
	aMap["key2"] = "{key2}"
	aMap["key1"] = "42"

	str := "{key1} == {key2} == true"

	str = props.ExpandPropsInString(aMap, str)
	require.Equal(t, "42 == {key2} == true", str)
}

func TestDeleteUnexpandedPropsFromString(t *testing.T) {
	aMap := make(map[string]string)
	aMap["key1"] = "42"
	aMap["key2"] = "{key1}"

	str := "{key1} == {key2} == {key3} == true"

	str = props.ExpandPropsInString(aMap, str)
	str, err := props.DeleteUnexpandedPropsFromString(str)
	NoError(t, err)
	require.Equal(t, "42 == 42 ==  == true", str)
}

func TestDeleteUnexpandedPropsFromString2(t *testing.T) {
	aMap := make(map[string]string)
	aMap["key2"] = "42"

	str := "{key1} == {key2} == {key3} == true"

	str = props.ExpandPropsInString(aMap, str)
	str, err := props.DeleteUnexpandedPropsFromString(str)
	NoError(t, err)
	require.Equal(t, " == 42 ==  == true", str)
}

func TestPropertiesRedBeearLabBoardsTxt(t *testing.T) {
	properties, err := props.Load(filepath.Join("props", "redbearlab_boards.txt"), i18n.HumanLogger{})

	NoError(t, err)

	require.Equal(t, 83, len(properties))
	require.Equal(t, "Blend", properties["blend.name"])
	require.Equal(t, "arduino:arduino", properties["blend.build.core"])
	require.Equal(t, "0x2404", properties["blendmicro16.pid.0"])

	ethernet := props.SubTree(properties, "blend")
	require.Equal(t, "arduino:arduino", ethernet[constants.BUILD_PROPERTIES_BUILD_CORE])
}

func TestPropertiesBroken(t *testing.T) {
	_, err := props.Load(filepath.Join("props", "broken.txt"), i18n.HumanLogger{})

	require.Error(t, err)
}
