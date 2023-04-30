package modify

func Float64(subject float64, action Action, arg float64, _ ModifierConf) float64 {
	var modified float64
	switch action {
	case ActionPrepend:
		modified = arg + subject
	case ActionAppend:
		modified = subject + arg
	case ActionReplace:
		modified = arg
	case ActionDelete:
	}
	return modified
}
