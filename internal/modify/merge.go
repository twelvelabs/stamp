package modify

// Merge merges src into dst and returns the result.
func Merge(dst, src any, conf ModifierConf) any {
	var result any

	switch dstCasted := dst.(type) {
	case map[string]any:
		if srcCasted, ok := src.(map[string]any); ok {
			// dst and src are both maps - merge
			result = MergeMap(dstCasted, srcCasted, conf)
		} else {
			// trying to merge a slice or scalar into a map - replace.
			result = src
		}
	case []any:
		if srcCasted, ok := src.([]any); ok {
			// dst and src are both slices - merge
			result = MergeSlice(dstCasted, srcCasted, conf)
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

// MergeMap recursively merges src into dst and returns the result.
func MergeMap(dst, src map[string]any, conf ModifierConf) map[string]any {
	for k, v := range src {
		dst[k] = Merge(dst[k], v, conf)
	}
	return dst
}

// MergeSlice merges src into dst and returns the result.
// Merge logic depends on conf.SliceMerge, which can be one of:
//
//   - SliceMergeConcat: concatenate the two slices (default).
//   - SliceMergeUpsert: insert source items only if they are not already present in dst.
//   - SliceMergeReplace: fully replace src with dst.
func MergeSlice(dst, src []any, conf ModifierConf) []any {
	// Slice merge behavior is configurable.
	// Defaults to `SliceMergeConcat`.
	switch conf.SliceMerge {
	case SliceMergeReplace:
		return src
	case SliceMergeUpsert:
		// Note: not just doing a set intersection because we don't want to remove
		// pre-existing dupes from the dst slice - just prevent new ones.
		set := NewSet(dst...)
		for _, item := range src {
			if !set.Contains(item) {
				dst = append(dst, item)
			}
		}
		return dst
	default:
		return append(dst, src...)
	}
}