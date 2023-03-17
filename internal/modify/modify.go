package modify

import (
	"github.com/spf13/cast"
)

// spell: words ohler55

// ModifierFunc is a callback for the modifying JSON Path expressions.
// See [github.com/ohler55/ojg/jp/Expr.Modify].
type ModifierFunc func(element any) (altered any, changed bool)

// Modifier returns a modifier function for the given action and arg.
func Modifier(action Action, arg any) ModifierFunc {
	return func(element any) (any, bool) {
		var altered any
		var changed bool

		switch v := element.(type) {
		case bool:
			altered = Bool(v, action, cast.ToBool(arg))
			changed = true
		case float64:
			altered = Float64(v, action, cast.ToFloat64(arg))
			changed = true
		case int64:
			altered = Int64(v, action, cast.ToInt64(arg))
			changed = true
		case map[string]any:
			altered = Map(v, action, cast.ToStringMap(arg))
			changed = true
		case []any:
			altered = Slice(v, action, arg)
			changed = true
		case string:
			altered = String(v, action, cast.ToString(arg))
			changed = true
		default:
			altered = element
			changed = false
		}

		return altered, changed
	}
}
