package test

import (
	"arduino.cc/builder/builder_utils"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func sleep(t *testing.T) {
	dur, err := time.ParseDuration("1s")
	NoError(t, err)
	time.Sleep(dur)
}

func tempFile(t *testing.T, prefix string) string {
	file, err := ioutil.TempFile("", prefix)
	NoError(t, err)
	return file.Name()
}

func TestObjFileIsUpToDateObjMissing(t *testing.T) {
	sourceFile := tempFile(t, "source")
	defer os.RemoveAll(sourceFile)

	upToDate, err := builder_utils.ObjFileIsUpToDate(sourceFile, "", "")
	NoError(t, err)
	require.False(t, upToDate)
}

func TestObjFileIsUpToDateDepMissing(t *testing.T) {
	sourceFile := tempFile(t, "source")
	defer os.RemoveAll(sourceFile)

	objFile := tempFile(t, "obj")
	defer os.RemoveAll(objFile)

	upToDate, err := builder_utils.ObjFileIsUpToDate(sourceFile, objFile, "")
	NoError(t, err)
	require.False(t, upToDate)
}

func TestObjFileIsUpToDateObjOlder(t *testing.T) {
	objFile := tempFile(t, "obj")
	defer os.RemoveAll(objFile)
	depFile := tempFile(t, "dep")
	defer os.RemoveAll(depFile)

	sleep(t)

	sourceFile := tempFile(t, "source")
	defer os.RemoveAll(sourceFile)

	upToDate, err := builder_utils.ObjFileIsUpToDate(sourceFile, objFile, depFile)
	NoError(t, err)
	require.False(t, upToDate)
}

func TestObjFileIsUpToDateObjNewer(t *testing.T) {
	sourceFile := tempFile(t, "source")
	defer os.RemoveAll(sourceFile)

	sleep(t)

	objFile := tempFile(t, "obj")
	defer os.RemoveAll(objFile)
	depFile := tempFile(t, "dep")
	defer os.RemoveAll(depFile)

	upToDate, err := builder_utils.ObjFileIsUpToDate(sourceFile, objFile, depFile)
	NoError(t, err)
	require.True(t, upToDate)
}

func TestObjFileIsUpToDateDepIsNewer(t *testing.T) {
	sourceFile := tempFile(t, "source")
	defer os.RemoveAll(sourceFile)

	sleep(t)

	objFile := tempFile(t, "obj")
	defer os.RemoveAll(objFile)
	depFile := tempFile(t, "dep")
	defer os.RemoveAll(depFile)

	sleep(t)

	headerFile := tempFile(t, "header")
	defer os.RemoveAll(headerFile)

	ioutil.WriteFile(depFile, []byte(objFile+": \\\n\t"+sourceFile+" \\\n\t"+headerFile), os.FileMode(0644))

	upToDate, err := builder_utils.ObjFileIsUpToDate(sourceFile, objFile, depFile)
	NoError(t, err)
	require.False(t, upToDate)
}

func TestObjFileIsUpToDateDepIsOlder(t *testing.T) {
	sourceFile := tempFile(t, "source")
	defer os.RemoveAll(sourceFile)

	headerFile := tempFile(t, "header")
	defer os.RemoveAll(headerFile)

	sleep(t)

	objFile := tempFile(t, "obj")
	defer os.RemoveAll(objFile)
	depFile := tempFile(t, "dep")
	defer os.RemoveAll(depFile)

	ioutil.WriteFile(depFile, []byte(objFile+": \\\n\t"+sourceFile+" \\\n\t"+headerFile), os.FileMode(0644))

	upToDate, err := builder_utils.ObjFileIsUpToDate(sourceFile, objFile, depFile)
	NoError(t, err)
	require.True(t, upToDate)
}

func TestObjFileIsUpToDateDepIsWrong(t *testing.T) {
	sourceFile := tempFile(t, "source")
	defer os.RemoveAll(sourceFile)

	sleep(t)

	objFile := tempFile(t, "obj")
	defer os.RemoveAll(objFile)
	depFile := tempFile(t, "dep")
	defer os.RemoveAll(depFile)

	sleep(t)

	headerFile := tempFile(t, "header")
	defer os.RemoveAll(headerFile)

	ioutil.WriteFile(depFile, []byte(sourceFile+": \\\n\t"+sourceFile+" \\\n\t"+headerFile), os.FileMode(0644))

	upToDate, err := builder_utils.ObjFileIsUpToDate(sourceFile, objFile, depFile)
	NoError(t, err)
	require.False(t, upToDate)
}
