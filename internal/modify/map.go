package modify

func Map(subject map[string]any, action Action, arg map[string]any, conf ModifierConf) map[string]any {
	modified := map[string]any{}

	switch action {
	case ActionPrepend:
		modified = MergeMap(arg, subject, conf)
	case ActionAppend:
		modified = MergeMap(subject, arg, conf)
	case ActionReplace:
		for k, v := range arg {
			modified[k] = v
		}
	case ActionDelete:
		modified = nil
	}

	return modified
}
