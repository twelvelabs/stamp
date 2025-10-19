package pkg

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/hashicorp/go-getter"
)

const (
	DefaultMetaFile string = "package.yaml"
)

type Store struct {
	BasePath string
	MetaFile string
	getter   Getter
}

func NewStore(path string) *Store {
	return &Store{
		BasePath: path,
		MetaFile: DefaultMetaFile,
		getter:   DefaultGetter,
	}
}

// WithGetter returns the receiver with getter set to g.
func (s *Store) WithGetter(g Getter) *Store {
	s.getter = g
	return s
}

// WithMetaFile returns the receiver with MetaFile set to filename.
func (s *Store) WithMetaFile(filename string) *Store {
	s.MetaFile = filename
	return s
}

// Returns the named package from the store.
func (s *Store) Load(name string) (*Package, error) {
	var pkg *Package
	var err error

	// The name may be a direct path to a package on the filesystem.
	pkg, err = LoadPackage(name, s.MetaFile)
	nfErr := NewNotFoundError(s.MetaFile)
	if err != nil && !errors.Is(err, nfErr) {
		return nil, err
	}

	// If that doesn't return a result, then it must be
	// a named package in the store.
	// Convert the name to a path and load.
	if pkg == nil {
		var pkgPath string
		pkgPath, err = s.path(name)
		if err != nil {
			return nil, err
		}
		pkg, err = LoadPackage(pkgPath, s.MetaFile)
	}

	return pkg, err
}

// Returns all valid packages in the store.
// Silently ignores any packages that fail to load.
func (s *Store) LoadAll() ([]*Package, error) {
	return LoadPackages(s.BasePath, s.MetaFile)
}

type CleanupFunc func()

func (s *Store) Stage(src string) (*Package, CleanupFunc, error) {
	// Create a temp staging dir.
	stagingRoot, err := os.MkdirTemp("", "pkg-")
	if err != nil {
		return nil, func() {}, fmt.Errorf("staging error: %w", err)
	}
	cleanup := func() {
		_ = os.RemoveAll(stagingRoot)
	}

	// Normalize `src` into a fully qualified url (i.e. "." to "file:///${PWD}").
	// Needed so the update logic can work properly.
	pwd, _ := os.Getwd()
	src, err = getter.Detect(src, pwd, getter.Detectors)
	if err != nil {
		return nil, cleanup, fmt.Errorf("staging error: %w", err)
	}

	// Copy `src` into the staging dir.
	pkgPath := path.Join(stagingRoot, "staged")
	err = s.getter(context.Background(), src, pkgPath)
	if err != nil {
		return nil, cleanup, fmt.Errorf("staging error: %w", err)
	}

	// See if it's a loadable package.
	pkg, err := LoadPackage(pkgPath, s.MetaFile)
	if err != nil {
		return nil, cleanup, fmt.Errorf("staging error: %w", err)
	}

	// Store the source url (used to update the package later).
	pkg.SetOrigin(src)
	err = StorePackage(pkg)
	if err != nil {
		return nil, cleanup, fmt.Errorf("staging error: %w", err)
	}

	return pkg, cleanup, nil
}

// Install copies the package from `src` to the base path of the store.
func (s *Store) Install(src string) (*Package, error) {
	pkg, cleanup, err := s.Stage(src)
	if err != nil {
		return nil, err
	}
	defer cleanup()

	// Make sure it's not already installed
	if _, err = s.Load(pkg.Name()); err == nil {
		return nil, ErrPkgExists
	}

	// Move it to the store dir.
	pkgPath, err := s.path(pkg.Name())
	if err != nil {
		return nil, fmt.Errorf("install error: %w", err)
	}
	err = MovePackage(pkg, pkgPath)
	if err != nil {
		return nil, fmt.Errorf("install error: %w", err)
	}

	return pkg, nil
}

// Uninstall removes the package from the store.
func (s *Store) Uninstall(name string) (*Package, error) {
	pkg, err := s.Load(name)
	if err != nil {
		return nil, err
	}

	err = RemovePackage(pkg)
	if err != nil {
		return nil, err
	}

	return pkg, nil
}

// Update re-installs the named package from the source.
func (s *Store) Update(name string) (*Package, error) {
	pkg, err := s.Load(name)
	if err != nil {
		return nil, err
	}

	origin := pkg.Origin()
	if origin == "" {
		return nil, fmt.Errorf("origin missing: %s", pkg.Name())
	}

	err = RemovePackage(pkg)
	if err != nil {
		return nil, err
	}

	return s.Install(origin)
}

func (s *Store) path(name string) (string, error) {
	return PackagePath(s.BasePath, name)
}
