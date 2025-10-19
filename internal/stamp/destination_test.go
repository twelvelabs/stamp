package stamp

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twelvelabs/termite/render"
	"github.com/twelvelabs/termite/testutil"
)

func TestDestination_ForPath(t *testing.T) {
	var mode os.FileMode = 0655

	dst := Destination{
		ContentTypeTpl: *render.MustCompile(`text`),
		ModeTpl:        *render.MustCompile(`0655`),
		PathTpl:        *render.MustCompile(`templates/`),
	}
	values := map[string]any{
		"DstPath": "testdata",
	}

	err := dst.SetValues(values)
	assert.NoError(t, err)
	assert.Equal(t, "templates", filepath.Base(dst.Path()))
	assert.Equal(t, FileTypeText, dst.ContentType())
	assert.Equal(t, mode, dst.Mode())
	assert.Equal(t, true, dst.IsDir())

	dst, err = dst.ForPath("templates/valid.json", values)
	assert.NoError(t, err)

	assert.Equal(t, "valid.json", filepath.Base(dst.Path()))
	assert.Equal(t, FileTypeText, dst.ContentType())
	assert.Equal(t, mode, dst.Mode())
	assert.Equal(t, false, dst.IsDir())

	dst, err = dst.ForPath("../../foo", values)
	assert.ErrorContains(t, err, "attempted to traverse outside of")
}

func TestDestination_FilesystemMethods(t *testing.T) {
	var err error

	dst := Destination{}
	values := map[string]any{}

	dst, err = dst.ForPath("unknown", values)
	assert.NoError(t, err)
	assert.Equal(t, false, dst.Exists())
	assert.Equal(t, false, dst.IsDir())

	dst, err = dst.ForPath("testdata", values)
	assert.NoError(t, err)
	assert.Equal(t, true, dst.Exists())
	assert.Equal(t, true, dst.IsDir())

	dst, err = dst.ForPath("testdata/templates/valid.txt", values)
	assert.NoError(t, err)
	assert.Equal(t, true, dst.Exists())
	assert.Equal(t, false, dst.IsDir())
}

func TestDestination_SetValues(t *testing.T) {
	tests := []struct {
		desc   string
		dest   Destination
		values map[string]any
		setup  func(t *testing.T, dir string)
		assert func(t *testing.T, d Destination)
		err    string
	}{
		/**
		* PATH
		**/
		{
			desc: "should return error if path is missing",
			dest: Destination{
				PathTpl: *render.MustCompile(`{{ .Foo }}`),
			},
			values: map[string]any{
				"Foo": "",
			},
			err: "evaluated to an empty string",
		},
		{
			desc: "should return error if path traverses outside of destination dir",
			dest: Destination{
				PathTpl: *render.MustCompile(`../nope`),
			},
			err: "dst path validate: ../nope attempted to traverse outside",
		},
		{
			desc: "should parse a valid path",
			dest: Destination{
				PathTpl: *render.MustCompile(`example.txt`),
			},
			assert: func(t *testing.T, d Destination) {
				t.Helper()
				assert.Equal(t, "example.txt", filepath.Base(d.Path()))
			},
		},

		/**
		* CONTENT TYPE
		**/
		{
			desc: "should return error if unable to render content_type",
			dest: Destination{
				PathTpl:        *render.MustCompile(`example.txt`),
				ContentTypeTpl: *render.MustCompile(`{{ fail "boom" }}`),
			},
			err: "boom",
		},
		{
			desc: "should return error if unable to parse content_type",
			dest: Destination{
				PathTpl:        *render.MustCompile(`example.txt`),
				ContentTypeTpl: *render.MustCompile(`unknown`),
			},
			err: "unknown is not a valid",
		},
		{
			desc: "should parse a valid content_type",
			dest: Destination{
				PathTpl:        *render.MustCompile(`example.txt`),
				ContentTypeTpl: *render.MustCompile(`json`),
			},
			assert: func(t *testing.T, d Destination) {
				t.Helper()
				assert.Equal(t, FileTypeJson, d.ContentType())
			},
		},
		{
			desc: "should infer missing content type from path",
			dest: Destination{
				PathTpl:        *render.MustCompile(`example.json`),
				ContentTypeTpl: *render.MustCompile(``),
			},
			assert: func(t *testing.T, d Destination) {
				t.Helper()
				assert.Equal(t, FileTypeJson, d.ContentType())
			},
		},

		/**
		* MODE
		**/
		{
			desc: "should not require mode",
			dest: Destination{
				PathTpl: *render.MustCompile(`example.json`),
				ModeTpl: *render.MustCompile(``),
			},
		},
		{
			desc: "should return error if unable to render mode",
			dest: Destination{
				PathTpl: *render.MustCompile(`example.json`),
				ModeTpl: *render.MustCompile(`{{ fail "boom" }}`),
			},
			err: "boom",
		},
		{
			desc: "should return error if unable to parse mode",
			dest: Destination{
				PathTpl: *render.MustCompile(`example.json`),
				ModeTpl: *render.MustCompile(`not a file mode`),
			},
			err: "invalid syntax",
		},
		{
			desc: "should parse a valid mode",
			dest: Destination{
				PathTpl: *render.MustCompile(`example.json`),
				ModeTpl: *render.MustCompile(`0755`),
			},
			assert: func(t *testing.T, d Destination) {
				t.Helper()
				assert.Equal(t, os.FileMode(0755), d.Mode())
			},
		},
		{
			desc: "should set a default mode for create",
			dest: Destination{
				PathTpl: *render.MustCompile(`example.json`),
			},
			assert: func(t *testing.T, d Destination) {
				t.Helper()
				assert.Equal(t, DstFileMode, d.Mode())
			},
		},
		{
			desc: "should not set a default mode when updating",
			setup: func(t *testing.T, dir string) {
				t.Helper()
				testutil.WritePaths(t, dir, map[string]any{
					"example.json": `{}`,
				})
			},
			dest: Destination{
				PathTpl: *render.MustCompile(`example.json`),
			},
			assert: func(t *testing.T, d Destination) {
				t.Helper()
				assert.Equal(t, fs.FileMode(0), d.Mode())
			},
		},

		/**
		* CONTENT
		**/
		{
			desc: "should return error if unable to parse content",
			setup: func(t *testing.T, dir string) {
				t.Helper()
				testutil.WritePaths(t, dir, map[string]any{
					"example.json": `{,}`,
				})
			},
			dest: Destination{
				PathTpl: *render.MustCompile(`example.json`),
			},
			err: "json decode",
		},
		{
			desc: "should parse content",
			setup: func(t *testing.T, dir string) {
				t.Helper()
				testutil.WritePaths(t, dir, map[string]any{
					"example.json": `{"foo": "bar"}`,
				})
			},
			dest: Destination{
				PathTpl: *render.MustCompile(`example.json`),
			},
			assert: func(t *testing.T, d Destination) {
				t.Helper()
				assert.Equal(t, map[string]any{
					"foo": "bar",
				}, d.Content())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			testutil.InTempDir(t, func(dir string) {
				if tt.values == nil {
					tt.values = map[string]any{}
				}
				tt.values["DstPath"] = "."

				if tt.setup != nil {
					tt.setup(t, dir)
				}

				err := tt.dest.SetValues(tt.values)

				if tt.err == "" {
					assert.NoError(t, err)
				} else {
					assert.ErrorContains(t, err, tt.err)
				}

				if tt.assert != nil {
					tt.assert(t, tt.dest)
				}
			})
		})
	}
}

// cspell: words unmarshalable
type unmarshalable struct {
}

func (u unmarshalable) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("boom")
}

func TestDestination_Write(t *testing.T) {
	tests := []struct {
		desc   string
		dest   Destination
		values map[string]any
		data   any
		setup  func(t *testing.T, dir string)
		assert func(t *testing.T, dir string, d Destination)
		err    string
	}{
		{
			desc: "should return an error when unable to encode",
			dest: Destination{
				PathTpl: *render.MustCompile(`example.json`),
			},
			data: unmarshalable{},
			assert: func(t *testing.T, dir string, d Destination) {
				t.Helper()
				assert.Nil(t, d.Content())
			},
			err: "boom",
		},
		{
			desc: "should create a new file when one does not yet exist",
			dest: Destination{
				PathTpl: *render.MustCompile(`example.txt`),
			},
			data: []byte("Hello"),
			assert: func(t *testing.T, dir string, d Destination) {
				t.Helper()
				testutil.AssertPaths(t, dir, map[string]any{
					"example.txt": []any{
						"Hello",
						DstFileMode,
					},
				})
				assert.Equal(t, []byte("Hello"), d.Content())
			},
		},
		{
			desc: "should update an existing file with new content",
			setup: func(t *testing.T, dir string) {
				t.Helper()
				testutil.WritePaths(t, dir, map[string]any{
					"example.txt": "Hello",
				})
			},
			dest: Destination{
				PathTpl: *render.MustCompile(`example.txt`),
				ModeTpl: *render.MustCompile(`0600`),
			},
			data: []byte("Goodbye"),
			assert: func(t *testing.T, dir string, d Destination) {
				t.Helper()
				testutil.AssertPaths(t, dir, map[string]any{
					"example.txt": []any{
						"Goodbye",
						0600,
					},
				})
				assert.Equal(t, []byte("Goodbye"), d.Content())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			testutil.InTempDir(t, func(dir string) {
				if tt.values == nil {
					tt.values = map[string]any{}
				}
				tt.values["DstPath"] = "."

				if tt.setup != nil {
					tt.setup(t, dir)
				}

				err := tt.dest.SetValues(tt.values)
				assert.NoError(t, err)

				err = tt.dest.Write(tt.data)
				if tt.err == "" {
					assert.NoError(t, err)
				} else {
					assert.ErrorContains(t, err, tt.err)
				}

				if tt.assert != nil {
					tt.assert(t, dir, tt.dest)
				}
			})
		})
	}
}

func TestDestination_Delete(t *testing.T) {
	testutil.InTempDir(t, func(dir string) {
		testutil.WritePaths(t, dir, map[string]any{
			"example.txt": "Hello",
		})

		dest := Destination{
			PathTpl: *render.MustCompile(`example.txt`),
		}
		err := dest.SetValues(map[string]any{
			"DstPath": ".",
		})
		assert.NoError(t, err)
		assert.FileExists(t, dest.Path())
		assert.Equal(t, []byte("Hello"), dest.Content())

		err = dest.Delete()
		assert.NoError(t, err)
		assert.NoFileExists(t, dest.Path())
	})
}
