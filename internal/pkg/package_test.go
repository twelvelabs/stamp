package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPackageTreeMethods(t *testing.T) {
	nested, err := LoadPackage(packageFixtureDir("nested"), DEFAULT_META_FILE)

	assert.NoError(t, err)
	if assert.NotNil(t, nested) {
		children, err := nested.Children()
		assert.NoError(t, err)
		if assert.Len(t, children, 4) {
			assert.Equal(t, "nested:aaa", children[0].Name())
			assert.Equal(t, "nested:aaa:111", children[1].Name())
			assert.Equal(t, "nested:bbb", children[2].Name())
			assert.Equal(t, "nested:ccc", children[3].Name())
		}
		assert.Nil(t, nested.Parent())
		assert.Equal(t, nested, nested.Root())
	}

	aaa, err := LoadPackage(packageFixtureDir("nested/aaa"), DEFAULT_META_FILE)

	assert.NoError(t, err)
	if assert.NotNil(t, aaa) {
		children, err := aaa.Children()
		assert.NoError(t, err)
		if assert.Len(t, children, 1) {
			assert.Equal(t, "nested:aaa:111", children[0].Name())
		}
		parent := aaa.Parent()
		if assert.NotNil(t, parent) {
			assert.Equal(t, nested.Name(), parent.Name())
		}
		assert.Equal(t, nested, aaa.Root())
	}
}
