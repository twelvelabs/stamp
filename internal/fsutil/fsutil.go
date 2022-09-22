package fsutil

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	// DefaultDirMode grants `rwx------`.
	DefaultDirMode = 0700
	// DefaultFileMode grants `rw-------`.
	DefaultFileMode = 0600
)

type FsUtil struct {
}

func NewFsUtil() *FsUtil {
	return &FsUtil{}
}

func (fsu *FsUtil) NormalizePath(name string) (string, error) {
	return NormalizePath(name)
}

// NormalizePath ensures that name is an absolute path.
// Environment variables (and the ~ string) are expanded.
func NormalizePath(name string) (string, error) {
	normalized := strings.TrimSpace(name)

	// Replace ENV vars
	normalized = os.ExpandEnv(normalized)

	// Replace ~
	if strings.HasPrefix(normalized, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("unable to normalize %s: %w", name, err)
		}
		normalized = home + strings.TrimPrefix(normalized, "~")
	}

	// Ensure abs path
	normalized, err := filepath.Abs(normalized)
	if err != nil {
		return "", fmt.Errorf("unable to normalize %s: %w", name, err)
	}

	return normalized, nil
}

func (fsu *FsUtil) EnsureDirWritable(name string) error {
	// Ensures dir exists (and IsDir).
	err := os.MkdirAll(name, DefaultDirMode)
	if err != nil {
		return fmt.Errorf("unable to create %s: %w", name, err)
	}

	f := path.Join(name, ".touch")
	if err := os.WriteFile(f, []byte(""), DefaultFileMode); err != nil {
		return fmt.Errorf("unable to write to %s: %w", name, err)
	}
	defer os.Remove(f) //nolint:errcheck

	return nil
}
