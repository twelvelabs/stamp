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
				"src":  "/some/src/path",
				"dst":  123, // not a string
			},
			Task: nil,
			Err:  "'Dst' expected type 'string', got unconvertible type 'int'",
		},
		{
			Name: "returns an error if validation fails",
			TaskData: map[string]any{
				"type": "create",
				"src":  "/some/src/path",
				"dst":  "", // not present
			},
			Task: nil,
			Err:  "Dst is a required field",
		},
		{
			Name: "returns the correct struct for the given type",
			TaskData: map[string]any{
				"type": "create",
				"src":  "/some/src/path",
				"dst":  "/some/dst/path",
			},
			Task: &CreateTask{
				Common: Common{
					If:     "true",
					Each:   "",
					DryRun: false,
				},
				Src:      "/some/src/path",
				Dst:      "/some/dst/path",
				Mode:     "0666",
				Conflict: "prompt",
			},
			Err: "",
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
