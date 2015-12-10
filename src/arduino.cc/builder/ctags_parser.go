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
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

const KIND_PROTOTYPE = "prototype"
const KIND_FUNCTION = "function"

//const KIND_PROTOTYPE_MODIFIERS = "prototype_modifiers"

const TEMPLATE = "template"
const STATIC = "static"

var KNOWN_TAG_KINDS = map[string]bool{
	"prototype": true,
	"function":  true,
}

type CTagsParser struct{}

func (s *CTagsParser) Run(context map[string]interface{}) error {
	rows := strings.Split(context[constants.CTX_CTAGS_OUTPUT].(string), "\n")

	rows = removeEmpty(rows)

	var tags []*types.CTag
	for _, row := range rows {
		tags = append(tags, parseTag(row))
	}

	skipTagsWhere(tags, tagIsUnknown, context)
	skipTagsWhere(tags, tagIsUnhandled, context)
	skipTagsWhere(tags, signatureContainsDefaultArg, context)
	addPrototypes(tags)
	removeDefinedProtypes(tags, context)
	removeDuplicate(tags)
	skipTagsWhere(tags, prototypeAndCodeDontMatch, context)

	context[constants.CTX_CTAGS_OF_PREPROC_SOURCE] = tags

	return nil
}

func addPrototypes(tags []*types.CTag) {
	for _, tag := range tags {
		if !tag.SkipMe {
			addPrototype(tag)
		}
	}
}

func addPrototype(tag *types.CTag) {
	if strings.Index(tag.Returntype, TEMPLATE) == 0 || strings.Index(tag.Code, TEMPLATE) == 0 {
		code := tag.Code
		if strings.Contains(code, "{") {
			code = code[:strings.Index(code, "{")]
		} else {
			code = code[:strings.LastIndex(code, ")")+1]
		}
		tag.Prototype = code + ";"
		return
	}

	tag.Prototype = tag.Returntype + " " + tag.FunctionName + tag.Signature + ";"

	tag.PrototypeModifiers = ""
	if strings.Index(tag.Code, STATIC+" ") != -1 {
		tag.PrototypeModifiers = tag.PrototypeModifiers + " " + STATIC
	}
	tag.PrototypeModifiers = strings.TrimSpace(tag.PrototypeModifiers)
}

func removeDefinedProtypes(tags []*types.CTag, context map[string]interface{}) {
	definedPrototypes := make(map[string]bool)
	for _, tag := range tags {
		if tag.Kind == KIND_PROTOTYPE {
			definedPrototypes[tag.Prototype] = true
		}
	}

	for _, tag := range tags {
		if definedPrototypes[tag.Prototype] {
			if utils.DebugLevel(context) >= 10 {
				utils.Logger(context).Fprintln(os.Stderr, constants.MSG_SKIPPING_TAG_ALREADY_DEFINED, tag.FunctionName)
			}
			tag.SkipMe = true
		}
	}
}

func removeDuplicate(tags []*types.CTag) {
	definedPrototypes := make(map[string]bool)

	for _, tag := range tags {
		if !definedPrototypes[tag.Prototype] {
			definedPrototypes[tag.Prototype] = true
		} else {
			tag.SkipMe = true
		}
	}
}

type skipFuncType func(tag *types.CTag) bool

func skipTagsWhere(tags []*types.CTag, skipFunc skipFuncType, context map[string]interface{}) {
	for _, tag := range tags {
		if !tag.SkipMe {
			skip := skipFunc(tag)
			if skip && utils.DebugLevel(context) >= 10 {
				utils.Logger(context).Fprintln(os.Stderr, constants.MSG_SKIPPING_TAG_WITH_REASON, tag.FunctionName, runtime.FuncForPC(reflect.ValueOf(skipFunc).Pointer()).Name())
			}
			tag.SkipMe = skip
		}
	}
}

func signatureContainsDefaultArg(tag *types.CTag) bool {
	return strings.Contains(tag.Signature, "=")
}

func prototypeAndCodeDontMatch(tag *types.CTag) bool {
	if tag.SkipMe {
		return true
	}

	code := removeSpacesAndTabs(tag.Code)
	prototype := removeSpacesAndTabs(tag.Prototype)
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

func tagIsUnhandled(tag *types.CTag) bool {
	return !isHandled(tag)
}

func isHandled(tag *types.CTag) bool {
	if tag.Class != "" {
		return false
	}
	if tag.Struct != "" {
		return false
	}
	if tag.Namespace != "" {
		return false
	}
	return true
}

func tagIsUnknown(tag *types.CTag) bool {
	return !KNOWN_TAG_KINDS[tag.Kind]
}

func parseTag(row string) *types.CTag {
	tag := &types.CTag{}
	parts := strings.Split(row, "\t")

	tag.FunctionName = parts[0]
	tag.Filename = parts[1]

	parts = parts[2:]

	for _, part := range parts {
		if strings.Contains(part, ":") {
			colon := strings.Index(part, ":")
			field := part[:colon]
			value := strings.TrimSpace(part[colon+1:])
			switch field {
			case "kind":
				tag.Kind = value
			case "line":
				val, _ := strconv.Atoi(value)
				// TODO: Check err from strconv.Atoi
				tag.Line = val
			case "typeref":
				tag.Typeref = value
			case "signature":
				tag.Signature = value
			case "returntype":
				tag.Returntype = value
			case "class":
				tag.Class = value
			case "struct":
				tag.Struct = value
			case "namespace":
				tag.Namespace = value
			}
		}
	}

	if strings.Contains(row, "/^") && strings.Contains(row, "$/;") {
		tag.Code = row[strings.Index(row, "/^")+2 : strings.Index(row, "$/;")]
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
