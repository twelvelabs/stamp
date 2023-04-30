package modify

func Bool(subject bool, action Action, arg bool, _ ModifierConf) bool {
	var modified bool
	switch action {
	case ActionPrepend:
		modified = arg && subject
	case ActionAppend:
		modified = subject && arg
	case ActionReplace:
		modified = arg
	case ActionDelete:
	}
	return modified
}
