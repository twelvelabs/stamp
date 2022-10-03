package pkg

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/twelvelabs/stamp/internal/testutil"
)

func TestPackagePath(t *testing.T) {
	var tests = []struct {
		PackageName string
		PackagePath string
		Err         string
	}{
		{
			PackageName: "foo",
			PackagePath: "/packages/foo",
			Err:         "",
		},
		{
			PackageName: "foo:bar:baz",
			PackagePath: "/packages/foo/bar/baz",
			Err:         "",
		},
		{
			PackageName: "",
			PackagePath: "",
			Err:         "invalid package name",
		},
		{
			PackageName: "foo:../../..:/etc/sudoers",
			PackagePath: "",
			Err:         "invalid package name",
		},
	}

	for _, test := range tests {
		t.Run(test.PackageName, func(t *testing.T) {
			path, err := PackagePath("/packages", test.PackageName)

			assert.Equal(t, test.PackagePath, path)
			if test.Err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, test.Err)
			}
		})
	}
}

func TestLoadPackage(t *testing.T) {
	var tests = []struct {
		PackageName string
		PackagePath string
		Err         string
	}{
		{
			PackageName: "non-existent",
			PackagePath: packageFixtureDir("non-existent"),
			Err:         "package not found",
		},
		{
			PackageName: "empty",
			PackagePath: packageFixtureDir("empty"),
			Err:         "no such file or directory",
		},
		{
			PackageName: "non-parsable",
			PackagePath: packageFixtureDir("non-parsable"),
			Err:         "unmarshal errors",
		},
		{
			PackageName: "minimal",
			PackagePath: packageFixtureDir("minimal"),
			Err:         "",
		},
		{
			PackageName: "nested",
			PackagePath: packageFixtureDir("nested"),
			Err:         "",
		},
	}

	for _, test := range tests {
		t.Run(test.PackageName, func(t *testing.T) {
			pkg, err := LoadPackage(test.PackagePath, DefaultMetaFile)

			if test.Err == "" {
				if assert.NotNil(t, pkg) {
					assert.Equal(t, test.PackageName, pkg.Name())
					assert.Equal(t, test.PackagePath, pkg.Path())
				}
				assert.NoError(t, err)
			} else {
				assert.Nil(t, pkg)
				assert.ErrorContains(t, err, test.Err)
			}
		})
	}
}

func TestStorePackage(t *testing.T) {
	defer testutil.Cleanup()

	rootPath := testutil.MkdirTemp()
	pkgPath, pkgMetaPath := createPackage(t, rootPath, "foo")

	pkg, err := LoadPackage(pkgPath, DefaultMetaFile)
	if assert.NotNil(t, pkg) {
		assert.Equal(t, "foo", pkg.Name())
	}
	assert.NoError(t, err)

	pkg.SetName("bar")

	err = StorePackage(pkg)
	assert.Equal(t, "bar", pkg.Name())
	assert.NoError(t, err)

	data, err := os.ReadFile(pkgMetaPath)
	assert.Equal(t, "Name: bar\n", string(data))
	assert.NoError(t, err)
}

func TestMovePackage(t *testing.T) {
	defer testutil.Cleanup()

	rootPath := testutil.MkdirTemp()
	fooPath, _ := createPackage(t, rootPath, "foo")
	barPath, _ := createPackage(t, rootPath, "bar")
	renamedPath := path.Join(rootPath, "renamed")

	pkg, err := LoadPackage(fooPath, DefaultMetaFile)
	assert.NotNil(t, pkg)
	assert.NoError(t, err)

	// Can't replace an existing package
	err = MovePackage(pkg, barPath)
	assert.ErrorIs(t, err, ErrPkgExists)

	err = MovePackage(pkg, renamedPath)
	assert.NoError(t, err)
	assert.Equal(t, renamedPath, pkg.Path())

	_, err = os.Stat(fooPath)
	assert.ErrorIs(t, err, os.ErrNotExist)

	_, err = os.Stat(renamedPath)
	assert.NoError(t, err)
}

func TestRemovePackage(t *testing.T) {
	defer testutil.Cleanup()

	rootPath := testutil.MkdirTemp()
	pkgPath, _ := createPackage(t, rootPath, "foo")

	pkg, err := LoadPackage(pkgPath, DefaultMetaFile)
	assert.NotNil(t, pkg)
	assert.NoError(t, err)

	err = RemovePackage(pkg)
	assert.NoError(t, err)

	err = RemovePackage(pkg)
	assert.ErrorContains(t, err, "package not found")

	_, err = os.Stat(pkgPath)
	assert.ErrorIs(t, err, os.ErrNotExist)
}

// Creates a new package directory and metadata file in the root
// dir and returns the paths to both.
func createPackage(t *testing.T, root string, name string) (string, string) {
	t.Helper()
	pkgPath := path.Join(root, strings.ReplaceAll(name, ":", "/"))
	pkgMetaPath := path.Join(pkgPath, DefaultMetaFile)

	err := os.MkdirAll(pkgPath, 0755)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(pkgMetaPath, []byte("Name: "+name), 0600)
	if err != nil {
		t.Fatal(err)
	}

	return pkgPath, pkgMetaPath
}

func workingDir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

// Returns the absolute path to `./testdata/packages`.
func packageFixturesDir() string {
	dir, err := filepath.Abs(path.Join(workingDir(), "..", "..", "testdata", "packages"))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

// Returns the absolute path to the named package.
func packageFixtureDir(name string) string {
	return path.Join(packageFixturesDir(), name)
}

// // Returns the absolute path to metadata file of the named package.
// func packageFixturePath(name string) string {
// 	return path.Join(packageFixturesDir(), name, DEFAULT_META_FILE)
// }
