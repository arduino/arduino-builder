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

package ctags

import (
	"arduino.cc/builder/constants"
	"arduino.cc/builder/types"
	"bufio"
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
const EXTERN = "extern \"C\""

var KNOWN_TAG_KINDS = map[string]bool{
	"prototype": true,
	"function":  true,
}

type CTagsParser struct{}

func (s *CTagsParser) Run(ctx *types.Context) error {
	rows := strings.Split(ctx.CTagsOutput, "\n")

	rows = removeEmpty(rows)

	var tags []*types.CTag
	for _, row := range rows {
		tags = append(tags, parseTag(row))
	}

	skipTagsWhere(tags, tagIsUnknown, ctx)
	skipTagsWhere(tags, tagIsUnhandled, ctx)
	addPrototypes(tags)
	removeDefinedProtypes(tags, ctx)
	removeDuplicate(tags)
	skipTagsWhere(tags, prototypeAndCodeDontMatch, ctx)

	ctx.CTagsOfPreprocessedSource = tags

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
	if strings.Index(tag.Prototype, TEMPLATE) == 0 || strings.Index(tag.Code, TEMPLATE) == 0 {
		code := tag.Code
		if strings.Contains(code, "{") {
			code = code[:strings.Index(code, "{")]
		} else {
			code = code[:strings.LastIndex(code, ")")+1]
		}
		tag.Prototype = code + ";"
		return
	}

	tag.PrototypeModifiers = ""
	if strings.Index(tag.Code, STATIC+" ") != -1 {
		tag.PrototypeModifiers = tag.PrototypeModifiers + " " + STATIC
	}
	if strings.Index(tag.Code, EXTERN+" ") != -1 {
		tag.PrototypeModifiers = tag.PrototypeModifiers + " " + EXTERN
	}
	tag.PrototypeModifiers = strings.TrimSpace(tag.PrototypeModifiers)
}

func removeDefinedProtypes(tags []*types.CTag, ctx *types.Context) {
	definedPrototypes := make(map[string]bool)
	for _, tag := range tags {
		if tag.Kind == KIND_PROTOTYPE {
			definedPrototypes[tag.Prototype] = true
		}
	}

	for _, tag := range tags {
		if definedPrototypes[tag.Prototype] {
			if ctx.DebugLevel >= 10 {
				ctx.GetLogger().Fprintln(os.Stdout, constants.LOG_LEVEL_DEBUG, constants.MSG_SKIPPING_TAG_ALREADY_DEFINED, tag.FunctionName)
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

func skipTagsWhere(tags []*types.CTag, skipFunc skipFuncType, ctx *types.Context) {
	for _, tag := range tags {
		if !tag.SkipMe {
			skip := skipFunc(tag)
			if skip && ctx.DebugLevel >= 10 {
				ctx.GetLogger().Fprintln(os.Stdout, constants.LOG_LEVEL_DEBUG, constants.MSG_SKIPPING_TAG_WITH_REASON, tag.FunctionName, runtime.FuncForPC(reflect.ValueOf(skipFunc).Pointer()).Name())
			}
			tag.SkipMe = skip
		}
	}
}

func prototypeAndCodeDontMatch(tag *types.CTag) bool {
	if tag.SkipMe {
		return true
	}

	code := removeSpacesAndTabs(tag.Code)

	// original code is multi-line, which tags doesn't have - could we find this code in the
	// original source file, for purposes of checking here?
	if strings.Index(code, ")") == -1 {
		file, err := os.Open(tag.Filename)
		if err == nil {
			defer file.Close()

			scanner := bufio.NewScanner(file)
			line := 1

			// skip lines until we get to the start of this tag
			for scanner.Scan() && line < tag.Line {
				line++
			}

			// read up to 10 lines in search of a closing paren
			newcode := scanner.Text()
			for scanner.Scan() && line < (tag.Line+10) && strings.Index(newcode, ")") == -1 {
				newcode += scanner.Text()
			}

			// don't bother replacing the code text if we haven't found a closing paren
			if strings.Index(newcode, ")") != -1 {
				code = removeSpacesAndTabs(newcode)
			}
		}
	}

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

	signature := ""
	returntype := ""
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
				signature = value
			case "returntype":
				returntype = value
			case "class":
				tag.Class = value
			case "struct":
				tag.Struct = value
			case "namespace":
				tag.Namespace = value
			}
		}
	}
	tag.Prototype = returntype + " " + tag.FunctionName + signature + ";"

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
