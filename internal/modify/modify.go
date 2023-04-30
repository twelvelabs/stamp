package modify

import (
	"github.com/spf13/cast"
)

// spell: words ohler55

// ModifierFunc is a callback for modifying JSON Path expressions.
// See [github.com/ohler55/ojg/jp/Expr.Modify].
type ModifierFunc func(element any) (altered any, changed bool)

// Modifier returns a modifier function for the given action and arg.
// When the modifier function is called with an element, it will perform
// the appropriate action (i.e. prepend, append, replace, or delete) with
// the arg and return the altered element and a bool representing
// whether the element has changed (unsupported types will be unchanged).
//
// For most types the arg is expected to match (or be cast-able to) the
// type of the element. The exception is when the element is `[]any`,
// in which case the arg can be either `[]any` or `any`.
//
// Example:
//
//	// Creates a modifier func that appends `3`.
//	modify := Modifier(ActionAppend, 3)
//
//	// Modify an int
//	modified, changed := modify(2)
//	print(modified) // => 5
//	print(changed) // => true
//
//	// Modify a slice
//	modified, changed = modify([]any{1,2})
//	print(modified) // => []any{1,2,3}
//	print(changed) // => true
//
//	// Attempt to modify an unsupported type
//	modified, changed = modify(someStruct)
//	print(modified) // => someStruct
//	print(changed) // => false
func Modifier(action Action, arg any, opts ...ModifierOpt) ModifierFunc {
	conf := ModifierConf{}
	for _, opt := range opts {
		conf = opt(conf)
	}

	return func(element any) (any, bool) {
		var altered any
		var changed bool

		switch v := element.(type) {
		case bool:
			altered = Bool(v, action, cast.ToBool(arg), conf)
			changed = true
		case float32:
			altered = Float64(float64(v), action, cast.ToFloat64(arg), conf)
			changed = true
		case float64:
			altered = Float64(v, action, cast.ToFloat64(arg), conf)
			changed = true
		case int:
			altered = Int64(int64(v), action, cast.ToInt64(arg), conf)
			changed = true
		case int32:
			altered = Int64(int64(v), action, cast.ToInt64(arg), conf)
			changed = true
		case int64:
			altered = Int64(v, action, cast.ToInt64(arg), conf)
			changed = true
		case map[string]any:
			altered = Map(v, action, cast.ToStringMap(arg), conf)
			changed = true
		case []any:
			altered = Slice(v, action, arg, conf)
			changed = true
		case string:
			altered = String(v, action, cast.ToString(arg), conf)
			changed = true
		default:
			altered = element
			changed = false
		}

		return altered, changed
	}
}

type ModifierConf struct {
	Upsert bool
}

type ModifierOpt func(conf ModifierConf) ModifierConf

// WithUpsert returns a ModifierOpt that configures the upsert option.
// When upsert is enabled, Slice appends and prepends only occur
// if the value is not already present.
func WithUpsert(upsert bool) ModifierOpt {
	return func(c ModifierConf) ModifierConf {
		c.Upsert = upsert
		return c
	}
}
