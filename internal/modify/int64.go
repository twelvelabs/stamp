package modify

func Int64(dst int64, action Action, src int64, _ ModifierConf) int64 {
	var result int64

	switch action {
	case ActionPrepend:
		result = src + dst
	case ActionAppend:
		result = dst + src
	case ActionReplace:
		result = src
	case ActionDelete:
	}

	return result
}
