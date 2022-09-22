package test_util

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	paths = []string{}
)

// Cleanup removes all paths created by [test_util].
// Should be called via defer at the beginning of any test cases
// using [test_util] functions.
func Cleanup() {
	errs := []error{}
	for _, p := range paths {
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
	paths = append(paths, dir)
	return dir
}
