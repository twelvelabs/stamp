package pkg

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/edwardrf/symwalk" // cspell:disable-line
	yaml "gopkg.in/yaml.v3"
)

var (
	ErrPkgExists      = errors.New("package already installed")
	ErrPkgNameInvalid = errors.New("invalid package name")
	ErrUnknown        = errors.New("unexpected error")

	pkgNameRegexp = regexp.MustCompile(`^[\w\-.:]+$`)
)

func NewNotFoundError(metaFile string) NotFoundError {
	return NotFoundError{
		metaFile: metaFile,
	}
}

type NotFoundError struct {
	metaFile string
}

func (e NotFoundError) Error() string {
	// This assumes that the meta filename describes what the package is.
	// i.e. if your package is a "widget", then name the file widget.yaml
	kind := strings.TrimSuffix(e.metaFile, filepath.Ext(e.metaFile))
	return fmt.Sprintf("%s not found", kind)
}

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

// LoadPackage parses and returns the package at `pkgPath`.
func LoadPackage(pkgPath string, metaFile string) (*Package, error) {
	// Ensure package path exists.
	if _, err := os.Stat(pkgPath); os.IsNotExist(err) {
		return nil, NewNotFoundError(metaFile)
	}

	// Read the package metadata file.
	pkgMeta, err := os.ReadFile(filepath.Join(pkgPath, metaFile))
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
	_ = symwalk.Walk(root, func(name string, info os.FileInfo, err error) error {
		if err != nil {
			return nil //nolint:nilerr
		}

		filename := filepath.Base(name)
		if filename != metaFile {
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

	err := os.MkdirAll(filepath.Dir(newPath), 0755)
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
	err = os.WriteFile(name, data, 0755)
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
	return NewNotFoundError(pkg.MetaFile())
}
