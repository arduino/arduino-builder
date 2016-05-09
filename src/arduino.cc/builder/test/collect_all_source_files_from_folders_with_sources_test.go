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
	"path/filepath"
	"sort"
	"testing"
)

func TestCollectAllSourceFilesFromFoldersWithSources(t *testing.T) {
	ctx := &types.Context{}

	sourceFiles := &types.UniqueStringQueue{}
	ctx.CollectedSourceFiles = sourceFiles
	foldersWithSources := &types.UniqueSourceFolderQueue{}
	foldersWithSources.Push(types.SourceFolder{Folder: Abs(t, "sketch_with_config"), Recurse: true})
	ctx.FoldersWithSourceFiles = foldersWithSources

	commands := []types.Command{
		&builder.CollectAllSourceFilesFromFoldersWithSources{},
	}

	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}

	require.Equal(t, 1, len(*sourceFiles))
	require.Equal(t, 0, len(*foldersWithSources))
	sort.Strings(*sourceFiles)

	require.Equal(t, Abs(t, filepath.Join("sketch_with_config", "includes", "de bug.cpp")), sourceFiles.Pop())
	require.Equal(t, 0, len(*sourceFiles))
}

func TestCollectAllSourceFilesFromFoldersWithSourcesOfLibrary(t *testing.T) {
	ctx := &types.Context{}

	sourceFiles := &types.UniqueStringQueue{}
	ctx.CollectedSourceFiles = sourceFiles
	foldersWithSources := &types.UniqueSourceFolderQueue{}
	foldersWithSources.Push(types.SourceFolder{Folder: Abs(t, filepath.Join("downloaded_libraries", "Bridge")), Recurse: true})
	ctx.FoldersWithSourceFiles = foldersWithSources

	commands := []types.Command{
		&builder.CollectAllSourceFilesFromFoldersWithSources{},
	}

	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}

	require.Equal(t, 9, len(*sourceFiles))
	require.Equal(t, 0, len(*foldersWithSources))
	sort.Strings(*sourceFiles)

	require.Equal(t, Abs(t, filepath.Join("downloaded_libraries", "Bridge", "src", "Bridge.cpp")), sourceFiles.Pop())
	require.Equal(t, Abs(t, filepath.Join("downloaded_libraries", "Bridge", "src", "BridgeClient.cpp")), sourceFiles.Pop())
	require.Equal(t, Abs(t, filepath.Join("downloaded_libraries", "Bridge", "src", "BridgeServer.cpp")), sourceFiles.Pop())
	require.Equal(t, Abs(t, filepath.Join("downloaded_libraries", "Bridge", "src", "BridgeUdp.cpp")), sourceFiles.Pop())
	require.Equal(t, Abs(t, filepath.Join("downloaded_libraries", "Bridge", "src", "Console.cpp")), sourceFiles.Pop())
	require.Equal(t, Abs(t, filepath.Join("downloaded_libraries", "Bridge", "src", "FileIO.cpp")), sourceFiles.Pop())
	require.Equal(t, Abs(t, filepath.Join("downloaded_libraries", "Bridge", "src", "HttpClient.cpp")), sourceFiles.Pop())
	require.Equal(t, Abs(t, filepath.Join("downloaded_libraries", "Bridge", "src", "Mailbox.cpp")), sourceFiles.Pop())
	require.Equal(t, Abs(t, filepath.Join("downloaded_libraries", "Bridge", "src", "Process.cpp")), sourceFiles.Pop())
	require.Equal(t, 0, len(*sourceFiles))
}

func TestCollectAllSourceFilesFromFoldersWithSourcesOfOldLibrary(t *testing.T) {
	ctx := &types.Context{}

	sourceFiles := &types.UniqueStringQueue{}
	ctx.CollectedSourceFiles = sourceFiles
	foldersWithSources := &types.UniqueSourceFolderQueue{}
	foldersWithSources.Push(types.SourceFolder{Folder: Abs(t, filepath.Join("libraries", "ShouldNotRecurseWithOldLibs")), Recurse: false})
	foldersWithSources.Push(types.SourceFolder{Folder: Abs(t, filepath.Join("libraries", "ShouldNotRecurseWithOldLibs", "utility")), Recurse: false})
	foldersWithSources.Push(types.SourceFolder{Folder: Abs(t, "non existent folder"), Recurse: false})
	ctx.FoldersWithSourceFiles = foldersWithSources

	commands := []types.Command{
		&builder.CollectAllSourceFilesFromFoldersWithSources{},
	}

	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}

	require.Equal(t, 2, len(*sourceFiles))
	require.Equal(t, 0, len(*foldersWithSources))
	sort.Strings(*sourceFiles)

	require.Equal(t, Abs(t, filepath.Join("libraries", "ShouldNotRecurseWithOldLibs", "ShouldNotRecurseWithOldLibs.cpp")), sourceFiles.Pop())
	require.Equal(t, Abs(t, filepath.Join("libraries", "ShouldNotRecurseWithOldLibs", "utility", "utils.cpp")), sourceFiles.Pop())
	require.Equal(t, 0, len(*sourceFiles))
}
