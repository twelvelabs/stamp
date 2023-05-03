package modify

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModify_Bool(t *testing.T) {
	tests := []struct {
		subject  bool
		action   Action
		arg      bool
		expected bool
	}{
		{
			subject:  true,
			action:   "prepend",
			arg:      true,
			expected: true,
		},
		{
			subject:  true,
			action:   "append",
			arg:      false,
			expected: false,
		},
		{
			subject:  false,
			action:   "replace",
			arg:      true,
			expected: true,
		},
		{
			subject:  true,
			action:   "delete",
			arg:      true,
			expected: false,
		},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("%v/%v/%v", tt.subject, tt.action, tt.arg)
		t.Run(name, func(t *testing.T) {
			modifier := Modifier(tt.action, tt.arg)
			altered, changed := modifier(tt.subject)
			assert.Equal(t, tt.expected, altered)
			assert.Equal(t, true, changed)
		})
	}
}

func TestModify_Float64(t *testing.T) {
	tests := []struct {
		subject  float64
		action   Action
		arg      float64
		expected float64
	}{
		{
			subject:  111.0,
			action:   "prepend",
			arg:      222.0,
			expected: 333.0,
		},
		{
			subject:  111.0,
			action:   "append",
			arg:      222.0,
			expected: 333.0,
		},
		{
			subject:  111.0,
			action:   "replace",
			arg:      222.0,
			expected: 222.0,
		},
		{
			subject:  111.0,
			action:   "delete",
			arg:      222.0,
			expected: 0.0,
		},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("%v/%v/%v", tt.subject, tt.action, tt.arg)
		t.Run(name, func(t *testing.T) {
			modifier := Modifier(tt.action, tt.arg)
			altered, changed := modifier(tt.subject)
			assert.Equal(t, tt.expected, altered)
			assert.Equal(t, true, changed)
		})
	}
}

func TestModify_Int64(t *testing.T) {
	tests := []struct {
		subject  int64
		action   Action
		arg      int64
		expected int64
	}{
		{
			subject:  111,
			action:   "prepend",
			arg:      222,
			expected: 333,
		},
		{
			subject:  111,
			action:   "append",
			arg:      222,
			expected: 333,
		},
		{
			subject:  111,
			action:   "replace",
			arg:      222,
			expected: 222,
		},
		{
			subject:  111,
			action:   "delete",
			arg:      222,
			expected: 0,
		},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("%v/%v/%v", tt.subject, tt.action, tt.arg)
		t.Run(name, func(t *testing.T) {
			modifier := Modifier(tt.action, tt.arg)
			altered, changed := modifier(tt.subject)
			assert.Equal(t, tt.expected, altered)
			assert.Equal(t, true, changed)
		})
	}
}

func TestModify_Map(t *testing.T) {
	tests := []struct {
		subject  map[string]any
		action   Action
		arg      map[string]any
		expected map[string]any
	}{
		{
			subject: map[string]any{
				"aaa": 111,
				"bbb": 111,
				"ddd": []any{1, 1},
			},
			action: "prepend",
			arg: map[string]any{
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
			subject: map[string]any{
				"aaa": 111,
				"bbb": 111,
				"ddd": []any{1, 1},
			},
			action: "append",
			arg: map[string]any{
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
			subject: map[string]any{
				"aaa": 111,
				"bbb": 111,
			},
			action: "replace",
			arg: map[string]any{
				"bbb": 222,
				"ccc": 222,
			},
			expected: map[string]any{
				"bbb": 222,
				"ccc": 222,
			},
		},
		{
			subject: map[string]any{
				"aaa": 111,
				"bbb": 111,
			},
			action:   "delete",
			expected: nil,
		},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("%v/%v/%v", tt.subject, tt.action, tt.arg)
		t.Run(name, func(t *testing.T) {
			modifier := Modifier(tt.action, tt.arg)
			altered, changed := modifier(tt.subject)
			assert.Equal(t, tt.expected, altered)
			assert.Equal(t, true, changed)
		})
	}
}

func TestModify_Slice(t *testing.T) {
	tests := []struct {
		subject  []any
		action   Action
		arg      any
		upsert   bool
		expected []any
	}{
		{
			subject:  []any{1, 1, 1},
			action:   "prepend",
			arg:      []any{2, 2, 2},
			expected: []any{2, 2, 2, 1, 1, 1},
		},
		{
			subject:  []any{1, 1, 1},
			action:   "prepend",
			arg:      222,
			expected: []any{222, 1, 1, 1},
		},
		{
			subject:  []any{1, 1, 1},
			action:   "append",
			arg:      []any{2, 2, 2},
			expected: []any{1, 1, 1, 2, 2, 2},
		},
		{
			subject:  []any{1, 1, 1},
			action:   "append",
			arg:      222,
			expected: []any{1, 1, 1, 222},
		},
		{
			subject:  []any{1, 2, 3},
			action:   "append",
			arg:      2,
			expected: []any{1, 2, 3, 2},
		},
		{
			subject:  []any{1, 2, 3},
			action:   "append",
			arg:      2,
			upsert:   true,
			expected: []any{1, 2, 3},
		},
		{
			subject:  []any{[]any{1}, []any{2}, []any{3}},
			action:   "append",
			arg:      []any{[]any{2}},
			upsert:   true,
			expected: []any{[]any{1}, []any{2}, []any{3}},
		},
		{
			subject:  []any{1, 1, 1},
			action:   "replace",
			arg:      []any{2, 2, 2},
			expected: []any{2, 2, 2},
		},
		{
			subject:  []any{1, 1, 1},
			action:   "replace",
			arg:      222,
			expected: []any{222},
		},
		{
			subject:  []any{1, 1, 1},
			action:   "delete",
			arg:      []any{2, 2, 2},
			expected: []any(nil),
		},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("%v/%v/%v", tt.subject, tt.action, tt.arg)
		t.Run(name, func(t *testing.T) {
			merge := SliceMergeConcat
			if tt.upsert {
				merge = SliceMergeUpsert
			}
			modifier := Modifier(tt.action, tt.arg, WithSliceMerge(merge))
			altered, changed := modifier(tt.subject)
			assert.Equal(t, tt.expected, altered)
			assert.Equal(t, true, changed)
		})
	}
}

func TestModify_String(t *testing.T) {
	tests := []struct {
		subject  string
		action   Action
		arg      string
		expected string
	}{
		{
			subject:  "111",
			action:   "prepend",
			arg:      "222",
			expected: "222111",
		},
		{
			subject:  "111",
			action:   "append",
			arg:      "222",
			expected: "111222",
		},
		{
			subject:  "111",
			action:   "replace",
			arg:      "222",
			expected: "222",
		},
		{
			subject:  "111",
			action:   "delete",
			arg:      "222",
			expected: "",
		},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("%v/%v/%v", tt.subject, tt.action, tt.arg)
		t.Run(name, func(t *testing.T) {
			modifier := Modifier(tt.action, tt.arg)
			altered, changed := modifier(tt.subject)
			assert.Equal(t, tt.expected, altered)
			assert.Equal(t, true, changed)
		})
	}
}

func TestModify_AliasTypes(t *testing.T) {
	modifier := Modifier(ActionAppend, int(222))
	altered, changed := modifier(int(111))
	assert.Equal(t, int64(333), altered)
	assert.Equal(t, true, changed)

	modifier = Modifier(ActionAppend, int32(222))
	altered, changed = modifier(int32(111))
	assert.Equal(t, int64(333), altered)
	assert.Equal(t, true, changed)

	modifier = Modifier(ActionAppend, float32(222.0))
	altered, changed = modifier(float32(111.0))
	assert.Equal(t, float64(333.0), altered)
	assert.Equal(t, true, changed)
}

func TestModify_Unsupported(t *testing.T) {
	subject := struct{}{}
	arg := struct{}{}

	modifier := Modifier(ActionReplace, arg)
	altered, changed := modifier(subject)
	assert.Equal(t, subject, altered)
	assert.Equal(t, false, changed)
}
