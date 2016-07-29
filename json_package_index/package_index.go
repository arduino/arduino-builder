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
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/arduino/arduino-builder/constants"
	_ "github.com/arduino/arduino-builder/i18n"
	properties "github.com/arduino/go-properties-map"
)

type core struct {
	Architecture string `json:"architecture"`
	Version      string `json:"version"`
	URL          string `json:"url"`
	Maintainer   string `json:"maintainer"`
	Name         string `json:"archiveFileName"`
	Checksum     string `json:"checksum"`
	destination  string
	Dependencies []struct {
		Packager string `json:"packager"`
		Name     string `json:"name"`
		Version  string `json:"version"`
	} `json:"toolsDependencies"`
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

func PackageIndexFoldersToPropertiesMap(folders []string) (map[string]properties.Map, error) {

	var paths []string
	for _, folder := range folders {
		folder, err := filepath.Abs(folder)
		if err != nil {
			break
		}
		files, _ := ioutil.ReadDir(folder)
		for _, f := range files {
			if strings.HasPrefix(f.Name(), "package") && strings.HasSuffix(f.Name(), "index.json") {
				paths = append(paths, filepath.Join(folder, f.Name()))
			}
		}
	}
	return PackageIndexesToPropertiesMap(paths)
}

func PackageIndexesToPropertiesMap(urls []string) (map[string]properties.Map, error) {

	globalProperties = make(map[string]properties.Map)

	data, err := PackageIndexesToGlobalIndex(urls)

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
			globalProperties[p.Name+":"+a.Architecture+":"+a.Version] = localProperties.Clone()
		}
	}
	return globalProperties, err
}

func PackageIndexesToGlobalIndex(urls []string) (index, error) {

	// firststub of arduino-pdpm
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
	return data, err
}
