package stamp

import (
	"embed"
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

type Store struct {
	*pkg.Store
}

func NewStore(root string) *Store {
	return &Store{
		Store: pkg.NewStore(root).WithMetaFile(metaFileName),
	}
}

// AsGenerator returns pkg wrapped in a Generator type or err.
// Useful when calling [pkg.Store] methods that normally return a [pkg.Package].
func (s *Store) AsGenerator(pkg *pkg.Package, err error) (*Generator, error) {
	if err != nil {
		return nil, err
	}
	return NewGenerator(s, pkg)
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
	defaultGenPath := filepath.Join(s.Store.BasePath, "generator")
	if fsutil.PathExists(defaultGenPath) {
		return nil
	}
	return os.CopyFS(s.Store.BasePath, defaultGen)
}

func (s *Store) Load(name string) (*Generator, error) {
	return s.AsGenerator(s.Store.Load(name))
}

func (s *Store) LoadAll() ([]*Generator, error) {
	return s.AsGenerators(s.Store.LoadAll())
}
