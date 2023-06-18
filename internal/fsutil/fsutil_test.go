package fsutil

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twelvelabs/termite/testutil"
)

func TestNoPathExists(t *testing.T) {
	testutil.InTempDir(t, func(dir string) {
		assert.NoFileExists(t, "foo.txt")
		assert.Equal(t, true, NoPathExists("foo.txt"))

		testutil.WriteFile(t, "foo.txt", []byte(""), 0600)
		assert.Equal(t, false, NoPathExists("foo.txt"))
	})
}

func TestPathExists(t *testing.T) {
	testutil.InTempDir(t, func(dir string) {
		assert.NoFileExists(t, "foo.txt")
		assert.Equal(t, false, PathExists("foo.txt"))

		testutil.WriteFile(t, "foo.txt", []byte(""), 0600)
		assert.Equal(t, true, PathExists("foo.txt"))
	})
}

func TestPathIsDir(t *testing.T) {
	testutil.InTempDir(t, func(dir string) {
		assert.NoDirExists(t, "foo")
		assert.Equal(t, false, PathIsDir("foo"))

		testutil.WriteFile(t, "foo", []byte(""), 0600)
		assert.Equal(t, false, PathIsDir("foo"))
		testutil.RemoveAll(t, "foo")

		testutil.MkdirAll(t, "foo", 0777)
		assert.Equal(t, true, PathIsDir("foo"))
	})
}

func TestNormalizePath(t *testing.T) {
	homeDir, _ := os.UserHomeDir()
	workingDir, _ := filepath.Abs(".")

	tests := []struct {
		Desc    string
		EnvVars map[string]string
		Input   string
		Output  string
		Err     string
	}{
		{
			Desc:   "is a noop when passed an empty string",
			Input:  "",
			Output: "",
			Err:    "",
		},
		{
			Desc:   "expands env vars",
			Input:  filepath.Join(".", "${FOO}-dir", "$BAR"),
			Output: filepath.Join(workingDir, "aaa-dir", "bbb"),
			EnvVars: map[string]string{
				"FOO": "aaa",
				"BAR": "bbb",
			},
			Err: "",
		},
		{
			Desc:   "expands tilde",
			Input:  "~",
			Output: homeDir,
			Err:    "",
		},
		{
			Desc:   "expands tilde when prefix",
			Input:  filepath.Join("~", "foo"),
			Output: filepath.Join(homeDir, "foo"),
			Err:    "",
		},
		{
			Desc:   "returns an absolute path",
			Input:  ".",
			Output: workingDir,
			Err:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Desc, func(t *testing.T) {
			if tt.EnvVars != nil {
				for k, v := range tt.EnvVars {
					t.Setenv(k, v)
				}
			}

			actual, err := NormalizePath(tt.Input)

			assert.Equal(t, tt.Output, actual)
			if tt.Err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.Err)
			}
		})
	}
}

func TestEnsureDirWritable(t *testing.T) {
	testutil.InTempDir(t, func(tmpDir string) {
		dir := filepath.Join(tmpDir, "foo")
		err := EnsureDirWritable(dir)
		assert.NoError(t, err)
		assert.DirExists(t, dir, "dir should exist")

		dirEntry := filepath.Join(dir, "bar")
		testutil.WriteFile(t, dirEntry, []byte(""), 0600)
		assert.FileExists(t, dirEntry, "dir should be writable")
	})
}

func TestEnsurePathRelativeToRoot(t *testing.T) {
	rootDir, _ := filepath.Abs(filepath.Join("testdata", "aaa"))

	tests := []struct {
		desc     string
		path     string
		root     string
		expected string
		err      string
	}{
		{
			desc:     "returns path to dirs in root",
			path:     "bbb",
			root:     rootDir,
			expected: filepath.Join(rootDir, "bbb"),
		},
		{
			desc:     "returns path to files in root",
			path:     "bbb/ccc.txt",
			root:     rootDir,
			expected: filepath.Join(rootDir, "bbb", "ccc.txt"),
		},
		{
			desc:     "returns path to files in root even if not present",
			path:     "bbb/missing.txt",
			root:     rootDir,
			expected: filepath.Join(rootDir, "bbb", "missing.txt"),
		},
		{
			desc:     "returns path to symlinks in root",
			path:     "bbb/ddd.txt",
			root:     rootDir,
			expected: filepath.Join(rootDir, "bbb", "ccc.txt"),
		},
		{
			desc:     "handles relative roots without error",
			path:     "bbb",
			root:     "./testdata/aaa",
			expected: filepath.Join(rootDir, "bbb"),
		},
		{
			desc:     "handles absolute paths without mangling",
			path:     filepath.Join(rootDir, "bbb"),
			root:     rootDir,
			expected: filepath.Join(rootDir, "bbb"),
		},
		{
			desc: "returns error if path traverses outside of root",
			path: "../protected.txt",
			root: rootDir,
			err:  "attempted to traverse",
		},
		{
			desc: "returns error if path traverses outside of root regardless of presence",
			path: "../missing.txt",
			root: rootDir,
			err:  "attempted to traverse",
		},
		{
			desc: "returns error if symlink traverses outside of root",
			path: "bbb/eee.txt",
			root: rootDir,
			err:  "attempted to traverse",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			actual, err := EnsurePathRelativeToRoot(tt.path, tt.root)

			assert.Equal(t, tt.expected, actual)

			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.err)
			}
		})
	}
}

func TestIsSubDir(t *testing.T) {
	tests := []struct {
		path     string
		dir      string
		expected bool
	}{
		{
			path:     "/foo/bar/baz",
			dir:      "/foo",
			expected: true,
		},
		{
			path:     "/foo/bar/baz",
			dir:      "/",
			expected: true,
		},
		{
			path:     "/foo/bar/baz",
			dir:      "/bar",
			expected: false,
		},
	}
	for _, tt := range tests {
		desc := fmt.Sprintf(":%s:%s", tt.path, tt.dir)
		t.Run(desc, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsSubDir(tt.path, tt.dir))
		})
	}
}
