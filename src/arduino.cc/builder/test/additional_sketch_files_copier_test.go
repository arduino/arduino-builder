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
	"arduino.cc/builder/gohasissues"
	"arduino.cc/builder/types"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"
)

type ByFileInfoName []os.FileInfo

func (s ByFileInfoName) Len() int {
	return len(s)
}
func (s ByFileInfoName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByFileInfoName) Less(i, j int) bool {
	return s[i].Name() < s[j].Name()
}

func TestCopyOtherFiles(t *testing.T) {
	ctx := &types.Context{
		SketchLocation: filepath.Join("sketch1", "sketch.ino"),
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	commands := []types.Command{
		&builder.AddAdditionalEntriesToContext{},
		&builder.SketchLoader{},
		&builder.AdditionalSketchFilesCopier{},
	}

	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}

	_, err1 := os.Stat(filepath.Join(buildPath, constants.FOLDER_SKETCH, "header.h"))
	NoError(t, err1)

	files, err1 := gohasissues.ReadDir(filepath.Join(buildPath, constants.FOLDER_SKETCH))
	NoError(t, err1)
	require.Equal(t, 3, len(files))

	sort.Sort(ByFileInfoName(files))
	require.Equal(t, "header.h", files[0].Name())
	require.Equal(t, "s_file.S", files[1].Name())
	require.Equal(t, "subfolder", files[2].Name())

	files, err1 = gohasissues.ReadDir(filepath.Join(buildPath, constants.FOLDER_SKETCH, "subfolder"))
	NoError(t, err1)
	require.Equal(t, 1, len(files))
	require.Equal(t, "helper.h", files[0].Name())
}

func TestCopyOtherFilesOnlyIfChanged(t *testing.T) {
	ctx := &types.Context{
		SketchLocation: filepath.Join("sketch1", "sketch.ino"),
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	commands := []types.Command{
		&builder.AddAdditionalEntriesToContext{},
		&builder.SketchLoader{},
		&builder.AdditionalSketchFilesCopier{},
	}

	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}

	headerStatBefore, err := os.Stat(filepath.Join(buildPath, constants.FOLDER_SKETCH, "header.h"))
	NoError(t, err)

	time.Sleep(2 * time.Second)

	ctx = &types.Context{
		SketchLocation: filepath.Join("sketch1", "sketch.ino"),
		BuildPath:      buildPath,
	}

	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}

	headerStatAfter, err := os.Stat(filepath.Join(buildPath, constants.FOLDER_SKETCH, "header.h"))
	NoError(t, err)

	require.Equal(t, headerStatBefore.ModTime().Unix(), headerStatAfter.ModTime().Unix())
}
