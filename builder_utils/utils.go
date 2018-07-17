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

package builder_utils

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/arduino/arduino-builder/constants"
	"github.com/arduino/arduino-builder/i18n"
	"github.com/arduino/arduino-builder/types"
	"github.com/arduino/arduino-builder/utils"
	"github.com/arduino/go-properties-map"
)

func PrintProgressIfProgressEnabledAndMachineLogger(ctx *types.Context) {

	if !ctx.Progress.PrintEnabled {
		return
	}

	log := ctx.GetLogger()
	if log.Name() == "machine" {
		log.Println(constants.LOG_LEVEL_INFO, constants.MSG_PROGRESS, strconv.FormatFloat(ctx.Progress.Progress, 'f', 2, 32))
		ctx.Progress.Progress += ctx.Progress.Steps
	}
}

func CompileFilesRecursive(ctx *types.Context, objectFiles []string, sourcePath string, buildPath string, buildProperties properties.Map, includes []string) ([]string, error) {
	objectFiles, err := CompileFiles(ctx, objectFiles, sourcePath, false, buildPath, buildProperties, includes)
	if err != nil {
		return nil, i18n.WrapError(err)
	}

	folders, err := utils.ReadDirFiltered(sourcePath, utils.FilterDirs)
	if err != nil {
		return nil, i18n.WrapError(err)
	}

	for _, folder := range folders {
		objectFiles, err = CompileFilesRecursive(ctx, objectFiles, filepath.Join(sourcePath, folder.Name()), filepath.Join(buildPath, folder.Name()), buildProperties, includes)
		if err != nil {
			return nil, i18n.WrapError(err)
		}
	}

	return objectFiles, nil
}

func CompileFiles(ctx *types.Context, objectFiles []string, sourcePath string, recurse bool, buildPath string, buildProperties properties.Map, includes []string) ([]string, error) {
	objectFiles, err := compileFilesWithExtensionWithRecipe(ctx, objectFiles, sourcePath, recurse, buildPath, buildProperties, includes, ".S", constants.RECIPE_S_PATTERN)
	if err != nil {
		return nil, i18n.WrapError(err)
	}
	objectFiles, err = compileFilesWithExtensionWithRecipe(ctx, objectFiles, sourcePath, recurse, buildPath, buildProperties, includes, ".c", constants.RECIPE_C_PATTERN)
	if err != nil {
		return nil, i18n.WrapError(err)
	}
	objectFiles, err = compileFilesWithExtensionWithRecipe(ctx, objectFiles, sourcePath, recurse, buildPath, buildProperties, includes, ".cpp", constants.RECIPE_CPP_PATTERN)
	if err != nil {
		return nil, i18n.WrapError(err)
	}
	return objectFiles, nil
}

func compileFilesWithExtensionWithRecipe(ctx *types.Context, objectFiles []string, sourcePath string, recurse bool, buildPath string, buildProperties properties.Map, includes []string, extension string, recipe string) ([]string, error) {
	sources, err := findFilesInFolder(sourcePath, extension, recurse)
	if err != nil {
		return nil, i18n.WrapError(err)
	}
	return compileFilesWithRecipe(ctx, objectFiles, sourcePath, sources, buildPath, buildProperties, includes, recipe)
}

func findFilesInFolder(sourcePath string, extension string, recurse bool) ([]string, error) {
	files, err := utils.ReadDirFiltered(sourcePath, utils.FilterFilesWithExtensions(extension))
	if err != nil {
		return nil, i18n.WrapError(err)
	}
	var sources []string
	for _, file := range files {
		sources = append(sources, filepath.Join(sourcePath, file.Name()))
	}

	if recurse {
		folders, err := utils.ReadDirFiltered(sourcePath, utils.FilterDirs)
		if err != nil {
			return nil, i18n.WrapError(err)
		}

		for _, folder := range folders {
			otherSources, err := findFilesInFolder(filepath.Join(sourcePath, folder.Name()), extension, recurse)
			if err != nil {
				return nil, i18n.WrapError(err)
			}
			sources = append(sources, otherSources...)
		}
	}

	return sources, nil
}

func findAllFilesInFolder(sourcePath string, recurse bool) ([]string, error) {
	files, err := utils.ReadDirFiltered(sourcePath, utils.FilterFiles())
	if err != nil {
		return nil, i18n.WrapError(err)
	}
	var sources []string
	for _, file := range files {
		sources = append(sources, filepath.Join(sourcePath, file.Name()))
	}

	if recurse {
		folders, err := utils.ReadDirFiltered(sourcePath, utils.FilterDirs)
		if err != nil {
			return nil, i18n.WrapError(err)
		}

		for _, folder := range folders {
			otherSources, err := findAllFilesInFolder(filepath.Join(sourcePath, folder.Name()), recurse)
			if err != nil {
				return nil, i18n.WrapError(err)
			}
			sources = append(sources, otherSources...)
		}
	}

	return sources, nil
}

func compileFilesWithRecipe(ctx *types.Context, objectFiles []string, sourcePath string, sources []string, buildPath string, buildProperties properties.Map, includes []string, recipe string) ([]string, error) {
	if len(sources) == 0 {
		return objectFiles, nil
	}
	objectFilesChan := make(chan string)
	errorsChan := make(chan error)
	doneChan := make(chan struct{})

	ctx.Progress.Steps = ctx.Progress.Steps / float64(len(sources))
	var wg sync.WaitGroup
	wg.Add(len(sources))

	for _, source := range sources {
		go func(source string) {
			defer wg.Done()
			PrintProgressIfProgressEnabledAndMachineLogger(ctx)
			objectFile, err := compileFileWithRecipe(ctx, sourcePath, source, buildPath, buildProperties, includes, recipe)
			if err != nil {
				errorsChan <- err
			} else {
				objectFilesChan <- objectFile
			}
		}(source)
	}

	go func() {
		wg.Wait()
		doneChan <- struct{}{}
	}()

	for {
		select {
		case objectFile := <-objectFilesChan:
			objectFiles = append(objectFiles, objectFile)
		case err := <-errorsChan:
			return nil, i18n.WrapError(err)
		case <-doneChan:
			close(objectFilesChan)
			for objectFile := range objectFilesChan {
				objectFiles = append(objectFiles, objectFile)
			}
			sort.Strings(objectFiles)
			return objectFiles, nil
		}
	}
}

func compileFileWithRecipe(ctx *types.Context, sourcePath string, source string, buildPath string, buildProperties properties.Map, includes []string, recipe string) (string, error) {
	logger := ctx.GetLogger()
	properties := buildProperties.Clone()
	properties[constants.BUILD_PROPERTIES_COMPILER_WARNING_FLAGS] = properties[constants.BUILD_PROPERTIES_COMPILER_WARNING_FLAGS+"."+ctx.WarningsLevel]
	properties[constants.BUILD_PROPERTIES_INCLUDES] = strings.Join(includes, constants.SPACE)
	properties[constants.BUILD_PROPERTIES_SOURCE_FILE] = source
	relativeSource, err := filepath.Rel(sourcePath, source)
	if err != nil {
		return "", i18n.WrapError(err)
	}
	properties[constants.BUILD_PROPERTIES_OBJECT_FILE] = filepath.Join(buildPath, relativeSource+".o")

	err = utils.EnsureFolderExists(filepath.Dir(properties[constants.BUILD_PROPERTIES_OBJECT_FILE]))
	if err != nil {
		return "", i18n.WrapError(err)
	}

	objIsUpToDate, err := ObjFileIsUpToDate(ctx, properties[constants.BUILD_PROPERTIES_SOURCE_FILE], properties[constants.BUILD_PROPERTIES_OBJECT_FILE], filepath.Join(buildPath, relativeSource+".d"))
	if err != nil {
		return "", i18n.WrapError(err)
	}

	if !objIsUpToDate {
		_, _, err = ExecRecipe(ctx, properties, recipe, false /* stdout */, utils.ShowIfVerbose /* stderr */, utils.Show)
		if err != nil {
			return "", i18n.WrapError(err)
		}
	} else if ctx.Verbose {
		logger.Println(constants.LOG_LEVEL_INFO, constants.MSG_USING_PREVIOUS_COMPILED_FILE, properties[constants.BUILD_PROPERTIES_OBJECT_FILE])
	}

	return properties[constants.BUILD_PROPERTIES_OBJECT_FILE], nil
}

func ObjFileIsUpToDate(ctx *types.Context, sourceFile, objectFile, dependencyFile string) (bool, error) {
	sourceFile = filepath.Clean(sourceFile)
	objectFile = filepath.Clean(objectFile)
	dependencyFile = filepath.Clean(dependencyFile)
	logger := ctx.GetLogger()
	debugLevel := ctx.DebugLevel

	if debugLevel >= 20 {
		logger.Fprintln(os.Stdout, constants.LOG_LEVEL_DEBUG, "Checking previous results for {0} (result = {1}, dep = {2})", sourceFile, objectFile, dependencyFile)
	}

	sourceFileStat, err := os.Stat(sourceFile)
	if err != nil {
		return false, i18n.WrapError(err)
	}

	objectFileStat, err := os.Stat(objectFile)
	if err != nil {
		if os.IsNotExist(err) {
			if debugLevel >= 20 {
				logger.Fprintln(os.Stdout, constants.LOG_LEVEL_DEBUG, "Not found: {0}", objectFile)
			}
			return false, nil
		} else {
			return false, i18n.WrapError(err)
		}
	}

	dependencyFileStat, err := os.Stat(dependencyFile)
	if err != nil {
		if os.IsNotExist(err) {
			if debugLevel >= 20 {
				logger.Fprintln(os.Stdout, constants.LOG_LEVEL_DEBUG, "Not found: {0}", dependencyFile)
			}
			return false, nil
		} else {
			return false, i18n.WrapError(err)
		}
	}

	if sourceFileStat.ModTime().After(objectFileStat.ModTime()) {
		if debugLevel >= 20 {
			logger.Fprintln(os.Stdout, constants.LOG_LEVEL_DEBUG, "{0} newer than {1}", sourceFile, objectFile)
		}
		return false, nil
	}
	if sourceFileStat.ModTime().After(dependencyFileStat.ModTime()) {
		if debugLevel >= 20 {
			logger.Fprintln(os.Stdout, constants.LOG_LEVEL_DEBUG, "{0} newer than {1}", sourceFile, dependencyFile)
		}
		return false, nil
	}

	rows, err := utils.ReadFileToRows(dependencyFile)
	if err != nil {
		return false, i18n.WrapError(err)
	}

	rows = utils.Map(rows, removeEndingBackSlash)
	rows = utils.Map(rows, strings.TrimSpace)
	rows = utils.Map(rows, unescapeDep)
	rows = utils.Filter(rows, nonEmptyString)

	if len(rows) == 0 {
		return true, nil
	}

	firstRow := rows[0]
	if !strings.HasSuffix(firstRow, ":") {
		if debugLevel >= 20 {
			logger.Fprintln(os.Stdout, constants.LOG_LEVEL_DEBUG, "No colon in first line of depfile")
		}
		return false, nil
	}
	objFileInDepFile := firstRow[:len(firstRow)-1]
	if objFileInDepFile != objectFile {
		if debugLevel >= 20 {
			logger.Fprintln(os.Stdout, constants.LOG_LEVEL_DEBUG, "Depfile is about different file: {0}", objFileInDepFile)
		}
		return false, nil
	}

	rows = rows[1:]
	for _, row := range rows {
		depStat, err := os.Stat(row)
		if err != nil && !os.IsNotExist(err) {
			// There is probably a parsing error of the dep file
			// Ignore the error and trigger a full rebuild anyway
			if debugLevel >= 20 {
				logger.Fprintln(os.Stdout, constants.LOG_LEVEL_DEBUG, "Failed to read: {0}", row)
				logger.Fprintln(os.Stdout, constants.LOG_LEVEL_DEBUG, i18n.WrapError(err).Error())
			}
			return false, nil
		}
		if os.IsNotExist(err) {
			if debugLevel >= 20 {
				logger.Fprintln(os.Stdout, constants.LOG_LEVEL_DEBUG, "Not found: {0}", row)
			}
			return false, nil
		}
		if depStat.ModTime().After(objectFileStat.ModTime()) {
			if debugLevel >= 20 {
				logger.Fprintln(os.Stdout, constants.LOG_LEVEL_DEBUG, "{0} newer than {1}", row, objectFile)
			}
			return false, nil
		}
	}

	return true, nil
}

func unescapeDep(s string) string {
	s = strings.Replace(s, "\\ ", " ", -1)
	s = strings.Replace(s, "\\\t", "\t", -1)
	s = strings.Replace(s, "\\#", "#", -1)
	s = strings.Replace(s, "$$", "$", -1)
	s = strings.Replace(s, "\\\\", "\\", -1)
	return s
}

func removeEndingBackSlash(s string) string {
	if strings.HasSuffix(s, "\\") {
		s = s[:len(s)-1]
	}
	return s
}

func nonEmptyString(s string) bool {
	return s != constants.EMPTY_STRING
}

func CoreOrReferencedCoreHasChanged(corePath, targetCorePath, targetFile string) bool {

	targetFileStat, err := os.Stat(targetFile)
	if err == nil {
		files, err := findAllFilesInFolder(corePath, true)
		if err != nil {
			return true
		}
		for _, file := range files {
			fileStat, err := os.Stat(file)
			if err != nil || fileStat.ModTime().After(targetFileStat.ModTime()) {
				return true
			}
		}
		if targetCorePath != constants.EMPTY_STRING && !strings.EqualFold(corePath, targetCorePath) {
			return CoreOrReferencedCoreHasChanged(targetCorePath, constants.EMPTY_STRING, targetFile)
		}
		return false
	}
	return true
}

func TXTBuildRulesHaveChanged(corePath, targetCorePath, targetFile string) bool {

	targetFileStat, err := os.Stat(targetFile)
	if err == nil {
		files, err := findAllFilesInFolder(corePath, true)
		if err != nil {
			return true
		}
		for _, file := range files {
			// report changes only for .txt files
			if filepath.Ext(file) != ".txt" {
				continue
			}
			fileStat, err := os.Stat(file)
			if err != nil || fileStat.ModTime().After(targetFileStat.ModTime()) {
				return true
			}
		}
		if targetCorePath != constants.EMPTY_STRING && !strings.EqualFold(corePath, targetCorePath) {
			return TXTBuildRulesHaveChanged(targetCorePath, constants.EMPTY_STRING, targetFile)
		}
		return false
	}
	return true
}

func ArchiveCompiledFiles(ctx *types.Context, buildPath string, archiveFile string, objectFiles []string, buildProperties properties.Map) (string, error) {
	logger := ctx.GetLogger()
	archiveFilePath := filepath.Join(buildPath, archiveFile)

	rebuildArchive := false

	if archiveFileStat, err := os.Stat(archiveFilePath); err == nil {

		for _, objectFile := range objectFiles {
			objectFileStat, _ := os.Stat(objectFile)
			if objectFileStat.ModTime().After(archiveFileStat.ModTime()) {
				// need to rebuild the archive
				rebuildArchive = true
				break
			}
		}

		// something changed, rebuild the core archive
		if rebuildArchive {
			err = os.Remove(archiveFilePath)
			if err != nil {
				return "", i18n.WrapError(err)
			}
		} else {
			if ctx.Verbose {
				logger.Println(constants.LOG_LEVEL_INFO, constants.MSG_USING_PREVIOUS_COMPILED_FILE, archiveFilePath)
			}
			return archiveFilePath, nil
		}
	}

	for _, objectFile := range objectFiles {
		properties := buildProperties.Clone()
		properties[constants.BUILD_PROPERTIES_ARCHIVE_FILE] = filepath.Base(archiveFilePath)
		properties[constants.BUILD_PROPERTIES_ARCHIVE_FILE_PATH] = archiveFilePath
		properties[constants.BUILD_PROPERTIES_OBJECT_FILE] = objectFile

		_, _, err := ExecRecipe(ctx, properties, constants.RECIPE_AR_PATTERN, false /* stdout */, utils.ShowIfVerbose /* stderr */, utils.Show)
		if err != nil {
			return "", i18n.WrapError(err)
		}
	}

	return archiveFilePath, nil
}

// See util.ExecCommand for stdout/stderr arguments
func ExecRecipe(ctx *types.Context, buildProperties properties.Map, recipe string, removeUnsetProperties bool, stdout int, stderr int) ([]byte, []byte, error) {
	command, err := PrepareCommandForRecipe(ctx, buildProperties, recipe, removeUnsetProperties)
	if err != nil {
		return nil, nil, i18n.WrapError(err)
	}

	return utils.ExecCommand(ctx, command, stdout, stderr)
}

const COMMANDLINE_LIMIT = 30000

func PrepareCommandForRecipe(ctx *types.Context, buildProperties properties.Map, recipe string, removeUnsetProperties bool) (*exec.Cmd, error) {
	logger := ctx.GetLogger()
	pattern := buildProperties[recipe]
	if pattern == constants.EMPTY_STRING {
		return nil, i18n.ErrorfWithLogger(logger, constants.MSG_PATTERN_MISSING, recipe)
	}

	var err error
	commandLine := buildProperties.ExpandPropsInString(pattern)
	if removeUnsetProperties {
		commandLine = properties.DeleteUnexpandedPropsFromString(commandLine)
	}

	relativePath := ""

	if len(commandLine) > COMMANDLINE_LIMIT {
		relativePath = buildProperties[constants.BUILD_PROPERTIES_BUILD_PATH]
	}

	command, err := utils.PrepareCommand(commandLine, logger, relativePath)
	if err != nil {
		return nil, i18n.WrapError(err)
	}

	return command, nil
}

// CopyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

// GetCachedCoreArchiveFileName returns the filename to be used to store
// the global cached core.a.
func GetCachedCoreArchiveFileName(fqbn, coreFolder string) string {
	fqbnToUnderscore := strings.Replace(fqbn, ":", "_", -1)
	fqbnToUnderscore = strings.Replace(fqbnToUnderscore, "=", "_", -1)
	if absCoreFolder, err := filepath.Abs(coreFolder); err == nil {
		coreFolder = absCoreFolder
	} // silently continue if absolute path can't be detected
	hash := utils.MD5Sum([]byte(coreFolder))
	return "core_" + fqbnToUnderscore + "_" + hash + ".a"
}
