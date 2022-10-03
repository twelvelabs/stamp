package fsutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/twelvelabs/stamp/internal/testutil"
)

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
	defer testutil.Cleanup()
	testutil.InTempDir(t, func(tmpDir string) {
		dir := filepath.Join(tmpDir, "foo")
		err := EnsureDirWritable(dir)
		assert.NoError(t, err)
		assert.DirExists(t, dir, "dir should exist")

		dirEntry := filepath.Join(dir, "bar")
		testutil.WriteFile(dirEntry, []byte(""), 0600)
		assert.FileExists(t, dirEntry, "dir should be writable")
	})
}
