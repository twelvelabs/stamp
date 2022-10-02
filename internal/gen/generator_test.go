package gen

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twelvelabs/stamp/internal/pkg"
)

func NewTestStore() *Store {
	return NewStore(filepath.Join("..", "..", "testdata", "generators"))
}

func TestNewGenerator(t *testing.T) {
	store := NewTestStore()

	_, err := NewGenerator(nil, nil)
	assert.ErrorContains(t, err, "nil store")

	_, err = NewGenerator(store, nil)
	assert.ErrorContains(t, err, "nil package")

	p := &pkg.Package{
		Metadata: map[string]any{
			"values": []any{
				map[string]any{
					"key": 123, // key should be a string
				},
			},
		},
	}
	_, err = NewGenerator(store, p)
	assert.ErrorContains(t, err, "generator metadata invalid")

	p = &pkg.Package{
		Metadata: map[string]any{
			"tasks": []any{
				map[string]any{
					"type": "unknown", // unknown type
				},
			},
		},
	}
	_, err = NewGenerator(store, p)
	assert.ErrorContains(t, err, "generator metadata invalid")
}

func TestNewGenerators(t *testing.T) {
	store := NewTestStore()

	items, err := NewGenerators(nil, []*pkg.Package{})
	assert.Nil(t, items)
	assert.ErrorContains(t, err, "nil store")

	items, err = NewGenerators(store, []*pkg.Package{nil})
	assert.Nil(t, items)
	assert.ErrorContains(t, err, "nil package")

	items, err = NewGenerators(store, []*pkg.Package{})
	assert.Equal(t, []*Generator{}, items)
	assert.NoError(t, err)

	p1 := &pkg.Package{}
	p2 := &pkg.Package{}
	items, err = NewGenerators(store, []*pkg.Package{p1, p2})
	assert.Len(t, items, 2)
	assert.NoError(t, err)
}
