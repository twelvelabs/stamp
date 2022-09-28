package render

import (
	"bytes"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/gobuffalo/flect"
)

// FuncMap returns a map of functions for the template engine.
func FuncMap() template.FuncMap {
	funcs := sprig.FuncMap()

	// See: https://pkg.go.dev/github.com/gobuffalo/flect
	funcs["camelize"] = flect.Camelize
	funcs["capitalize"] = flect.Capitalize
	funcs["dasherize"] = flect.Dasherize
	funcs["humanize"] = flect.Humanize
	funcs["ordinalize"] = flect.Ordinalize
	funcs["pascalize"] = flect.Pascalize
	funcs["pluralize"] = flect.Pluralize
	funcs["singularize"] = flect.Singularize
	funcs["titleize"] = flect.Titleize
	funcs["underscore"] = flect.Underscore

	return funcs
}

// RenderFile renders the named file using the data in values.
func RenderFile(path string, values map[string]any) (string, error) {
	name := filepath.Base(path)
	t, err := template.New(name).Funcs(FuncMap()).ParseFiles(path)
	if err != nil {
		return "", err
	}
	return execute(t, values)
}

// RenderString renders the template string using the data in values.
func RenderString(s string, values map[string]any) (string, error) {
	t, err := template.New("render-string").Funcs(FuncMap()).Parse(s)
	if err != nil {
		return "", err
	}
	return execute(t, values)
}

// executes template t with values and returns the rendered string.
func execute(t *template.Template, values map[string]any) (string, error) {
	buf := bytes.Buffer{}
	err := t.Execute(&buf, values)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
