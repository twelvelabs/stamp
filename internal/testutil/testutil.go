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

	// Abstracted to make testing easier.
	chdirFunc      = os.Chdir
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

// AddCleanupPath registers a path to be removed on Cleanup().
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

// InTempDir executes handler in a temp dir, restoring the working dir
// back to it's original location once handler exits.
func InTempDir(tb testing.TB, handler func(tmpDir string)) {
	tb.Helper()
	current := WorkingDir()
	tmp := MkdirTemp()
	if err := chdirFunc(tmp); err != nil {
		panic(err)
	}
	handler(tmp)
	if err := chdirFunc(current); err != nil {
		panic(err)
	}
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

func AssertPaths(tb testing.TB, base string, files map[string]any) {
	tb.Helper()
	if files == nil {
		return
	}
	for _, key := range sortedKeys(files) {
		value := files[key]
		path, isDir := joinPath(base, key)
		if isDir {
			AssertDirPath(tb, path, value)
		} else {
			AssertFilePath(tb, path, value)
		}
	}
}

// AssertDirPath asserts a directory path does not exist if
// value is false, otherwise asserts that it does.
func AssertDirPath(tb testing.TB, path string, value any) {
	tb.Helper()
	if exists, ok := value.(bool); ok && !exists {
		assert.NoDirExists(tb, path)
	} else {
		assert.DirExists(tb, path)
	}
}

// AssertFilePath asserts a file path does not exist if
// value is false, otherwise asserts that it does.
//
// Additionally:
//   - If value is a string, then the file contents should match.
//   - If value is an int, then file permissions should match.
func AssertFilePath(tb testing.TB, path string, value any) {
	tb.Helper()
	if exists, ok := value.(bool); ok && !exists {
		// value is `false`, file _should not_ be there
		assert.NoFileExists(tb, path)
	} else {
		// file _should_ be there
		assert.FileExists(tb, path)
		if content, ok := value.(string); ok {
			// value is a string, file content should match
			buf, _ := os.ReadFile(path)
			assert.Equal(tb, content, string(buf))
		}
		if perm, ok := value.(int); ok {
			// value is an int, file permissions should match
			info, _ := os.Stat(path)
			assert.Equal(tb, perm, int(info.Mode().Perm()))
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
func joinPath(base string, path string) (string, bool) {
	isDir := strings.HasSuffix(path, "/")
	joined := filepath.Join(base, path)
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
