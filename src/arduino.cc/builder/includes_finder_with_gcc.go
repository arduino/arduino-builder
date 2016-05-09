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

// XXX: Obsolete?

import (
	"arduino.cc/builder/builder_utils"
	"arduino.cc/builder/constants"
	"arduino.cc/builder/i18n"
	"arduino.cc/builder/types"
	"arduino.cc/builder/utils"
	"strings"
)

type IncludesFinderWithGCC struct {
	SourceFile string
}

func (s *IncludesFinderWithGCC) Run(ctx *types.Context) error {
	buildProperties := ctx.BuildProperties.Clone()
	verbose := ctx.Verbose
	logger := ctx.GetLogger()

	includes := utils.Map(ctx.IncludeFolders, utils.WrapWithHyphenI)
	includesParams := strings.Join(includes, " ")

	properties := buildProperties.Clone()
	properties[constants.BUILD_PROPERTIES_SOURCE_FILE] = s.SourceFile
	properties[constants.BUILD_PROPERTIES_INCLUDES] = includesParams
	builder_utils.RemoveHyphenMDDFlagFromGCCCommandLine(properties)

	if properties[constants.RECIPE_PREPROC_INCLUDES] == "" {
		//generate RECIPE_PREPROC_INCLUDES from RECIPE_CPP_PATTERN
		properties[constants.RECIPE_PREPROC_INCLUDES] = GeneratePreprocIncludePatternFromCompile(properties[constants.RECIPE_CPP_PATTERN])
	}

	output, err := builder_utils.ExecRecipe(properties, constants.RECIPE_PREPROC_INCLUDES, true, verbose, false, logger)
	if err != nil {
		return i18n.WrapError(err)
	}

	ctx.OutputGccMinusM = string(output)

	return nil
}

func GeneratePreprocIncludePatternFromCompile(compilePattern string) string {
	// add {preproc.includes.flags}
	// remove -o "{object_file}"
	returnString := compilePattern
	returnString = strings.Replace(returnString, "{compiler.cpp.flags}", "{compiler.cpp.flags} {preproc.includes.flags}", 1)
	returnString = strings.Replace(returnString, "-o {object_file}", "", 1)
	return returnString
}
