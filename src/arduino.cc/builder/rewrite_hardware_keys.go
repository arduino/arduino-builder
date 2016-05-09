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
	"arduino.cc/builder/types"
)

type RewriteHardwareKeys struct{}

func (s *RewriteHardwareKeys) Run(ctx *types.Context) error {
	if ctx.PlatformKeyRewrites.Empty() {
		return nil
	}

	packages := ctx.Hardware
	platformKeysRewrite := ctx.PlatformKeyRewrites
	hardwareRewriteResults := ctx.HardwareRewriteResults

	for _, aPackage := range packages.Packages {
		for _, platform := range aPackage.Platforms {
			if platform.Properties[constants.REWRITING] != constants.REWRITING_DISABLED {
				for _, rewrite := range platformKeysRewrite.Rewrites {
					if platform.Properties[rewrite.Key] != constants.EMPTY_STRING && platform.Properties[rewrite.Key] == rewrite.OldValue {
						platform.Properties[rewrite.Key] = rewrite.NewValue
						appliedRewrites := rewritesAppliedToPlatform(platform, hardwareRewriteResults)
						appliedRewrites = append(appliedRewrites, rewrite)
						hardwareRewriteResults[platform] = appliedRewrites
					}
				}
			}
		}
	}

	return nil
}

func rewritesAppliedToPlatform(platform *types.Platform, hardwareRewriteResults map[*types.Platform][]types.PlatforKeyRewrite) []types.PlatforKeyRewrite {
	if hardwareRewriteResults[platform] == nil {
		hardwareRewriteResults[platform] = []types.PlatforKeyRewrite{}
	}
	return hardwareRewriteResults[platform]
}
