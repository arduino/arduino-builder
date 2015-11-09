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

const FIELD_KIND = "kind"
const FIELD_LINE = "line"
const FIELD_SIGNATURE = "signature"
const FIELD_RETURNTYPE = "returntype"
const FIELD_CODE = "code"
const FIELD_FUNCTION_NAME = "functionName"
const FIELD_CLASS = "class"
const FIELD_STRUCT = "struct"
const FIELD_NAMESPACE = "namespace"
const FIELD_SKIP = "skipMe"

const KIND_PROTOTYPE = "prototype"
const KIND_FUNCTION = "function"
const KIND_PROTOTYPE_MODIFIERS = "prototype_modifiers"

const TEMPLATE = "template"
const STATIC = "static"
const TRUE = "true"

var FIELDS = map[string]bool{"kind": true, "line": true, "typeref": true, "signature": true, "returntype": true, "class": true, "struct": true, "namespace": true}
var KNOWN_TAG_KINDS = map[string]bool{"prototype": true, "function": true}
var FIELDS_MARKING_UNHANDLED_TAGS = []string{FIELD_CLASS, FIELD_STRUCT, FIELD_NAMESPACE}

type CTagsParser struct {
	PrototypesField string
}

func (s *CTagsParser) Run(context map[string]interface{}) error {
	rows := strings.Split(context[constants.CTX_CTAGS_OUTPUT].(string), "\n")

	rows = removeEmpty(rows)

	var tags []map[string]string
	for _, row := range rows {
		tags = append(tags, parseTag(row))
	}

	skipTagsWhere(tags, tagIsUnknown)
	skipTagsWithField(tags, FIELDS_MARKING_UNHANDLED_TAGS)
	skipTagsWhere(tags, signatureContainsDefaultArg)
	addPrototypes(tags)
	removeDefinedProtypes(tags)
	removeDuplicate(tags)
	skipTagsWhere(tags, prototypeAndCodeDontMatch)

	lineWhereToInsertPrototypes, err := findLineWhereToInsertPrototypes(tags)
	if err != nil {
		return utils.WrapError(err)
	}
	if lineWhereToInsertPrototypes != -1 {
		context[constants.CTX_LINE_WHERE_TO_INSERT_PROTOTYPES] = lineWhereToInsertPrototypes
	}

	prototypes := toPrototypes(tags)

	context[s.PrototypesField] = prototypes

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
		if !tagIsUnknown(tag) && !tagHasAtLeastOneField(tag, FIELDS_MARKING_UNHANDLED_TAGS) && tag[FIELD_KIND] == KIND_FUNCTION {
			return strconv.Atoi(tag[FIELD_LINE])
		}
	}
	return -1, nil
}

func toPrototypes(tags []map[string]string) []*types.Prototype {
	prototypes := []*types.Prototype{}
	for _, tag := range tags {
		if tag[FIELD_SKIP] != TRUE {
			ctag := types.Prototype{FunctionName: tag[FIELD_FUNCTION_NAME], Prototype: tag[KIND_PROTOTYPE], Modifiers: tag[KIND_PROTOTYPE_MODIFIERS], Fields: tag}
			prototypes = append(prototypes, &ctag)
		}
	}
	return prototypes
}

func addPrototypes(tags []map[string]string) {
	for _, tag := range tags {
		if tag[FIELD_SKIP] != TRUE {
			addPrototype(tag)
		}
	}
}

func addPrototype(tag map[string]string) {
	if strings.Index(tag[FIELD_RETURNTYPE], TEMPLATE) == 0 || strings.Index(tag[FIELD_CODE], TEMPLATE) == 0 {
		code := tag[FIELD_CODE]
		if strings.Contains(code, "{") {
			code = code[:strings.Index(code, "{")]
		} else {
			code = code[:strings.LastIndex(code, ")")+1]
		}
		tag[KIND_PROTOTYPE] = code + ";"
		return
	}

	tag[KIND_PROTOTYPE] = tag[FIELD_RETURNTYPE] + " " + tag[FIELD_FUNCTION_NAME] + tag[FIELD_SIGNATURE] + ";"

	tag[KIND_PROTOTYPE_MODIFIERS] = ""
	if strings.Index(tag[FIELD_CODE], STATIC+" ") != -1 {
		tag[KIND_PROTOTYPE_MODIFIERS] = tag[KIND_PROTOTYPE_MODIFIERS] + " " + STATIC
	}
	tag[KIND_PROTOTYPE_MODIFIERS] = strings.TrimSpace(tag[KIND_PROTOTYPE_MODIFIERS])
}

func removeDefinedProtypes(tags []map[string]string) {
	definedPrototypes := make(map[string]bool)
	for _, tag := range tags {
		if tag[FIELD_KIND] == KIND_PROTOTYPE {
			definedPrototypes[tag[KIND_PROTOTYPE]] = true
		}
	}

	for _, tag := range tags {
		if definedPrototypes[tag[KIND_PROTOTYPE]] {
			tag[FIELD_SKIP] = TRUE
		}
	}
}

func removeDuplicate(tags []map[string]string) {
	definedPrototypes := make(map[string]bool)

	for _, tag := range tags {
		if !definedPrototypes[tag[KIND_PROTOTYPE]] {
			definedPrototypes[tag[KIND_PROTOTYPE]] = true
		} else {
			tag[FIELD_SKIP] = TRUE
		}
	}
}

type skipFuncType func(tag map[string]string) bool

func skipTagsWhere(tags []map[string]string, skipFuncs ...skipFuncType) {
	for _, tag := range tags {
		if tag[FIELD_SKIP] != TRUE {
			skip := skipFuncs[0](tag)
			for _, skipFunc := range skipFuncs[1:] {
				skip = skip || skipFunc(tag)
			}
			tag[FIELD_SKIP] = strconv.FormatBool(skip)
		}
	}
}

func signatureContainsDefaultArg(tag map[string]string) bool {
	return strings.Contains(tag[FIELD_SIGNATURE], "=")
}

func prototypeAndCodeDontMatch(tag map[string]string) bool {
	if tag[FIELD_SKIP] == TRUE {
		return true
	}

	code := removeSpacesAndTabs(tag[FIELD_CODE])
	prototype := removeSpacesAndTabs(tag[KIND_PROTOTYPE])
	prototype = removeTralingSemicolon(prototype)

	return strings.Index(code, prototype) == -1
}

func removeTralingSemicolon(s string) string {
	return s[0 : len(s)-1]
}

func removeSpacesAndTabs(s string) string {
	s = strings.Replace(s, " ", "", -1)
	s = strings.Replace(s, "\t", "", -1)
	return s
}

func skipTagsWithField(tags []map[string]string, fields []string) {
	for _, tag := range tags {
		if tagHasAtLeastOneField(tag, fields) {
			tag[FIELD_SKIP] = TRUE
		}
	}
}

func tagHasAtLeastOneField(tag map[string]string, fields []string) bool {
	for _, field := range fields {
		if tag[field] != constants.EMPTY_STRING {
			return true
		}
	}
	return false
}

func tagIsUnknown(tag map[string]string) bool {
	return !KNOWN_TAG_KINDS[tag[FIELD_KIND]]
}

func parseTag(row string) map[string]string {
	tag := make(map[string]string)
	parts := strings.Split(row, "\t")

	tag[FIELD_FUNCTION_NAME] = parts[0]
	parts = parts[1:]

	for _, part := range parts {
		if strings.Contains(part, ":") {
			field := part[:strings.Index(part, ":")]
			if FIELDS[field] {
				tag[field] = strings.TrimSpace(part[strings.Index(part, ":")+1:])
			}
		}
	}

	if strings.Contains(row, "/^") && strings.Contains(row, "$/;") {
		tag[FIELD_CODE] = row[strings.Index(row, "/^")+2 : strings.Index(row, "$/;")]
	}

	return tag
}

func removeEmpty(rows []string) []string {
	var newRows []string
	for _, row := range rows {
		row = strings.TrimSpace(row)
		if len(row) > 0 {
			newRows = append(newRows, row)
		}
	}

	return newRows
}
