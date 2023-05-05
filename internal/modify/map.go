package modify

func Map(dst map[string]any, action Action, src map[string]any, conf ModifierConf) map[string]any {
	result := map[string]any{}

	switch action {
	case ActionPrepend:
		result = PrependMap(dst, src, conf)
	case ActionAppend:
		result = AppendMap(dst, src, conf)
	case ActionReplace:
		for k, v := range src {
			result[k] = v
		}
	case ActionDelete:
		result = nil
	}

	return result
}

func PrependMap(dst, src map[string]any, conf ModifierConf) map[string]any {
	for k, v := range dst {
		src[k] = appendMapValue(src[k], v, conf)
	}
	return src
}

func AppendMap(dst, src map[string]any, conf ModifierConf) map[string]any {
	for k, v := range src {
		dst[k] = appendMapValue(dst[k], v, conf)
	}
	return dst
}

// merges src into dst and returns the result.
func appendMapValue(dst, src any, conf ModifierConf) any {
	var result any

	switch dstCasted := dst.(type) {
	case map[string]any:
		if srcCasted, ok := src.(map[string]any); ok {
			// dst and src are both maps - append
			result = AppendMap(dstCasted, srcCasted, conf)
		} else {
			// trying to merge a slice or scalar into a map - replace.
			result = src
		}
	case []any:
		if srcCasted, ok := src.([]any); ok {
			// dst and src are both slices - append
			result = AppendSlice(dstCasted, srcCasted, conf)
		} else {
			// trying to merge a map or scalar into a slice - replace.
			result = src
		}
	default:
		// merging into a scalar always replaces.
		result = src
	}

	return result
}
