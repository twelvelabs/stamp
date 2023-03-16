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
				"content": "bar",
			},
			Task: &UpdateTask{
				Common: Common{
					If:   "true",
					Each: "",
				},
				Dst:     "example.txt",
				Missing: "ignore",
				Pattern: "foo",
				Action:  "replace",
				Content: "bar",
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
		{
			Desc: "returns an error if pattern evaluates to empty string",
			StartFiles: map[string]any{
				"README.md": "Hello World\n",
			},
			TaskData: map[string]any{
				"type":    "update",
				"dst":     "./README.md",
				"pattern": "{{ .Empty }}",
			},
			Values: map[string]any{
				"Empty": "",
			},
			EndFiles: map[string]any{
				"README.md": "Hello World\n",
			},
			Err: "pattern: '{{ .Empty }}' evaluated to an empty string",
		},
		{
			Desc: "returns an error if pattern can not be compiled to regexp",
			StartFiles: map[string]any{
				"README.md": "Hello World\n",
			},
			TaskData: map[string]any{
				"type":    "update",
				"dst":     "./README.md",
				"pattern": "(.}",
			},
			EndFiles: map[string]any{
				"README.md": "Hello World\n",
			},
			Err: "error parsing regexp",
		},

		{
			Desc: "updates a path",
			StartFiles: map[string]any{
				"README.md": "Hello World\n",
			},
			TaskData: map[string]any{
				"type":    "update",
				"dst":     "./README.md",
				"pattern": "Hello (\\w+)",
				"action":  "replace",
				"content": "Goodbye $1",
			},
			EndFiles: map[string]any{
				"README.md": "Goodbye World\n",
			},
		},
		{
			Desc: "updates a path and changes file mode",
			StartFiles: map[string]any{
				"README.md": "Hello World\n",
			},
			TaskData: map[string]any{
				"type":    "update",
				"dst":     "./README.md",
				"mode":    "0755",
				"pattern": "Hello (\\w+)",
				"action":  "replace",
				"content": "Goodbye $1",
			},
			EndFiles: map[string]any{
				"README.md": []any{
					"Goodbye World\n",
					0o755,
				},
			},
		},
		{
			Desc:   "does not update paths during a dry run",
			DryRun: true,
			StartFiles: map[string]any{
				"README.md": "Hello World\n",
			},
			TaskData: map[string]any{
				"type":    "update",
				"dst":     "./README.md",
				"pattern": "Hello (\\w+)",
				"action":  "replace",
				"content": "Goodbye $1",
			},
			EndFiles: map[string]any{
				"README.md": "Hello World\n",
			},
		},

		{
			Desc: "[missing:ignore] ignores missing paths",
			TaskData: map[string]any{
				"type":    "update",
				"dst":     "./README.md",
				"pattern": "foo",
			},
			EndFiles: map[string]any{
				"README.md": false,
			},
		},
		{
			Desc: "[missing:error] returns an error when path is missing",
			TaskData: map[string]any{
				"type":    "update",
				"dst":     "./README.md",
				"missing": "error",
				"pattern": "foo",
			},
			EndFiles: map[string]any{
				"README.md": false,
			},
			Err: "path not found",
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
