package pkg

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	// cspell:disable-line
	yaml "gopkg.in/yaml.v3"

	"github.com/twelvelabs/stamp/internal/fsutil"
)

var (
	ErrNotFound       = errors.New("package not found")
	ErrPkgExists      = errors.New("package already installed")
	ErrPkgNameInvalid = errors.New("invalid package name")
	ErrUnknown        = errors.New("unexpected error")

	pkgNameRegexp = regexp.MustCompile(`^[\w\-.:]+$`)
)

// PackagePath returns the absolute path to a package.
// Returns an error if `name` is empty or invalid.
// For example:
//
//	// "/packages/some/nested/name"
//	PackagePath("/packages", "some:nested:name")
func PackagePath(root string, name string) (string, error) {
	name = strings.TrimSpace(name)
	if !pkgNameRegexp.MatchString(name) {
		return "", ErrPkgNameInvalid
	}

	segments := strings.Split(name, ":")
	rel := filepath.Join(segments...)

	abs, err := filepath.Abs(filepath.Join(root, rel))
	if err != nil {
		return "", fmt.Errorf("package path: %w", err)
	}

	return abs, nil
}

// IsPackagePath returns true if pkgPath contains a metadata file.
func IsPackagePath(pkgPath string, metaFile string) bool {
	return fsutil.PathExists(filepath.Join(pkgPath, metaFile))
}

// LoadPackage parses and returns the package at `pkgPath`.
func LoadPackage(pkgPath string, metaFile string) (*Package, error) {
	if !IsPackagePath(pkgPath, metaFile) {
		return nil, ErrNotFound
	}

	// Read the package metadata file.
	pkgMetaPath := filepath.Join(pkgPath, metaFile)
	pkgMeta, err := os.ReadFile(pkgMetaPath) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("load package: %w", err)
	}

	pkg := &Package{
		path:     pkgPath,
		metaFile: metaFile,
	}

	// Parse the package metadata.
	err = yaml.Unmarshal(pkgMeta, &pkg.Metadata)
	if err != nil {
		return nil, fmt.Errorf("load package: %w", err)
	}

	return pkg, nil
}

// LoadPackages returns all valid packages found in `root`.
// Includes nested sub-packages.
func LoadPackages(root string, metaFile string) ([]*Package, error) {
	found := []*Package{}

	// Ensure the root path is readable.
	if _, err := os.Stat(root); err != nil {
		return nil, fmt.Errorf("load packages: %w", err)
	}

	// Note: We're intentionally ignoring FS or package parse errors
	//       here so that one bad package doesn't break everything.
	_ = filepath.Walk(root, func(name string, info os.FileInfo, err error) error {
		if err != nil {
			return nil //nolint:nilerr
		}

		basename := filepath.Base(name)
		if info.IsDir() && strings.HasPrefix(basename, "_") {
			return filepath.SkipDir // Ignore underscore (i.e. private) dirs
		}
		if basename != metaFile {
			return nil // Ignore non-package files
		}

		dir := filepath.Dir(name)
		if dir == root {
			return nil // Ignore the root path
		}

		pkg, _ := LoadPackage(dir, metaFile)
		if pkg != nil {
			found = append(found, pkg)
		}

		return nil
	})

	// Since we're only processing the metadata files, there's a chance
	// that the slice may be slightly out of order.
	// For example:
	//  - /a
	//  - /a/b
	//  - /a/b/package.yaml
	//  - /a/package.yaml
	// Results in:
	//   - []string{"a:b", "a"}
	// Resorting to ensure the packages are sorted lexically by name.
	sort.Slice(found, func(i, j int) bool {
		return found[i].Name() < found[j].Name()
	})

	return found, nil
}

// MovePackage moves the package to the given path.
func MovePackage(pkg *Package, newPath string) error {
	if _, err := os.Stat(newPath); !os.IsNotExist(err) {
		return ErrPkgExists
	}

	err := os.MkdirAll(filepath.Dir(newPath), 0755) //nolint:gosec
	if err != nil {
		return err
	}

	err = os.Rename(pkg.Path(), newPath)
	if err != nil {
		return err
	}

	pkg.SetPath(newPath)
	return nil
}

// StorePackage updates the package metadata file.
func StorePackage(pkg *Package) error {
	data, err := yaml.Marshal(&pkg.Metadata)
	if err != nil {
		return err
	}

	name := pkg.MetaPath()
	err = os.WriteFile(name, data, 0600)
	if err != nil {
		return err
	}

	return nil
}

// RemovePackage deletes the package from the filesystem.
func RemovePackage(pkg *Package) error {
	if _, err := os.Stat(pkg.Path()); !os.IsNotExist(err) {
		return os.RemoveAll(pkg.Path())
	}
	return ErrNotFound
}
