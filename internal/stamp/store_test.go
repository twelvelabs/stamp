package stamp

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twelvelabs/termite/testutil"

	"github.com/twelvelabs/stamp/internal/pkg"
)

// NewTestStore returns a new store pointing to ./testdata/generators.
// Note that since we're starting with a relative path, this function
// needs to be called outside any calls to InTempDir().
func NewTestStore() *Store {
	storePath, _ := filepath.Abs(filepath.Join("testdata", "generators"))
	return NewStore(storePath)
}

func TestStore_Init(t *testing.T) {
	testutil.InTempDir(t, func(tmpDir string) {
		store := NewStore(tmpDir)

		genDir := filepath.Join(tmpDir, "generator")
		assert.NoDirExists(t, genDir)

		err := store.Init()
		assert.NoError(t, err)
		assert.DirExists(t, genDir)
		assert.DirExists(t, filepath.Join(genDir, "_src"))
		assert.FileExists(t, filepath.Join(genDir, "generator.yaml"))
	})
}

func TestStore_Init_WhenAlreadyExists(t *testing.T) {
	testutil.InTempDir(t, func(tmpDir string) {
		store := NewStore(tmpDir)

		genDir := filepath.Join(tmpDir, "generator")
		_ = os.Mkdir(genDir, 0777) //nolint:gosec
		_ = os.WriteFile(filepath.Join(genDir, "hello.txt"), []byte("hi"), 0600)

		err := store.Init()
		assert.NoError(t, err)
		assert.DirExists(t, genDir)
		assert.FileExists(t, filepath.Join(genDir, "hello.txt"))
		assert.NoDirExists(t, filepath.Join(genDir, "_src"))
		assert.NoFileExists(t, filepath.Join(genDir, "generator.yaml"))
	})
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

func TestStore_Stage(t *testing.T) {
	getter := pkg.NewMockGetter(func(ctx context.Context, src, dst string) error {
		return errors.New("boom")
	})

	store := NewStore("/some/path/that/does/not/exist")
	store.WithGetter(getter.Get)

	pkg, cleanup, err := store.Stage("https://github.com/example/repo")
	defer cleanup()

	assert.Nil(t, pkg)
	assert.ErrorContains(t, err, "boom")
}
