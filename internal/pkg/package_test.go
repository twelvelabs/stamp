package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPackageTreeMethods(t *testing.T) {
	nested, err := LoadPackage(packageFixtureDir("nested"), DefaultMetaFile)

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

	aaa, err := LoadPackage(packageFixtureDir("nested/aaa"), DefaultMetaFile)

	assert.NoError(t, err)
	if assert.NotNil(t, aaa) {
		children, err := aaa.Children()
		assert.NoError(t, err)
		if assert.Len(t, children, 1) {
			assert.Equal(t, "nested:aaa:111", children[0].Name())
		}

		all, err := aaa.All()
		assert.NoError(t, err)
		if assert.Len(t, all, 2) {
			assert.Equal(t, "nested:aaa", all[0].Name())
			assert.Equal(t, "nested:aaa:111", all[1].Name())
		}

		parent := aaa.Parent()
		if assert.NotNil(t, parent) {
			assert.Equal(t, nested.Name(), parent.Name())
		}
		assert.Equal(t, nested, aaa.Root())
	}
}

func TestPackage_Name(t *testing.T) {
	p := &Package{
		Metadata: map[string]any{},
	}
	assert.Equal(t, "", p.Name())
	p.SetName("foo")
	assert.Equal(t, "foo", p.Name())
}

func TestPackage_Description(t *testing.T) {
	p := &Package{
		Metadata: map[string]any{},
	}
	assert.Equal(t, "", p.Description())
	p.SetDescription("foo")
	assert.Equal(t, "foo", p.Description())
}

func TestPackage_ShortDescription(t *testing.T) {
	tests := []struct {
		desc     string
		given    string
		expected string
	}{
		{
			desc:     "empty string is a noop",
			given:    "",
			expected: "",
		},
		{
			desc:     "single line is a noop",
			given:    "Example description",
			expected: "Example description",
		},
		{
			desc:     "otherwise first line is returned",
			given:    "Example description\nExtended info",
			expected: "Example description",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			p := &Package{
				Metadata: map[string]any{
					"description": tt.given,
				},
			}

			actual := p.ShortDescription()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestPackage_Origin(t *testing.T) {
	p := &Package{
		Metadata: map[string]any{},
	}
	assert.Equal(t, "", p.Origin())
	p.SetOrigin("~/packages/foo")
	assert.Equal(t, "~/packages/foo", p.Origin())
}

func TestPackage_MetadataLookup(t *testing.T) {
	tests := []struct {
		Desc     string
		Metadata map[string]any
		Key      string
		Value    any
	}{
		{
			Desc:     "returns nil for unset keys",
			Metadata: map[string]any{},
			Key:      "unknown",
			Value:    nil,
		},
		{
			Desc: "returns the value for key if present",
			Metadata: map[string]any{
				"foo_bar": "baz",
			},
			Key:   "foo_bar",
			Value: "baz",
		},
		{
			Desc: "returns the value for Pascal-key if present",
			Metadata: map[string]any{
				"FooBar": "baz",
			},
			Key:   "foo_bar",
			Value: "baz",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Desc, func(t *testing.T) {
			p := &Package{
				Metadata: tt.Metadata,
			}
			assert.Equal(t, tt.Value, p.MetadataLookup(tt.Key))
		})
	}
}

func TestPackage_MetadataSlice(t *testing.T) {
	tests := []struct {
		Desc     string
		Metadata map[string]any
		Key      string
		Value    any
		Panics   bool
	}{
		{
			Desc:     "returns empty slice for unset keys",
			Metadata: map[string]any{},
			Key:      "items",
			Value:    []any{},
		},
		{
			Desc: "returns the value for key if present",
			Metadata: map[string]any{
				"items": []any{1, 2, 3},
			},
			Key:   "items",
			Value: []any{1, 2, 3},
		},
		{
			Desc: "panics if value is not a slice",
			Metadata: map[string]any{
				"items": "non-slice",
			},
			Key:    "items",
			Panics: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Desc, func(t *testing.T) {
			p := &Package{
				Metadata: tt.Metadata,
			}

			if tt.Panics {
				assert.Panics(t, func() {
					p.MetadataSlice(tt.Key)
				})
			} else {
				assert.Equal(t, tt.Value, p.MetadataSlice(tt.Key))
			}
		})
	}
}

func TestPackage_MetadataMapSlice(t *testing.T) {
	tests := []struct {
		Desc     string
		Metadata map[string]any
		Key      string
		Value    any
		Panics   bool
	}{
		{
			Desc: "returns empty slice for unset keys",
			Metadata: map[string]any{
				"items": nil,
			},
			Key:   "items",
			Value: []map[string]any{},
		},
		{
			Desc: "returns the value for key if present",
			Metadata: map[string]any{
				"items": []any{
					map[string]any{
						"key":  "aaa",
						"type": "string",
					},
					map[string]any{
						"key":  "bbb",
						"type": "int",
					},
				},
			},
			Key: "items",
			Value: []map[string]any{
				{
					"key":  "aaa",
					"type": "string",
				},
				{
					"key":  "bbb",
					"type": "int",
				},
			},
		},
		{
			Desc: "panics if slice does not contain a map",
			Metadata: map[string]any{
				"items": []any{1, 2},
			},
			Key:    "items",
			Panics: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Desc, func(t *testing.T) {
			p := &Package{
				Metadata: tt.Metadata,
			}

			if tt.Panics {
				assert.Panics(t, func() {
					p.MetadataMapSlice(tt.Key)
				})
			} else {
				assert.Equal(t, tt.Value, p.MetadataMapSlice(tt.Key))
			}
		})
	}
}

func TestPackage_MetadataString(t *testing.T) {
	tests := []struct {
		Desc     string
		Metadata map[string]any
		Key      string
		Value    any
		Panics   bool
	}{
		{
			Desc:     "returns empty string for unset keys",
			Metadata: map[string]any{},
			Key:      "name",
			Value:    "",
		},
		{
			Desc: "returns the value for key if present",
			Metadata: map[string]any{
				"name": "foo",
			},
			Key:   "name",
			Value: "foo",
		},
		{
			Desc: "panics if value is not a slice",
			Metadata: map[string]any{
				"name": 123,
			},
			Key:    "name",
			Panics: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Desc, func(t *testing.T) {
			p := &Package{
				Metadata: tt.Metadata,
			}

			if tt.Panics {
				assert.Panics(t, func() {
					p.MetadataString(tt.Key)
				})
			} else {
				assert.Equal(t, tt.Value, p.MetadataString(tt.Key))
			}
		})
	}
}
