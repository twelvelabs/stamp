package modify

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	type args struct {
		dst    string
		action Action
		src    string
		conf   ModifierConf
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "prepend: concatenates src and dst slices",
			args: args{
				dst:    "foo bar",
				action: ActionPrepend,
				src:    "baz ",
			},
			want: "baz foo bar",
		},
		{
			name: "prepend(upsert): concatenates when dst does not contain src item",
			args: args{
				dst:    "foo bar",
				action: ActionPrepend,
				src:    "baz ",
				conf: ModifierConf{
					MergeType: MergeTypeUpsert,
				},
			},
			want: "baz foo bar",
		},
		{
			name: "prepend(upsert): is noop when src item already present",
			args: args{
				dst:    "baz foo bar",
				action: ActionPrepend,
				src:    "baz ",
				conf: ModifierConf{
					MergeType: MergeTypeUpsert,
				},
			},
			want: "baz foo bar",
		},
		{
			name: "prepend(replace): replaces dst with src",
			args: args{
				dst:    "foo bar",
				action: ActionPrepend,
				src:    "baz",
				conf: ModifierConf{
					MergeType: MergeTypeReplace,
				},
			},
			want: "baz",
		},

		{
			name: "append: concatenates dst and src slices",
			args: args{
				dst:    "foo bar",
				action: ActionAppend,
				src:    " baz",
			},
			want: "foo bar baz",
		},
		{
			name: "append: concatenates dst and src even if src is not a slice",
			args: args{
				dst:    "foo bar",
				action: ActionAppend,
				src:    " baz",
			},
			want: "foo bar baz",
		},
		{
			name: "append(upsert): concatenates when dst does not contain src item",
			args: args{
				dst:    "foo bar",
				action: ActionAppend,
				src:    " baz",
				conf: ModifierConf{
					MergeType: MergeTypeUpsert,
				},
			},
			want: "foo bar baz",
		},
		{
			name: "append(upsert): is noop when src item already present",
			args: args{
				dst:    "foo bar baz",
				action: ActionAppend,
				src:    " baz",
				conf: ModifierConf{
					MergeType: MergeTypeUpsert,
				},
			},
			want: "foo bar baz",
		},
		{
			name: "append(replace): replaces dst with src",
			args: args{
				dst:    "foo bar",
				action: ActionAppend,
				src:    "baz",
				conf: ModifierConf{
					MergeType: MergeTypeReplace,
				},
			},
			want: "baz",
		},

		{
			name: "replace: replaces dst with src",
			args: args{
				dst:    "foo bar",
				action: ActionReplace,
				src:    "baz",
			},
			want: "baz",
		},

		{
			name: "delete: deletes dst entirely and ignores src",
			args: args{
				dst:    "foo bar",
				action: ActionDelete,
				src:    "baz",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, String(tt.args.dst, tt.args.action, tt.args.src, tt.args.conf))
		})
	}
}
