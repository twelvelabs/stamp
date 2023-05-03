package modify

func Slice(subject []any, action Action, arg any, conf ModifierConf) []any {
	// ensure arg is a slice
	var argSlice []any
	if a, ok := arg.([]any); ok {
		argSlice = a
	} else {
		argSlice = append(argSlice, arg)
	}

	var modified []any
	switch action {
	case ActionPrepend:
		modified = MergeSlice(argSlice, subject, conf)
	case ActionAppend:
		modified = MergeSlice(subject, argSlice, conf)
	case ActionReplace:
		modified = argSlice
	case ActionDelete:
	}

	return modified
}
