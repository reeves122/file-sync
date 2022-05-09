package file

import (
	"github.com/champ-oss/file-sync/pkg/common"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func Test_copyFile_Success(t *testing.T) {
	// Write a test file to a test source directory
	sourceDir, _ := ioutil.TempDir("", "source")
	defer common.RemoveDir(sourceDir)
	sourceFile := filepath.Join(sourceDir, "test.txt")
	err := ioutil.WriteFile(sourceFile, []byte("test"), 0644)
	assert.NoError(t, err)

	// Copy file and check if it exists in destination
	assert.NoError(t, copyFile(sourceFile, "test.txt"))
	defer os.Remove("test.txt")
	_, err = os.Stat("test.txt")
	assert.NoError(t, err)
}

func Test_copyFile_Bad_Source(t *testing.T) {
	assert.Error(t, copyFile("/foo/invalid.txt", "/foo/invalid.txt"))
}

func Test_copyFile_Bad_Destination(t *testing.T) {
	// Write a test file to a test source directory
	sourceDir, _ := ioutil.TempDir("", "source")
	defer common.RemoveDir(sourceDir)
	sourceFile := filepath.Join(sourceDir, "test.txt")
	err := ioutil.WriteFile(sourceFile, []byte("test"), 0644)
	assert.NoError(t, err)
	_, err = os.Stat(sourceFile)
	assert.NoError(t, err)

	// Assert an error when copying to a bad destination
	assert.Error(t, copyFile(sourceFile, "/foo/invalid.txt"))
}

func Test_copyFile_Create_Dir(t *testing.T) {
	// Write a test file to a test source directory
	sourceDir, _ := ioutil.TempDir("", "source")
	defer common.RemoveDir(sourceDir)
	sourceFile := filepath.Join(sourceDir, "test.txt")
	err := ioutil.WriteFile(sourceFile, []byte("test"), 0644)
	assert.NoError(t, err)
	_, err = os.Stat(sourceFile)
	assert.NoError(t, err)

	// Create a test destination directory
	destDir, _ := ioutil.TempDir("", "dest")
	defer common.RemoveDir(sourceDir)

	// Use a nested directory that does not exist
	destFile := filepath.Join(destDir, "somedirectory", "test.txt")
	// Copy file and check if it exists in destination
	assert.NoError(t, copyFile(sourceFile, destFile))
	_, err = os.Stat(destFile)
	assert.NoError(t, err)
}

func Test_copySourceFiles_Success(t *testing.T) {
	// Write a test file to a test source directory
	sourceDir, _ := ioutil.TempDir("", "source")
	defer common.RemoveDir(sourceDir)
	sourceFile := filepath.Join(sourceDir, "test.txt")
	err := ioutil.WriteFile(sourceFile, []byte("test"), 0644)
	assert.NoError(t, err)
	_, err = os.Stat(sourceFile)
	assert.NoError(t, err)

	// Create a test destination directory
	destDir, _ := ioutil.TempDir("", "dest")
	defer common.RemoveDir(sourceDir)

	assert.NoError(t, copySourceFiles([]string{"test.txt"}, sourceDir, destDir))
}

func Test_copySourceFiles_Error(t *testing.T) {
	assert.Error(t, copySourceFiles([]string{"test.txt"}, "/foo", "/foo"))
}
