package value

import (
	"fmt"
	"strings"

	//cspell:disable
	"github.com/gobuffalo/flect"
	"github.com/spf13/cast"
	"github.com/twelvelabs/stamp/internal/fsutil"
	//cspell:enable
)

var (
	transformers = map[string]Transformer{}
)

func init() {
	RegisterTransformer("trim", StringTransformer(strings.TrimSpace))
	RegisterTransformer("uppercase", StringTransformer(strings.ToUpper))
	RegisterTransformer("lowercase", StringTransformer(strings.ToLower))
	RegisterTransformer("dasherize", StringTransformer(flect.Dasherize))
	RegisterTransformer("pascalize", StringTransformer(flect.Pascalize))
	RegisterTransformer("underscore", StringTransformer(flect.Underscore))
	RegisterTransformer("expand-path", ExpandPath)
}

func Transform(key string, value any, rule string) (any, error) {
	if rule == "" {
		return value, nil
	}
	transformed := value

	ts, err := ParseTransformRule(key, rule)
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

// Transformer is a function used to process value data.
type Transformer func(any) (any, error)

func GetTransformer(name string) Transformer {
	t, _ := transformers[name]
	return t
}
func RegisterTransformer(name string, t Transformer) {
	transformers[name] = t
}

func ParseTransformRule(key string, rule string) ([]Transformer, error) {
	ts := []Transformer{}
	rules := strings.Split(strings.TrimSpace(rule), ",")
	for _, rule := range rules {
		rule = strings.TrimSpace(rule)
		t := GetTransformer(rule)
		if t == nil {
			return nil, fmt.Errorf("undefined transform [%s: %s]", key, rule)
		}
		ts = append(ts, t)
	}
	return ts, nil
}

// StringTransformer accepts a string function and returns a Transformer
// that delegates to it.
func StringTransformer(f func(s string) string) Transformer {
	return func(data any) (any, error) {
		return f(cast.ToString(data)), nil
	}
}

// ExpandPath ensures that data is an absolute path.
// Environment variables (and the ~ string) are expanded.
func ExpandPath(data any) (any, error) {
	return fsutil.NormalizePath(cast.ToString(data))
}
