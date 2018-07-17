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

package phases

import (
	"path/filepath"

	"github.com/arduino/arduino-builder/builder_utils"
	"github.com/arduino/arduino-builder/constants"
	"github.com/arduino/arduino-builder/i18n"
	"github.com/arduino/arduino-builder/types"
	"github.com/arduino/arduino-builder/utils"
	"github.com/arduino/go-properties-map"
)

type CoreBuilder struct{}

func (s *CoreBuilder) Run(ctx *types.Context) error {
	coreBuildPath := ctx.CoreBuildPath
	coreBuildCachePath := ctx.CoreBuildCachePath
	buildProperties := ctx.BuildProperties

	err := utils.EnsureFolderExists(coreBuildPath)
	if err != nil {
		return i18n.WrapError(err)
	}

	if coreBuildCachePath != "" {
		err := utils.EnsureFolderExists(coreBuildCachePath)
		if err != nil {
			return i18n.WrapError(err)
		}
	}

	archiveFile, objectFiles, err := compileCore(ctx, coreBuildPath, coreBuildCachePath, buildProperties)
	if err != nil {
		return i18n.WrapError(err)
	}

	ctx.CoreArchiveFilePath = archiveFile
	ctx.CoreObjectsFiles = objectFiles

	return nil
}

func compileCore(ctx *types.Context, buildPath string, buildCachePath string, buildProperties properties.Map) (string, []string, error) {
	logger := ctx.GetLogger()
	coreFolder := buildProperties[constants.BUILD_PROPERTIES_BUILD_CORE_PATH]
	variantFolder := buildProperties[constants.BUILD_PROPERTIES_BUILD_VARIANT_PATH]

	targetCoreFolder := buildProperties[constants.BUILD_PROPERTIES_RUNTIME_PLATFORM_PATH]

	includes := []string{}
	includes = append(includes, coreFolder)
	if variantFolder != constants.EMPTY_STRING {
		includes = append(includes, variantFolder)
	}
	includes = utils.Map(includes, utils.WrapWithHyphenI)

	var err error

	variantObjectFiles := []string{}
	if variantFolder != constants.EMPTY_STRING {
		variantObjectFiles, err = builder_utils.CompileFiles(ctx, variantObjectFiles, variantFolder, true, buildPath, buildProperties, includes)
		if err != nil {
			return "", nil, i18n.WrapError(err)
		}
	}

	// Recreate the archive if ANY of the core files (including platform.txt) has changed
	realCoreFolder := utils.GetParentFolder(coreFolder, 2)

	var targetArchivedCore string
	if buildCachePath != "" {
		archivedCoreName := builder_utils.GetCachedCoreArchiveFileName(buildProperties[constants.BUILD_PROPERTIES_FQBN], realCoreFolder)
		targetArchivedCore = filepath.Join(buildCachePath, archivedCoreName)
		canUseArchivedCore := !builder_utils.CoreOrReferencedCoreHasChanged(realCoreFolder, targetCoreFolder, targetArchivedCore)

		if canUseArchivedCore {
			// use archived core
			if ctx.Verbose {
				logger.Println(constants.LOG_LEVEL_INFO, "Using precompiled core: {0}", targetArchivedCore)
			}
			return targetArchivedCore, variantObjectFiles, nil
		}
	}

	coreObjectFiles, err := builder_utils.CompileFiles(ctx, []string{}, coreFolder, true, buildPath, buildProperties, includes)
	if err != nil {
		return "", nil, i18n.WrapError(err)
	}

	archiveFile, err := builder_utils.ArchiveCompiledFiles(ctx, buildPath, "core.a", coreObjectFiles, buildProperties)
	if err != nil {
		return "", nil, i18n.WrapError(err)
	}

	// archive core.a
	if targetArchivedCore != "" {
		if ctx.Verbose {
			logger.Println(constants.LOG_LEVEL_INFO, constants.MSG_ARCHIVING_CORE_CACHE, targetArchivedCore)
		}
		builder_utils.CopyFile(archiveFile, targetArchivedCore)
	}

	return archiveFile, variantObjectFiles, nil
}
