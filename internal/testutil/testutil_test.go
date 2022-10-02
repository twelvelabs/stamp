package testutil

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanup(t *testing.T) {
	f := removeAllFunc
	defer func() {
		removeAllFunc = f
	}()

	removed := []string{}
	removeAllFunc = func(path string) error {
		removed = append(removed, path)
		return nil
	}

	AddCleanupPath("/path/aaa")
	AddCleanupPath("/path/bbb")
	assert.Equal(t, []string{"/path/aaa", "/path/bbb"}, CleanupPaths())

	Cleanup()

	assert.Equal(t, []string{}, CleanupPaths())
	assert.Equal(t, []string{"/path/aaa", "/path/bbb"}, removed)
}

func TestRemoveAll(t *testing.T) {
	f := removeAllFunc
	defer func() {
		removeAllFunc = f
	}()

	// should delegate to removeAllFunc
	removeAllFunc = func(path string) error {
		return nil
	}
	assert.NotPanics(t, func() {
		RemoveAll("/some/path")
	})

	// should panic if removeAllFunc returns an error
	removeAllFunc = func(path string) error {
		return errors.New("boom")
	}
	assert.Panics(t, func() {
		RemoveAll("/some/path")
	})
}

func TestMkdirAll(t *testing.T) {
	f := mkdirAllFunc
	defer func() {
		mkdirAllFunc = f
		ClearCleanupPaths()
	}()

	assert.Equal(t, []string{}, CleanupPaths())

	mkdirAllFunc = func(path string, perm os.FileMode) error {
		return nil
	}
	assert.NotPanics(t, func() {
		MkdirAll("/some/path", 0755)
	})
	assert.Equal(t, []string{"/some/path"}, CleanupPaths())
	ClearCleanupPaths()

	mkdirAllFunc = func(path string, perm os.FileMode) error {
		return errors.New("boom")
	}
	assert.Panics(t, func() {
		MkdirAll("/some/path", 0755)
	})
	assert.Equal(t, []string{}, CleanupPaths())
}

func TestMkdirTemp(t *testing.T) {
	f := mkdirTempFunc
	defer func() {
		mkdirTempFunc = f
		ClearCleanupPaths()
	}()

	assert.Equal(t, []string{}, CleanupPaths())

	// should delegate to mkdirTempFunc
	mkdirTempFunc = func(dir, pattern string) (string, error) {
		return "/some/path", nil
	}
	assert.Equal(t, "/some/path", MkdirTemp())
	assert.Equal(t, []string{"/some/path"}, CleanupPaths())
	ClearCleanupPaths()

	// should panic if mkdirTempFunc returns an error
	mkdirTempFunc = func(dir, pattern string) (string, error) {
		return "", errors.New("boom")
	}
	assert.Panics(t, func() {
		MkdirTemp()
	})
	assert.Equal(t, []string{}, CleanupPaths())
}

func TestInTempDir(t *testing.T) {
	f := chdirFunc
	defer func() {
		chdirFunc = f
		ClearCleanupPaths()
	}()

	actualDirs := []string{}
	chdirFunc = func(dir string) error {
		actualDirs = append(actualDirs, dir)
		return nil
	}

	expectedDirs := []string{}
	InTempDir(func(tmpDir string) {
		expectedDirs = append(expectedDirs, tmpDir)
	})
	expectedDirs = append(expectedDirs, WorkingDir())

	assert.Equal(t, expectedDirs, actualDirs)

	// error on first chdir call
	chdirFunc = func(dir string) error {
		return errors.New("boom")
	}
	assert.Panics(t, func() {
		InTempDir(func(tmpDir string) {})
	})

	// error on second chdir call... for the coverage :/
	chdirFunc = func(dir string) error {
		return nil
	}
	assert.Panics(t, func() {
		InTempDir(func(tmpDir string) {
			chdirFunc = func(dir string) error {
				return errors.New("boom")
			}
		})
	})
}

func TestWorkingDir(t *testing.T) {
	f := workingDirFunc
	defer func() {
		workingDirFunc = f
	}()

	// should delegate to workingDirFunc
	workingDirFunc = func() (dir string, err error) {
		return "/some/path", nil
	}
	assert.Equal(t, "/some/path", WorkingDir())

	// should panic if workingDirFunc returns an error
	workingDirFunc = func() (dir string, err error) {
		return "", errors.New("boom")
	}
	assert.Panics(t, func() {
		WorkingDir()
	})
}

func TestWriteFile(t *testing.T) {
	f := writeFileFunc
	defer func() {
		writeFileFunc = f
		ClearCleanupPaths()
	}()

	assert.Equal(t, []string{}, CleanupPaths())

	// should delegate to writeFileFunc
	writeFileFunc = func(name string, data []byte, perm os.FileMode) error {
		return nil
	}
	assert.NotPanics(t, func() {
		WriteFile("/some/path", []byte(""), 0666)
	})
	assert.Equal(t, []string{"/some/path"}, CleanupPaths())
	ClearCleanupPaths()

	// should panic if writeFileFunc returns an error
	writeFileFunc = func(name string, data []byte, perm os.FileMode) error {
		return errors.New("boom")
	}
	assert.Panics(t, func() {
		WriteFile("/some/path", []byte(""), 0666)
	})
	assert.Equal(t, []string{}, CleanupPaths())
}

func TestCreateFiles(t *testing.T) {
	defer Cleanup()
	tmpDir := MkdirTemp()

	// Handles nil
	CreatePaths(tmpDir, nil)
	assert.Equal(t, []string{tmpDir}, CleanupPaths())

	CreatePaths(tmpDir, map[string]any{
		"aaa/":       true,
		"bbb/ccc/":   true,
		"bin/aaa.sh": "aaa",
		"bin/bbb.sh": "bbb",
		"hello.txt":  "hello",
	})

	var data []byte

	assert.DirExists(t, filepath.Join(tmpDir, "aaa"))
	assert.DirExists(t, filepath.Join(tmpDir, "bbb", "ccc"))

	assert.FileExists(t, filepath.Join(tmpDir, "bin", "aaa.sh"))
	data, _ = os.ReadFile(filepath.Join(tmpDir, "bin", "aaa.sh"))
	assert.Equal(t, "aaa", string(data))

	assert.FileExists(t, filepath.Join(tmpDir, "bin", "bbb.sh"))
	data, _ = os.ReadFile(filepath.Join(tmpDir, "bin", "bbb.sh"))
	assert.Equal(t, "bbb", string(data))

	assert.FileExists(t, filepath.Join(tmpDir, "hello.txt"))
	data, _ = os.ReadFile(filepath.Join(tmpDir, "hello.txt"))
	assert.Equal(t, "hello", string(data))

	assert.Equal(t, []string{
		tmpDir,
		filepath.Join(tmpDir, "aaa"),
		filepath.Join(tmpDir, "bbb", "ccc"),
		filepath.Join(tmpDir, "bin"),
		filepath.Join(tmpDir, "bin", "aaa.sh"),
		filepath.Join(tmpDir, "bin", "bbb.sh"),
		filepath.Join(tmpDir, "hello.txt"),
	}, CleanupPaths())
}

func TestAssertFiles(t *testing.T) {
	defer Cleanup()
	tmpDir := MkdirTemp()

	// Handles nil
	AssertPaths(t, tmpDir, nil)

	CreatePaths(tmpDir, map[string]any{
		"aaa/":       true,
		"bbb/ccc/":   true,
		"bin/aaa.sh": "aaa",
		"bin/bbb.sh": "bbb",
		"hello.txt":  "hello",
	})

	os.Chmod(filepath.Join(tmpDir, "bin", "aaa.sh"), 0600)
	os.Chmod(filepath.Join(tmpDir, "bin", "bbb.sh"), 0600)

	AssertPaths(t, tmpDir, map[string]any{
		"aaa/":        true,    // dir exists
		"bbb/ccc/":    true,    // dir exists
		"bin/aaa.sh":  0600,    // file exists, perms match
		"bin/bbb.sh":  0600,    // file exists, perms match
		"hello.txt":   "hello", // file exists, content matches
		"unknown/":    false,   // dir should not exist
		"unknown.txt": false,   // file should not exist
	})
}
