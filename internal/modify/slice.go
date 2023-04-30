package modify

func Slice(subject []any, action Action, arg any, conf ModifierConf) []any {
	// ensure arg is a slice
	var argSlice []any
	if a, ok := arg.([]any); ok {
		argSlice = a
	} else {
		argSlice = append(argSlice, arg)
	}

	// lookup map
	lookup := map[any]struct{}{}
	for _, item := range subject {
		lookup[item] = struct{}{}
	}

	shouldPerformAction := func(arg any) bool {
		_, exists := lookup[arg]
		return !conf.Upsert || (conf.Upsert && !exists)
	}

	var modified []any
	switch action {
	case ActionPrepend:
		for _, a := range argSlice {
			if shouldPerformAction(a) {
				modified = append(modified, a)
			}
		}
		modified = append(modified, subject...)
	case ActionAppend:
		modified = append(modified, subject...)
		for _, a := range argSlice {
			if shouldPerformAction(a) {
				modified = append(modified, a)
			}
		}
	case ActionReplace:
		modified = argSlice
	case ActionDelete:
	}

	return modified
}
