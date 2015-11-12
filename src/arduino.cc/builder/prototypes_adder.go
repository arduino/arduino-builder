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
	"arduino.cc/builder/utils"
	"strconv"
	"strings"
)

type PrototypesAdder struct{}

func (s *PrototypesAdder) Run(context map[string]interface{}) error {
	source := context[constants.CTX_SOURCE].(string)
	sourceRows := strings.Split(source, "\n")

	if !utils.MapHas(context, constants.CTX_LINE_WHERE_TO_INSERT_PROTOTYPES) {
		return nil
	}

	firstFunctionLine := context[constants.CTX_LINE_WHERE_TO_INSERT_PROTOTYPES].(int)
	if firstFunctionOutsideOfSource(firstFunctionLine, sourceRows) {
		return nil
	}

	firstFunctionChar := len(strings.Join(sourceRows[:firstFunctionLine+context[constants.CTX_LINE_OFFSET].(int)-1], "\n")) + 1
	prototypeSection := composePrototypeSection(firstFunctionLine, context[constants.CTX_PROTOTYPES].([]*types.Prototype))
	context[constants.CTX_PROTOTYPE_SECTION] = prototypeSection
	source = source[:firstFunctionChar] + prototypeSection + source[firstFunctionChar:]

	context[constants.CTX_SOURCE] = source

	return nil
}

func composePrototypeSection(line int, prototypes []*types.Prototype) string {
	if len(prototypes) == 0 {
		return constants.EMPTY_STRING
	}

	str := joinPrototypes(prototypes)
	str += "\n#line "
	str += strconv.Itoa(line)
	str += "\n"

	return str
}

func joinPrototypes(prototypes []*types.Prototype) string {
	prototypesSlice := []string{}
	for _, proto := range prototypes {
		prototypesSlice = append(prototypesSlice, "#line "+proto.Line)
		prototypeParts := []string{}
		if proto.Modifiers != "" {
			prototypeParts = append(prototypeParts, proto.Modifiers)
		}
		prototypeParts = append(prototypeParts, proto.Prototype)
		prototypesSlice = append(prototypesSlice, strings.Join(prototypeParts, " "))
	}
	return strings.Join(prototypesSlice, "\n")
}

func firstFunctionOutsideOfSource(firstFunctionLine int, sourceRows []string) bool {
	return firstFunctionLine > len(sourceRows)-1
}
