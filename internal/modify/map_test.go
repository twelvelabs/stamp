package modify

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) { //nolint:maintidx
	type args struct {
		dst    map[string]any
		action Action
		src    map[string]any
		conf   ModifierConf
	}
	tests := []struct {
		name string
		args args
		want map[string]any
	}{
		{
			name: "prepend: concatenates the src and dst maps",
			args: args{
				dst: map[string]any{
					"aaa": "aaa.dst",
					"bbb": []any{1, 2, 3, 4},
					"ccc": map[string]any{
						"111": "222.dst",
						"222": "222.dst",
					},
					"dst": "dst",
				},
				action: ActionPrepend,
				src: map[string]any{
					"aaa": "aaa.src",
					"bbb": []any{2, 3, 4, 5},
					"ccc": map[string]any{
						"111": "222.src",
						"222": "222.src",
					},
					"src": "src",
				},
			},
			want: map[string]any{
				"aaa": "aaa.dst",
				"bbb": []any{2, 3, 4, 5, 1, 2, 3, 4},
				"ccc": map[string]any{
					"111": "222.dst",
					"222": "222.dst",
				},
				"dst": "dst",
				"src": "src",
			},
		},
		{
			name: "prepend(upsert): concatenates the src and dst maps without slice duplicates",
			args: args{
				dst: map[string]any{
					"aaa": "aaa.dst",
					"bbb": []any{1, 2, 3, 4},
					"ccc": map[string]any{
						"111": "222.dst",
						"222": "222.dst",
					},
					"dst": "dst",
				},
				action: ActionPrepend,
				src: map[string]any{
					"aaa": "aaa.src",
					"bbb": []any{2, 3, 4, 5},
					"ccc": map[string]any{
						"111": "222.src",
						"222": "222.src",
					},
					"src": "src",
				},
				conf: ModifierConf{
					MergeType: MergeTypeUpsert,
				},
			},
			want: map[string]any{
				"aaa": "aaa.dst",
				"bbb": []any{2, 3, 4, 5, 1},
				"ccc": map[string]any{
					"111": "222.dst",
					"222": "222.dst",
				},
				"dst": "dst",
				"src": "src",
			},
		},
		{
			name: "prepend(replace): concatenates the src and dst maps, replacing slice values",
			args: args{
				dst: map[string]any{
					"aaa": "aaa.dst",
					"bbb": []any{1, 2, 3, 4},
					"ccc": map[string]any{
						"111": "222.dst",
						"222": "222.dst",
					},
					"dst": "dst",
				},
				action: ActionPrepend,
				src: map[string]any{
					"aaa": "aaa.src",
					"bbb": []any{2, 3, 4, 5},
					"ccc": map[string]any{
						"111": "222.src",
						"222": "222.src",
					},
					"src": "src",
				},
				conf: ModifierConf{
					MergeType: MergeTypeReplace,
				},
			},
			want: map[string]any{
				"aaa": "aaa.dst",
				"bbb": []any{1, 2, 3, 4},
				"ccc": map[string]any{
					"111": "222.dst",
					"222": "222.dst",
				},
				"dst": "dst",
				"src": "src",
			},
		},

		{
			name: "append: concatenates the dst and src maps",
			args: args{
				dst: map[string]any{
					"aaa": "aaa.dst",
					"bbb": []any{1, 2, 3, 4},
					"ccc": map[string]any{
						"111": "222.dst",
						"222": "222.dst",
					},
					"dst": "dst",
				},
				action: ActionAppend,
				src: map[string]any{
					"aaa": "aaa.src",
					"bbb": []any{2, 3, 4, 5},
					"ccc": map[string]any{
						"111": "222.src",
						"222": "222.src",
					},
					"src": "src",
				},
			},
			want: map[string]any{
				"aaa": "aaa.src",
				"bbb": []any{1, 2, 3, 4, 2, 3, 4, 5},
				"ccc": map[string]any{
					"111": "222.src",
					"222": "222.src",
				},
				"dst": "dst",
				"src": "src",
			},
		},
		{
			name: "append: recursively merges maps",
			args: args{
				dst: map[string]any{
					"aaa": map[string]any{
						"111": map[string]any{
							"aaa": "aaa.dst",
							"bbb": "bbb.dst",
							"dst": "dst",
						},
					},
				},
				action: ActionAppend,
				src: map[string]any{
					"aaa": map[string]any{
						"111": map[string]any{
							"aaa": "aaa.src",
							"bbb": "bbb.src",
							"src": "src",
						},
					},
				},
			},
			want: map[string]any{
				"aaa": map[string]any{
					"111": map[string]any{
						"aaa": "aaa.src",
						"bbb": "bbb.src",
						"dst": "dst",
						"src": "src",
					},
				},
			},
		},
		{
			name: "append: replaces complex types when src value does not have the same type",
			args: args{
				dst: map[string]any{
					"aaa": "aaa.dst",
					"bbb": []any{1, 2, 3, 4},
					"ccc": map[string]any{
						"111": "222.dst",
						"222": "222.dst",
					},
					"dst": "dst",
				},
				action: ActionAppend,
				src: map[string]any{
					"aaa": "aaa.src",
					"bbb": "bbb.src",
					"ccc": "ccc.src",
					"src": "src",
				},
			},
			want: map[string]any{
				"aaa": "aaa.src",
				"bbb": "bbb.src",
				"ccc": "ccc.src",
				"dst": "dst",
				"src": "src",
			},
		},
		{
			name: "append(upsert): concatenates the dst and src maps without slice duplicates",
			args: args{
				dst: map[string]any{
					"aaa": "aaa.dst",
					"bbb": []any{1, 2, 3, 4},
					"ccc": map[string]any{
						"111": "222.dst",
						"222": "222.dst",
					},
					"dst": "dst",
				},
				action: ActionAppend,
				src: map[string]any{
					"aaa": "aaa.src",
					"bbb": []any{2, 3, 4, 5},
					"ccc": map[string]any{
						"111": "222.src",
						"222": "222.src",
					},
					"src": "src",
				},
				conf: ModifierConf{
					MergeType: MergeTypeUpsert,
				},
			},
			want: map[string]any{
				"aaa": "aaa.src",
				"bbb": []any{1, 2, 3, 4, 5},
				"ccc": map[string]any{
					"111": "222.src",
					"222": "222.src",
				},
				"dst": "dst",
				"src": "src",
			},
		},
		{
			name: "append(replace): concatenates the dst and src maps, replacing slice values",
			args: args{
				dst: map[string]any{
					"aaa": "aaa.dst",
					"bbb": []any{1, 2, 3, 4},
					"ccc": map[string]any{
						"111": "222.dst",
						"222": "222.dst",
					},
					"dst": "dst",
				},
				action: ActionAppend,
				src: map[string]any{
					"aaa": "aaa.src",
					"bbb": []any{2, 3, 4, 5},
					"ccc": map[string]any{
						"111": "222.src",
						"222": "222.src",
					},
					"src": "src",
				},
				conf: ModifierConf{
					MergeType: MergeTypeReplace,
				},
			},
			want: map[string]any{
				"aaa": "aaa.src",
				"bbb": []any{2, 3, 4, 5},
				"ccc": map[string]any{
					"111": "222.src",
					"222": "222.src",
				},
				"dst": "dst",
				"src": "src",
			},
		},

		{
			name: "replace: replaces dst with src map",
			args: args{
				dst: map[string]any{
					"aaa": "aaa.dst",
					"bbb": []any{1, 2, 3, 4},
					"ccc": map[string]any{
						"111": "222.dst",
						"222": "222.dst",
					},
					"dst": "dst",
				},
				action: ActionReplace,
				src: map[string]any{
					"aaa": "aaa.src",
					"bbb": []any{2, 3, 4, 5},
					"ccc": map[string]any{
						"111": "222.src",
						"222": "222.src",
					},
					"src": "src",
				},
			},
			want: map[string]any{
				"aaa": "aaa.src",
				"bbb": []any{2, 3, 4, 5},
				"ccc": map[string]any{
					"111": "222.src",
					"222": "222.src",
				},
				"src": "src",
			},
		},

		{
			name: "delete: deletes dst and ignores src map",
			args: args{
				dst: map[string]any{
					"aaa": "aaa.dst",
					"bbb": []any{1, 2, 3, 4},
					"ccc": map[string]any{
						"111": "222.dst",
						"222": "222.dst",
					},
					"dst": "dst",
				},
				action: ActionDelete,
				src: map[string]any{
					"aaa": "aaa.src",
					"bbb": []any{2, 3, 4, 5},
					"ccc": map[string]any{
						"111": "222.src",
						"222": "222.src",
					},
					"src": "src",
				},
			},
			want: map[string]any(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Map(tt.args.dst, tt.args.action, tt.args.src, tt.args.conf))
		})
	}
}

func TestMap_DoesNotMutateArguments(t *testing.T) {
	dst := map[string]any{
		"aaa": "aaa.dst",
		"bbb": []any{1, 2},
	}
	src := map[string]any{
		"bbb": []any{3, 4},
		"ccc": "ccc.src",
	}

	result := Map(dst, ActionAppend, src, ModifierConf{})

	assert.Equal(t, map[string]any{
		"aaa": "aaa.dst",
		"bbb": []any{1, 2, 3, 4},
		"ccc": "ccc.src",
	}, result)
	assert.Equal(t, map[string]any{
		"aaa": "aaa.dst",
		"bbb": []any{1, 2},
	}, dst, "dst should not have changed")
	assert.Equal(t, map[string]any{
		"bbb": []any{3, 4},
		"ccc": "ccc.src",
	}, src, "src should not have changed")
}
