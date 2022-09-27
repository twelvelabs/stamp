package test_util

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMkdirTemp(t *testing.T) {
	tmpDir := MkdirTemp(t)
	defer os.RemoveAll(tmpDir)

	assert.DirExists(t, tmpDir)
	assert.Equal(t, []string{tmpDir}, Paths)
}

func TestCleanup(t *testing.T) {
	tmpDir := MkdirTemp(t)
	assert.DirExists(t, tmpDir)

	Cleanup()

	assert.NoDirExists(t, tmpDir)
}

func TestWorkingDir(t *testing.T) {
	dir, _ := os.Getwd()
	assert.Equal(t, dir, WorkingDir())
}

func TestFixturesDir(t *testing.T) {
	dataDir, _ := filepath.Abs(filepath.Join(WorkingDir(), "..", "..", "testdata"))

	assert.Equal(t, (dataDir + "/foo"), FixturesDir("foo"))
	assert.Equal(t, (dataDir + "/foo/bar/baz"), FixturesDir("foo", "bar", "baz"))
}
