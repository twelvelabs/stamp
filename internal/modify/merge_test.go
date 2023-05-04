package modify

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMerge(t *testing.T) {
	tests := []struct {
		desc     string
		dst      any
		src      any
		expected any
		conf     ModifierConf
	}{
		// scalar -> scalar: src always replaces dst
		{dst: "old", src: "new", expected: "new"},
		{dst: 111, src: 222, expected: 222},
		{dst: true, src: false, expected: false},

		// Mixed types: src always replaces dst
		{dst: 123, src: "string", expected: "string"},
		{dst: []any{1, 2, 3}, src: true, expected: true},
		{dst: map[string]any{"foo": "bar"}, src: 123, expected: 123},

		{
			desc:     "slices: concatenate by default",
			dst:      []any{1, 2, 3},
			src:      []any{3, 4, 5},
			expected: []any{1, 2, 3, 3, 4, 5},
		},
		{
			desc:     "slices: can be configured to replace",
			dst:      []any{1, 2, 3},
			src:      []any{3, 4, 5},
			expected: []any{3, 4, 5},
			conf: ModifierConf{
				MergeType: MergeTypeReplace,
			},
		},
		{
			desc:     "slices: can be configured to upsert",
			dst:      []any{1, 2, 3},
			src:      []any{3, 4, 5},
			expected: []any{1, 2, 3, 4, 5},
			conf: ModifierConf{
				MergeType: MergeTypeUpsert,
			},
		},
		{
			desc:     "slices: upsert should not remove existing dupes from dst",
			dst:      []any{1, 1, 1},
			src:      []any{3, 4, 5},
			expected: []any{1, 1, 1, 3, 4, 5},
			conf: ModifierConf{
				MergeType: MergeTypeUpsert,
			},
		},

		{
			desc: "maps: when mixed value types",
			dst: map[string]any{
				"111": "foo",
				"222": []any{1, 2, 3},
				"333": map[string]any{
					"333.111": true,
				},
				"444": "untouched",
			},
			src: map[string]any{
				"111": 123,
				"222": "replacement",
				"333": []any{"foo", "bar"},
				"555": "new",
			},
			expected: map[string]any{
				"111": 123,
				"222": "replacement",
				"333": []any{"foo", "bar"},
				"444": "untouched",
				"555": "new",
			},
		},
		{
			desc: "maps: when common value types",
			dst: map[string]any{
				"111": "original",
				"222": []any{1, 2, 3},
				"333": map[string]any{
					"333.111": true,
					"333.222": "original",
				},
				"444": "untouched",
			},
			src: map[string]any{
				"111": "replaced",
				"222": []any{3, 4, 3, 5},
				"333": map[string]any{
					"333.222": "replaced",
					"333.333": "new",
				},
				"555": "new",
			},
			expected: map[string]any{
				"111": "replaced",
				"222": []any{1, 2, 3, 4, 5},
				"333": map[string]any{
					"333.111": true,
					"333.222": "replaced",
					"333.333": "new",
				},
				"444": "untouched",
				"555": "new",
			},
			conf: ModifierConf{
				MergeType: MergeTypeUpsert,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert.Equal(t, tt.expected, Merge(tt.dst, tt.src, tt.conf))
		})
	}
}
