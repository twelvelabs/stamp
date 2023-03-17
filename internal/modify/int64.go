package modify

func Int64(subject int64, action Action, arg int64) int64 {
	var modified int64
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
