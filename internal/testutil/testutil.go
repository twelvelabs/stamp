package testutil

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	paths = []string{}

	// Abstracted to make testing easier
	removeAllFunc  = os.RemoveAll
	mkdirTempFunc  = os.MkdirTemp
	mkdirAllFunc   = os.MkdirAll
	workingDirFunc = os.Getwd
	writeFileFunc  = os.WriteFile
)

func CleanupPaths() []string {
	return paths
}

func ClearCleanupPaths() {
	paths = []string{}
}

// AddCleanupPath registers a path to be removed on Cleanup()
func AddCleanupPath(path string) {
	paths = append(paths, path)
}

// Cleanup removes all paths created by [testutil].
// Should be called via defer at the beginning of any test cases
// using [testutil] functions.
func Cleanup() {
	for _, p := range CleanupPaths() {
		RemoveAll(p)
	}
	ClearCleanupPaths()
}

// RemoveAll calls os.RemoveAll with path and panics on error.
func RemoveAll(path string) {
	err := removeAllFunc(path)
	if err != nil {
		panic(err)
	}
}

func MkdirAll(path string, perm fs.FileMode) {
	err := mkdirAllFunc(path, perm)
	if err != nil {
		panic(err)
	}
	AddCleanupPath(path)
}

// MkdirTemp calls os.MkdirTemp and panics on error.
// The temp dir is registered to be removed on Cleanup().
func MkdirTemp() string {
	dir, err := mkdirTempFunc("", "testutil")
	if err != nil {
		panic(err)
	}
	AddCleanupPath(dir)
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

func WriteFile(name string, data []byte, perm fs.FileMode) {
	err := writeFileFunc(name, data, perm)
	if err != nil {
		panic(err)
	}
	AddCleanupPath(name)
}

func AssertPaths(t *testing.T, base string, files map[string]any) {
	if files == nil {
		return
	}
	for _, key := range sortedKeys(files) {
		value := files[key]
		path, isDir := joinPath(base, key)
		if isDir {
			AssertDirPath(t, path, value)
		} else {
			AssertFilePath(t, path, value)
		}
	}
}

// AssertDirPath asserts a directory path does not exist if
// value is false, otherwise asserts that it does.
func AssertDirPath(t *testing.T, path string, value any) {
	if exists, ok := value.(bool); ok && !exists {
		assert.NoDirExists(t, path)
	} else {
		assert.DirExists(t, path)
	}
}

// AssertFilePath asserts a file path does not exist if
// value is false, otherwise asserts that it does.
//
// Additionally:
//   - If value is a string, then the file contents should match.
//   - If value is an int, then file permissions should match.
func AssertFilePath(t *testing.T, path string, value any) {
	if exists, ok := value.(bool); ok && !exists {
		// value is `false`, file _should not_ be there
		assert.NoFileExists(t, path)
	} else {
		// file _should_ be there
		assert.FileExists(t, path)
		if content, ok := value.(string); ok {
			// value is a string, file content should match
			buf, _ := os.ReadFile(path)
			assert.Equal(t, content, string(buf))
		}
		if perm, ok := value.(int); ok {
			// value is an int, file permissions should match
			info, _ := os.Stat(path)
			assert.Equal(t, perm, int(info.Mode().Perm()))
		}
	}
}

func CreatePaths(base string, files map[string]any) {
	if files == nil {
		return
	}
	for _, key := range sortedKeys(files) {
		value := files[key]
		path, isDir := joinPath(base, key)
		if isDir {
			MkdirAll(path, 0755)
		} else {
			parent := filepath.Dir(path)
			if _, err := os.Stat(parent); os.IsNotExist(err) {
				MkdirAll(parent, 0755)
			}
			WriteFile(path, []byte(value.(string)), 0666)
		}
	}
}

// Helper that returns base+path and whether path is supposed
// to be a directory (indicated by it ending in "/").
func joinPath(base string, path string) (joined string, isDir bool) {
	isDir = strings.HasSuffix(path, "/")
	joined = filepath.Join(base, path)
	return joined, isDir
}

// Helper that returns a sorted set of map keys.
func sortedKeys(data map[string]any) []string {
	keys := []string{}
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
