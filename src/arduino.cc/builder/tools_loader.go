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

package builder

import (
	"arduino.cc/builder/constants"
	"arduino.cc/builder/gohasissues"
	"arduino.cc/builder/i18n"
	"arduino.cc/builder/types"
	"arduino.cc/builder/utils"
	"os"
	"path/filepath"
	"strings"
)

type ToolsLoader struct{}

func (s *ToolsLoader) Run(ctx *types.Context) error {
	folders := ctx.ToolsFolders

	tools := []*types.Tool{}

	for _, folder := range folders {
		builtinToolsVersionsFile, err := findBuiltinToolsVersionsFile(folder)
		if err != nil {
			return i18n.WrapError(err)
		}
		if builtinToolsVersionsFile != constants.EMPTY_STRING {
			err = loadToolsFrom(&tools, builtinToolsVersionsFile)
			if err != nil {
				return i18n.WrapError(err)
			}
		} else {
			subfolders, err := collectAllToolsFolders(folder)
			if err != nil {
				return i18n.WrapError(err)
			}

			for _, subfolder := range subfolders {
				err = loadToolsFromFolderStructure(&tools, subfolder)
				if err != nil {
					return i18n.WrapError(err)
				}
			}
		}
	}

	ctx.Tools = tools

	return nil
}

func collectAllToolsFolders(from string) ([]string, error) {
	folders := []string{}
	walkFunc := func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if !info.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(from, currentPath)
		if err != nil {
			return i18n.WrapError(err)
		}
		depth := len(strings.Split(rel, string(os.PathSeparator)))

		if info.Name() == constants.FOLDER_TOOLS && depth == 2 {
			folders = append(folders, currentPath)
		} else if depth > 2 {
			return filepath.SkipDir
		}

		return nil
	}
	err := gohasissues.Walk(from, walkFunc)

	if len(folders) == 0 {
		folders = append(folders, from)
	}

	return folders, i18n.WrapError(err)
}

func toolsSliceContains(tools *[]*types.Tool, name, version string) bool {
	for _, tool := range *tools {
		if name == tool.Name && version == tool.Version {
			return true
		}
	}
	return false
}

func loadToolsFrom(tools *[]*types.Tool, builtinToolsVersionsFilePath string) error {
	rows, err := utils.ReadFileToRows(builtinToolsVersionsFilePath)
	if err != nil {
		return i18n.WrapError(err)
	}

	folder, err := filepath.Abs(filepath.Dir(builtinToolsVersionsFilePath))
	if err != nil {
		return i18n.WrapError(err)
	}

	for _, row := range rows {
		row = strings.TrimSpace(row)
		if row != constants.EMPTY_STRING {
			rowParts := strings.Split(row, "=")
			toolName := strings.Split(rowParts[0], ".")[1]
			toolVersion := rowParts[1]
			if !toolsSliceContains(tools, toolName, toolVersion) {
				*tools = append(*tools, &types.Tool{Name: toolName, Version: toolVersion, Folder: folder})
			}
		}
	}

	return nil
}

func findBuiltinToolsVersionsFile(folder string) (string, error) {
	builtinToolsVersionsFilePath := constants.EMPTY_STRING
	findBuiltInToolsVersionsTxt := func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if builtinToolsVersionsFilePath != constants.EMPTY_STRING {
			return nil
		}
		if filepath.Base(currentPath) == constants.FILE_BUILTIN_TOOLS_VERSIONS_TXT {
			builtinToolsVersionsFilePath = currentPath
		}
		return nil
	}
	err := gohasissues.Walk(folder, findBuiltInToolsVersionsTxt)
	return builtinToolsVersionsFilePath, i18n.WrapError(err)
}

func loadToolsFromFolderStructure(tools *[]*types.Tool, folder string) error {
	toolsNames, err := utils.ReadDirFiltered(folder, utils.FilterDirs)
	if err != nil {
		return i18n.WrapError(err)
	}
	for _, toolName := range toolsNames {
		toolVersions, err := utils.ReadDirFiltered(filepath.Join(folder, toolName.Name()), utils.FilterDirs)
		if err != nil {
			return i18n.WrapError(err)
		}
		for _, toolVersion := range toolVersions {
			toolFolder, err := filepath.Abs(filepath.Join(folder, toolName.Name(), toolVersion.Name()))
			if err != nil {
				return i18n.WrapError(err)
			}
			if !toolsSliceContains(tools, toolName.Name(), toolVersion.Name()) {
				*tools = append(*tools, &types.Tool{Name: toolName.Name(), Version: toolVersion.Name(), Folder: toolFolder})
			}
		}
	}

	return nil
}
