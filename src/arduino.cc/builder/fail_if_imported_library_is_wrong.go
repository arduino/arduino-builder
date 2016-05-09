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
	"arduino.cc/builder/i18n"
	"arduino.cc/builder/types"
	"os"
	"path/filepath"
)

type FailIfImportedLibraryIsWrong struct{}

func (s *FailIfImportedLibraryIsWrong) Run(ctx *types.Context) error {
	if len(ctx.ImportedLibraries) == 0 {
		return nil
	}

	logger := ctx.GetLogger()

	for _, library := range ctx.ImportedLibraries {
		if !library.IsLegacy {
			if stat, err := os.Stat(filepath.Join(library.Folder, constants.LIBRARY_FOLDER_ARCH)); err == nil && stat.IsDir() {
				return i18n.ErrorfWithLogger(logger, constants.MSG_ARCH_FOLDER_NOT_SUPPORTED)
			}
			for _, propName := range LIBRARY_MANDATORY_PROPERTIES {
				if _, ok := library.Properties[propName]; !ok {
					return i18n.ErrorfWithLogger(logger, constants.MSG_PROP_IN_LIBRARY, propName, library.Folder)
				}
			}
			if library.Layout == types.LIBRARY_RECURSIVE {
				if stat, err := os.Stat(filepath.Join(library.Folder, constants.LIBRARY_FOLDER_UTILITY)); err == nil && stat.IsDir() {
					return i18n.ErrorfWithLogger(logger, constants.MSG_LIBRARY_CAN_USE_SRC_AND_UTILITY_FOLDERS, library.Folder)
				}
			}
		}
	}

	return nil
}
