package modify

func Map(subject map[string]any, action Action, arg map[string]any, _ ModifierConf) map[string]any {
	modified := map[string]any{}

	switch action {
	case ActionPrepend:
		for k, v := range arg {
			modified[k] = v
		}
		for k, v := range subject {
			modified[k] = v
		}
	case ActionAppend:
		for k, v := range subject {
			modified[k] = v
		}
		for k, v := range arg {
			modified[k] = v
		}
	case ActionReplace:
		for k, v := range arg {
			modified[k] = v
		}
	case ActionDelete:
		modified = nil
	}

	return modified
}
