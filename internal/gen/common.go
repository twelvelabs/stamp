package gen

import (
	"strings"

	"github.com/spf13/cast"

	"github.com/twelvelabs/stamp/internal/render"
)

type Common struct {
	If   string `default:"true"`
	Each string

	DryRun bool
}

func (c *Common) Iterator(values map[string]any) []any {
	if c.Each == "" {
		return nil
	}

	rendered := c.Render(c.Each, values)
	trimmed := []any{}
	for _, item := range strings.Split(rendered, ",") {
		trimmed = append(trimmed, strings.TrimSpace(item))
	}

	return trimmed
}

func (c *Common) Render(tpl string, values map[string]any) string {
	rendered, err := render.String(tpl, values)
	if err != nil {
		return tpl
	}
	return rendered
}

func (c *Common) ShouldExecute(values map[string]any) bool {
	rendered := c.Render(c.If, values)
	return cast.ToBool(rendered)
}
