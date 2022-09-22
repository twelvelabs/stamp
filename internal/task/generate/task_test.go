package generate_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twelvelabs/stamp/internal/task"
	"github.com/twelvelabs/stamp/internal/task/common"
	"github.com/twelvelabs/stamp/internal/task/generate"
)

func TestGenerateTask(t *testing.T) {
	tests := []struct {
		Name   string
		Input  map[string]any
		Output interface{}
		Err    string
	}{
		{
			Name: "empty",
			Input: map[string]any{
				"type": "generate",
			},
			Output: nil,
			Err:    "Src is a required field, Dst is a required field",
		},
		{
			Name: "only src",
			Input: map[string]any{
				"type": "generate",
				"src":  "example.tpl",
			},
			Output: nil,
			Err:    "Dst is a required field",
		},
		{
			Name: "only dst",
			Input: map[string]any{
				"type": "generate",
				"dst":  "example.txt",
			},
			Output: nil,
			Err:    "Src is a required field",
		},
		{
			Name: "valid",
			Input: map[string]any{
				"type": "generate",
				"src":  "example.tpl",
				"dst":  "example.txt",
			},
			Output: &generate.Task{
				Common: common.Common{
					If:   "true",
					Each: "",
				},
				Src:      "example.tpl",
				Dst:      "example.txt",
				Conflict: "prompt",
				Mode:     "0666",
			},
			Err: "",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			actual, err := task.NewTask(test.Input)

			assert.Equal(t, test.Output, actual)
			if test.Err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, test.Err)
			}
		})
	}
}
