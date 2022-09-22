package pkg

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twelvelabs/stamp/internal/test_util"
)

func TestNewStore(t *testing.T) {
	store := NewStore("/some/path")
	assert.Equal(t, "/some/path", store.BasePath)
	assert.Equal(t, DEFAULT_META_FILE, store.MetaFile)

	store = NewStore("/some/path").WithMetaFile("custom.yml")
	assert.Equal(t, "/some/path", store.BasePath)
	assert.Equal(t, "custom.yml", store.MetaFile)
}

func TestLoadingIndividualPackages(t *testing.T) {
	store := NewStore(packageFixturesDir())

	pkg, err := store.Load("minimal")
	if assert.NotNil(t, pkg) {
		assert.Equal(t, packageFixtureDir("minimal"), pkg.Path())
		assert.Equal(t, "minimal", pkg.Name())
	}
	assert.NoError(t, err)

	pkg, err = store.Load("nested:aaa:111")
	if assert.NotNil(t, pkg) {
		assert.Equal(t, packageFixtureDir("nested/aaa/111"), pkg.Path())
		assert.Equal(t, "nested:aaa:111", pkg.Name())
	}
	assert.NoError(t, err)
}

func TestLoadingAllPackages(t *testing.T) {
	store := NewStore(packageFixturesDir())
	items, err := store.LoadAll()

	if assert.Len(t, items, 6) {
		assert.Equal(t, "minimal", items[0].Name())
		assert.Equal(t, "nested", items[1].Name())
		assert.Equal(t, "nested:aaa", items[2].Name())
		assert.Equal(t, "nested:aaa:111", items[3].Name())
		assert.Equal(t, "nested:bbb", items[4].Name())
		assert.Equal(t, "nested:ccc", items[5].Name())
	}
	assert.NoError(t, err)
}

func TestLoadingAllPackagesWhenBasePathMissing(t *testing.T) {
	store := NewStore(packageFixtureDir("non-existent"))
	items, err := store.LoadAll()

	assert.Len(t, items, 0)
	assert.ErrorContains(t, err, "no such file or directory")
}

func TestInstallingPackages(t *testing.T) {
	defer test_util.Cleanup()

	storeDir := test_util.MkdirTemp(t)
	pkgDir := path.Join(storeDir, "minimal")

	store := NewStore(storeDir)
	pkg, err := store.Install(packageFixtureDir("minimal"))

	// Should have been returned the installed package
	if assert.NotNil(t, pkg) {
		assert.Equal(t, "minimal", pkg.Name())
		assert.Equal(t, pkgDir, pkg.Path())
	}
	assert.NoError(t, err)

	// And the package should have been copied to the store dir
	info, err := os.Stat(pkgDir)
	assert.NoError(t, err)
	if assert.NotNil(t, info) {
		assert.Equal(t, true, info.IsDir())
	}

	// And it should be loadable
	item, err := store.Load("minimal")
	assert.Equal(t, pkg, item)
	assert.NoError(t, err)
	items, err := store.LoadAll()
	if assert.Len(t, items, 1) {
		assert.Equal(t, "minimal", items[0].Name())
	}
	assert.NoError(t, err)

	// Should not be able to install another w/ the same name
	pkg, err = store.Install(packageFixtureDir("minimal"))
	assert.Nil(t, pkg)
	assert.ErrorIs(t, err, ErrPkgExists)
}

func TestInstallingPackagesWhenGetError(t *testing.T) {
	getter := NewMockGetter(func(ctx context.Context, src, dst string) error {
		return ErrUnknown
	})
	store := NewStore(packageFixturesDir()).WithGetter(getter.Get)
	pkg, err := store.Install("https://github.com/example/repo")

	assert.Nil(t, pkg)
	assert.ErrorIs(t, err, ErrUnknown)

	// Mock should have been called properly
	assert.Equal(t, true, getter.Called)
	assert.Equal(t, "https://github.com/example/repo", getter.Src)
	// The staging dir should have been cleaned up
	_, err = os.Stat(getter.Dst)
	assert.Error(t, err)
	assert.Equal(t, true, errors.Is(err, fs.ErrNotExist))
}
