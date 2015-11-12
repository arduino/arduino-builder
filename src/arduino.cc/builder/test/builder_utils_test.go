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
	"arduino.cc/builder/builder_utils"
	"arduino.cc/builder/utils"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func sleep(t *testing.T) {
	dur, err := time.ParseDuration("1s")
	NoError(t, err)
	time.Sleep(dur)
}

func tempFile(t *testing.T, prefix string) string {
	file, err := ioutil.TempFile("", prefix)
	NoError(t, err)
	return file.Name()
}

func TestObjFileIsUpToDateObjMissing(t *testing.T) {
	sourceFile := tempFile(t, "source")
	defer os.RemoveAll(sourceFile)

	upToDate, err := builder_utils.ObjFileIsUpToDate(sourceFile, "", "")
	NoError(t, err)
	require.False(t, upToDate)
}

func TestObjFileIsUpToDateDepMissing(t *testing.T) {
	sourceFile := tempFile(t, "source")
	defer os.RemoveAll(sourceFile)

	objFile := tempFile(t, "obj")
	defer os.RemoveAll(objFile)

	upToDate, err := builder_utils.ObjFileIsUpToDate(sourceFile, objFile, "")
	NoError(t, err)
	require.False(t, upToDate)
}

func TestObjFileIsUpToDateObjOlder(t *testing.T) {
	objFile := tempFile(t, "obj")
	defer os.RemoveAll(objFile)
	depFile := tempFile(t, "dep")
	defer os.RemoveAll(depFile)

	sleep(t)

	sourceFile := tempFile(t, "source")
	defer os.RemoveAll(sourceFile)

	upToDate, err := builder_utils.ObjFileIsUpToDate(sourceFile, objFile, depFile)
	NoError(t, err)
	require.False(t, upToDate)
}

func TestObjFileIsUpToDateObjNewer(t *testing.T) {
	sourceFile := tempFile(t, "source")
	defer os.RemoveAll(sourceFile)

	sleep(t)

	objFile := tempFile(t, "obj")
	defer os.RemoveAll(objFile)
	depFile := tempFile(t, "dep")
	defer os.RemoveAll(depFile)

	upToDate, err := builder_utils.ObjFileIsUpToDate(sourceFile, objFile, depFile)
	NoError(t, err)
	require.True(t, upToDate)
}

func TestObjFileIsUpToDateDepIsNewer(t *testing.T) {
	sourceFile := tempFile(t, "source")
	defer os.RemoveAll(sourceFile)

	sleep(t)

	objFile := tempFile(t, "obj")
	defer os.RemoveAll(objFile)
	depFile := tempFile(t, "dep")
	defer os.RemoveAll(depFile)

	sleep(t)

	headerFile := tempFile(t, "header")
	defer os.RemoveAll(headerFile)

	utils.WriteFile(depFile, objFile+": \\\n\t"+sourceFile+" \\\n\t"+headerFile)

	upToDate, err := builder_utils.ObjFileIsUpToDate(sourceFile, objFile, depFile)
	NoError(t, err)
	require.False(t, upToDate)
}

func TestObjFileIsUpToDateDepIsOlder(t *testing.T) {
	sourceFile := tempFile(t, "source")
	defer os.RemoveAll(sourceFile)

	headerFile := tempFile(t, "header")
	defer os.RemoveAll(headerFile)

	sleep(t)

	objFile := tempFile(t, "obj")
	defer os.RemoveAll(objFile)
	depFile := tempFile(t, "dep")
	defer os.RemoveAll(depFile)

	utils.WriteFile(depFile, objFile+": \\\n\t"+sourceFile+" \\\n\t"+headerFile)

	upToDate, err := builder_utils.ObjFileIsUpToDate(sourceFile, objFile, depFile)
	NoError(t, err)
	require.True(t, upToDate)
}

func TestObjFileIsUpToDateDepIsWrong(t *testing.T) {
	sourceFile := tempFile(t, "source")
	defer os.RemoveAll(sourceFile)

	sleep(t)

	objFile := tempFile(t, "obj")
	defer os.RemoveAll(objFile)
	depFile := tempFile(t, "dep")
	defer os.RemoveAll(depFile)

	sleep(t)

	headerFile := tempFile(t, "header")
	defer os.RemoveAll(headerFile)

	utils.WriteFile(depFile, sourceFile+": \\\n\t"+sourceFile+" \\\n\t"+headerFile)

	upToDate, err := builder_utils.ObjFileIsUpToDate(sourceFile, objFile, depFile)
	NoError(t, err)
	require.False(t, upToDate)
}
