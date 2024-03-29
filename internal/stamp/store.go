package stamp

import (
	"github.com/twelvelabs/stamp/internal/pkg"
)

type Store struct {
	*pkg.Store
}

func NewStore(root string) *Store {
	return &Store{
		Store: pkg.NewStore(root).WithMetaFile(metaFileName),
	}
}

func (s *Store) Load(name string) (*Generator, error) {
	pkg, err := s.Store.Load(name)
	if err != nil {
		return nil, err
	}
	return NewGenerator(s, pkg)
}

func (s *Store) LoadAll() ([]*Generator, error) {
	packages, err := s.Store.LoadAll()
	if err != nil {
		return nil, err
	}
	return NewGenerators(s, packages)
}
