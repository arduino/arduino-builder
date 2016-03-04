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
 * Copyright 2015 Matthijs Kooijman
 */

package builder

import (
	"arduino.cc/builder/constants"
	"arduino.cc/builder/types"
	"regexp"
	"strings"
)

type SketchSourceMerger struct{}

func (s *SketchSourceMerger) Run(context map[string]interface{}) error {
	sketch := context[constants.CTX_SKETCH].(*types.Sketch)

	lineOffset := 0
	includeSection := ""
	if !sketchIncludesArduinoH(&sketch.MainFile) {
		includeSection += "#include <Arduino.h>\n"
		lineOffset++
	}
	includeSection += "#line 1 \"" + strings.Replace((&sketch.MainFile).Name, "\\", "\\\\", -1) + "\"\n"
	lineOffset++
	context[constants.CTX_INCLUDE_SECTION] = includeSection

	source := includeSection
	source += addSourceWrappedWithLineDirective(&sketch.MainFile)
	lineOffset += 1
	for _, file := range sketch.OtherSketchFiles {
		source += addSourceWrappedWithLineDirective(&file)
	}

	context[constants.CTX_LINE_OFFSET] = lineOffset
	context[constants.CTX_SOURCE] = source

	return nil
}

func sketchIncludesArduinoH(sketch *types.SketchFile) bool {
	if matched, err := regexp.MatchString("(?m)^\\s*#\\s*include\\s*[<\"]Arduino\\.h[>\"]", sketch.Source); err != nil {
		panic(err)
	} else {
		return matched
	}
}

func addSourceWrappedWithLineDirective(sketch *types.SketchFile) string {
	source := "#line 1 \"" + strings.Replace(sketch.Name, "\\", "\\\\", -1) + "\"\n"
	source += sketch.Source
	source += "\n"

	return source
}
