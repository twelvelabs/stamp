package stamp

import (
	"errors"
	"strings"

	"github.com/spf13/cast"
	"github.com/twelvelabs/termite/render"
)

var (
	ErrPathNotFound = errors.New("path not found")
)

type Common struct {
	IfTpl   render.Template `mapstructure:"if" title:"If" default:"true" examples:"[\"true\", \"{{ .SomeBool }}\"]" description:"Determines whether the task should be executed. The value must be [coercible](https://pkg.go.dev/strconv#ParseBool) to a boolean."`                                   //nolint:lll
	EachTpl render.Template `mapstructure:"each" title:"Each" examples:"[\"foo, bar, baz\", \"{{ .SomeList | join \\\",\\\" }}\"]" description:"Set to a comma separated value and the task will be executued once per-item. On each iteration, the _Item and _Index values will be set accordingly."` //nolint:lll
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
