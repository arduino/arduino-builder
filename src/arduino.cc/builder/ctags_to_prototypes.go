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

type CTagsToPrototypes struct{}

func (s *CTagsToPrototypes) Run(context map[string]interface{}) error {
	tags := context[constants.CTX_COLLECTED_CTAGS].([]map[string]string)

	lineWhereToInsertPrototypes, err := findLineWhereToInsertPrototypes(tags)
	if err != nil {
		return utils.WrapError(err)
	}
	if lineWhereToInsertPrototypes != -1 {
		context[constants.CTX_LINE_WHERE_TO_INSERT_PROTOTYPES] = lineWhereToInsertPrototypes
	}

	prototypes := toPrototypes(tags)
	context[constants.CTX_PROTOTYPES] = prototypes

	return nil
}

func findLineWhereToInsertPrototypes(tags []map[string]string) (int, error) {
	firstFunctionLine, err := firstFunctionAtLine(tags)
	if err != nil {
		return -1, utils.WrapError(err)
	}
	firstFunctionPointerAsArgument, err := firstFunctionPointerUsedAsArgument(tags)
	if err != nil {
		return -1, utils.WrapError(err)
	}
	if firstFunctionLine != -1 && firstFunctionPointerAsArgument != -1 {
		if firstFunctionLine < firstFunctionPointerAsArgument {
			return firstFunctionLine, nil
		} else {
			return firstFunctionPointerAsArgument, nil
		}
	} else if firstFunctionLine == -1 {
		return firstFunctionPointerAsArgument, nil
	} else {
		return firstFunctionLine, nil
	}
}

func firstFunctionPointerUsedAsArgument(tags []map[string]string) (int, error) {
	functionNames := collectFunctionNames(tags)
	for _, tag := range tags {
		if functionNameUsedAsFunctionPointerIn(tag, functionNames) {
			return strconv.Atoi(tag[FIELD_LINE])
		}
	}
	return -1, nil
}

func functionNameUsedAsFunctionPointerIn(tag map[string]string, functionNames []string) bool {
	for _, functionName := range functionNames {
		if strings.Index(tag[FIELD_CODE], "&"+functionName) != -1 {
			return true
		}
	}
	return false
}

func collectFunctionNames(tags []map[string]string) []string {
	names := []string{}
	for _, tag := range tags {
		if tag[FIELD_KIND] == KIND_FUNCTION {
			names = append(names, tag[FIELD_FUNCTION_NAME])
		}
	}
	return names
}

func firstFunctionAtLine(tags []map[string]string) (int, error) {
	for _, tag := range tags {
		_, tagHasAtLeastOneField := utils.TagHasAtLeastOneField(tag, FIELDS_MARKING_UNHANDLED_TAGS)
		if !tagIsUnknown(tag) && !tagHasAtLeastOneField && tag[FIELD_KIND] == KIND_FUNCTION {
			return strconv.Atoi(tag[FIELD_LINE])
		}
	}
	return -1, nil
}

func toPrototypes(tags []map[string]string) []*types.Prototype {
	prototypes := []*types.Prototype{}
	for _, tag := range tags {
		if tag[FIELD_SKIP] != TRUE {
			prototype := &types.Prototype{FunctionName: tag[FIELD_FUNCTION_NAME], File: tag[FIELD_FILENAME], Prototype: tag[KIND_PROTOTYPE], Modifiers: tag[KIND_PROTOTYPE_MODIFIERS], Line: tag[FIELD_LINE], Fields: tag}
			prototypes = append(prototypes, prototype)
		}
	}
	return prototypes
}
