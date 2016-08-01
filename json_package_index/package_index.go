/*
 * This file is part of Arduino Builder.
 *
 * Copyright 2016 Arduino LLC (http://www.arduino.cc/)
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
 */

package json_package_index

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/blang/semver"

	"github.com/arduino/arduino-builder/constants"
	_ "github.com/arduino/arduino-builder/i18n"
	"github.com/arduino/arduino-builder/types"
	"github.com/arduino/arduino-builder/utils"
	"github.com/arduino/go-properties-map"
)

type core struct {
	Architecture string `json:"architecture"`
	Version      string `json:"version"`
	URL          string `json:"url"`
	Maintainer   string `json:"maintainer"`
	Name         string `json:"archiveFileName"`
	Checksum     string `json:"checksum"`
	destination  string
	installed    bool
	Dependencies []struct {
		Packager string `json:"packager"`
		Name     string `json:"name"`
		Version  string `json:"version"`
	} `json:"toolsDependencies"`
	CoreDependencies []struct {
		Packager string `json:"packager"`
	} `json:"coreDependencies"`
}

type tool struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Systems []struct {
		Host     string `json:"host"`
		URL      string `json:"url"`
		Name     string `json:"archiveFileName"`
		Checksum string `json:"checksum"`
	} `json:"systems"`
	url         string
	destination string
}

type index struct {
	Packages []struct {
		Name       string `json:"name"`
		Maintainer string `json:"maintainer"`
		Platforms  []core `json:"platforms"`
		Tools      []tool `json:"tools"`
	} `json:"packages"`
}

var systems = map[string]string{
	"linuxamd64":  "x86_64-linux-gnu",
	"linux386":    "i686-linux-gnu",
	"darwinamd64": "apple-darwin",
	"windows386":  "i686-mingw32",
}

// globalProperties is a big map of properties maps in the form
// globalProperties["arduino:avr:1.6.12"] = usual properties Map
// at compile time, when de board is well defined, the relevant map
// should be merged with the "classic" map overriding its values

var globalProperties map[string]properties.Map

func PackageIndexFoldersToPropertiesMap(packages *types.Packages, folders []string, specifiedFilenames []string) (map[string]properties.Map, error) {

	var paths []string

	for _, folder := range folders {
		folder, err := filepath.Abs(folder)
		if err != nil {
			break
		}
		files, _ := ioutil.ReadDir(folder)
		for _, f := range files {
			if strings.HasPrefix(f.Name(), "package") && strings.HasSuffix(f.Name(), "index.json") {
				// if a list of required json has been provided only add them
				if specifiedFilenames != nil && len(specifiedFilenames) > 1 &&
					!utils.SliceContains(specifiedFilenames, f.Name()) {
					continue
				} else {
					paths = append(paths, filepath.Join(folder, f.Name()))
				}
			}
		}
	}
	return PackageIndexesToPropertiesMap(packages, paths)
}

func PackageIndexesToPropertiesMap(packages *types.Packages, urls []string) (map[string]properties.Map, error) {

	globalProperties = make(map[string]properties.Map)
	coreDependencyMap := make(map[string]string)

	data, err := PackageIndexesToGlobalIndex(packages, urls)

	for _, p := range data.Packages {
		for _, a := range p.Platforms {
			localProperties := make(properties.Map)
			for _, dep := range a.Dependencies {
				localProperties[constants.BUILD_PROPERTIES_RUNTIME_TOOLS_PREFIX+dep.Name+constants.BUILD_PROPERTIES_RUNTIME_TOOLS_SUFFIX] =
					"{" + constants.BUILD_PROPERTIES_RUNTIME_TOOLS_PREFIX + dep.Name + "-" + dep.Version + constants.BUILD_PROPERTIES_RUNTIME_TOOLS_SUFFIX + "}"
				if dep.Packager != p.Name {
					localProperties[constants.BUILD_PROPERTIES_RUNTIME_TOOLS_PREFIX+dep.Name+"-"+dep.Version+constants.BUILD_PROPERTIES_RUNTIME_TOOLS_SUFFIX] =
						"{" + constants.BUILD_PROPERTIES_RUNTIME_TOOLS_PREFIX + dep.Name + "-" + dep.Packager + "-" + dep.Version + constants.BUILD_PROPERTIES_RUNTIME_TOOLS_SUFFIX + "}"
				}
			}
			for _, coredep := range a.CoreDependencies {
				// inherit all the tools from latest coredep
				if err == nil {
					coreDependencyMap[p.Name+":"+a.Architecture+":"+a.Version] =
						coredep.Packager + ":" + a.Architecture
				}
			}
			globalProperties[p.Name+":"+a.Architecture+":"+a.Version] = localProperties.Clone()
		}
	}

	for idx, parentCore := range coreDependencyMap {
		version, err := findLatestInstalledCore(data, strings.Split(parentCore, ":")[0], strings.Split(parentCore, ":")[1])
		if err == nil {
			globalProperties[idx] = globalProperties[parentCore+":"+version].Clone()
		}
	}

	return globalProperties, err
}

func findLatestInstalledCore(data index, Packager string, Name string) (string, error) {
	latest, _ := semver.Make("0.0.0")
	for _, p := range data.Packages {
		for _, a := range p.Platforms {
			if p.Name == Packager && a.Architecture == Name {
				test, _ := semver.Make(a.Version)
				if test.GT(latest) && a.installed {
					latest = test
				}
			}
		}
	}
	var err error
	test, _ := semver.Make("0.0.0")
	if latest.EQ(test) {
		err = errors.New("No such core available")
	} else {
		err = nil
	}
	return latest.String(), err
}

func PackageIndexesToGlobalIndex(packages *types.Packages, urls []string) (index, error) {

	// first stub of arduino-pdpm
	var data index
	var err error

	for _, url := range urls {

		var body []byte
		var localdata index
		localpath, _ := filepath.Abs(url)
		_, err := os.Stat(localpath)

		if err != nil {
			resp, err := http.Get(url)
			if err == nil {
				defer resp.Body.Close()
				body, err = ioutil.ReadAll(resp.Body)
				if err != nil {
					break
				}
			}
		} else {
			body, err = ioutil.ReadFile(localpath)
			if err != nil {
				break
			}
		}
		json.Unmarshal(body, &localdata)
		for _, entry := range localdata.Packages {
			data.Packages = append(data.Packages, entry)
		}
	}

	for i, p := range data.Packages {
		for j, a := range p.Platforms {
			if packages != nil && packages.Packages[p.Name] != nil &&
				packages.Packages[p.Name].Platforms[a.Architecture] != nil &&
				packages.Packages[p.Name].Platforms[a.Architecture].Properties["version"] == a.Version {
				data.Packages[i].Platforms[j].installed = true
			}
		}
	}

	return data, err
}

func CompareVersions(fv string, sv string) int {
	v1, _ := semver.Make(fv)
	v2, _ := semver.Make(sv)
	if v1.EQ(v2) {
		return 0
	}
	if v1.GT(v2) {
		return 1
	} else {
		return -1
	}
}
