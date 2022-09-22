package gen

import "github.com/twelvelabs/stamp/internal/pkg"

type Store struct {
	*pkg.Store
}

func NewStore(root string) *Store {
	return &Store{
		Store: pkg.NewStore(root).WithMetaFile("generator.yaml"),
	}
}

func (s *Store) Load(name string) (*Generator, error) {
	pkg, err := s.Store.Load(name)
	if err != nil {
		return nil, err
	}
	return NewGenerator(pkg)
}

func (s *Store) LoadAll() ([]*Generator, error) {
	packages, err := s.Store.LoadAll()
	if err != nil {
		return nil, err
	}
	return NewGenerators(packages)
}
