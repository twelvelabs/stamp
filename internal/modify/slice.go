package modify

import "fmt"

func Slice(subject []any, action Action, arg any, conf ModifierConf) []any {
	// ensure arg is a slice
	var argSlice []any
	if a, ok := arg.([]any); ok {
		argSlice = a
	} else {
		argSlice = append(argSlice, arg)
	}

	// Simple lookup map for subject slice content.
	lookup := map[any]struct{}{}
	for _, item := range subject {
		// Maps and nested slices are not allowed as map keys.
		// Go-syntax representation probably isn't bulletproof,
		// but should be good enough for the relatively simple
		// data loaded from `generator.yaml` files.
		key := fmt.Sprintf("%#v", item)
		lookup[key] = struct{}{}
	}

	// Helper to decide whether to perform append/prepend based on upsert config.
	shouldPerformAction := func(item any) bool {
		key := fmt.Sprintf("%#v", item)
		_, exists := lookup[key]
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
