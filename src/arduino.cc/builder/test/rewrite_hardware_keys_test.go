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
	"arduino.cc/builder/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRewriteHardwareKeys(t *testing.T) {
	ctx := &types.Context{}

	packages := &types.Packages{}
	packages.Packages = make(map[string]*types.Package)
	aPackage := &types.Package{PackageId: "dummy"}
	packages.Packages["dummy"] = aPackage
	aPackage.Platforms = make(map[string]*types.Platform)

	platform := &types.Platform{PlatformId: "dummy"}
	aPackage.Platforms["dummy"] = platform
	platform.Properties = make(map[string]string)
	platform.Properties[constants.PLATFORM_NAME] = "A test platform"
	platform.Properties[constants.BUILD_PROPERTIES_COMPILER_PATH] = "{runtime.ide.path}/hardware/tools/avr/bin/"

	ctx.Hardware = packages

	rewrite := types.PlatforKeyRewrite{Key: constants.BUILD_PROPERTIES_COMPILER_PATH, OldValue: "{runtime.ide.path}/hardware/tools/avr/bin/", NewValue: "{runtime.tools.avr-gcc.path}/bin/"}
	platformKeysRewrite := types.PlatforKeysRewrite{Rewrites: []types.PlatforKeyRewrite{rewrite}}
	ctx.PlatformKeyRewrites = platformKeysRewrite

	commands := []types.Command{
		&builder.AddAdditionalEntriesToContext{},
		&builder.RewriteHardwareKeys{},
	}

	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}

	require.Equal(t, "{runtime.tools.avr-gcc.path}/bin/", platform.Properties[constants.BUILD_PROPERTIES_COMPILER_PATH])
}

func TestRewriteHardwareKeysWithRewritingDisabled(t *testing.T) {
	ctx := &types.Context{}

	packages := &types.Packages{}
	packages.Packages = make(map[string]*types.Package)
	aPackage := &types.Package{PackageId: "dummy"}
	packages.Packages["dummy"] = aPackage
	aPackage.Platforms = make(map[string]*types.Platform)

	platform := &types.Platform{PlatformId: "dummy"}
	aPackage.Platforms["dummy"] = platform
	platform.Properties = make(map[string]string)
	platform.Properties[constants.PLATFORM_NAME] = "A test platform"
	platform.Properties[constants.BUILD_PROPERTIES_COMPILER_PATH] = "{runtime.ide.path}/hardware/tools/avr/bin/"
	platform.Properties[constants.REWRITING] = constants.REWRITING_DISABLED

	ctx.Hardware = packages

	rewrite := types.PlatforKeyRewrite{Key: constants.BUILD_PROPERTIES_COMPILER_PATH, OldValue: "{runtime.ide.path}/hardware/tools/avr/bin/", NewValue: "{runtime.tools.avr-gcc.path}/bin/"}
	platformKeysRewrite := types.PlatforKeysRewrite{Rewrites: []types.PlatforKeyRewrite{rewrite}}

	ctx.PlatformKeyRewrites = platformKeysRewrite

	commands := []types.Command{
		&builder.AddAdditionalEntriesToContext{},
		&builder.RewriteHardwareKeys{},
	}

	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}

	require.Equal(t, "{runtime.ide.path}/hardware/tools/avr/bin/", platform.Properties[constants.BUILD_PROPERTIES_COMPILER_PATH])
}
