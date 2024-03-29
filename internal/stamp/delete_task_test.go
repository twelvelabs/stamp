package stamp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twelvelabs/termite/testutil"
)

func TestNewTask_WhenTypeIsDelete(t *testing.T) {
	tests := []struct {
		Name     string
		TaskData map[string]any
		Err      string
	}{
		{
			Name: "returns an error when missing field is invalid",
			TaskData: map[string]any{
				"type": "delete",
				"dst": map[string]any{
					"path":    "example.txt",
					"missing": "unknown",
				},
			},
			Err: "unknown is not a valid MissingConfig",
		},
		{
			Name: "returns the task when all fields are valid",
			TaskData: map[string]any{
				"type": "delete",
				"dst": map[string]any{
					"path": "example.txt",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			actual, err := NewTask(test.TaskData)

			if test.Err == "" {
				assert.NoError(t, err)
				assert.NotNil(t, actual)
			} else {
				assert.ErrorContains(t, err, test.Err)
				assert.Nil(t, actual)
			}
		})
	}
}

func TestDeleteTask_Execute(t *testing.T) {
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
				"type": "delete",
				"dst": map[string]any{
					"path": "{{ .Empty }}",
				},
			},
			Values: map[string]any{
				"Empty": "",
			},
			Err: "evaluated to an empty string",
		},

		{
			Desc: "deletes a path",
			StartFiles: map[string]any{
				"README.md": "Pre-existing content",
			},
			TaskData: map[string]any{
				"type": "delete",
				"dst": map[string]any{
					"path": "README.md",
				},
			},
			Values: map[string]any{
				"DstPath": ".",
			},
			EndFiles: map[string]any{
				"README.md": false,
			},
			Err: "",
		},
		{
			Desc:   "does not delete paths during a dry run",
			DryRun: true,
			StartFiles: map[string]any{
				"README.md": "Pre-existing content",
			},
			TaskData: map[string]any{
				"type": "delete",
				"dst": map[string]any{
					"path": "README.md",
				},
			},
			Values: map[string]any{
				"DstPath": ".",
			},
			EndFiles: map[string]any{
				"README.md": "Pre-existing content",
			},
			Err: "",
		},

		{
			Desc: "[missing:ignore] ignores missing paths",
			TaskData: map[string]any{
				"type": "delete",
				"dst": map[string]any{
					"path": "README.md",
				},
			},
			Values: map[string]any{
				"DstPath": ".",
			},
			EndFiles: map[string]any{
				"README.md": false,
			},
			Err: "",
		},

		{
			Desc: "[missing:error] returns an error when path is missing",
			TaskData: map[string]any{
				"type": "delete",
				"dst": map[string]any{
					"path":    "README.md",
					"missing": "error",
				},
			},
			Values: map[string]any{
				"DstPath": ".",
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
