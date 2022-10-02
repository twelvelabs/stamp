package gen

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twelvelabs/stamp/internal/iostreams"
	"github.com/twelvelabs/stamp/internal/testutil"
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
				Name:  "foo",
				Extra: map[string]any{},
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
	packagesDir := filepath.Join("..", "..", "testdata", "generators")
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
			Desc: "executes the named generator with values",
			TaskData: map[string]any{
				"type": "generator",
				"name": "hello",
			},
			Values: map[string]any{
				"Greeting": "hello, world!",
			},
			EndFiles: map[string]any{
				"hello.txt": "hello, world!",
			},
			Err: "",
		},

		{
			Desc: "executes the named generator with extras",
			TaskData: map[string]any{
				"type": "generator",
				"name": "hello",
				"extra": map[string]any{
					"Greeting": "hello, world!",
				},
			},
			Values: map[string]any{},
			EndFiles: map[string]any{
				"hello.txt": "hello, world!",
			},
			Err: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Desc, func(t *testing.T) {
			defer testutil.Cleanup()

			// Create a temp dir
			tmpDir := testutil.MkdirTemp()
			tt.Values["DstPath"] = tmpDir

			// Populate the temp dir w/ any initial files
			testutil.CreatePaths(tmpDir, tt.StartFiles)

			store := NewStore(packagesDir)

			task, err := NewTask(tt.TaskData)
			assert.NoError(t, err)

			ios := iostreams.Test()
			ctx := NewTaskContext(ios, tt.Prompter, store, false)
			err = task.Execute(ctx, tt.Values)

			// Ensure the expected files were generated
			testutil.AssertPaths(t, tmpDir, tt.EndFiles)

			if tt.Err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.Err)
			}
		})
	}
}
