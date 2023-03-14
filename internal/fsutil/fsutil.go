package fsutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	// DefaultDirMode grants `rwx------`.
	DefaultDirMode = 0700
	// DefaultFileMode grants `rw-------`.
	DefaultFileMode = 0600
)

func NoPathExists(path string) bool {
	_, err := os.Stat(path)
	return errors.Is(err, os.ErrNotExist)
}

func PathExists(path string) bool {
	return !NoPathExists(path)
}

// NormalizePath ensures that name is an absolute path.
// Environment variables (and the ~ string) are expanded.
func NormalizePath(name string) (string, error) {
	normalized := strings.TrimSpace(name)
	if normalized == "" {
		return "", nil
	}

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

// EnsureDirWritable ensures that path is a writable directory.
// Will attempt to create a new directory if path does not exist.
func EnsureDirWritable(path string) error {
	// Ensure dir exists (and IsDir).
	err := os.MkdirAll(path, DefaultDirMode)
	if err != nil {
		return fmt.Errorf("unable to create %s: %w", path, err)
	}

	f := filepath.Join(path, ".touch")
	if err := os.WriteFile(f, []byte(""), DefaultFileMode); err != nil {
		return fmt.Errorf("unable to write to %s: %w", path, err)
	}
	defer os.Remove(f)

	return nil
}
