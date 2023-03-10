package stamp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twelvelabs/termite/testutil"

	"github.com/twelvelabs/stamp/internal/value"
)

func TestNewTask_WhenTypeIsGenerator(t *testing.T) {
	tests := []struct {
		Desc     string
		TaskData map[string]any
		Task     interface{}
		Err      string
	}{
		{
			Desc: "returns an error when name is missing",
			TaskData: map[string]any{
				"type": "generator",
			},
			Task: nil,
			Err:  "Name is a required field",
		},
		{
			Desc: "returns the task when all fields are valid",
			TaskData: map[string]any{
				"type": "generator",
				"name": "foo",
			},
			Task: &GeneratorTask{
				Common: Common{
					If:   "true",
					Each: "",
				},
				Name:   "foo",
				Values: map[string]any{},
			},
			Err: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Desc, func(t *testing.T) {
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

func TestGeneratorTask_Execute(t *testing.T) {
	tests := []struct {
		Desc       string
		TaskData   map[string]any
		Values     map[string]any
		Prompter   *value.PrompterMock
		StartFiles map[string]any
		EndFiles   map[string]any
		Err        string
	}{
		{
			Desc: "returns an error when named generator is not found",
			TaskData: map[string]any{
				"type": "generator",
				"name": "unknown",
			},
			Values: map[string]any{},
			Err:    "generator not found",
		},

		{
			Desc: "executes the named generator",
			TaskData: map[string]any{
				"type": "generator",
				"name": "file",
			},
			Values: map[string]any{
				"FileName":    "hello.txt",
				"FileContent": "hello, world!",
			},
			EndFiles: map[string]any{
				"hello.txt": "hello, world!",
			},
			Err: "",
		},

		{
			Desc: "executes the named generator with value overrides",
			TaskData: map[string]any{
				"type": "generator",
				"name": "file",
				"values": map[string]any{
					"FileName":    "custom.txt",
					"FileContent": "custom content",
				},
			},
			Values: map[string]any{
				"FileName":    "hello.txt",
				"FileContent": "hello, world!",
			},
			EndFiles: map[string]any{
				"custom.txt": "custom content",
				"hello.txt":  false,
			},
			Err: "",
		},

		{
			Desc: "returns an error if unable to set custom values",
			TaskData: map[string]any{
				"type": "generator",
				"name": "file",
				"values": map[string]any{
					"FileName":    "custom.txt",
					"FileContent": "{{}",
				},
			},
			Values: map[string]any{},
			EndFiles: map[string]any{
				"custom.txt": false,
				"hello.txt":  false,
			},
			Err: "unexpected \"}\" in command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Desc, func(t *testing.T) {
			store := NewTestStore() // must call before changing dirs
			app := NewTestApp()
			app.Store = store

			testutil.InTempDir(t, func(tmpDir string) {
				// Populate the temp dir w/ any initial files
				testutil.WritePaths(t, tmpDir, tt.StartFiles)

				task, err := NewTask(tt.TaskData)
				assert.NoError(t, err)

				ctx := NewTaskContext(app, false)
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
