package pkg

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultGetter(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pkg-")
	if err != nil {
		assert.FailNow(t, "unable to create temp dir")
	}
	defer os.RemoveAll(tmpDir)

	pkgSrcPath := packageFixtureDir("minimal")
	pkgDstPath := path.Join(tmpDir, "package")
	pkgManifestPath := path.Join(pkgDstPath, DEFAULT_META_FILE)

	assert.NoDirExists(t, pkgDstPath)
	assert.NoFileExists(t, pkgManifestPath)

	err = DefaultGetter(context.Background(), pkgSrcPath, pkgDstPath)

	assert.NoError(t, err)
	assert.FileExists(t, pkgManifestPath)
}
