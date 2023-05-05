package modify

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModifier(t *testing.T) {
	tests := []struct {
		name      string
		dst       any
		action    Action
		src       any
		mergeType MergeType
		expected  any
	}{
		{
			dst:      true,
			action:   "prepend",
			src:      true,
			expected: true,
		},
		{
			dst:      true,
			action:   "append",
			src:      false,
			expected: false,
		},
		{
			dst:      false,
			action:   "replace",
			src:      true,
			expected: true,
		},
		{
			dst:      true,
			action:   "delete",
			src:      true,
			expected: false,
		},

		{
			dst:      111.0,
			action:   "prepend",
			src:      222.0,
			expected: 333.0,
		},
		{
			dst:      111.0,
			action:   "append",
			src:      222.0,
			expected: 333.0,
		},
		{
			dst:      111.0,
			action:   "replace",
			src:      222.0,
			expected: 222.0,
		},
		{
			dst:      111.0,
			action:   "delete",
			src:      222.0,
			expected: 0.0,
		},

		{
			dst:      int64(111),
			action:   "prepend",
			src:      int64(222),
			expected: int64(333),
		},
		{
			dst:      int64(111),
			action:   "append",
			src:      int64(222),
			expected: int64(333),
		},
		{
			dst:      int64(111),
			action:   "replace",
			src:      int64(222),
			expected: int64(222),
		},
		{
			dst:      int64(111),
			action:   "delete",
			src:      int64(222),
			expected: int64(0),
		},

		{
			dst: map[string]any{
				"aaa": 111,
				"bbb": 111,
				"ddd": []any{1, 1},
			},
			action: "prepend",
			src: map[string]any{
				"bbb": 222,
				"ccc": 222,
				"ddd": []any{2, 2},
			},
			expected: map[string]any{
				"aaa": 111,
				"bbb": 111,
				"ccc": 222,
				"ddd": []any{2, 2, 1, 1},
			},
		},
		{
			dst: map[string]any{
				"aaa": 111,
				"bbb": 111,
				"ddd": []any{1, 1},
			},
			action: "append",
			src: map[string]any{
				"bbb": 222,
				"ccc": 222,
				"ddd": []any{2, 2},
			},
			expected: map[string]any{
				"aaa": 111,
				"bbb": 222,
				"ccc": 222,
				"ddd": []any{1, 1, 2, 2},
			},
		},
		{
			dst: map[string]any{
				"aaa": 111,
				"bbb": 111,
			},
			action: "replace",
			src: map[string]any{
				"bbb": 222,
				"ccc": 222,
			},
			expected: map[string]any{
				"bbb": 222,
				"ccc": 222,
			},
		},
		{
			dst: map[string]any{
				"aaa": 111,
				"bbb": 111,
			},
			action: "delete",
			src: map[string]any{
				"bbb": 222,
				"ccc": 222,
			},
			expected: map[string]any(nil),
		},

		{
			dst:      []any{1, 1, 1},
			action:   "prepend",
			src:      []any{2, 2, 2},
			expected: []any{2, 2, 2, 1, 1, 1},
		},
		{
			dst:      []any{1, 1, 1},
			action:   "append",
			src:      []any{2, 2, 2},
			expected: []any{1, 1, 1, 2, 2, 2},
		},
		{
			dst:      []any{1, 1, 1},
			action:   "replace",
			src:      []any{2, 2, 2},
			expected: []any{2, 2, 2},
		},
		{
			dst:      []any{1, 1, 1},
			action:   "delete",
			src:      []any{2, 2, 2},
			expected: []any(nil),
		},

		{
			dst:      []byte("111"),
			action:   "prepend",
			src:      []byte("222"),
			expected: []byte("222111"),
		},
		{
			dst:      []byte("111"),
			action:   "append",
			src:      []byte("222"),
			expected: []byte("111222"),
		},
		{
			dst:      []byte("111"),
			action:   "replace",
			src:      []byte("222"),
			expected: []byte("222"),
		},
		{
			dst:      []byte("111"),
			action:   "delete",
			src:      []byte("222"),
			expected: []byte(nil),
		},

		{
			dst:      "111",
			action:   "prepend",
			src:      "222",
			expected: "222111",
		},
		{
			dst:      "111",
			action:   "append",
			src:      "222",
			expected: "111222",
		},
		{
			dst:      "111",
			action:   "replace",
			src:      "222",
			expected: "222",
		},
		{
			dst:      "111",
			action:   "delete",
			src:      "222",
			expected: "",
		},
	}
	for _, tt := range tests {
		if tt.name == "" {
			tt.name = fmt.Sprintf("%v/%v/%v", tt.dst, tt.action, tt.src)
		}
		t.Run(tt.name, func(t *testing.T) {
			modify := Modifier(tt.action, tt.src, WithMergeType(tt.mergeType))
			actual, changed := modify(tt.dst)
			assert.Equal(t, tt.expected, actual)
			assert.Equal(t, true, changed)
		})
	}
}

func TestModifier_ByteSliceConversion(t *testing.T) {
	// Merging into a byte slice should work if the arg is convertible...
	modify := Modifier(ActionAppend, "bar")
	altered, changed := modify([]byte("foo"))
	assert.Equal(t, []byte("foobar"), altered)
	assert.Equal(t, true, changed)

	// ... and noop if it is not.
	modify = Modifier(ActionAppend, 123)
	altered, changed = modify([]byte("foo"))
	assert.Equal(t, []byte("foo"), altered)
	assert.Equal(t, false, changed)
}

func TestModifier_AliasTypes(t *testing.T) {
	modify := Modifier(ActionAppend, int(222))
	altered, changed := modify(int(111))
	assert.Equal(t, int64(333), altered)
	assert.Equal(t, true, changed)

	modify = Modifier(ActionAppend, int32(222))
	altered, changed = modify(int32(111))
	assert.Equal(t, int64(333), altered)
	assert.Equal(t, true, changed)

	modify = Modifier(ActionAppend, float32(222.0))
	altered, changed = modify(float32(111.0))
	assert.Equal(t, float64(333.0), altered)
	assert.Equal(t, true, changed)
}

func TestModifier_UnsupportedTypes(t *testing.T) {
	subject := struct{}{}
	arg := struct{}{}

	modifier := Modifier(ActionReplace, arg)
	altered, changed := modifier(subject)
	assert.Equal(t, subject, altered)
	assert.Equal(t, false, changed)
}
