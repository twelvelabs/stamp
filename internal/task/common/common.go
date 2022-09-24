package common

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/huandu/xstrings"
	"github.com/spf13/cast"
	"github.com/twelvelabs/stamp/internal/iostreams"
)

const (
	ACTION_WIDTH = 10
)

type Status = int

const (
	StatusUnknown Status = iota
	StatusSuccess Status = iota
	StatusWarning Status = iota
	StatusFailure Status = iota
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
	parsed, err := template.New("render").Funcs(sprig.FuncMap()).Parse(tpl)
	if err != nil {
		return tpl
	}
	buf := &bytes.Buffer{}
	err = parsed.Execute(buf, values)
	if err != nil {
		return tpl
	}
	return buf.String()
}

func (c *Common) LogSuccess(ios *iostreams.IOStreams, action string, msg string) {
	c.LogStatus(ios, StatusSuccess, action, msg)
}

func (c *Common) LogWarning(ios *iostreams.IOStreams, action string, msg string) {
	c.LogStatus(ios, StatusWarning, action, msg)
}

func (c *Common) LogFailure(ios *iostreams.IOStreams, action string, msg string) {
	c.LogStatus(ios, StatusFailure, action, msg)
}

func (c *Common) LogStatus(ios *iostreams.IOStreams, status Status, action string, msg string) {
	cs := ios.Formatter()

	var icon string
	var color string
	switch status {
	case StatusSuccess:
		icon = cs.SuccessIcon()
		color = "green"
	case StatusWarning:
		icon = cs.WarningIcon()
		color = "yellow"
	case StatusFailure:
		icon = cs.FailureIcon()
		color = "red"
	}

	prefix := icon + " "
	if c.DryRun {
		prefix = prefix + "[DRY RUN]"
	}
	// need to justify _before_ adding color codes
	action = xstrings.RightJustify(action, ACTION_WIDTH, " ")
	action = cs.ColorFromString(color)(action)
	fmt.Fprintf(ios.Err, "%s[%s]: %s\n", prefix, action, msg)
}

func (c *Common) ShouldExecute(values map[string]any) bool {
	rendered := c.Render(c.If, values)
	return cast.ToBool(rendered)
}
