package modify

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlice(t *testing.T) {
	type args struct {
		dst    []any
		action Action
		src    any
		conf   ModifierConf
	}
	tests := []struct {
		name string
		args args
		want []any
	}{
		{
			name: "prepend: concatenates src and dst slices",
			args: args{
				dst:    []any{"foo", "bar"},
				action: ActionPrepend,
				src:    []any{"baz"},
			},
			want: []any{"baz", "foo", "bar"},
		},
		{
			name: "prepend(upsert): concatenates when dst does not contain src item",
			args: args{
				dst:    []any{"foo", "bar"},
				action: ActionPrepend,
				src:    []any{"baz"},
				conf: ModifierConf{
					MergeType: MergeTypeUpsert,
				},
			},
			want: []any{"baz", "foo", "bar"},
		},
		{
			name: "prepend(upsert): is noop when src item already present",
			args: args{
				dst:    []any{"baz", "foo", "bar"},
				action: ActionPrepend,
				src:    []any{"baz"},
				conf: ModifierConf{
					MergeType: MergeTypeUpsert,
				},
			},
			want: []any{"baz", "foo", "bar"},
		},
		{
			name: "prepend(replace): replaces dst with src",
			args: args{
				dst:    []any{"foo", "bar"},
				action: ActionPrepend,
				src:    []any{"baz"},
				conf: ModifierConf{
					MergeType: MergeTypeReplace,
				},
			},
			want: []any{"baz"},
		},

		{
			name: "append: concatenates dst and src slices",
			args: args{
				dst:    []any{"foo", "bar"},
				action: ActionAppend,
				src:    []any{"baz"},
			},
			want: []any{"foo", "bar", "baz"},
		},
		{
			name: "append: concatenates dst and src even if src is not a slice",
			args: args{
				dst:    []any{"foo", "bar"},
				action: ActionAppend,
				src:    "baz",
			},
			want: []any{"foo", "bar", "baz"},
		},
		{
			name: "append(upsert): concatenates when dst does not contain src item",
			args: args{
				dst:    []any{"foo", "bar"},
				action: ActionAppend,
				src:    []any{"baz"},
				conf: ModifierConf{
					MergeType: MergeTypeUpsert,
				},
			},
			want: []any{"foo", "bar", "baz"},
		},
		{
			name: "append(upsert): is noop when src item already present",
			args: args{
				dst:    []any{"foo", "bar", "baz"},
				action: ActionAppend,
				src:    []any{"baz"},
				conf: ModifierConf{
					MergeType: MergeTypeUpsert,
				},
			},
			want: []any{"foo", "bar", "baz"},
		},
		{
			name: "append(replace): replaces dst with src",
			args: args{
				dst:    []any{"foo", "bar"},
				action: ActionAppend,
				src:    []any{"baz"},
				conf: ModifierConf{
					MergeType: MergeTypeReplace,
				},
			},
			want: []any{"baz"},
		},

		{
			name: "replace: replaces dst with src",
			args: args{
				dst:    []any{"foo", "bar"},
				action: ActionReplace,
				src:    []any{"baz"},
			},
			want: []any{"baz"},
		},

		{
			name: "delete: deletes dst entirely and ignores src",
			args: args{
				dst:    []any{"foo", "bar"},
				action: ActionDelete,
				src:    []any{"baz"},
			},
			want: []any(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Slice(tt.args.dst, tt.args.action, tt.args.src, tt.args.conf))
		})
	}
}

func TestSlice_DoesNotMutateArguments(t *testing.T) {
	dst := []any{1, 2}
	src := []any{3, 4}

	result := Slice(dst, ActionAppend, src, ModifierConf{})

	assert.Equal(t, []any{1, 2, 3, 4}, result)
	assert.Equal(t, []any{1, 2}, dst, "dst should not have changed")
	assert.Equal(t, []any{3, 4}, src, "src should not have changed")
}
