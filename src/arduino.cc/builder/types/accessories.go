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

package types

type UniqueStringQueue []string

func (h UniqueStringQueue) Len() int           { return len(h) }
func (h UniqueStringQueue) Less(i, j int) bool { return false }
func (h UniqueStringQueue) Swap(i, j int)      { panic("Who called me?!?") }

func (h *UniqueStringQueue) Push(x interface{}) {
	if !sliceContains(*h, x.(string)) {
		*h = append(*h, x.(string))
	}
}

func (h *UniqueStringQueue) Pop() interface{} {
	old := *h
	x := old[0]
	*h = old[1:]
	return x
}

func (h *UniqueStringQueue) Empty() bool {
	return h.Len() == 0
}

// duplication of utils.SliceContains! Thanks golang! Why? Because you can't have import cycles, so types cannot import from utils because utils already imports from types
func sliceContains(slice []string, target string) bool {
	for _, value := range slice {
		if value == target {
			return true
		}
	}
	return false
}

type LibraryResolutionResult struct {
	Library          *Library
	NotUsedLibraries []*Library
}
