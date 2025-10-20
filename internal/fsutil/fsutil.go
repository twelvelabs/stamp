package fsutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

const (
	// DefaultDirMode grants `rwx------`.
	DefaultDirMode = 0700
	// DefaultFileMode grants `rw-------`.
	DefaultFileMode = 0600
)

func NoPathExists(path string) bool {
	_, err := os.Stat(path)
	// for some reason os.ErrInvalid sometimes != syscall.EINVAL :shrug:
	if errors.Is(err, os.ErrNotExist) ||
		errors.Is(err, os.ErrInvalid) ||
		errors.Is(err, syscall.EINVAL) {
		return true
	}
	return false
}

func PathExists(path string) bool {
	return !NoPathExists(path)
}

func PathIsDir(path string) bool {
	info, _ := os.Stat(path)
	return info != nil && info.IsDir()
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

// EnsurePathRelativeToRoot ensures that the relative path exists inside the trusted root,
// and returns it's absolute path. Returns an error if the path traverses outside the root.
func EnsurePathRelativeToRoot(path string, root string) (string, error) {
	path = filepath.FromSlash(path)

	absRoot, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}

	// EvalSymlinks returns an lstat error if path does not exist.
	if PathExists(absRoot) {
		absRoot, err = filepath.EvalSymlinks(absRoot)
		if err != nil {
			return "", err
		}
	}

	absPath := path
	if !filepath.IsAbs(absPath) {
		absPath, err = filepath.Abs(filepath.Join(absRoot, path))
		if err != nil {
			return "", err
		}
	}

	// EvalSymlinks returns an lstat error if path does not exist.
	if PathExists(absPath) {
		absPath, err = filepath.EvalSymlinks(absPath)
		if err != nil {
			return "", err
		}
	}

	if !IsSubDir(absPath, absRoot) {
		return "", fmt.Errorf("%s attempted to traverse outside of %s", path, root)
	}

	return absPath, nil
}

func IsSubDir(path string, dir string) bool {
	separator := string(filepath.Separator)
	for path != separator {
		path = filepath.Dir(path)
		if path == dir {
			return true
		}
	}
	return false
}
