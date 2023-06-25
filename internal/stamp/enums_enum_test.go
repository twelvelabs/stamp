package stamp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/swaggest/jsonschema-go"
)

// dummy tests that exercise all the generated enum methods
// and ensure that coverage numbers don't take a hit.

func TestConflictConfig(t *testing.T) {
	name := ConflictConfigNames()[0]
	enum := ConflictConfig(name)

	assert.Equal(t, true, enum.IsValid())
	assert.Equal(t, name, enum.String())

	buf, err := enum.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, []byte(name), buf)

	err = (&enum).UnmarshalText(buf)
	assert.NoError(t, err)
	err = (&enum).UnmarshalText([]byte{})
	assert.Error(t, err)

	err = enum.PrepareJSONSchema(&jsonschema.Schema{})
	assert.NoError(t, err)
}

func TestMatchSource(t *testing.T) {
	name := MatchSourceNames()[0]
	enum := MatchSource(name)

	assert.Equal(t, true, enum.IsValid())
	assert.Equal(t, name, enum.String())

	buf, err := enum.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, []byte(name), buf)

	err = (&enum).UnmarshalText(buf)
	assert.NoError(t, err)
	err = (&enum).UnmarshalText([]byte{})
	assert.Error(t, err)

	err = enum.PrepareJSONSchema(&jsonschema.Schema{})
	assert.NoError(t, err)
}

func TestMissingConfig(t *testing.T) {
	name := MissingConfigNames()[0]
	enum := MissingConfig(name)

	assert.Equal(t, true, enum.IsValid())
	assert.Equal(t, name, enum.String())

	buf, err := enum.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, []byte(name), buf)

	err = (&enum).UnmarshalText(buf)
	assert.NoError(t, err)
	err = (&enum).UnmarshalText([]byte{})
	assert.Error(t, err)

	err = enum.PrepareJSONSchema(&jsonschema.Schema{})
	assert.NoError(t, err)
}

func TestFileType(t *testing.T) {
	name := FileTypeNames()[0]
	enum := FileType(name)

	assert.Equal(t, true, enum.IsValid())
	assert.Equal(t, name, enum.String())

	buf, err := enum.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, []byte(name), buf)

	err = (&enum).UnmarshalText(buf)
	assert.NoError(t, err)
	err = (&enum).UnmarshalText([]byte{})
	assert.Error(t, err)

	err = enum.PrepareJSONSchema(&jsonschema.Schema{})
	assert.NoError(t, err)
}

func TestParseFileTypeFromPath(t *testing.T) {
	tests := []struct {
		path      string
		expected  FileType
		assertion assert.ErrorAssertionFunc
	}{
		{
			path:      "example.json",
			expected:  FileTypeJson,
			assertion: assert.NoError,
		},
		{
			path:      "example.yaml",
			expected:  FileTypeYaml,
			assertion: assert.NoError,
		},
		{
			path:      "example.yml",
			expected:  FileTypeYaml,
			assertion: assert.NoError,
		},
		{
			path:      "example.text",
			expected:  FileTypeText,
			assertion: assert.NoError,
		},
		{
			path:      "example.nope",
			expected:  FileTypeText,
			assertion: assert.NoError,
		},
		{
			path:      "example",
			expected:  FileTypeText,
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			actual, err := ParseFileTypeFromPath(tt.path)
			tt.assertion(t, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
