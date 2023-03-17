package modify

func Slice(subject []any, action Action, arg any) []any {
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
		modified = append(modified, argSlice...)
		modified = append(modified, subject...)
	case ActionAppend:
		modified = append(modified, subject...)
		modified = append(modified, argSlice...)
	case ActionReplace:
		modified = argSlice
	case ActionDelete:
	}

	return modified
}
