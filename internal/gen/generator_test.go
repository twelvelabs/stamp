package gen

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twelvelabs/stamp/internal/pkg"
)

func TestNewGenerator(t *testing.T) {
	_, err := NewGenerator(nil)
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
	_, err = NewGenerator(p)
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
	_, err = NewGenerator(p)
	assert.ErrorContains(t, err, "generator metadata invalid")
}

func TestNewGenerators(t *testing.T) {
	items, err := NewGenerators([]*pkg.Package{})
	assert.Equal(t, []*Generator{}, items)
	assert.NoError(t, err)

	items, err = NewGenerators([]*pkg.Package{nil})
	assert.Nil(t, items)
	assert.ErrorContains(t, err, "nil package")

	p1 := &pkg.Package{}
	p2 := &pkg.Package{}
	items, err = NewGenerators([]*pkg.Package{p1, p2})
	assert.Len(t, items, 2)
	assert.NoError(t, err)
}
