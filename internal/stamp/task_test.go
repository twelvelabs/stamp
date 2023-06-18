package stamp

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTask(t *testing.T) {
	tests := []struct {
		Name        string
		TaskData    map[string]any
		Task        Task
		SetDefaults SetDefaultsFunc
		Err         string
	}{
		{
			Name:     "returns an error if type field is missing",
			TaskData: map[string]any{},
			Task:     nil,
			Err:      "undefined task type",
		},
		{
			Name: "returns an error if type field is unknown",
			TaskData: map[string]any{
				"type": "not-a-type",
			},
			Task: nil,
			Err:  "unknown task type: not-a-type",
		},
		{
			Name: "returns an error if unable to set defaults",
			TaskData: map[string]any{
				"type": "create",
			},
			Task: nil,
			SetDefaults: func(a any) error {
				return errors.New("boom")
			},
			Err: "boom",
		},
		{
			Name: "returns an error if decoding fails",
			TaskData: map[string]any{
				"type": "create",
				"src": map[string]any{
					"path": "/some/src/path",
				},
				"dst": map[string]any{
					"path": 123, // not a string
				},
			},
			Task: nil,
			Err:  "error(s) decoding",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.SetDefaults != nil {
				SetDefaults = tt.SetDefaults
				defer func() {
					SetDefaults = DefaultSetDefaultsFunc
				}()
			}
			actual, err := NewTask(tt.TaskData)
			assert.Equal(t, tt.Task, actual)
			if tt.Err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.Err)
			}
		})
	}
}
