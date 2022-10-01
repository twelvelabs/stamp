package testutil

import (
	"errors"
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

	RegisterForCleanup("/path/aaa")
	RegisterForCleanup("/path/bbb")
	assert.Equal(t, []string{"/path/aaa", "/path/bbb"}, paths)

	Cleanup()
	assert.Equal(t, []string{}, paths)
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

func TestMkdirTemp(t *testing.T) {
	f := mkdirTempFunc
	defer func() {
		mkdirTempFunc = f
	}()

	assert.Equal(t, []string{}, paths)

	// should delegate to mkdirTempFunc
	mkdirTempFunc = func(dir, pattern string) (string, error) {
		return "/some/path", nil
	}
	assert.Equal(t, "/some/path", MkdirTemp())
	assert.Equal(t, []string{"/some/path"}, paths, "temp dir should have been added to paths")

	// should panic if mkdirTempFunc returns an error
	mkdirTempFunc = func(dir, pattern string) (string, error) {
		return "", errors.New("boom")
	}
	assert.Panics(t, func() {
		_ = MkdirTemp()
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
		_ = WorkingDir()
	})
}
