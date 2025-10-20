package stamp

import (
	"embed"
	"errors"
	"os"
	"path/filepath"

	"github.com/twelvelabs/stamp/internal/fsutil"
	"github.com/twelvelabs/stamp/internal/pkg"
)

// Note: the `all:` prefix is required so that the
// dirs starting with `_` are included.
//
//go:embed all:generator
var defaultGen embed.FS

type CleanupFunc = pkg.CleanupFunc

type Store struct {
	*pkg.Store
}

func NewStore(root string) *Store {
	return &Store{
		Store: pkg.NewStore(root).WithMetaFile(metaFileName),
	}
}

// AsGenerator returns p wrapped in a Generator type or err.
// Useful when calling [pkg.Store] methods that normally return a [pkg.Package].
func (s *Store) AsGenerator(p *pkg.Package, err error) (*Generator, error) {
	if err != nil {
		if errors.Is(err, pkg.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return NewGenerator(s, p)
}

// AsGenerators returns packages wrapped in Generator types or err.
// Useful when calling [pkg.Store] methods that normally return a slice of [pkg.Package] types.
func (s *Store) AsGenerators(packages []*pkg.Package, err error) ([]*Generator, error) {
	if err != nil {
		return nil, err
	}
	return NewGenerators(s, packages)
}

func (s *Store) Init() error {
	defaultGenPath := filepath.Join(s.BasePath, "generator")
	if fsutil.PathExists(defaultGenPath) {
		return nil
	}
	return os.CopyFS(s.BasePath, defaultGen)
}

// Returns the named generator from the store.
func (s *Store) Load(name string) (*Generator, error) {
	return s.AsGenerator(s.Store.Load(name))
}

// Stage copies a generator from src into a temp dir.
// Returns the generator and a cleanup function that
// removes the temp dir.
func (s *Store) Stage(src string) (*Generator, CleanupFunc, error) {
	p, cleanup, err := s.Store.Stage(src)
	g, err := s.AsGenerator(p, err)
	return g, cleanup, err
}

// Returns all valid generators in the store.
// Silently ignores any generators that fail to load.
func (s *Store) LoadAll() ([]*Generator, error) {
	return s.AsGenerators(s.Store.LoadAll())
}
