package stamp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twelvelabs/termite/render"
)

func TestCommon_Iterator(t *testing.T) {
	tests := []struct {
		Name   string
		Each   render.Template
		Values map[string]any
		Output []any
	}{
		{
			Name:   "should return nil when Each is an empty string",
			Each:   *render.MustCompile(``),
			Values: map[string]any{},
			Output: nil,
		},
		{
			Name:   "should return a slice when Each is a comma separated string",
			Each:   *render.MustCompile(`foo, bar, baz`),
			Values: map[string]any{},
			Output: []any{"foo", "bar", "baz"},
		},
		{
			Name: "should render Each as a template value before processing",
			Each: *render.MustCompile(`{{ .Tags }}`),
			Values: map[string]any{
				"Tags": "foo, bar, baz",
			},
			Output: []any{"foo", "bar", "baz"},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			task := &Common{
				EachTpl: test.Each,
			}
			assert.Equal(t, test.Output, task.Iterator(test.Values))
		})
	}
}

func TestCommon_ShouldExecute(t *testing.T) {
	tests := []struct {
		Name   string
		Values map[string]any
		If     render.Template
		Output bool
	}{
		{
			Name:   "returns false when empty string",
			Values: map[string]any{},
			If:     *render.MustCompile(``),
			Output: false,
		},
		{
			Name:   "returns true when literal string true",
			Values: map[string]any{},
			If:     *render.MustCompile(`true`),
			Output: true,
		},
		{
			Name:   "returns false when literal string false",
			Values: map[string]any{},
			If:     *render.MustCompile(`false`),
			Output: false,
		},
		{
			Name: "returns template value if present",
			Values: map[string]any{
				"SomeBool": true,
			},
			If:     *render.MustCompile(`{{ .SomeBool }}`),
			Output: true,
		},
		{
			Name:   "returns false if template value missing",
			Values: map[string]any{},
			If:     *render.MustCompile(`{{ .SomeBool }}`),
			Output: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			task := &Common{
				IfTpl: test.If,
			}
			assert.Equal(t, test.Output, task.ShouldExecute(test.Values))
		})
	}
}
