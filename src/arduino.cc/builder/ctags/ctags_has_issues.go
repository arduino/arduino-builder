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
	"bufio"
	"os"
	"strings"

	"arduino.cc/builder/types"
)

func (p *CTagsParser) FixCLinkageTagsDeclarations(tags []*types.CTag) {

	linesMap := p.FindCLinkageLines(tags)
	for i, _ := range tags {

		if sliceContainsInt(linesMap[tags[i].Filename], tags[i].Line) &&
			!strings.Contains(tags[i].PrototypeModifiers, EXTERN) {
			tags[i].PrototypeModifiers = tags[i].PrototypeModifiers + " " + EXTERN
		}
	}
}

func sliceContainsInt(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (p *CTagsParser) prototypeAndCodeDontMatch(tag *types.CTag) bool {
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

/* This function scans the source files searching for "extern C" context
 * It save the line numbers in a map filename -> {lines...}
 */
func (p *CTagsParser) FindCLinkageLines(tags []*types.CTag) map[string][]int {

	lines := make(map[string][]int)

	for _, tag := range tags {

		if lines[tag.Filename] != nil {
			break
		}

		file, err := os.Open(tag.Filename)
		if err == nil {
			defer file.Close()

			lines[tag.Filename] = append(lines[tag.Filename], 0)

			scanner := bufio.NewScanner(file)

			// we can't remove the comments otherwise the line number will be wrong
			// there are three cases:
			// 1 - extern "C" void foo()
			// 2 - extern "C" {
			//		void foo();
			//		void bar();
			//	}
			// 3 - extern "C"
			//	{
			//		void foo();
			//		void bar();
			//	}
			// case 1 and 2 can be simply recognized with string matching and indent level count
			// case 3 needs specia attention: if the line ONLY contains `extern "C"` string, don't bail out on indent level = 0

			inScope := false
			enteringScope := false
			indentLevels := 0
			line := 0

			externCDecl := removeSpacesAndTabs(EXTERN)

			for scanner.Scan() {
				line++
				str := removeSpacesAndTabs(scanner.Text())

				if len(str) == 0 {
					continue
				}

				// check if we are on the first non empty line after externCDecl in case 3
				if enteringScope == true {
					enteringScope = false
				}

				// check if the line contains externCDecl
				if strings.Contains(str, externCDecl) {
					inScope = true
					if len(str) == len(externCDecl) {
						// case 3
						enteringScope = true
					}
				}
				if inScope == true {
					lines[tag.Filename] = append(lines[tag.Filename], line)
				}
				indentLevels += strings.Count(str, "{") - strings.Count(str, "}")

				// Bail out if indentLevel is zero and we are not in case 3
				if indentLevels == 0 && enteringScope == false {
					inScope = false
				}
			}
		}

	}
	return lines
}
