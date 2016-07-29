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
	"github.com/arduino/arduino-builder/constants"
	"github.com/arduino/arduino-builder/json_package_index"
	"github.com/arduino/arduino-builder/types"
)

type AddMissingBuildPropertiesFromParentPlatformTxtFiles struct{}

func (s *AddMissingBuildPropertiesFromParentPlatformTxtFiles) Run(ctx *types.Context) error {
	packages := ctx.Hardware
	targetPackage := ctx.TargetPackage
	buildProperties := ctx.BuildProperties

	newBuildProperties := packages.Properties.Clone()
	newBuildProperties.Merge(targetPackage.Properties)
	newBuildProperties.Merge(buildProperties)

	ctx.BuildProperties = newBuildProperties

	return nil
}

type OverridePropertiesWithJsonInfo struct{}

func (s *OverridePropertiesWithJsonInfo) Run(ctx *types.Context) error {

	if ctx.JsonFolders != nil {

		jsonProperties, err := json_package_index.PackageIndexFoldersToPropertiesMap(ctx.JsonFolders)

		if err != nil {
			// doesn't matter, log the broken package in verbose mode
		}

		newBuildProperties := jsonProperties[ctx.TargetPackage.PackageId+":"+ctx.TargetPlatform.PlatformId+":"+ctx.TargetPlatform.Properties["version"]]

		buildProperties := ctx.BuildProperties.Clone()

		buildProperties.Merge(newBuildProperties)

		// HACK!!! To overcome AVR core 1.6.12 lto problems, replace avr-gcc-4.8.1-arduino5 with
		// 4.9.2-atmel3.5.3-arduino2 if it exists
		if buildProperties[constants.HACK_PROPERTIES_AVR_GCC_NEW] != "" {
			buildProperties[constants.HACK_PROPERTIES_AVR_GCC_OLD] =
				"{" + constants.HACK_PROPERTIES_AVR_GCC_NEW + "}"
		}

		ctx.BuildProperties = buildProperties
	}

	return nil
}
