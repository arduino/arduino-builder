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

package utils

import (
	"arduino.cc/builder/constants"
	"arduino.cc/builder/gohasissues"
	"arduino.cc/builder/i18n"
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func KeysOfMapOfStringInterface(input map[string]interface{}) []string {
	var keys []string
	for key, _ := range input {
		keys = append(keys, key)
	}
	return keys
}

func KeysOfMapOfString(input map[string]string) []string {
	var keys []string
	for key, _ := range input {
		keys = append(keys, key)
	}
	return keys
}

func PrettyOSName() string {
	switch osName := runtime.GOOS; osName {
	case "darwin":
		return "macosx"
	case "linux":
		return "linux"
	case "windows":
		return "windows"
	default:
		return "other"
	}
}

func ParseCommandLine(input string, logger i18n.Logger) ([]string, error) {
	var parts []string
	escapingChar := constants.EMPTY_STRING
	escapedArg := constants.EMPTY_STRING
	for _, inputPart := range strings.Split(input, constants.SPACE) {
		inputPart = strings.TrimSpace(inputPart)
		if len(inputPart) == 0 {
			continue
		}

		if escapingChar == constants.EMPTY_STRING {
			if inputPart[0] != '"' && inputPart[0] != '\'' {
				parts = append(parts, inputPart)
				continue
			}

			escapingChar = string(inputPart[0])
			inputPart = inputPart[1:]
			escapedArg = constants.EMPTY_STRING
		}

		if inputPart[len(inputPart)-1] != '"' && inputPart[len(inputPart)-1] != '\'' {
			escapedArg = escapedArg + inputPart + " "
			continue
		}

		escapedArg = escapedArg + inputPart[:len(inputPart)-1]
		escapedArg = strings.TrimSpace(escapedArg)
		if len(escapedArg) > 0 {
			parts = append(parts, escapedArg)
		}
		escapingChar = constants.EMPTY_STRING
	}

	if escapingChar != constants.EMPTY_STRING {
		return nil, i18n.ErrorfWithLogger(logger, constants.MSG_INVALID_QUOTING, escapingChar)
	}

	return parts, nil
}

type filterFiles func([]os.FileInfo) []os.FileInfo

func ReadDirFiltered(folder string, fn filterFiles) ([]os.FileInfo, error) {
	files, err := gohasissues.ReadDir(folder)
	if err != nil {
		return nil, i18n.WrapError(err)
	}
	return fn(files), nil
}

func FilterDirs(files []os.FileInfo) []os.FileInfo {
	var filtered []os.FileInfo
	for _, info := range files {
		if info.IsDir() {
			filtered = append(filtered, info)
		}
	}
	return filtered
}

func FilterFilesWithExtension(extension string) filterFiles {
	return func(files []os.FileInfo) []os.FileInfo {
		var filtered []os.FileInfo
		for _, file := range files {
			if !file.IsDir() && filepath.Ext(file.Name()) == extension {
				filtered = append(filtered, file)
			}
		}
		return filtered
	}
}

var SOURCE_CONTROL_FOLDERS = map[string]bool{"CVS": true, "RCS": true, ".git": true, ".svn": true, ".hg": true, ".bzr": true}

func IsSCCSOrHiddenFile(file os.FileInfo) bool {
	return IsSCCSFile(file) || IsHiddenFile(file)
}

func IsHiddenFile(file os.FileInfo) bool {
	name := filepath.Base(file.Name())

	if name[0] == '.' {
		return true
	}

	return false
}

func IsSCCSFile(file os.FileInfo) bool {
	name := filepath.Base(file.Name())

	if SOURCE_CONTROL_FOLDERS[name] {
		return true
	}

	return false
}

func SliceContains(slice []string, target string) bool {
	for _, value := range slice {
		if value == target {
			return true
		}
	}
	return false
}

type mapFunc func(string) string

func Map(slice []string, fn mapFunc) []string {
	newSlice := []string{}
	for _, elem := range slice {
		newSlice = append(newSlice, fn(elem))
	}
	return newSlice
}

type filterFunc func(string) bool

func Filter(slice []string, fn filterFunc) []string {
	newSlice := []string{}
	for _, elem := range slice {
		if fn(elem) {
			newSlice = append(newSlice, elem)
		}
	}
	return newSlice
}

func WrapWithHyphenI(value string) string {
	return "\"-I" + value + "\""
}

func TrimSpace(value string) string {
	return strings.TrimSpace(value)
}

type argFilterFunc func(int, string, []string) bool

func PrepareCommandFilteredArgs(pattern string, filter argFilterFunc, logger i18n.Logger) (*exec.Cmd, error) {
	parts, err := ParseCommandLine(pattern, logger)
	if err != nil {
		return nil, i18n.WrapError(err)
	}
	command := parts[0]
	parts = parts[1:]
	var args []string
	for idx, part := range parts {
		if filter(idx, part, parts) {
			args = append(args, part)
		}
	}

	return exec.Command(command, args...), nil
}

func filterEmptyArg(_ int, arg string, _ []string) bool {
	return arg != constants.EMPTY_STRING
}

func PrepareCommand(pattern string, logger i18n.Logger) (*exec.Cmd, error) {
	return PrepareCommandFilteredArgs(pattern, filterEmptyArg, logger)
}

func MapHas(aMap map[string]interface{}, key string) bool {
	_, ok := aMap[key]
	return ok
}

func MapStringStringHas(aMap map[string]string, key string) bool {
	_, ok := aMap[key]
	return ok
}

func SliceToMapStringBool(keys []string, value bool) map[string]bool {
	aMap := make(map[string]bool)
	for _, key := range keys {
		aMap[key] = value
	}
	return aMap
}

func AbsolutizePaths(files []string) ([]string, error) {
	for idx, file := range files {
		absFile, err := filepath.Abs(file)
		if err != nil {
			return nil, i18n.WrapError(err)
		}
		files[idx] = absFile
	}

	return files, nil
}

func ReadFileToRows(file string) ([]string, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, i18n.WrapError(err)
	}
	txt := string(bytes)
	txt = strings.Replace(txt, "\r\n", "\n", -1)

	return strings.Split(txt, "\n"), nil
}

func TheOnlySubfolderOf(folder string) (string, error) {
	subfolders, err := ReadDirFiltered(folder, FilterDirs)
	if err != nil {
		return constants.EMPTY_STRING, i18n.WrapError(err)
	}

	if len(subfolders) != 1 {
		return constants.EMPTY_STRING, nil
	}

	return subfolders[0].Name(), nil
}

func FilterOutFoldersByNames(folders []os.FileInfo, names ...string) []os.FileInfo {
	filterNames := SliceToMapStringBool(names, true)

	var filtered []os.FileInfo
	for _, folder := range folders {
		if !filterNames[folder.Name()] {
			filtered = append(filtered, folder)
		}
	}

	return filtered
}

type CheckFilePathFunc func(filePath string) bool

func CollectAllReadableFiles(collector *[]string, test CheckFilePathFunc) filepath.WalkFunc {
	walkFunc := func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}
		if !test(currentPath) {
			return nil
		}
		currentFile, err := os.Open(currentPath)
		if err != nil {
			return nil
		}
		currentFile.Close()

		*collector = append(*collector, currentPath)
		return nil
	}
	return walkFunc
}

func AppendIfNotPresent(target []string, elements ...string) []string {
	for _, element := range elements {
		if !SliceContains(target, element) {
			target = append(target, element)
		}
	}
	return target
}

func EnsureFolderExists(folder string) error {
	return os.MkdirAll(folder, os.FileMode(0755))
}

func WriteFileBytes(targetFilePath string, data []byte) error {
	return ioutil.WriteFile(targetFilePath, data, os.FileMode(0644))
}

func WriteFile(targetFilePath string, data string) error {
	return WriteFileBytes(targetFilePath, []byte(data))
}

func TouchFile(targetFilePath string) error {
	return WriteFileBytes(targetFilePath, []byte{})
}

func NULLFile() string {
	if runtime.GOOS == "windows" {
		return "nul"
	}
	return "/dev/null"
}

func MD5Sum(data []byte) string {
	md5sumBytes := md5.Sum(data)
	return hex.EncodeToString(md5sumBytes[:])
}
