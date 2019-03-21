package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/sirupsen/logrus"
	assert "github.com/stretchr/testify/require"
)

var dummyData = []byte{0, 1, 2}

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	Log.Level = logrus.DebugLevel
}

type walkResult struct {
	output []string
	root   string
}

func (w *walkResult) accumulate(path string, info os.FileInfo, err error) error {
	if err == nil {
		relPath, relErr := filepath.Rel(w.root, path)
		if relErr != nil {
			return relErr
		}
		w.output = append(w.output, relPath)
	}
	return err
}

func (w *walkResult) sorted() []string {
	output := w.output
	sort.Strings(output)
	return output
}

func newWalkResult(root string) *walkResult {
	return &walkResult{root: root}
}

func mktemp(t *testing.T) string {
	tempDir, err := ioutil.TempDir("", "walk_test")
	assert.NoError(t, err)
	assert.NotEmpty(t, tempDir)
	return tempDir
}

func TestWalkFilterBlacklist(t *testing.T) {
	tempDir := mktemp(t)
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	blackListedFilename := filepath.Join(tempDir, "ios")
	assert.NoError(t, ioutil.WriteFile(blackListedFilename, dummyData, 0644))

	blackListedDir := filepath.Join(tempDir, ".git")
	assert.NoError(t, os.Mkdir(blackListedDir, 0755))
	blackListedSubFilename := filepath.Join(blackListedDir, "config")
	assert.NoError(t, ioutil.WriteFile(blackListedSubFilename, dummyData, 0644))

	whiteListedFilename := filepath.Join(tempDir, "whiteListedFile.dat")
	assert.NoError(t, ioutil.WriteFile(whiteListedFilename, dummyData, 0644))

	whiteListedDirectory := filepath.Join(tempDir, "whiteListedDir")
	assert.NoError(t, os.Mkdir(whiteListedDirectory, 0755))
	whiteListedSubFilename := filepath.Join(whiteListedDirectory, "whiteListedSubFilename")
	assert.NoError(t, ioutil.WriteFile(whiteListedSubFilename, dummyData, 0644))

	result := newWalkResult(tempDir)
	assert.NoError(t, filepath.Walk(tempDir, WalkFilterBlacklist(tempDir, result.accumulate)))
	output := result.sorted()
	assert.Len(t, output, 4)
	assert.Equal(t, ".", output[0])
	assert.Equal(t, "whiteListedDir", output[1])
	assert.Equal(t, "whiteListedDir/whiteListedSubFilename", output[2])
	assert.Equal(t, "whiteListedFile.dat", output[3])
}

func createSimpleFilesystem(tempDir string, t *testing.T) {
	file := filepath.Join(tempDir, "file.txt")
	assert.NoError(t, ioutil.WriteFile(file, dummyData, 0644))

	dir := filepath.Join(tempDir, "dir.dirext")
	assert.NoError(t, os.Mkdir(dir, 0755))
	subFile := filepath.Join(dir, "subfile.png")
	assert.NoError(t, ioutil.WriteFile(subFile, dummyData, 0644))
}

func TestWalkOnlyDirectory(t *testing.T) {
	tempDir := mktemp(t)
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()
	createSimpleFilesystem(tempDir, t)

	result := newWalkResult(tempDir)
	assert.NoError(t, filepath.Walk(tempDir, WalkOnlyDirectory(result.accumulate)))
	output := result.sorted()
	assert.Len(t, output, 2)
	assert.Equal(t, ".", output[0])
	assert.Equal(t, "dir.dirext", output[1])
}

func TestWalkOnlyFile(t *testing.T) {
	tempDir := mktemp(t)
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()
	createSimpleFilesystem(tempDir, t)

	result := newWalkResult(tempDir)
	assert.NoError(t, filepath.Walk(tempDir, WalkOnlyFile(result.accumulate)))
	output := result.sorted()
	assert.Len(t, output, 2)
	assert.Equal(t, "dir.dirext/subfile.png", output[0])
	assert.Equal(t, "file.txt", output[1])
}

func TestWalkFilterPrefix(t *testing.T) {
	tempDir := mktemp(t)
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()
	createSimpleFilesystem(tempDir, t)

	result := newWalkResult(tempDir)
	assert.NoError(t, filepath.Walk(tempDir, WalkFilterPrefix(result.accumulate, "dir")))
	output := result.sorted()
	assert.Len(t, output, 2)
	assert.Equal(t, ".", output[0])
	assert.Equal(t, "file.txt", output[1])
}

func TestWalkOnlyExtension(t *testing.T) {
	tempDir := mktemp(t)
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()
	createSimpleFilesystem(tempDir, t)

	result := newWalkResult(tempDir)
	assert.NoError(t, filepath.Walk(tempDir, WalkOnlyExtension(result.accumulate, "txt")))
	output := result.sorted()
	assert.Len(t, output, 3)
	assert.Equal(t, ".", output[0])
	assert.Equal(t, "dir.dirext", output[1])
	assert.Equal(t, "file.txt", output[2])
}
