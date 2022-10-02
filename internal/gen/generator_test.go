package gen

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twelvelabs/stamp/internal/iostreams"
	"github.com/twelvelabs/stamp/internal/pkg"
	"github.com/twelvelabs/stamp/internal/testutil"
)

func NewTestStore() *Store {
	storePath, _ := filepath.Abs(filepath.Join("..", "..", "testdata", "generators"))
	return NewStore(storePath)
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

func TestGenerator_AddsValuesFromDelegatedGenerators(t *testing.T) {
	defer testutil.Cleanup()

	store := NewTestStore() // store path is relative, can't be called in tmp dir

	testutil.InTempDir(func(tmpDir string) {
		gen, err := store.Load("delegating")
		assert.NotNil(t, gen)
		assert.NoError(t, err)

		values := gen.Values.GetAll()

		assert.Len(t, values, 2)
		assert.Equal(t, "customized.txt", values["FileName"])
		assert.Equal(t, "custom content", values["FileContent"])

		ctx := NewTaskContext(iostreams.Test(), nil, store, false)
		err = gen.Tasks.Execute(ctx, values)

		testutil.AssertPaths(t, tmpDir, map[string]any{
			"customized.txt": "custom content",
		})
	})

	testutil.InTempDir(func(tmpDir string) {
		gen, err := store.Load("delegating-dupe")
		assert.NotNil(t, gen)
		assert.NoError(t, err)

		values := gen.Values.GetAll()

		// should only be two values, even though the generator was referenced twice
		assert.Len(t, values, 2)
		// the defaults should be set by the last generator task in the list.
		assert.Equal(t, "untitled.txt", values["FileName"])
		assert.Equal(t, "", values["FileContent"])

		ctx := NewTaskContext(iostreams.Test(), nil, store, false)
		err = gen.Tasks.Execute(ctx, values)

		// Should have respected the `extra` attribute
		// and created two different filenames.
		testutil.AssertPaths(t, tmpDir, map[string]any{
			"customized.txt": "custom content",
			"untitled.txt":   "",
		})
	})
}
