package stamp

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// NewTestStore returns a new store pointing to ./testdata/generators.
// Note that since we're starting with a relative path, this function
// needs to be called outside any calls to InTempDir().
func NewTestStore() *Store {
	storePath, _ := filepath.Abs(filepath.Join("testdata", "generators"))
	return NewStore(storePath)
}

func TestStore_LoadAll(t *testing.T) {
	store := NewStore("/some/path/that/does/not/exist")
	items, err := store.LoadAll()
	assert.Nil(t, items)
	assert.Error(t, err)

	store = NewTestStore()
	items, err = store.LoadAll()

	assert.True(t, len(items) > 0)
	assert.NoError(t, err)
}
