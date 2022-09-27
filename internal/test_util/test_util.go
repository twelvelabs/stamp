package test_util

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	Paths = []string{}
)

// Cleanup removes all paths created by [test_util].
// Should be called via defer at the beginning of any test cases
// using [test_util] functions.
func Cleanup() {
	errs := []error{}
	for _, p := range Paths {
		err := os.RemoveAll(p)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
		panic(fmt.Sprintf("failed to cleanup %d paths", len(errs)))
	}
}

// MkdirTemp creates a new temp directory.
func MkdirTemp(t *testing.T) string {
	dir, err := os.MkdirTemp("", "test_util")
	if err != nil {
		assert.FailNow(t, "unable to create temp dir", err)
	}
	Paths = append(Paths, dir)
	return dir
}

// WorkingDir returns the current working dir.
func WorkingDir() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return dir
}

// FixturesDir returns the path to a dir in `testdata`.
func FixturesDir(args ...string) string {
	path := filepath.Join(args...)
	dir, err := filepath.Abs(filepath.Join("..", "..", "testdata", path))
	if err != nil {
		panic(err)
	}
	return dir
}
