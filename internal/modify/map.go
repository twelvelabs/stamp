package modify

import (
	"github.com/imdario/mergo"
)

func Map(subject map[string]any, action Action, arg map[string]any, _ ModifierConf) map[string]any {
	modified := map[string]any{}

	switch action {
	case ActionPrepend:
		for k, v := range arg {
			modified[k] = v
		}
		_ = mergo.Merge(&modified, subject, mergo.WithOverride, mergo.WithAppendSlice)
	case ActionAppend:
		for k, v := range subject {
			modified[k] = v
		}
		_ = mergo.Merge(&modified, arg, mergo.WithOverride, mergo.WithAppendSlice)
	case ActionReplace:
		for k, v := range arg {
			modified[k] = v
		}
	case ActionDelete:
		modified = nil
	}

	return modified
}
