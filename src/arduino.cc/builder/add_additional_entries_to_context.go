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
	"path/filepath"

	"arduino.cc/builder/constants"
	"arduino.cc/builder/i18n"
	"arduino.cc/builder/types"
)

type AddAdditionalEntriesToContext struct{}

func (s *AddAdditionalEntriesToContext) Run(ctx *types.Context) error {
	if ctx.BuildPath != "" {
		buildPath := ctx.BuildPath
		preprocPath, err := filepath.Abs(filepath.Join(buildPath, constants.FOLDER_PREPROC))
		if err != nil {
			return i18n.WrapError(err)
		}
		sketchBuildPath, err := filepath.Abs(filepath.Join(buildPath, constants.FOLDER_SKETCH))
		if err != nil {
			return i18n.WrapError(err)
		}
		librariesBuildPath, err := filepath.Abs(filepath.Join(buildPath, constants.FOLDER_LIBRARIES))
		if err != nil {
			return i18n.WrapError(err)
		}
		coreBuildPath, err := filepath.Abs(filepath.Join(buildPath, constants.FOLDER_CORE))
		if err != nil {
			return i18n.WrapError(err)
		}

		ctx.PreprocPath = preprocPath
		ctx.SketchBuildPath = sketchBuildPath
		ctx.LibrariesBuildPath = librariesBuildPath
		ctx.CoreBuildPath = coreBuildPath
	}

	if ctx.BuildCachePath != "" {
		coreBuildCachePath, err := filepath.Abs(filepath.Join(ctx.BuildCachePath, constants.FOLDER_CORE))
		if err != nil {
			return i18n.WrapError(err)
		}

		ctx.CoreBuildCachePath = coreBuildCachePath
	}

	if ctx.WarningsLevel == "" {
		ctx.WarningsLevel = DEFAULT_WARNINGS_LEVEL
	}

	ctx.CollectedSourceFiles = &types.UniqueSourceFileQueue{}

	ctx.LibrariesResolutionResults = make(map[string]types.LibraryResolutionResult)
	ctx.HardwareRewriteResults = make(map[*types.Platform][]types.PlatforKeyRewrite)

	return nil
}
