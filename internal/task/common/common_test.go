package common

import (
	"testing"

	//cspell: disable
	"github.com/creasty/defaults"
	"github.com/stretchr/testify/assert"
	"github.com/twelvelabs/stamp/internal/iostreams"
	//cspell: enable
)

func TestCommon_Defaults(t *testing.T) {
	// empty struct should have defaults set
	task := &Common{}
	defaults.Set(task)
	assert.Equal(t, "true", task.If)
	assert.Equal(t, "", task.Each)

	// existing values should not be overridden by defaults
	task = &Common{
		If:   "{{ .SomeBool }}",
		Each: "{{ .SomeList }}",
	}
	defaults.Set(task)
	assert.Equal(t, "{{ .SomeBool }}", task.If)
	assert.Equal(t, "{{ .SomeList }}", task.Each)
}

func TestCommon_Iterator(t *testing.T) {
	tests := []struct {
		Name   string
		Each   string
		Values map[string]any
		Output []any
	}{
		{
			Name:   "should return nil when Each is an empty string",
			Each:   "",
			Values: map[string]any{},
			Output: nil,
		},
		{
			Name:   "should return a slice when Each is a comma separated string",
			Each:   "foo, bar, baz",
			Values: map[string]any{},
			Output: []any{"foo", "bar", "baz"},
		},
		{
			Name: "should render Each as a template value before processing",
			Each: "{{ .Tags }}",
			Values: map[string]any{
				"Tags": "foo, bar, baz",
			},
			Output: []any{"foo", "bar", "baz"},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			task := &Common{
				Each: test.Each,
			}
			assert.Equal(t, test.Output, task.Iterator(test.Values))
		})
	}
}

func TestCommon_Render(t *testing.T) {
	tests := []struct {
		Name     string
		Template string
		Values   map[string]any
		Output   string
	}{
		{
			Name:     "treats empty string as a noop",
			Template: "",
			Values:   map[string]any{},
			Output:   "",
		},
		{
			Name:     "returns the template unchanged if unable to parse",
			Template: "{{}",
			Values:   map[string]any{},
			Output:   "{{}",
		},
		{
			Name:     "semi-gracefully handles missing values",
			Template: "Hello, {{ .Name }}.",
			Values:   map[string]any{},
			Output:   "Hello, <no value>.",
		},
		{
			Name:     "renders values when present",
			Template: "Hello, {{ .Name }}.",
			Values:   map[string]any{"Name": "World"},
			Output:   "Hello, World.",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			task := &Common{}
			assert.Equal(t, test.Output, task.Render(test.Template, test.Values))
		})
	}
}

func TestCommon_Report(t *testing.T) {
	tests := []struct {
		Name    string
		Action  string
		Message string
		DryRun  bool
		Output  string
	}{
		{
			Name:    "writes a formatted log line",
			Action:  "create",
			Message: "/something",
			DryRun:  false,
			Output:  "✓ [    create]: /something\n",
		},
		{
			Name:    "adds a DRY RUN prefix when configured",
			Action:  "create",
			Message: "/something",
			DryRun:  true,
			Output:  "✓ [DRY RUN][    create]: /something\n",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			ios, _, _, stderr := iostreams.Test()

			task := &Common{
				DryRun: test.DryRun,
			}
			task.LogSuccess(ios, test.Action, test.Message)
			assert.Equal(t, test.Output, stderr.String())
		})
	}
}

func TestCommon_DryRunMethods(t *testing.T) {
	task := &Common{}
	assert.Equal(t, false, task.IsDryRun())
	task.SetDryRun(true)
	assert.Equal(t, true, task.IsDryRun())

	task = &Common{DryRun: true}
	assert.Equal(t, true, task.IsDryRun())
}

func TestCommon_ShouldExecute(t *testing.T) {
	tests := []struct {
		Name   string
		Values map[string]any
		If     string
		Output bool
	}{
		{
			Name:   "returns false when empty string",
			Values: map[string]any{},
			If:     "",
			Output: false,
		},
		{
			Name:   "returns true when literal string true",
			Values: map[string]any{},
			If:     "true",
			Output: true,
		},
		{
			Name:   "returns false when literal string false",
			Values: map[string]any{},
			If:     "false",
			Output: false,
		},
		{
			Name: "returns template value if present",
			Values: map[string]any{
				"SomeBool": true,
			},
			If:     "{{ .SomeBool }}",
			Output: true,
		},
		{
			Name:   "returns false if template value missing",
			Values: map[string]any{},
			If:     "{{ .SomeBool }}",
			Output: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			task := &Common{
				If: test.If,
			}
			assert.Equal(t, test.Output, task.ShouldExecute(test.Values))
		})
	}
}
