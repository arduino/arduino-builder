package json_package_index

import (
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func TestPropertiesPackageIndex(t *testing.T) {

	var paths []string
	paths = append(paths, filepath.Join("testdata", "package_index.json"))

	p, err := PackageIndexesToPropertiesMap(nil, paths)

	require.NoError(t, err)

	require.Equal(t, "{runtime.tools.avr-gcc-4.9.2-atmel3.5.3-arduino2.path}", p["arduino:avr:1.6.12"]["runtime.tools.avr-gcc.path"])
}

func TestPropertiesPackageIndexRemote(t *testing.T) {

	var paths []string
	paths = append(paths, filepath.Join("testdata", "package_index.json"))
	paths = append(paths, "http://downloads.arduino.cc/packages/package_arduino.cc_linux_index.json")

	p, err := PackageIndexesToPropertiesMap(nil, paths)

	require.NoError(t, err)

	require.Equal(t, "{runtime.tools.avr-gcc-4.9.2-atmel3.5.3-arduino2.path}", p["arduino:avr:1.6.12"]["runtime.tools.avr-gcc.path"])
	require.Equal(t, "{runtime.tools.linuxuploader-1.2.path}", p["arduino:arm_cortexA:0.4.0"]["runtime.tools.linuxuploader.path"])
}

func TestPackageIndexToGlobalIndex(t *testing.T) {

	var paths []string
	paths = append(paths, filepath.Join("testdata", "package_index.json"))

	p, err := PackageIndexesToGlobalIndex(nil, paths)
	require.NoError(t, err)

	require.Equal(t, "Arduino", p.Packages[0].Maintainer)
}

func TestPackageIndexFoldersToPropertiesMap(t *testing.T) {
	var paths []string
	paths = append(paths, "testdata")

	p, err := PackageIndexFoldersToPropertiesMap(nil, paths)
	require.NoError(t, err)

	require.Equal(t, "{runtime.tools.avr-gcc-4.9.2-atmel3.5.3-arduino2.path}", p["arduino:avr:1.6.12"]["runtime.tools.avr-gcc.path"])
}
