package stamp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twelvelabs/termite/testutil"
)

func TestNewTask_WhenTypeIsUpdate(t *testing.T) {
	tests := []struct {
		Name     string
		TaskData map[string]any
		Task     Task
		Err      string
	}{
		{
			Name: "returns an error when dst field is missing",
			TaskData: map[string]any{
				"type":    "update",
				"dst":     "",
				"pattern": "foo",
			},
			Task: nil,
			Err:  "Dst is a required field",
		},
		{
			Name: "returns an error when pattern field is missing",
			TaskData: map[string]any{
				"type":    "update",
				"dst":     "example.txt",
				"pattern": "",
			},
			Task: nil,
			Err:  "Pattern is a required field",
		},
		{
			Name: "returns an error when missing field is invalid",
			TaskData: map[string]any{
				"type":    "update",
				"dst":     "example.txt",
				"missing": "unknown",
				"pattern": "foo",
			},
			Task: nil,
			Err:  "unknown is not a valid MissingConfig",
		},
		{
			Name: "returns the task when all fields are valid",
			TaskData: map[string]any{
				"type":    "update",
				"dst":     "example.txt",
				"pattern": "foo",
				"value":   "bar",
			},
			Task: &UpdateTask{
				Common: Common{
					If:   "true",
					Each: "",
				},
				Dst:     "example.txt",
				Missing: "ignore",
				Pattern: "foo",
				Value:   "bar",
			},
			Err: "",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			actual, err := NewTask(test.TaskData)

			assert.Equal(t, test.Task, actual)
			if test.Err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, test.Err)
			}
		})
	}
}

func TestUpdateTask_Execute(t *testing.T) {
	tests := []struct {
		Desc       string
		DryRun     bool
		TaskData   map[string]any
		Values     map[string]any
		StartFiles map[string]any
		EndFiles   map[string]any
		Setup      func(app *App)
		Err        string
	}{
		{
			Desc: "returns an error if dst evaluates to empty string",
			TaskData: map[string]any{
				"type":    "update",
				"dst":     "{{ .Empty }}",
				"pattern": "foo",
			},
			Values: map[string]any{
				"Empty": "",
			},
			Err: "dst: '{{ .Empty }}' evaluated to an empty string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Desc, func(t *testing.T) {
			testutil.InTempDir(t, func(tmpDir string) {
				// Populate the temp dir w/ any initial files
				testutil.WritePaths(t, tmpDir, tt.StartFiles)

				// Setup the app.
				app := NewTestApp()
				app.Config.DryRun = tt.DryRun
				if tt.Setup != nil {
					tt.Setup(app)
				}
				defer app.UI.VerifyStubs(t)

				// Create a new task and execute it.
				task, err := NewTask(tt.TaskData)
				require.NoError(t, err)
				ctx := NewTaskContext(app)
				err = task.Execute(ctx, tt.Values)

				// Ensure the expected files were generated
				testutil.AssertPaths(t, tmpDir, tt.EndFiles)

				if tt.Err == "" {
					assert.NoError(t, err)
				} else {
					assert.ErrorContains(t, err, tt.Err)
				}
			})
		})
	}
}
