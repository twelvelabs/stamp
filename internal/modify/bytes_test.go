package modify

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBytes(t *testing.T) {
	type args struct {
		dst    []byte
		action Action
		src    []byte
		conf   ModifierConf
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "prepend: concatenates src and dst bytes",
			args: args{
				dst:    []byte("foo bar"),
				action: ActionPrepend,
				src:    []byte("baz "),
			},
			want: []byte("baz foo bar"),
		},
		{
			name: "prepend(upsert): concatenates when dst does not contain src prefix",
			args: args{
				dst:    []byte("foo bar"),
				action: ActionPrepend,
				src:    []byte("baz "),
				conf: ModifierConf{
					MergeType: MergeTypeUpsert,
				},
			},
			want: []byte("baz foo bar"),
		},
		{
			name: "prepend(upsert): is noop when src prefix already present",
			args: args{
				dst:    []byte("baz foo bar"),
				action: ActionPrepend,
				src:    []byte("baz "),
				conf: ModifierConf{
					MergeType: MergeTypeUpsert,
				},
			},
			want: []byte("baz foo bar"),
		},
		{
			name: "prepend(replace): replaces dst with src",
			args: args{
				dst:    []byte("foo bar"),
				action: ActionPrepend,
				src:    []byte("baz"),
				conf: ModifierConf{
					MergeType: MergeTypeReplace,
				},
			},
			want: []byte("baz"),
		},

		{
			name: "append: concatenates dst and src bytes",
			args: args{
				dst:    []byte("foo bar"),
				action: ActionAppend,
				src:    []byte(" baz"),
			},
			want: []byte("foo bar baz"),
		},
		{
			name: "append(upsert): concatenates when dst does not contain src suffix",
			args: args{
				dst:    []byte("foo bar"),
				action: ActionAppend,
				src:    []byte(" baz"),
				conf: ModifierConf{
					MergeType: MergeTypeUpsert,
				},
			},
			want: []byte("foo bar baz"),
		},
		{
			name: "append(upsert): is noop when src suffix already present",
			args: args{
				dst:    []byte("foo bar baz"),
				action: ActionAppend,
				src:    []byte(" baz"),
				conf: ModifierConf{
					MergeType: MergeTypeUpsert,
				},
			},
			want: []byte("foo bar baz"),
		},
		{
			name: "append(replace): replaces dst with src",
			args: args{
				dst:    []byte("foo bar"),
				action: ActionAppend,
				src:    []byte("baz"),
				conf: ModifierConf{
					MergeType: MergeTypeReplace,
				},
			},
			want: []byte("baz"),
		},

		{
			name: "replace: replaces dst with src",
			args: args{
				dst:    []byte("foo"),
				action: ActionReplace,
				src:    []byte("bar"),
			},
			want: []byte("bar"),
		},

		{
			name: "delete: deletes dst entirely and ignores src",
			args: args{
				dst:    []byte("foo"),
				action: ActionDelete,
				src:    []byte("bar"),
			},
			want: []byte(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Bytes(tt.args.dst, tt.args.action, tt.args.src, tt.args.conf))
		})
	}
}
