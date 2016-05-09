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
	"arduino.cc/builder/gohasissues"
	"arduino.cc/builder/types"
	"arduino.cc/builder/utils"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestWipeoutBuildPathIfBuildOptionsChanged(t *testing.T) {
	ctx := &types.Context{}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	ctx.BuildOptionsJsonPrevious = "{ \"old\":\"old\" }"
	ctx.BuildOptionsJson = "{ \"new\":\"new\" }"

	utils.TouchFile(filepath.Join(buildPath, "should_be_deleted.txt"))

	commands := []types.Command{
		&builder.WipeoutBuildPathIfBuildOptionsChanged{},
	}

	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}

	_, err := os.Stat(buildPath)
	NoError(t, err)

	files, err := gohasissues.ReadDir(buildPath)
	NoError(t, err)
	require.Equal(t, 0, len(files))

	_, err = os.Stat(filepath.Join(buildPath, "should_be_deleted.txt"))
	require.Error(t, err)
}

func TestWipeoutBuildPathIfBuildOptionsChangedNoPreviousBuildOptions(t *testing.T) {
	ctx := &types.Context{}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	ctx.BuildOptionsJson = "{ \"new\":\"new\" }"

	utils.TouchFile(filepath.Join(buildPath, "should_not_be_deleted.txt"))

	commands := []types.Command{
		&builder.WipeoutBuildPathIfBuildOptionsChanged{},
	}

	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}

	_, err := os.Stat(buildPath)
	NoError(t, err)

	files, err := gohasissues.ReadDir(buildPath)
	NoError(t, err)
	require.Equal(t, 1, len(files))

	_, err = os.Stat(filepath.Join(buildPath, "should_not_be_deleted.txt"))
	NoError(t, err)
}

func TestWipeoutBuildPathIfBuildOptionsChangedBuildOptionsMatch(t *testing.T) {
	ctx := &types.Context{}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	ctx.BuildOptionsJsonPrevious = "{ \"old\":\"old\" }"
	ctx.BuildOptionsJson = "{ \"old\":\"old\" }"

	utils.TouchFile(filepath.Join(buildPath, "should_not_be_deleted.txt"))

	commands := []types.Command{
		&builder.WipeoutBuildPathIfBuildOptionsChanged{},
	}

	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}

	_, err := os.Stat(buildPath)
	NoError(t, err)

	files, err := gohasissues.ReadDir(buildPath)
	NoError(t, err)
	require.Equal(t, 1, len(files))

	_, err = os.Stat(filepath.Join(buildPath, "should_not_be_deleted.txt"))
	NoError(t, err)
}
