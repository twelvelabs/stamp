package modify

func Float64(dst float64, action Action, src float64, _ ModifierConf) float64 {
	var result float64

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
