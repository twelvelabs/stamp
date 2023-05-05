package modify

import "bytes"

func Bytes(dst []byte, action Action, src []byte, conf ModifierConf) []byte {
	var result []byte

	switch action {
	case ActionPrepend:
		result = PrependBytes(dst, src, conf)
	case ActionAppend:
		result = AppendBytes(dst, src, conf)
	case ActionReplace:
		result = append(result, src...)
	case ActionDelete:
	}

	return result
}

func PrependBytes(dst, src []byte, conf ModifierConf) []byte {
	switch conf.MergeType {
	case MergeTypeReplace:
		return src
	case MergeTypeUpsert:
		if bytes.HasPrefix(dst, src) {
			return dst
		}
		return append(src, dst...)
	default: // case MergeTypeConcat:
		return append(src, dst...)
	}
}

func AppendBytes(dst, src []byte, conf ModifierConf) []byte {
	switch conf.MergeType {
	case MergeTypeReplace:
		return src
	case MergeTypeUpsert:
		if bytes.HasSuffix(dst, src) {
			return dst
		}
		return append(dst, src...)
	default: // case MergeTypeConcat:
		return append(dst, src...)
	}
}
