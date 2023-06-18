package stamp

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twelvelabs/termite/render"
	"github.com/twelvelabs/termite/testutil"
)

func TestNewSourceWithValues(t *testing.T) {
	values := map[string]any{
		"foo": "bar",
	}
	src, err := NewSourceWithValues("{{ .foo }}", values)
	assert.NoError(t, err)
	assert.Equal(t, "bar", filepath.Base(src.Path()))

	src, err = NewSourceWithValues("{{}", values)
	assert.ErrorContains(t, err, `unexpected "}" in command`)

	src, err = NewSourceWithValues("../../foo", values)
	assert.ErrorContains(t, err, "attempted to traverse outside of")
}

func TestSource(t *testing.T) {
	tests := []struct {
		desc   string
		src    Source
		values map[string]any
		setup  func(t *testing.T, dir string)
		assert func(t *testing.T, s Source)
		err    string
	}{
		/**
		* PATH
		**/
		{
			desc: "should return error if unable to render path",
			src: Source{
				PathTpl: *render.MustCompile(`{{ fail "boom" }}`),
			},
			err: "boom",
		},
		{
			desc: "should return error if both path and content are present",
			src: Source{
				PathTpl:       *render.MustCompile(`example.txt`),
				InlineContent: "inline content",
			},
			err: "path and content are mutually exclusive",
		},
		{
			desc: "should return error if path traverses outside of destination dir",
			src: Source{
				PathTpl: *render.MustCompile(`../nope`),
			},
			err: "../nope attempted to traverse outside",
		},
		{
			desc: "should parse a valid path",
			setup: func(t *testing.T, dir string) {
				t.Helper()
			},
			src: Source{
				PathTpl: *render.MustCompile(`example.txt`),
			},
			assert: func(t *testing.T, s Source) {
				t.Helper()
				assert.Equal(t, "example.txt", filepath.Base(s.Path()))
			},
		},

		/**
		* CONTENT TYPE
		**/
		{
			desc: "should return error if unable to render content_type",
			src: Source{
				PathTpl:        *render.MustCompile(`example.txt`),
				ContentTypeTpl: *render.MustCompile(`{{ fail "boom" }}`),
			},
			err: "boom",
		},
		{
			desc: "should return error if unable to parse content_type",
			src: Source{
				PathTpl:        *render.MustCompile(`example.txt`),
				ContentTypeTpl: *render.MustCompile(`unknown`),
			},
			err: "unknown is not a valid",
		},
		{
			desc: "should parse a valid content_type",
			src: Source{
				PathTpl:        *render.MustCompile(`example.txt`),
				ContentTypeTpl: *render.MustCompile(`json`),
			},
			assert: func(t *testing.T, s Source) {
				t.Helper()
				assert.Equal(t, FileTypeJson, s.ContentType())
			},
		},
		{
			desc: "should infer missing content type from path",
			src: Source{
				PathTpl:        *render.MustCompile(`example.json`),
				ContentTypeTpl: *render.MustCompile(``),
			},
			assert: func(t *testing.T, s Source) {
				t.Helper()
				assert.Equal(t, FileTypeJson, s.ContentType())
			},
		},

		/**
		* INLINE CONTENT
		**/
		{
			desc: "should return error if unable to parse inline content",
			src: Source{
				InlineContent: `{{ fail "boom" }}`,
			},
			err: "boom",
		},
		{
			desc: "should parse inline content",
			src: Source{
				InlineContent: map[string]any{
					"foo": "bar",
				},
			},
			assert: func(t *testing.T, s Source) {
				t.Helper()
				assert.Equal(t, map[string]any{
					"foo": "bar",
				}, s.Content())
			},
		},

		/**
		* PATH CONTENT
		**/
		{
			desc: "should return error if unable to render path content",
			setup: func(t *testing.T, dir string) {
				t.Helper()
				testutil.WritePaths(t, dir, map[string]any{
					"example.json": `{{ fail "boom" }}`,
				})
			},
			src: Source{
				PathTpl: *render.MustCompile(`example.json`),
			},
			err: "boom",
		},
		{
			desc: "should return error if unable to parse path content",
			setup: func(t *testing.T, dir string) {
				t.Helper()
				testutil.WritePaths(t, dir, map[string]any{
					"example.json": `{,}`,
				})
			},
			src: Source{
				PathTpl: *render.MustCompile(`example.json`),
			},
			err: "json decode",
		},
		{
			desc: "should render and parse path content",
			setup: func(t *testing.T, dir string) {
				t.Helper()
				testutil.WritePaths(t, dir, map[string]any{
					"example.json": `{"foo": "{{ .FooVal }}"}`,
				})
			},
			src: Source{
				PathTpl: *render.MustCompile(`example.json`),
			},
			values: map[string]any{
				"FooVal": "bar",
			},
			assert: func(t *testing.T, s Source) {
				t.Helper()
				assert.Equal(t, map[string]any{
					"foo": "bar",
				}, s.Content())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			testutil.InTempDir(t, func(dir string) {
				if tt.values == nil {
					tt.values = map[string]any{}
				}
				tt.values["SrcPath"] = "."

				if tt.setup != nil {
					tt.setup(t, dir)
				}

				err := tt.src.SetValues(tt.values)

				if tt.err == "" {
					assert.NoError(t, err)
				} else {
					assert.ErrorContains(t, err, tt.err)
				}

				if tt.assert != nil {
					tt.assert(t, tt.src)
				}
			})
		})
	}
}
