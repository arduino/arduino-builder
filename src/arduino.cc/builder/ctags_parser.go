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
	"arduino.cc/builder/utils"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

const FIELD_KIND = "kind"
const FIELD_LINE = "line"
const FIELD_SIGNATURE = "signature"
const FIELD_RETURNTYPE = "returntype"
const FIELD_CODE = "code"
const FIELD_CLASS = "class"
const FIELD_STRUCT = "struct"
const FIELD_NAMESPACE = "namespace"
const FIELD_FILENAME = "filename"
const FIELD_SKIP = "skipMe"
const FIELD_FUNCTION_NAME = "functionName"

const KIND_PROTOTYPE = "prototype"
const KIND_FUNCTION = "function"
const KIND_PROTOTYPE_MODIFIERS = "prototype_modifiers"

const TEMPLATE = "template"
const STATIC = "static"
const TRUE = "true"

var FIELDS = map[string]bool{"kind": true, "line": true, "typeref": true, "signature": true, "returntype": true, "class": true, "struct": true, "namespace": true}
var KNOWN_TAG_KINDS = map[string]bool{"prototype": true, "function": true}
var FIELDS_MARKING_UNHANDLED_TAGS = []string{FIELD_CLASS, FIELD_STRUCT, FIELD_NAMESPACE}

type CTagsParser struct{}

func (s *CTagsParser) Run(context map[string]interface{}) error {
	rows := strings.Split(context[constants.CTX_CTAGS_OUTPUT].(string), "\n")

	rows = removeEmpty(rows)

	var tags []map[string]string
	for _, row := range rows {
		tags = append(tags, parseTag(row))
	}

	skipTagsWhere(tags, tagIsUnknown, context)
	skipTagsWithField(tags, FIELDS_MARKING_UNHANDLED_TAGS, context)
	skipTagsWhere(tags, signatureContainsDefaultArg, context)
	addPrototypes(tags)
	removeDefinedProtypes(tags, context)
	removeDuplicate(tags)
	skipTagsWhere(tags, prototypeAndCodeDontMatch, context)

	context[constants.CTX_CTAGS_OF_PREPROC_SOURCE] = tags

	return nil
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

func removeDefinedProtypes(tags []map[string]string, context map[string]interface{}) {
	definedPrototypes := make(map[string]bool)
	for _, tag := range tags {
		if tag[FIELD_KIND] == KIND_PROTOTYPE {
			definedPrototypes[tag[KIND_PROTOTYPE]] = true
		}
	}

	for _, tag := range tags {
		if definedPrototypes[tag[KIND_PROTOTYPE]] {
			if utils.DebugLevel(context) >= 10 {
				utils.Logger(context).Fprintln(os.Stderr, constants.MSG_SKIPPING_TAG_ALREADY_DEFINED, tag[FIELD_FUNCTION_NAME])
			}
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

func skipTagsWhere(tags []map[string]string, skipFunc skipFuncType, context map[string]interface{}) {
	for _, tag := range tags {
		if tag[FIELD_SKIP] != TRUE {
			skip := skipFunc(tag)
			if skip && utils.DebugLevel(context) >= 10 {
				utils.Logger(context).Fprintln(os.Stderr, constants.MSG_SKIPPING_TAG_WITH_REASON, tag[FIELD_FUNCTION_NAME], runtime.FuncForPC(reflect.ValueOf(skipFunc).Pointer()).Name())
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

func skipTagsWithField(tags []map[string]string, fields []string, context map[string]interface{}) {
	for _, tag := range tags {
		if field, skip := utils.TagHasAtLeastOneField(tag, fields); skip {
			if utils.DebugLevel(context) >= 10 {
				utils.Logger(context).Fprintln(os.Stderr, constants.MSG_SKIPPING_TAG_BECAUSE_HAS_FIELD, field)
			}
			tag[FIELD_SKIP] = TRUE
		}
	}
}

func tagIsUnknown(tag map[string]string) bool {
	return !KNOWN_TAG_KINDS[tag[FIELD_KIND]]
}

func parseTag(row string) map[string]string {
	tag := make(map[string]string)
	parts := strings.Split(row, "\t")

	tag[FIELD_FUNCTION_NAME] = parts[0]
	tag[FIELD_FILENAME] = parts[1]

	parts = parts[2:]

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
