package stamp

import (
	"errors"
	"strings"

	"github.com/spf13/cast"
	"github.com/swaggest/jsonschema-go"
	"github.com/twelvelabs/termite/render"
)

var (
	ErrPathNotFound = errors.New("path not found")
)

type Common struct {
	IfTpl   render.Template `mapstructure:"if" default:"true"`
	EachTpl render.Template `mapstructure:"each"`
}

// PrepareJSONSchema implements the jsonschema.Preparer interface.
func (c Common) PrepareJSONSchema(schema *jsonschema.Schema) error {
	if prop, ok := schema.Properties["if"]; ok {
		prop.TypeObjectEns().
			WithTitle("If").
			WithDescription(
				"Determines whether the task should be executed. "+
					"The value must be [coercible](https://pkg.go.dev/strconv#ParseBool) "+
					"to a boolean.",
			).
			WithExamples(
				"true",
				"{{ .SomeBool }}",
			)
	}
	if prop, ok := schema.Properties["each"]; ok {
		prop.TypeObjectEns().
			WithTitle("Each").
			WithDescription(
				"Set to a comma separated value and the task will be executued once per-item. "+
					"On each iteration, the `_Item` and `_Index` values will be set accordingly.",
			).
			WithExamples(
				"foo, bar, baz",
				"{{ .SomeList | join \",\" }}",
			)
	}
	return nil
}

func (c *Common) Iterator(values map[string]any) []any {
	rendered, _ := c.EachTpl.Render(values)
	if rendered == "" {
		return nil
	}

	trimmed := []any{}
	for _, item := range strings.Split(rendered, ",") {
		trimmed = append(trimmed, strings.TrimSpace(item))
	}

	return trimmed
}

func (c *Common) ShouldExecute(values map[string]any) bool {
	rendered, _ := c.IfTpl.Render(values)
	return cast.ToBool(rendered)
}
