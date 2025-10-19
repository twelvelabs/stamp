package pkg

import (
	"context"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultGetter(t *testing.T) {
	tmpDir := t.TempDir()

	pkgSrcPath := packageFixtureDir("minimal")
	pkgDstPath := path.Join(tmpDir, "package")
	pkgManifestPath := path.Join(pkgDstPath, DefaultMetaFile)

	assert.NoDirExists(t, pkgDstPath)
	assert.NoFileExists(t, pkgManifestPath)

	err := DefaultGetter(context.Background(), pkgSrcPath, pkgDstPath)

	assert.NoError(t, err)
	assert.FileExists(t, pkgManifestPath)
}
