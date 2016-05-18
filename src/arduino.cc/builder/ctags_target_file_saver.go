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
	"arduino.cc/builder/i18n"
	"arduino.cc/builder/types"
	"arduino.cc/builder/utils"
	"path/filepath"
	"strconv"
	"strings"
)

type CTagsTargetFileSaver struct {
	Source         *string
	TargetFileName string
}

func (s *CTagsTargetFileSaver) Run(ctx *types.Context) error {
	source := *s.Source

	preprocPath := ctx.PreprocPath
	err := utils.EnsureFolderExists(preprocPath)
	if err != nil {
		return i18n.WrapError(err)
	}

	// drop every line which is not part of user sketch or a define/ifdef/endif/etc
	var searchSlice []string
	searchSlice = append(searchSlice, filepath.Dir(ctx.SketchLocation))
	searchSlice = append(searchSlice, filepath.Dir(ctx.BuildPath))
	source = saveLinesContainingDirectivesAndSketch(source, searchSlice)

	ctagsTargetFilePath := filepath.Join(preprocPath, s.TargetFileName)
	err = utils.WriteFile(ctagsTargetFilePath, source)
	if err != nil {
		return i18n.WrapError(err)
	}

	ctx.CTagsTargetFile = ctagsTargetFilePath

	return nil
}

func saveLinesContainingDirectivesAndSketch(src string, tofind []string) string {
	lines := strings.Split(src, "\n")

	saveLine := false
	minimizedString := ""

	for _, line := range lines {
		if saveLine || startsWithHashtag(line) {
			minimizedString += line + "\n"
		}
		if containsAny(line, tofind) && isLineMarker(line) {
			saveLine = true
		}
		if saveLine && !containsAny(line, tofind) && isLineMarker(line) {
			saveLine = false
		}
	}
	return minimizedString
}

func containsAny(src string, tofind []string) bool {
	for _, str := range tofind {
		if strings.Contains(src, str) {
			return true
		}
	}
	return false
}

func startsWithHashtag(src string) bool {
	trimmedStr := strings.TrimSpace(src)
	if len(trimmedStr) > 0 && trimmedStr[0] == '#' {
		return true
	}
	return false
}

func isLineMarker(src string) bool {
	trimmedStr := strings.TrimSpace(src)
	splittedStr := strings.Split(trimmedStr, " ")
	if len(splittedStr) > 2 && splittedStr[0][0] == '#' {
		_, err := strconv.Atoi(splittedStr[1])
		if err == nil {
			return true
		}
	}
	return false
}
