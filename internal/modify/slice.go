package modify

func Slice(dst []any, action Action, src any, conf ModifierConf) []any {
	var result []any

	// The source can be either a single element, or a slice of elements.
	// Lazily wrap single elements in a slice so we can use that same
	// logic for both.
	var srcSlice []any
	if a, ok := src.([]any); ok {
		srcSlice = a
	} else {
		srcSlice = append(srcSlice, src)
	}

	switch action {
	case ActionPrepend:
		result = PrependSlice(dst, srcSlice, conf)
	case ActionAppend:
		result = AppendSlice(dst, srcSlice, conf)
	case ActionReplace:
		result = srcSlice
	case ActionDelete:
	}

	return result
}

func PrependSlice(dst, src []any, conf ModifierConf) []any {
	switch conf.MergeType {
	case MergeTypeReplace:
		return src
	case MergeTypeUpsert:
		dstSet := NewSet(dst...)
		head := []any{}
		for _, item := range src {
			if !dstSet.Contains(item) {
				head = append(head, item)
			}
		}
		return append(head, dst...)
	default: // case MergeTypeConcat:
		return append(src, dst...)
	}
}

func AppendSlice(dst, src []any, conf ModifierConf) []any {
	switch conf.MergeType {
	case MergeTypeReplace:
		return src
	case MergeTypeUpsert:
		set := NewSet(dst...)
		for _, item := range src {
			if !set.Contains(item) {
				dst = append(dst, item)
			}
		}
		return dst
	default: // case MergeTypeConcat:
		return append(dst, src...)
	}
}
