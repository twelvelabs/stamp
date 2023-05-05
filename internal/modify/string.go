package modify

import "strings"

func String(dst string, action Action, src string, conf ModifierConf) string {
	var result string

	switch action {
	case ActionPrepend:
		result = PrependString(dst, src, conf)
	case ActionAppend:
		result = AppendString(dst, src, conf)
	case ActionReplace:
		result = src
	case ActionDelete:
	}

	return result
}

func AppendString(dst, src string, conf ModifierConf) string {
	switch conf.MergeType {
	case MergeTypeReplace:
		return src
	case MergeTypeUpsert:
		if strings.HasSuffix(dst, src) {
			return dst
		}
		return dst + src
	default: // case MergeTypeConcat:
		return dst + src
	}
}

func PrependString(dst, src string, conf ModifierConf) string {
	switch conf.MergeType {
	case MergeTypeReplace:
		return src
	case MergeTypeUpsert:
		if strings.HasPrefix(dst, src) {
			return dst
		}
		return src + dst
	default: // case MergeTypeConcat:
		return src + dst
	}
}
