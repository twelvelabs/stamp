package render

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderFile(t *testing.T) {
	templatesDir := filepath.Join("..", "..", "testdata", "templates")
	tests := []struct {
		Desc     string
		Values   map[string]any
		Template string
		Rendered string
		Err      string
	}{
		{
			Desc:     "returns an error if unable to parse",
			Template: filepath.Join(templatesDir, "parse-error.txt"),
			Rendered: "",
			Err:      "unexpected \"}\" in command",
		},
		{
			Desc:     "returns an error if unable to render",
			Template: filepath.Join(templatesDir, "render-error.txt"),
			Rendered: "",
			Err:      "boom",
		},
		{
			Desc:     "does not error if template values are missing",
			Values:   map[string]any{},
			Template: filepath.Join(templatesDir, "valid.txt"),
			Rendered: "Hello, <no value>",
			Err:      "",
		},
		{
			Desc: "renders template values successfully",
			Values: map[string]any{
				"Name": "World",
			},
			Template: filepath.Join(templatesDir, "valid.txt"),
			Rendered: "Hello, World",
			Err:      "",
		},
		{
			Desc: "renders template functions successfully",
			Values: map[string]any{
				"Name": "Some Name",
			},
			Template: filepath.Join(templatesDir, "valid-func.txt"),
			Rendered: "Hello, some-name",
			Err:      "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Desc, func(t *testing.T) {
			actual, err := RenderFile(tt.Template, tt.Values)
			assert.Equal(t, tt.Rendered, actual)
			if tt.Err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.Err)
			}
		})
	}
}

func TestRenderString(t *testing.T) {
	tests := []struct {
		Desc     string
		Values   map[string]any
		Template string
		Rendered string
		Err      string
	}{
		{
			Desc:     "returns an error if unable to parse",
			Template: "{{}",
			Rendered: "",
			Err:      "unexpected \"}\" in command",
		},
		{
			Desc:     "returns an error if unable to render",
			Template: "{{ fail \"boom\" }}",
			Rendered: "",
			Err:      "boom",
		},
		{
			Desc:     "does not error if template values are missing",
			Values:   map[string]any{},
			Template: "Hello, {{ .NotFound }}",
			Rendered: "Hello, <no value>",
			Err:      "",
		},
		{
			Desc: "renders template values successfully",
			Values: map[string]any{
				"Name": "World",
			},
			Template: "Hello, {{ .Name }}",
			Rendered: "Hello, World",
			Err:      "",
		},
		{
			Desc: "renders template functions successfully",
			Values: map[string]any{
				"Name": "Some Name",
			},
			Template: "Hello, {{ .Name | dasherize }}",
			Rendered: "Hello, some-name",
			Err:      "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Desc, func(t *testing.T) {
			actual, err := RenderString(tt.Template, tt.Values)
			assert.Equal(t, tt.Rendered, actual)
			if tt.Err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.Err)
			}
		})
	}
}
