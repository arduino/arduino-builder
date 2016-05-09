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
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestUnusedCompiledLibrariesRemover(t *testing.T) {
	temp, err := ioutil.TempDir("", "test")
	NoError(t, err)
	defer os.RemoveAll(temp)

	NoError(t, os.MkdirAll(filepath.Join(temp, "SPI"), os.FileMode(0755)))
	NoError(t, os.MkdirAll(filepath.Join(temp, "Bridge"), os.FileMode(0755)))
	NoError(t, ioutil.WriteFile(filepath.Join(temp, "dummy_file"), []byte{}, os.FileMode(0644)))

	ctx := &types.Context{}
	ctx.LibrariesBuildPath = temp
	ctx.ImportedLibraries = []*types.Library{&types.Library{Name: "Bridge"}}

	cmd := builder.UnusedCompiledLibrariesRemover{}
	err = cmd.Run(ctx)
	NoError(t, err)

	_, err = os.Stat(filepath.Join(temp, "SPI"))
	require.Error(t, err)
	require.True(t, os.IsNotExist(err))
	_, err = os.Stat(filepath.Join(temp, "Bridge"))
	NoError(t, err)
	_, err = os.Stat(filepath.Join(temp, "dummy_file"))
	NoError(t, err)
}

func TestUnusedCompiledLibrariesRemoverLibDoesNotExist(t *testing.T) {
	ctx := &types.Context{}
	ctx.LibrariesBuildPath = filepath.Join(os.TempDir(), "test")
	ctx.ImportedLibraries = []*types.Library{&types.Library{Name: "Bridge"}}

	cmd := builder.UnusedCompiledLibrariesRemover{}
	err := cmd.Run(ctx)
	NoError(t, err)
}

func TestUnusedCompiledLibrariesRemoverNoUsedLibraries(t *testing.T) {
	temp, err := ioutil.TempDir("", "test")
	NoError(t, err)
	defer os.RemoveAll(temp)

	NoError(t, os.MkdirAll(filepath.Join(temp, "SPI"), os.FileMode(0755)))
	NoError(t, os.MkdirAll(filepath.Join(temp, "Bridge"), os.FileMode(0755)))
	NoError(t, ioutil.WriteFile(filepath.Join(temp, "dummy_file"), []byte{}, os.FileMode(0644)))

	ctx := &types.Context{}
	ctx.LibrariesBuildPath = temp
	ctx.ImportedLibraries = []*types.Library{}

	cmd := builder.UnusedCompiledLibrariesRemover{}
	err = cmd.Run(ctx)
	NoError(t, err)

	_, err = os.Stat(filepath.Join(temp, "SPI"))
	require.Error(t, err)
	require.True(t, os.IsNotExist(err))
	_, err = os.Stat(filepath.Join(temp, "Bridge"))
	require.Error(t, err)
	require.True(t, os.IsNotExist(err))
	_, err = os.Stat(filepath.Join(temp, "dummy_file"))
	NoError(t, err)
}
