package stamp

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cast"
	"github.com/twelvelabs/termite/render"
)

var (
	ErrPathNotFound = errors.New("path not found")
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

func (c *Common) RenderRequired(key, tpl string, values map[string]any) (string, error) {
	rendered := c.Render(tpl, values)
	if rendered == "" {
		return "", fmt.Errorf("%s: '%s' evaluated to an empty string", key, tpl)
	}
	return rendered, nil
}

func (c *Common) RenderMode(tpl string, values map[string]any) (os.FileMode, error) {
	rendered := c.Render(tpl, values)
	parsed, err := strconv.ParseUint(rendered, 8, 32)
	if err != nil {
		return 0, err
	}
	return os.FileMode(parsed), nil
}

func (c *Common) ShouldExecute(values map[string]any) bool {
	rendered := c.Render(c.If, values)
	return cast.ToBool(rendered)
}
