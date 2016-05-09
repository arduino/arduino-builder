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
	"arduino.cc/builder/types"
	"strings"
)

type CTagsToPrototypes struct{}

func (s *CTagsToPrototypes) Run(ctx *types.Context) error {
	tags := ctx.CTagsCollected

	lineWhereToInsertPrototypes := findLineWhereToInsertPrototypes(tags)
	if lineWhereToInsertPrototypes != -1 {
		ctx.PrototypesLineWhereToInsert = lineWhereToInsertPrototypes
	}

	ctx.Prototypes = toPrototypes(tags)
	return nil
}

func findLineWhereToInsertPrototypes(tags []*types.CTag) int {
	firstFunctionLine := firstFunctionAtLine(tags)
	firstFunctionPointerAsArgument := firstFunctionPointerUsedAsArgument(tags)
	if firstFunctionLine != -1 && firstFunctionPointerAsArgument != -1 {
		if firstFunctionLine < firstFunctionPointerAsArgument {
			return firstFunctionLine
		} else {
			return firstFunctionPointerAsArgument
		}
	} else if firstFunctionLine == -1 {
		return firstFunctionPointerAsArgument
	} else {
		return firstFunctionLine
	}
}

func firstFunctionPointerUsedAsArgument(tags []*types.CTag) int {
	functionNames := collectFunctionNames(tags)
	for _, tag := range tags {
		if functionNameUsedAsFunctionPointerIn(tag, functionNames) {
			return tag.Line
		}
	}
	return -1
}

func functionNameUsedAsFunctionPointerIn(tag *types.CTag, functionNames []string) bool {
	for _, functionName := range functionNames {
		if strings.Index(tag.Code, "&"+functionName) != -1 {
			return true
		}
	}
	return false
}

func collectFunctionNames(tags []*types.CTag) []string {
	names := []string{}
	for _, tag := range tags {
		if tag.Kind == KIND_FUNCTION {
			names = append(names, tag.FunctionName)
		}
	}
	return names
}

func firstFunctionAtLine(tags []*types.CTag) int {
	for _, tag := range tags {
		if !tagIsUnknown(tag) && isHandled(tag) && tag.Kind == KIND_FUNCTION {
			return tag.Line
		}
	}
	return -1
}

func toPrototypes(tags []*types.CTag) []*types.Prototype {
	prototypes := []*types.Prototype{}
	for _, tag := range tags {
		if strings.TrimSpace(tag.Prototype) == "" {
			continue
		}
		if !tag.SkipMe {
			prototype := &types.Prototype{
				FunctionName: tag.FunctionName,
				File:         tag.Filename,
				Prototype:    tag.Prototype,
				Modifiers:    tag.PrototypeModifiers,
				Line:         tag.Line,
				//Fields:       tag,
			}
			prototypes = append(prototypes, prototype)
		}
	}
	return prototypes
}
