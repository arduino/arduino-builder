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

package test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/arduino/arduino-builder/constants"
	"github.com/arduino/arduino-builder/gohasissues"
	"github.com/arduino/arduino-builder/i18n"
	"github.com/arduino/arduino-builder/utils"
	"github.com/arduino/go-properties-map"
	"github.com/go-errors/errors"
)

const HARDWARE_FOLDER = "downloaded_hardware"
const BOARD_MANAGER_FOLDER = "downloaded_board_manager_stuff"
const TOOLS_FOLDER = "downloaded_tools"
const LIBRARIES_FOLDER = "downloaded_libraries"
const PATCHES_FOLDER = "downloaded_stuff_patches"

type Tool struct {
	Name    string
	Package string
	Version string
	OsUrls  []OsUrl
}

type OsUrl struct {
	Os  string
	Url string
}

type Library struct {
	Name                   string
	Version                string
	VersionInLibProperties string
	Url                    string
}

type Core struct {
	Maintainer string
	Arch       string
	Version    string
	Url        string
}

func DownloadCoresAndToolsAndLibraries(t *testing.T) {
	cores := []Core{
		Core{Maintainer: "arduino", Arch: "avr", Version: "1.6.10"},
		Core{Maintainer: "arduino", Arch: "sam", Version: "1.6.7"},
	}

	boardsManagerCores := []Core{
		Core{Maintainer: "arduino", Arch: "samd", Version: "1.6.5"},
	}

	boardsManagerRedBearCores := []Core{
		Core{Maintainer: "RedBearLab", Arch: "avr", Version: "1.0.0", Url: "https://redbearlab.github.io/arduino/Blend/blend_boards.zip"},
	}

	toolsMultipleVersions := []Tool{
		Tool{Name: "bossac", Version: "1.6.1-arduino"},
		Tool{Name: "bossac", Version: "1.5-arduino"},
	}

	tools := []Tool{
		Tool{Name: "avrdude", Version: "6.0.1-arduino5"},
		Tool{Name: "avr-gcc", Version: "4.8.1-arduino5"},
		Tool{Name: "arm-none-eabi-gcc", Version: "4.8.3-2014q1"},
		Tool{Name: "ctags", Version: "5.8-arduino11",
			OsUrls: []OsUrl{
				OsUrl{Os: "i686-pc-linux-gnu", Url: "http://downloads.arduino.cc/tools/ctags-5.8-arduino11-i686-pc-linux-gnu.tar.bz2"},
				OsUrl{Os: "x86_64-pc-linux-gnu", Url: "http://downloads.arduino.cc/tools/ctags-5.8-arduino11-x86_64-pc-linux-gnu.tar.bz2"},
				OsUrl{Os: "i686-mingw32", Url: "http://downloads.arduino.cc/tools/ctags-5.8-arduino11-i686-mingw32.zip"},
				OsUrl{Os: "x86_64-apple-darwin", Url: "http://downloads.arduino.cc/tools/ctags-5.8-arduino11-x86_64-apple-darwin.zip"},
				OsUrl{Os: "arm-linux-gnueabihf", Url: "http://downloads.arduino.cc/tools/ctags-5.8-arduino11-armv6-linux-gnueabihf.tar.bz2"},
				OsUrl{Os: "aarch64-linux-gnu", Url: "http://downloads.arduino.cc/tools/ctags-5.8-arduino11-aarch64-linux-gnu.tar.bz2"},
			},
		},
		Tool{Name: "arduino-preprocessor", Version: "0.1.5",
			OsUrls: []OsUrl{
				OsUrl{Os: "i686-pc-linux-gnu", Url: "https://github.com/arduino/arduino-preprocessor/releases/download/0.1.5/arduino-preprocessor-0.1.5-i686-pc-linux-gnu.tar.bz2"},
				OsUrl{Os: "x86_64-pc-linux-gnu", Url: "https://github.com/arduino/arduino-preprocessor/releases/download/0.1.5/arduino-preprocessor-0.1.5-x86_64-pc-linux-gnu.tar.bz2"},
				OsUrl{Os: "i686-mingw32", Url: "https://github.com/arduino/arduino-preprocessor/releases/download/0.1.5/arduino-preprocessor-0.1.5-i686-w64-mingw32.tar.bz2"},
				OsUrl{Os: "x86_64-apple-darwin", Url: "https://github.com/arduino/arduino-preprocessor/releases/download/0.1.5/arduino-preprocessor-0.1.5-x86_64-apple-darwin11.tar.bz2"},
				OsUrl{Os: "arm-linux-gnueabihf", Url: "https://github.com/arduino/arduino-preprocessor/releases/download/0.1.5/arduino-preprocessor-0.1.5-arm-linux-gnueabihf.tar.bz2"},
				OsUrl{Os: "aarch64-linux-gnu", Url: "https://github.com/arduino/arduino-preprocessor/releases/download/0.1.5/arduino-preprocessor-0.1.5-aarch64-linux-gnu.tar.bz2"},
			},
		},
	}

	boardsManagerTools := []Tool{
		Tool{Name: "openocd", Version: "0.9.0-arduino", Package: "arduino"},
		Tool{Name: "CMSIS", Version: "4.0.0-atmel", Package: "arduino"},
	}

	boardsManagerRFduinoTools := []Tool{
		Tool{Name: "arm-none-eabi-gcc", Version: "4.8.3-2014q1", Package: "RFduino"},
	}

	libraries := []Library{
		Library{Name: "Audio", Version: "1.0.4"},
		Library{Name: "Adafruit PN532", Version: "1.0.0"},
		Library{Name: "Bridge", Version: "1.6.1"},
		Library{Name: "CapacitiveSensor", Version: "0.5.0", VersionInLibProperties: "0.5"},
		Library{Name: "Ethernet", Version: "1.1.1"},
		Library{Name: "Robot IR Remote", Version: "1.0.2"},
		Library{Name: "FastLED", Version: "3.1.0"},
	}

	download(t, cores, boardsManagerCores, boardsManagerRedBearCores, tools, toolsMultipleVersions, boardsManagerTools, boardsManagerRFduinoTools, libraries)

	patchFiles(t)
}

func patchFiles(t *testing.T) {
	err := utils.EnsureFolderExists(PATCHES_FOLDER)
	NoError(t, err)
	files, err := ioutil.ReadDir(PATCHES_FOLDER)
	NoError(t, err)

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".patch" {
			panic("Patching for downloaded tools is not available! (see https://github.com/arduino/arduino-builder/issues/147)")
			// XXX: find an alternative to x/codereview/patch
			// https://github.com/arduino/arduino-builder/issues/147
			/*
				data, err := ioutil.ReadFile(Abs(t, filepath.Join(PATCHES_FOLDER, file.Name())))
				NoError(t, err)
				patchSet, err := patch.Parse(data)
				NoError(t, err)
				operations, err := patchSet.Apply(ioutil.ReadFile)
				for _, op := range operations {
					utils.WriteFileBytes(op.Dst, op.Data)
				}
			*/
		}
	}
}

func download(t *testing.T, cores, boardsManagerCores, boardsManagerRedBearCores []Core, tools, toolsMultipleVersions, boardsManagerTools, boardsManagerRFduinoTools []Tool, libraries []Library) {
	allCoresDownloaded, err := allCoresAlreadyDownloadedAndUnpacked(HARDWARE_FOLDER, cores)
	NoError(t, err)
	if allCoresDownloaded &&
		allBoardsManagerCoresAlreadyDownloadedAndUnpacked(BOARD_MANAGER_FOLDER, boardsManagerCores) &&
		allBoardsManagerCoresAlreadyDownloadedAndUnpacked(BOARD_MANAGER_FOLDER, boardsManagerRedBearCores) &&
		allBoardsManagerToolsAlreadyDownloadedAndUnpacked(BOARD_MANAGER_FOLDER, boardsManagerTools) &&
		allBoardsManagerToolsAlreadyDownloadedAndUnpacked(BOARD_MANAGER_FOLDER, boardsManagerRFduinoTools) &&
		allToolsAlreadyDownloadedAndUnpacked(TOOLS_FOLDER, tools) &&
		allToolsAlreadyDownloadedAndUnpacked(TOOLS_FOLDER, toolsMultipleVersions) &&
		allLibrariesAlreadyDownloadedAndUnpacked(LIBRARIES_FOLDER, libraries) {
		return
	}

	index, err := downloadIndex("http://downloads.arduino.cc/packages/package_index.json")
	NoError(t, err)

	err = downloadCores(cores, index)
	NoError(t, err)

	err = downloadBoardManagerCores(boardsManagerCores, index)
	NoError(t, err)

	err = downloadTools(tools, index)
	NoError(t, err)

	err = downloadToolsMultipleVersions(toolsMultipleVersions, index)
	NoError(t, err)

	err = downloadBoardsManagerTools(boardsManagerTools, index)
	NoError(t, err)

	rfduinoIndex, err := downloadIndex("http://downloads.arduino.cc/packages/test_package_rfduino_index.json")
	NoError(t, err)

	err = downloadBoardsManagerTools(boardsManagerRFduinoTools, rfduinoIndex)
	NoError(t, err)

	err = downloadBoardManagerCores(boardsManagerRedBearCores, nil)
	NoError(t, err)

	librariesIndex, err := downloadIndex("http://downloads.arduino.cc/libraries/library_index.json")
	NoError(t, err)

	err = downloadLibraries(libraries, librariesIndex)
	NoError(t, err)
}

func downloadIndex(url string) (map[string]interface{}, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body.Close()

	index := make(map[string]interface{})
	err = json.Unmarshal(bytes, &index)
	if err != nil {
		return nil, err
	}

	return index, nil
}

func downloadCores(cores []Core, index map[string]interface{}) error {
	for _, core := range cores {
		url, err := findCoreUrl(index, core)
		if err != nil {
			return i18n.WrapError(err)
		}
		err = downloadAndUnpackCore(core, url, HARDWARE_FOLDER)
		if err != nil {
			return i18n.WrapError(err)
		}
	}
	return nil
}

func downloadBoardManagerCores(cores []Core, index map[string]interface{}) error {
	for _, core := range cores {
		url, err := findCoreUrl(index, core)
		if err != nil {
			return i18n.WrapError(err)
		}
		err = downloadAndUnpackBoardManagerCore(core, url, BOARD_MANAGER_FOLDER)
		if err != nil {
			return i18n.WrapError(err)
		}
	}
	return nil
}

func findCoreUrl(index map[string]interface{}, core Core) (string, error) {
	if core.Url != "" {
		return core.Url, nil
	}
	packages := index["packages"].([]interface{})
	for _, p := range packages {
		pack := p.(map[string]interface{})
		if pack[constants.PACKAGE_NAME].(string) == core.Maintainer {
			packagePlatforms := pack["platforms"].([]interface{})
			for _, pt := range packagePlatforms {
				packagePlatform := pt.(map[string]interface{})
				if packagePlatform[constants.PLATFORM_ARCHITECTURE] == core.Arch && packagePlatform[constants.PLATFORM_VERSION] == core.Version {
					return packagePlatform[constants.PLATFORM_URL].(string), nil
				}
			}
		}
	}

	return constants.EMPTY_STRING, errors.Errorf("Unable to find tool " + core.Maintainer + " " + core.Arch + " " + core.Version)
}

func downloadTools(tools []Tool, index map[string]interface{}) error {
	host := translateGOOSGOARCHToPackageIndexValue()

	for _, tool := range tools {
		url, err := findToolUrl(index, tool, host)
		if err != nil {
			return i18n.WrapError(err)
		}
		err = downloadAndUnpackTool(tool, url, TOOLS_FOLDER, true)
		if err != nil {
			return i18n.WrapError(err)
		}
	}

	return nil
}

func downloadToolsMultipleVersions(tools []Tool, index map[string]interface{}) error {
	host := translateGOOSGOARCHToPackageIndexValue()

	for _, tool := range tools {
		if !toolAlreadyDownloadedAndUnpacked(TOOLS_FOLDER, tool) {
			_, err := os.Stat(filepath.Join(TOOLS_FOLDER, tool.Name))
			if err == nil {
				err = os.RemoveAll(filepath.Join(TOOLS_FOLDER, tool.Name))
				if err != nil {
					return i18n.WrapError(err)
				}
			}
		}
	}

	for _, tool := range tools {
		url, err := findToolUrl(index, tool, host)
		if err != nil {
			return i18n.WrapError(err)
		}
		err = downloadAndUnpackTool(tool, url, TOOLS_FOLDER, false)
		if err != nil {
			return i18n.WrapError(err)
		}
	}

	return nil
}

func downloadBoardsManagerTools(tools []Tool, index map[string]interface{}) error {
	host := translateGOOSGOARCHToPackageIndexValue()

	for _, tool := range tools {
		url, err := findToolUrl(index, tool, host)
		if err != nil {
			return i18n.WrapError(err)
		}
		err = downloadAndUnpackBoardsManagerTool(tool, url, BOARD_MANAGER_FOLDER)
		if err != nil {
			return i18n.WrapError(err)
		}
	}

	return nil
}

func allBoardsManagerCoresAlreadyDownloadedAndUnpacked(targetPath string, cores []Core) bool {
	for _, core := range cores {
		if !boardManagerCoreAlreadyDownloadedAndUnpacked(targetPath, core) {
			return false
		}
	}
	return true
}

func boardManagerCoreAlreadyDownloadedAndUnpacked(targetPath string, core Core) bool {
	_, err := os.Stat(filepath.Join(targetPath, core.Maintainer, "hardware", core.Arch, core.Version))
	return !os.IsNotExist(err)
}

func allCoresAlreadyDownloadedAndUnpacked(targetPath string, cores []Core) (bool, error) {
	for _, core := range cores {
		alreadyDownloaded, err := coreAlreadyDownloadedAndUnpacked(targetPath, core)
		if err != nil {
			return false, i18n.WrapError(err)
		}
		if !alreadyDownloaded {
			return false, nil
		}
	}
	return true, nil
}

func coreAlreadyDownloadedAndUnpacked(targetPath string, core Core) (bool, error) {
	corePath := filepath.Join(targetPath, core.Maintainer, core.Arch)

	_, err := os.Stat(corePath)
	if os.IsNotExist(err) {
		return false, nil
	}
	platform, err := properties.Load(filepath.Join(corePath, "platform.txt"))
	if err != nil {
		return false, i18n.WrapError(err)
	}

	if core.Version != platform["version"] {
		err := os.RemoveAll(corePath)
		return false, i18n.WrapError(err)
	}

	return true, nil
}

func allBoardsManagerToolsAlreadyDownloadedAndUnpacked(targetPath string, tools []Tool) bool {
	for _, tool := range tools {
		if !boardManagerToolAlreadyDownloadedAndUnpacked(targetPath, tool) {
			return false
		}
	}
	return true
}

func boardManagerToolAlreadyDownloadedAndUnpacked(targetPath string, tool Tool) bool {
	_, err := os.Stat(filepath.Join(targetPath, tool.Package, constants.FOLDER_TOOLS, tool.Name, tool.Version))
	return !os.IsNotExist(err)
}

func allToolsAlreadyDownloadedAndUnpacked(targetPath string, tools []Tool) bool {
	for _, tool := range tools {
		if !toolAlreadyDownloadedAndUnpacked(targetPath, tool) {
			return false
		}
	}
	return true
}

func toolAlreadyDownloadedAndUnpacked(targetPath string, tool Tool) bool {
	_, err := os.Stat(filepath.Join(targetPath, tool.Name, tool.Version))
	return !os.IsNotExist(err)
}

func allLibrariesAlreadyDownloadedAndUnpacked(targetPath string, libraries []Library) bool {
	for _, library := range libraries {
		if !libraryAlreadyDownloadedAndUnpacked(targetPath, library) {
			return false
		}
	}
	return true
}

func libraryAlreadyDownloadedAndUnpacked(targetPath string, library Library) bool {
	_, err := os.Stat(filepath.Join(targetPath, strings.Replace(library.Name, " ", "_", -1)))
	if os.IsNotExist(err) {
		return false
	}

	libProps, err := properties.Load(filepath.Join(targetPath, strings.Replace(library.Name, " ", "_", -1), "library.properties"))
	if err != nil {
		return false
	}
	return libProps["version"] == library.Version || libProps["version"] == library.VersionInLibProperties
}

func downloadAndUnpackCore(core Core, url string, targetPath string) error {
	alreadyDownloaded, err := coreAlreadyDownloadedAndUnpacked(targetPath, core)
	if err != nil {
		return i18n.WrapError(err)
	}
	if alreadyDownloaded {
		return nil
	}

	targetPath, err = filepath.Abs(targetPath)
	if err != nil {
		return i18n.WrapError(err)
	}

	unpackFolder, files, err := downloadAndUnpack(url)
	if err != nil {
		return i18n.WrapError(err)
	}
	defer os.RemoveAll(unpackFolder)

	_, err = os.Stat(filepath.Join(targetPath, core.Maintainer, core.Arch))
	if err == nil {
		err = os.RemoveAll(filepath.Join(targetPath, core.Maintainer, core.Arch))
		if err != nil {
			return i18n.WrapError(err)
		}
	}

	if len(files) == 1 && files[0].IsDir() {
		err = utils.EnsureFolderExists(filepath.Join(targetPath, core.Maintainer))
		if err != nil {
			return i18n.WrapError(err)
		}
		err = copyRecursive(filepath.Join(unpackFolder, files[0].Name()), filepath.Join(targetPath, core.Maintainer, core.Arch))
		if err != nil {
			return i18n.WrapError(err)
		}
	} else {
		err = utils.EnsureFolderExists(filepath.Join(targetPath, core.Maintainer, core.Arch))
		if err != nil {
			return i18n.WrapError(err)
		}
		for _, file := range files {
			err = copyRecursive(filepath.Join(unpackFolder, file.Name()), filepath.Join(targetPath, core.Maintainer, core.Arch, file.Name()))
			if err != nil {
				return i18n.WrapError(err)
			}
		}
	}

	return nil
}

func downloadAndUnpackBoardManagerCore(core Core, url string, targetPath string) error {
	if boardManagerCoreAlreadyDownloadedAndUnpacked(targetPath, core) {
		return nil
	}

	targetPath, err := filepath.Abs(targetPath)
	if err != nil {
		return i18n.WrapError(err)
	}

	unpackFolder, files, err := downloadAndUnpack(url)
	if err != nil {
		return i18n.WrapError(err)
	}
	defer os.RemoveAll(unpackFolder)

	_, err = os.Stat(filepath.Join(targetPath, core.Maintainer, "hardware", core.Arch))
	if err == nil {
		err = os.RemoveAll(filepath.Join(targetPath, core.Maintainer, "hardware", core.Arch))
		if err != nil {
			return i18n.WrapError(err)
		}
	}

	if len(files) == 1 && files[0].IsDir() {
		err = utils.EnsureFolderExists(filepath.Join(targetPath, core.Maintainer, "hardware", core.Arch))
		if err != nil {
			return i18n.WrapError(err)
		}
		err = copyRecursive(filepath.Join(unpackFolder, files[0].Name()), filepath.Join(targetPath, core.Maintainer, "hardware", core.Arch, core.Version))
		if err != nil {
			return i18n.WrapError(err)
		}
	} else {
		err = utils.EnsureFolderExists(filepath.Join(targetPath, core.Maintainer, "hardware", core.Arch, core.Version))
		if err != nil {
			return i18n.WrapError(err)
		}
		for _, file := range files {
			err = copyRecursive(filepath.Join(unpackFolder, file.Name()), filepath.Join(targetPath, core.Maintainer, "hardware", core.Arch, core.Version, file.Name()))
			if err != nil {
				return i18n.WrapError(err)
			}
		}
	}

	return nil
}

func downloadAndUnpackBoardsManagerTool(tool Tool, url string, targetPath string) error {
	if boardManagerToolAlreadyDownloadedAndUnpacked(targetPath, tool) {
		return nil
	}

	targetPath, err := filepath.Abs(targetPath)
	if err != nil {
		return i18n.WrapError(err)
	}

	unpackFolder, files, err := downloadAndUnpack(url)
	if err != nil {
		return i18n.WrapError(err)
	}
	defer os.RemoveAll(unpackFolder)

	if len(files) == 1 && files[0].IsDir() {
		err = utils.EnsureFolderExists(filepath.Join(targetPath, tool.Package, constants.FOLDER_TOOLS, tool.Name))
		if err != nil {
			return i18n.WrapError(err)
		}
		err = copyRecursive(filepath.Join(unpackFolder, files[0].Name()), filepath.Join(targetPath, tool.Package, constants.FOLDER_TOOLS, tool.Name, tool.Version))
		if err != nil {
			return i18n.WrapError(err)
		}
	} else {
		err = utils.EnsureFolderExists(filepath.Join(targetPath, tool.Package, constants.FOLDER_TOOLS, tool.Name, tool.Version))
		if err != nil {
			return i18n.WrapError(err)
		}
		for _, file := range files {
			err = copyRecursive(filepath.Join(unpackFolder, file.Name()), filepath.Join(targetPath, tool.Package, constants.FOLDER_TOOLS, tool.Name, tool.Version, file.Name()))
			if err != nil {
				return i18n.WrapError(err)
			}
		}
	}

	return nil
}

func downloadAndUnpackTool(tool Tool, url string, targetPath string, deleteIfMissing bool) error {
	if toolAlreadyDownloadedAndUnpacked(targetPath, tool) {
		return nil
	}

	targetPath, err := filepath.Abs(targetPath)
	if err != nil {
		return i18n.WrapError(err)
	}

	unpackFolder, files, err := downloadAndUnpack(url)
	if err != nil {
		return i18n.WrapError(err)
	}
	defer os.RemoveAll(unpackFolder)

	if deleteIfMissing {
		_, err = os.Stat(filepath.Join(targetPath, tool.Name))
		if err == nil {
			err = os.RemoveAll(filepath.Join(targetPath, tool.Name))
			if err != nil {
				return i18n.WrapError(err)
			}
		}
	}

	if len(files) == 1 && files[0].IsDir() {
		err = utils.EnsureFolderExists(filepath.Join(targetPath, tool.Name))
		if err != nil {
			return i18n.WrapError(err)
		}
		err = copyRecursive(filepath.Join(unpackFolder, files[0].Name()), filepath.Join(targetPath, tool.Name, tool.Version))
		if err != nil {
			return i18n.WrapError(err)
		}
	} else {
		err = utils.EnsureFolderExists(filepath.Join(targetPath, tool.Name, tool.Version))
		if err != nil {
			return i18n.WrapError(err)
		}
		for _, file := range files {
			err = copyRecursive(filepath.Join(unpackFolder, file.Name()), filepath.Join(targetPath, tool.Name, tool.Version, file.Name()))
			if err != nil {
				return i18n.WrapError(err)
			}
		}
	}

	return nil
}

func downloadAndUnpack(url string) (string, []os.FileInfo, error) {
	fmt.Fprintln(os.Stderr, "Downloading "+url)

	unpackFolder, err := ioutil.TempDir(constants.EMPTY_STRING, "arduino-builder-tool")
	if err != nil {
		return constants.EMPTY_STRING, nil, i18n.WrapError(err)
	}

	urlParts := strings.Split(url, "/")
	archiveFileName := urlParts[len(urlParts)-1]
	archiveFilePath := filepath.Join(unpackFolder, archiveFileName)

	res, err := http.Get(url)
	if err != nil {
		return constants.EMPTY_STRING, nil, i18n.WrapError(err)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return constants.EMPTY_STRING, nil, i18n.WrapError(err)
	}
	res.Body.Close()

	utils.WriteFileBytes(archiveFilePath, bytes)

	cmd := buildUnpackCmd(archiveFilePath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return constants.EMPTY_STRING, nil, i18n.WrapError(err)
	}
	if len(out) > 0 {
		fmt.Println(string(out))
	}

	os.Remove(archiveFilePath)

	files, err := gohasissues.ReadDir(unpackFolder)
	if err != nil {
		return constants.EMPTY_STRING, nil, i18n.WrapError(err)
	}

	return unpackFolder, files, nil
}

func buildUnpackCmd(file string) *exec.Cmd {
	var cmd *exec.Cmd
	if strings.HasSuffix(file, "zip") {
		cmd = exec.Command("unzip", "-qq", filepath.Base(file))
	} else {
		cmd = exec.Command("tar", "xf", filepath.Base(file))
	}
	cmd.Dir = filepath.Dir(file)
	return cmd
}

func translateGOOSGOARCHToPackageIndexValue() []string {
	switch value := runtime.GOOS + "-" + runtime.GOARCH; value {
	case "linux-amd64":
		return []string{"x86_64-pc-linux-gnu", "x86_64-linux-gnu"}
	case "linux-386":
		return []string{"i686-pc-linux-gnu", "i686-linux-gnu"}
	case "windows-amd64":
		return []string{"i686-mingw32", "i686-cygwin"}
	case "windows-386":
		return []string{"i686-mingw32", "i686-cygwin"}
	case "darwin-amd64":
		return []string{"i386-apple-darwin11", "x86_64-apple-darwin"}
	case "linux-arm":
		return []string{"arm-linux-gnueabihf"}
	default:
		panic("Unknown OS: " + value)
	}
}

func findToolUrl(index map[string]interface{}, tool Tool, host []string) (string, error) {
	if len(tool.OsUrls) > 0 {
		for _, osUrl := range tool.OsUrls {
			if utils.SliceContains(host, osUrl.Os) {
				return osUrl.Url, nil
			}
		}
	} else {
		packages := index["packages"].([]interface{})
		for _, p := range packages {
			pack := p.(map[string]interface{})
			packageTools := pack[constants.PACKAGE_TOOLS].([]interface{})
			for _, pt := range packageTools {
				packageTool := pt.(map[string]interface{})
				name := packageTool[constants.TOOL_NAME].(string)
				version := packageTool[constants.TOOL_VERSION].(string)
				if name == tool.Name && version == tool.Version {
					systems := packageTool["systems"].([]interface{})
					for _, s := range systems {
						system := s.(map[string]interface{})
						if utils.SliceContains(host, system["host"].(string)) {
							return system[constants.TOOL_URL].(string), nil
						}
					}
				}
			}
		}
	}

	return constants.EMPTY_STRING, errors.Errorf("Unable to find tool " + tool.Name + " " + tool.Version)
}

func downloadLibraries(libraries []Library, index map[string]interface{}) error {
	for _, library := range libraries {
		url, err := findLibraryUrl(index, library)
		if err != nil {
			return i18n.WrapError(err)
		}
		err = downloadAndUnpackLibrary(library, url, LIBRARIES_FOLDER)
		if err != nil {
			return i18n.WrapError(err)
		}
	}

	return nil
}

func findLibraryUrl(index map[string]interface{}, library Library) (string, error) {
	if library.Url != "" {
		return library.Url, nil
	}
	libs := index["libraries"].([]interface{})
	for _, l := range libs {
		lib := l.(map[string]interface{})
		if library.Name == lib["name"].(string) && library.Version == lib["version"].(string) {
			return lib["url"].(string), nil
		}
	}

	return constants.EMPTY_STRING, errors.Errorf("Unable to find library " + library.Name + " " + library.Version)
}

func downloadAndUnpackLibrary(library Library, url string, targetPath string) error {
	if libraryAlreadyDownloadedAndUnpacked(targetPath, library) {
		return nil
	}

	targetPath, err := filepath.Abs(targetPath)
	if err != nil {
		return i18n.WrapError(err)
	}

	unpackFolder, files, err := downloadAndUnpack(url)
	if err != nil {
		return i18n.WrapError(err)
	}
	defer os.RemoveAll(unpackFolder)

	_, err = os.Stat(filepath.Join(targetPath, strings.Replace(library.Name, " ", "_", -1)))
	if err == nil {
		err = os.RemoveAll(filepath.Join(targetPath, strings.Replace(library.Name, " ", "_", -1)))
		if err != nil {
			return i18n.WrapError(err)
		}
	}

	err = copyRecursive(filepath.Join(unpackFolder, files[0].Name()), filepath.Join(targetPath, strings.Replace(library.Name, " ", "_", -1)))
	if err != nil {
		return i18n.WrapError(err)
	}

	return nil
}

func copyRecursive(from, to string) error {
	copyFunc := func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(from, currentPath)
		if err != nil {
			return i18n.WrapError(err)
		}
		targetPath := filepath.Join(to, rel)
		if info.IsDir() {
			err := os.MkdirAll(targetPath, info.Mode())
			if err != nil {
				return i18n.WrapError(err)
			}
		} else if info.Mode().IsRegular() {
			fromFile, err := os.Open(currentPath)
			if err != nil {
				return i18n.WrapError(err)
			}
			defer fromFile.Close()
			targetFile, err := os.Create(targetPath)
			if err != nil {
				return i18n.WrapError(err)
			}
			defer targetFile.Close()
			_, err = io.Copy(targetFile, fromFile)
			if err != nil {
				return i18n.WrapError(err)
			}
			err = os.Chmod(targetPath, info.Mode())
			if err != nil {
				return i18n.WrapError(err)
			}
		} else if info.Mode()&os.ModeSymlink == os.ModeSymlink {
			linkedFile, err := os.Readlink(currentPath)
			if err != nil {
				return i18n.WrapError(err)
			}
			fromFile := filepath.Join(filepath.Dir(targetPath), linkedFile)
			err = os.Symlink(fromFile, targetPath)
			if err != nil {
				return i18n.WrapError(err)
			}
		} else {
			return errors.Errorf("unable to copy file " + currentPath)
		}

		return nil
	}
	err := gohasissues.Walk(from, copyFunc)
	return i18n.WrapError(err)
}
