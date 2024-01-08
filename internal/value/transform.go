package value

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/gobuffalo/flect"
	"github.com/spf13/cast"

	"github.com/twelvelabs/stamp/internal/fsutil"
)

var (
	ErrUnknownTransformer = errors.New("undefined transform")
	transformers          = map[string]Transformer{}
)

func init() {
	RegisterTransformer(Transformer{
		Name:        "trim",
		Description: "Removes all leading and trailing whitespace.",
		Func:        StringTransformerFunc(strings.TrimSpace),
	})
	RegisterTransformer(Transformer{
		Name:        "uppercase",
		Description: "Converts to `UPPERCASE`.",
		Func:        StringTransformerFunc(strings.ToUpper),
	})
	RegisterTransformer(Transformer{
		Name:        "lowercase",
		Description: "Converts to `lowercase`.",
		Func:        StringTransformerFunc(strings.ToLower),
	})
	RegisterTransformer(Transformer{
		Name:        "dasherize",
		Description: "Converts to `kebab-case`.",
		Func:        StringTransformerFunc(flect.Dasherize),
	})
	RegisterTransformer(Transformer{
		Name:        "pascalize",
		Description: "Converts to `PascalCase`.",
		Func:        StringTransformerFunc(flect.Pascalize),
	})
	RegisterTransformer(Transformer{
		Name:        "underscore",
		Description: "Converts to snake_case.",
		Func:        StringTransformerFunc(flect.Underscore),
	})
	RegisterTransformer(Transformer{
		Name:        "expand-path",
		Description: "Converts to abs file path; Expands env vars and tilde.",
		Func:        expandPath,
	})
}

func Transform(key string, value any, rule string) (any, error) {
	if rule == "" {
		return value, nil
	}
	transformed := value

	ts, err := parseTransformRule(key, rule)
	if err != nil {
		return nil, err
	}

	for _, t := range ts {
		transformed, err = t(transformed)
		if err != nil {
			return nil, err
		}
	}

	return transformed, nil
}

// TransformerFunc is a function used to process value data.
type TransformerFunc func(any) (any, error)

type Transformer struct {
	Name        string
	Description string
	Func        TransformerFunc
}

// GetTransformer returns the transformer registered for name.
// If name is not found, returns ErrUnknownTransformer.
func GetTransformer(name string) (Transformer, error) {
	if t, ok := transformers[name]; ok {
		return t, nil
	}
	return Transformer{}, ErrUnknownTransformer
}

func RegisterTransformer(t Transformer) {
	if _, ok := transformers[t.Name]; ok {
		panic("transformer already registered for name: " + t.Name)
	}
	transformers[t.Name] = t
}

func UnregisterTransformer(t Transformer) {
	delete(transformers, t.Name)
}

func RegisteredTransformers() []Transformer {
	ts := []Transformer{}
	for _, t := range transformers {
		ts = append(ts, t)
	}
	sort.Slice(ts, func(i, j int) bool {
		return ts[i].Name < ts[j].Name
	})
	return ts
}

func parseTransformRule(key string, rule string) ([]TransformerFunc, error) {
	tfs := []TransformerFunc{}
	rules := strings.Split(strings.TrimSpace(rule), ",")
	for _, rule := range rules {
		rule = strings.TrimSpace(rule)
		t, err := GetTransformer(rule)
		if err != nil {
			return nil, fmt.Errorf("%w [%s: %s]", err, key, rule)
		}
		tfs = append(tfs, t.Func)
	}
	return tfs, nil
}

// StringTransformerFunc accepts a string function and returns a TransformerFunc
// that delegates to it.
func StringTransformerFunc(f func(s string) string) TransformerFunc {
	return func(data any) (any, error) {
		return f(cast.ToString(data)), nil
	}
}

func expandPath(data any) (any, error) {
	return fsutil.NormalizePath(cast.ToString(data))
}
