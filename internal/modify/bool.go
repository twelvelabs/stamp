package modify

func Bool(dst bool, action Action, src bool, _ ModifierConf) bool {
	var result bool

	switch action {
	case ActionPrepend:
		result = src && dst
	case ActionAppend:
		result = dst && src
	case ActionReplace:
		result = src
	case ActionDelete:
	}

	return result
}
