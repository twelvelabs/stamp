package stamp

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twelvelabs/termite/testutil"

	"github.com/twelvelabs/stamp/internal/pkg"
)

func TestNewGenerator(t *testing.T) {
	store := NewTestStore()

	_, err := NewGenerator(nil, nil)
	assert.ErrorContains(t, err, "nil store")

	_, err = NewGenerator(store, nil)
	assert.ErrorContains(t, err, "nil package")

	var p *pkg.Package

	p = &pkg.Package{
		Metadata: map[string]any{
			"visibility": "hidden",
		},
	}
	_, err = NewGenerator(store, p)
	assert.NoError(t, err)

	p = &pkg.Package{
		Metadata: map[string]any{
			"visibility": "purple",
		},
	}
	_, err = NewGenerator(store, p)
	assert.ErrorContains(t, err, "generator metadata invalid")

	p = &pkg.Package{
		Metadata: map[string]any{
			"values": []any{
				map[string]any{
					"key": 123, // key should be a string
				},
			},
		},
	}
	_, err = NewGenerator(store, p)
	assert.ErrorContains(t, err, "generator metadata invalid")

	p = &pkg.Package{
		Metadata: map[string]any{
			"tasks": []any{
				map[string]any{
					"type": "unknown", // unknown type
				},
			},
		},
	}
	_, err = NewGenerator(store, p)
	assert.ErrorContains(t, err, "generator metadata invalid")
}

func TestNewGenerators(t *testing.T) {
	store := NewTestStore()

	items, err := NewGenerators(nil, []*pkg.Package{})
	assert.Nil(t, items)
	assert.ErrorContains(t, err, "nil store")

	items, err = NewGenerators(store, []*pkg.Package{nil})
	assert.Nil(t, items)
	assert.ErrorContains(t, err, "nil package")

	items, err = NewGenerators(store, []*pkg.Package{})
	assert.Equal(t, []*Generator{}, items)
	assert.NoError(t, err)

	p1 := &pkg.Package{}
	p2 := &pkg.Package{}
	items, err = NewGenerators(store, []*pkg.Package{p1, p2})
	assert.Len(t, items, 2)
	assert.NoError(t, err)
}

func TestGenerator_AddsValuesFromDelegatedGenerators(t *testing.T) {
	store := NewTestStore() // must call before changing dirs
	app := NewTestApp()
	app.Store = store

	testutil.InTempDir(t, func(tmpDir string) {
		gen, err := store.Load("delegating")
		assert.NotNil(t, gen)
		assert.NoError(t, err)

		values := gen.Values.GetAll()

		assert.Len(t, values, 3)
		assert.Equal(t, "customized.txt", values["FileName"])
		assert.Equal(t, "custom content", values["FileContent"])

		ctx := NewTaskContext(app)
		err = gen.Tasks.Execute(ctx, values)
		assert.NoError(t, err)

		testutil.AssertPaths(t, tmpDir, map[string]any{
			"customized.txt": "custom content",
		})
	})

	testutil.InTempDir(t, func(tmpDir string) {
		gen, err := store.Load("delegating-dupe")
		assert.NotNil(t, gen)
		assert.NoError(t, err)

		values := gen.Values.GetAll()

		// should only be three values, even though the generator was referenced twice
		assert.Len(t, values, 3)
		// the defaults should be set by the last generator task in the list.
		assert.Equal(t, "untitled.txt", values["FileName"])
		assert.Equal(t, "", values["FileContent"])

		ctx := NewTaskContext(app)
		err = gen.Tasks.Execute(ctx, values)
		assert.NoError(t, err)

		// Should have respected the `values` attribute
		// and created two different filenames.
		testutil.AssertPaths(t, tmpDir, map[string]any{
			"customized.txt": "custom content",
			"untitled.txt":   "",
		})
	})
}

func TestGenerator_Description(t *testing.T) {
	gen := &Generator{
		Package: &pkg.Package{
			Metadata: map[string]any{
				"description": "a test generator",
			},
		},
	}
	assert.Equal(t, "a test generator", gen.Description())
}

func TestGenerator_ShortDescription(t *testing.T) {
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
			gen := &Generator{
				Package: &pkg.Package{
					Metadata: map[string]any{
						"description": tt.given,
					},
				},
			}

			actual := gen.ShortDescription()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestGenerator_SrcPath(t *testing.T) {
	store := NewTestStore()
	gen, err := store.Load("file")

	srcPath := filepath.Join(gen.Path(), "_src")

	assert.NoError(t, err)
	assert.Equal(t, srcPath, gen.SrcPath())
}

func TestGeneratorName(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T, tmpDir string)
		path  string
		want  string
	}{
		{
			name: "returns empty string if missing path",
			path: "",
			want: "",
		},
		{
			name: "returns empty string if not a generator",
			setup: func(t *testing.T, tmpDir string) {
				t.Helper()
				testutil.WritePaths(t, tmpDir, map[string]any{
					"aaa/": "",
				})
			},
			path: "aaa",
			want: "",
		},
		{
			name: "returns non-nested names",
			setup: func(t *testing.T, tmpDir string) {
				t.Helper()
				testutil.WritePaths(t, tmpDir, map[string]any{
					"aaa/":               "",
					"aaa/generator.yaml": "name: 'aaa'",
				})
			},
			path: "aaa",
			want: "aaa",
		},
		{
			name: "returns nested names",
			setup: func(t *testing.T, tmpDir string) {
				t.Helper()
				testutil.WritePaths(t, tmpDir, map[string]any{
					"aaa/":                       "",
					"aaa/bbb/":                   "",
					"aaa/bbb/ccc/":               "",
					"aaa/bbb/ccc/generator.yaml": "name: 'aaa:bbb:ccc'",
				})
			},
			path: "aaa/bbb/ccc",
			want: "aaa:bbb:ccc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutil.InTempDir(t, func(tmpDir string) {
				if tt.setup != nil {
					tt.setup(t, tmpDir)
				}
				assert.Equal(t, tt.want, GeneratorName(tt.path))
			})
		})
	}
}

func TestGeneratorNameForCreate(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T, tmpDir string)
		path  string
		want  string
	}{
		// {
		// 	name: "returns empty string if missing path",
		// 	path: "",
		// 	want: "",
		// },
		{
			name: "returns immediate dirname when no generator in ancestor tree",
			setup: func(t *testing.T, tmpDir string) {
				t.Helper()
				testutil.WritePaths(t, tmpDir, map[string]any{
					"aaa/":         "",
					"aaa/bbb/":     "",
					"aaa/bbb/ccc/": "",
				})
			},
			path: "aaa/bbb/ccc",
			want: "ccc",
		},
		{
			name: "returns name prefixed with parent generator",
			setup: func(t *testing.T, tmpDir string) {
				t.Helper()
				testutil.WritePaths(t, tmpDir, map[string]any{
					"aaa/":                   "",
					"aaa/bbb/":               "",
					"aaa/bbb/generator.yaml": "name: 'bbb'",
					"aaa/bbb/ccc/":           "",
				})
			},
			path: "aaa/bbb/ccc",
			want: "bbb:ccc",
		},
		{
			name: "returns name prefixed with furthest parent generator",
			setup: func(t *testing.T, tmpDir string) {
				t.Helper()
				testutil.WritePaths(t, tmpDir, map[string]any{
					"aaa/":                   "",
					"aaa/generator.yaml":     "name: 'aaa'",
					"aaa/bbb/":               "",
					"aaa/bbb/generator.yaml": "name: 'aaa:bbb'",
					"aaa/bbb/ccc/":           "",
				})
			},
			path: "aaa/bbb/ccc",
			want: "aaa:bbb:ccc",
		},
		{
			name: "returns name prefixed with furthest parent generator even when gaps",
			setup: func(t *testing.T, tmpDir string) {
				t.Helper()
				testutil.WritePaths(t, tmpDir, map[string]any{
					"aaa/":               "",
					"aaa/generator.yaml": "name: 'aaa'",
					"aaa/bbb/":           "",
					"aaa/bbb/ccc/":       "",
				})
			},
			path: "aaa/bbb/ccc",
			want: "aaa:bbb:ccc",
		},
		{
			name: "returns name prefixed with furthest parent generator even when custom name",
			setup: func(t *testing.T, tmpDir string) {
				t.Helper()
				testutil.WritePaths(t, tmpDir, map[string]any{
					"aaa/":               "",
					"aaa/generator.yaml": "name: 'foo'",
					"aaa/bbb/":           "",
					"aaa/bbb/ccc/":       "",
				})
			},
			path: "aaa/bbb/ccc",
			want: "foo:bbb:ccc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutil.InTempDir(t, func(tmpDir string) {
				if tt.setup != nil {
					tt.setup(t, tmpDir)
				}
				assert.Equal(t, tt.want, GeneratorNameForCreate(tt.path))
			})
		})
	}
}
