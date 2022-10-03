package pkg

import (
	"context"
	"fmt"
	"os"
	"path"
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
	pkgPath, err := s.path(name)
	if err != nil {
		return nil, err
	}
	return LoadPackage(pkgPath, s.MetaFile)
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
		os.RemoveAll(stagingRoot)
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
