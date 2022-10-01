package testutil

import (
	"os"
)

var (
	paths = []string{}

	// Abstracted to make testing easier
	removeAllFunc  = os.RemoveAll
	mkdirTempFunc  = os.MkdirTemp
	workingDirFunc = os.Getwd
)

// RegisterForCleanup registers a path to be removed on Cleanup()
func RegisterForCleanup(path string) {
	paths = append(paths, path)
}

// Cleanup removes all paths created by [testutil].
// Should be called via defer at the beginning of any test cases
// using [testutil] functions.
func Cleanup() {
	for _, p := range paths {
		RemoveAll(p)
	}
	paths = []string{}
}

// RemoveAll calls os.RemoveAll with path and panics on error.
func RemoveAll(path string) {
	err := removeAllFunc(path)
	if err != nil {
		panic(err)
	}
}

// MkdirTemp calls os.MkdirTemp and panics on error.
// The temp dir is registered to be removed on Cleanup().
func MkdirTemp() string {
	dir, err := mkdirTempFunc("", "testutil")
	if err != nil {
		panic(err)
	}
	RegisterForCleanup(dir)
	return dir
}

// WorkingDir calls os.Getwd and panics on error.
func WorkingDir() string {
	dir, err := workingDirFunc()
	if err != nil {
		panic(err)
	}
	return dir
}
