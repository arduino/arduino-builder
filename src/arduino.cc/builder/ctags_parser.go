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

const TEMPLATE = "template"

var FIELDS = map[string]bool{"kind": true, "line": true, "typeref": true, "signature": true, "returntype": true, "class": true, "struct": true, "namespace": true}
var KNOWN_TAG_KINDS = map[string]bool{"prototype": true, "function": true}

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

	tags = filterOutUnknownTags(tags)
	tags = filterOutTagsWithField(tags, FIELD_CLASS)
	tags = filterOutTagsWithField(tags, FIELD_STRUCT)
	tags = filterOutTagsWithField(tags, FIELD_NAMESPACE)
	tags = skipTagsWhere(tags, signatureContainsDefaultArg)
	tags = addPrototypes(tags)
	tags = removeDefinedProtypes(tags)
	tags = removeDuplicate(tags)
	tags = skipTagsWhere(tags, prototypeAndCodeDontMatch)

	if len(tags) > 0 {
		line, err := strconv.Atoi(tags[0][FIELD_LINE])
		if err != nil {
			return utils.WrapError(err)
		}
		context[constants.CTX_FIRST_FUNCTION_AT_LINE] = line
	}

	prototypes := toPrototypes(tags)

	context[s.PrototypesField] = prototypes

	return nil
}

func toPrototypes(tags []map[string]string) []*types.Prototype {
	prototypes := []*types.Prototype{}
	for _, tag := range tags {
		if tag[FIELD_SKIP] != "true" {
			ctag := types.Prototype{FunctionName: tag[FIELD_FUNCTION_NAME], Prototype: tag[KIND_PROTOTYPE], Fields: tag}
			prototypes = append(prototypes, &ctag)
		}
	}
	return prototypes
}

func addPrototypes(tags []map[string]string) []map[string]string {
	for _, tag := range tags {
		addPrototype(tag)
	}
	return tags
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
	} else {
		tag[KIND_PROTOTYPE] = tag[FIELD_RETURNTYPE] + " " + tag[FIELD_FUNCTION_NAME] + tag[FIELD_SIGNATURE] + ";"
	}
}

func removeDefinedProtypes(tags []map[string]string) []map[string]string {
	definedPrototypes := make(map[string]bool)
	for _, tag := range tags {
		if tag[FIELD_KIND] == KIND_PROTOTYPE {
			definedPrototypes[tag[KIND_PROTOTYPE]] = true
		}
	}

	var newTags []map[string]string
	for _, tag := range tags {
		if !definedPrototypes[tag[KIND_PROTOTYPE]] {
			newTags = append(newTags, tag)
		}
	}
	return newTags
}

func removeDuplicate(tags []map[string]string) []map[string]string {
	definedPrototypes := make(map[string]bool)

	var newTags []map[string]string
	for _, tag := range tags {
		if !definedPrototypes[tag[KIND_PROTOTYPE]] {
			newTags = append(newTags, tag)
			definedPrototypes[tag[KIND_PROTOTYPE]] = true
		}
	}
	return newTags
}

type skipFuncType func(tag map[string]string) bool

func skipTagsWhere(tags []map[string]string, skipFuncs ...skipFuncType) []map[string]string {
	for _, tag := range tags {
		skip := skipFuncs[0](tag)
		for _, skipFunc := range skipFuncs[1:] {
			skip = skip || skipFunc(tag)
		}
		tag[FIELD_SKIP] = strconv.FormatBool(skip)
	}
	return tags
}

func signatureContainsDefaultArg(tag map[string]string) bool {
	return strings.Contains(tag[FIELD_SIGNATURE], "=")
}

func prototypeAndCodeDontMatch(tag map[string]string) bool {
	if tag[FIELD_SKIP] == "true" {
		return true
	}

	code := removeSpacesAndTabs(tag[FIELD_CODE])
	prototype := removeSpacesAndTabs(tag[KIND_PROTOTYPE])
	prototype = prototype[0 : len(prototype)-1]

	return strings.Index(code, prototype) == -1
}

func removeSpacesAndTabs(s string) string {
	s = strings.Replace(s, " ", "", -1)
	s = strings.Replace(s, "\t", "", -1)
	return s
}

func filterOutTagsWithField(tags []map[string]string, field string) []map[string]string {
	var newTags []map[string]string
	for _, tag := range tags {
		if tag[field] == constants.EMPTY_STRING {
			newTags = append(newTags, tag)
		}
	}
	return newTags
}

func filterOutUnknownTags(tags []map[string]string) []map[string]string {
	var newTags []map[string]string
	for _, tag := range tags {
		if KNOWN_TAG_KINDS[tag[FIELD_KIND]] {
			newTags = append(newTags, tag)
		}
	}
	return newTags
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
