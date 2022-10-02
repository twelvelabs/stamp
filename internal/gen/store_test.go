package gen

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStore_LoadAll(t *testing.T) {
	store := NewStore("/some/path/that/does/not/exist")
	items, err := store.LoadAll()
	assert.Nil(t, items)
	assert.Error(t, err)

	storeDir := filepath.Join("..", "..", "testdata", "generators")
	store = NewStore(storeDir)
	items, err = store.LoadAll()

	assert.Len(t, items, 1)
	assert.NoError(t, err)
}
