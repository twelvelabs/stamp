package stamp

import (
	"path/filepath"
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
				"type":  "update",
				"dst":   "",
				"match": "foo",
			},
			Task: nil,
			Err:  "Dst is a required field",
		},
		{
			Name: "returns an error when missing field is invalid",
			TaskData: map[string]any{
				"type":    "update",
				"dst":     "example.txt",
				"missing": "unknown",
				"match":   "foo",
			},
			Task: nil,
			Err:  "unknown is not a valid MissingConfig",
		},
		{
			Name: "returns the task when all fields are valid",
			TaskData: map[string]any{
				"type":  "update",
				"dst":   "example.txt",
				"match": "foo",
			},
			Task: &UpdateTask{
				Common: Common{
					If:   "true",
					Each: "",
				},
				Dst:     "example.txt",
				Missing: "ignore",
				Match:   "foo",
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

func TestUpdateTask_Execute(t *testing.T) { //nolint: maintidx
	// Note: have to do this here since sub-tests are run in a temp dir.
	srcPath, _ := filepath.Abs(filepath.Join("testdata", "templates"))
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
			Desc: "returns an error if src evaluates to empty string",
			TaskData: map[string]any{
				"type": "update",
				"src":  "{{ .Empty }}",
				"dst":  "./README.md",
			},
			Values: map[string]any{
				"Empty": "",
			},
			Err: "src: '{{ .Empty }}' evaluated to an empty string",
		},
		{
			Desc: "returns an error if src can not be rendered",
			TaskData: map[string]any{
				"type": "update",
				"src":  "./render-error.txt",
				"dst":  "./README.md",
			},
			Values: map[string]any{
				"SrcPath": srcPath,
			},
			Err: "boom",
		},
		{
			Desc: "returns an error if src can not be parsed",
			TaskData: map[string]any{
				"type": "update",
				"src":  "./invalid.json",
				"dst":  "./something.json",
			},
			Values: map[string]any{
				"SrcPath": srcPath,
			},
			Err: "unexpected character",
		},
		{
			Desc: "returns an error if dst evaluates to empty string",
			TaskData: map[string]any{
				"type": "update",
				"dst":  "{{ .Empty }}",
			},
			Values: map[string]any{
				"Empty": "",
			},
			Err: "dst: '{{ .Empty }}' evaluated to an empty string",
		},
		{
			Desc: "returns an error if match can not be decoded into config",
			TaskData: map[string]any{
				"type": "update",
				"dst":  "./README.md",
				"match": map[string]any{
					"pattern": true,
				},
			},
			Err: "expected type 'string', got unconvertible type 'bool'",
		},
		{
			Desc: "returns an error if pattern can not be compiled to regexp",
			StartFiles: map[string]any{
				"README.md": "Hello World\n",
			},
			TaskData: map[string]any{
				"type":  "update",
				"dst":   "./README.md",
				"match": "(.}",
			},
			EndFiles: map[string]any{
				"README.md": "Hello World\n",
			},
			Err: "error parsing regexp",
		},
		{
			Desc: "returns an error if replacement can not be cast to string",
			StartFiles: map[string]any{
				"README.md": "Hello World\n",
			},
			TaskData: map[string]any{
				"type":  "update",
				"dst":   "./README.md",
				"match": "World",
				"src":   struct{}{},
			},
			EndFiles: map[string]any{
				"README.md": "Hello World\n",
			},
			Err: "unable to cast struct",
		},
		{
			Desc: "returns an error if action can not be parsed",
			TaskData: map[string]any{
				"type":   "update",
				"dst":    "./README.md",
				"action": "unknown",
			},
			Err: "unknown is not a valid Action",
		},
		{
			Desc: "returns an error if file_type can not be parsed",
			TaskData: map[string]any{
				"type":      "update",
				"dst":       "./README.md",
				"file_type": "unknown",
			},
			Err: "unknown is not a valid FileType",
		},
		{
			Desc: "returns an error if mode can not be parsed",
			TaskData: map[string]any{
				"type": "update",
				"dst":  "README.md",
				"mode": "unknown",
			},
			Err: "invalid syntax",
		},
		{
			Desc: "returns an error if dst can not be touched",
			TaskData: map[string]any{
				"type":    "update",
				"dst":     "invalid\x00file.txt",
				"missing": "touch",
				"src":     "Hello",
			},
			Err: "invalid argument",
		},

		{
			Desc: "prepends a string in dst",
			StartFiles: map[string]any{
				"README.md": "aaa bbb ccc\n",
			},
			TaskData: map[string]any{
				"type":   "update",
				"dst":    "./README.md",
				"match":  "aaa",
				"action": "prepend",
				"src":    "000 ",
			},
			EndFiles: map[string]any{
				"README.md": "000 aaa bbb ccc\n",
			},
		},
		{
			Desc: "appends a string in dst",
			StartFiles: map[string]any{
				"README.md": "aaa bbb ccc\n",
			},
			TaskData: map[string]any{
				"type":        "update",
				"dst":         "./README.md",
				"match":       "ccc",
				"action":      "append",
				"src":         " ddd",
				"description": "append ddd",
			},
			EndFiles: map[string]any{
				"README.md": "aaa bbb ccc ddd\n",
			},
		},
		{
			Desc: "replaces a string in dst",
			StartFiles: map[string]any{
				"README.md": "Hello World\n",
			},
			TaskData: map[string]any{
				"type":   "update",
				"dst":    "./README.md",
				"match":  "Hello (\\w+)",
				"action": "replace",
				"src":    "Goodbye $1",
			},
			EndFiles: map[string]any{
				"README.md": "Goodbye World\n",
			},
		},
		{
			Desc: "deletes a string in dst",
			StartFiles: map[string]any{
				"README.md": "Hello World\n",
			},
			TaskData: map[string]any{
				"type":   "update",
				"dst":    "./README.md",
				"match":  "(?m)\\s*(\\w+)$",
				"action": "delete",
			},
			EndFiles: map[string]any{
				"README.md": "Hello\n",
			},
		},

		{
			Desc: "matches all text if pattern evaluates to empty string",
			StartFiles: map[string]any{
				"README.md": "Hello World\n",
			},
			TaskData: map[string]any{
				"type":   "update",
				"dst":    "./README.md",
				"action": "append",
				"src":    "Goodbye\n",
			},
			EndFiles: map[string]any{
				"README.md": "Hello World\nGoodbye\n",
			},
		},
		{
			Desc: "creates and updates a path when missing set to touch",
			TaskData: map[string]any{
				"type":    "update",
				"missing": "touch",
				"dst":     "README.md",
				"src":     "Howdy",
			},
			EndFiles: map[string]any{
				"README.md": "Howdy",
			},
		},
		{
			Desc: "updates a path and changes file mode",
			StartFiles: map[string]any{
				"README.md": "Hello World\n",
			},
			TaskData: map[string]any{
				"type":   "update",
				"dst":    "./README.md",
				"mode":   "0755",
				"match":  "Hello (\\w+)",
				"action": "replace",
				"src":    "Goodbye $1",
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
				"type":   "update",
				"dst":    "./README.md",
				"match":  "Hello (\\w+)",
				"action": "replace",
				"src":    "Goodbye $1",
			},
			EndFiles: map[string]any{
				"README.md": "Hello World\n",
			},
		},
		{
			Desc: "updates a path using content from src field",
			StartFiles: map[string]any{
				"README.md": "Hello World\n",
			},
			TaskData: map[string]any{
				"type":   "update",
				"src":    "./valid.txt",
				"dst":    "./README.md",
				"action": "replace",
			},
			Values: map[string]any{
				"SrcPath": srcPath,
				"Name":    "Some Person",
			},
			EndFiles: map[string]any{
				"README.md": "Hello, Some Person",
			},
		},
		{
			Desc: "matches patterns by line by default",
			StartFiles: map[string]any{
				"README.md": `
the first line might be foo
the second line might be bar
the third line might be baz
`,
			},
			TaskData: map[string]any{
				"type": "update",
				"dst":  "README.md",
				"match": map[string]any{
					"pattern": `^the (\w+) line might be (\w+)$`,
				},
				"src": "the ${1} line IS ${2}",
			},
			Values: map[string]any{
				"SrcPath": srcPath,
			},
			EndFiles: map[string]any{
				"README.md": `
the first line IS foo
the second line IS bar
the third line IS baz
`,
			},
		},
		{
			Desc: "matches patterns multiline if configured",
			StartFiles: map[string]any{
				"README.md": `
the first line might be foo
the second line might be bar
the third line might be baz
`,
			},
			TaskData: map[string]any{
				"type": "update",
				"dst":  "README.md",
				"match": map[string]any{
					"pattern": `(?ms) second(.+)baz$`,
					"source":  "file",
				},
				"action": "delete",
			},
			Values: map[string]any{
				"SrcPath": srcPath,
			},
			EndFiles: map[string]any{
				"README.md": `
the first line might be foo
the
`,
			},
		},

		{
			Desc: "prepends JSON data in dst",
			StartFiles: map[string]any{
				"example.json": `{"foo":[1,2,3]}`,
			},
			TaskData: map[string]any{
				"type":   "update",
				"dst":    "example.json",
				"match":  "$.foo",
				"action": "prepend",
				"src":    []any{4, 5},
			},
			EndFiles: map[string]any{
				"example.json": `{
    "foo": [
        4,
        5,
        1,
        2,
        3
    ]
}`,
			},
		},
		{
			Desc: "appends JSON data in dst",
			StartFiles: map[string]any{
				"example.json": `{"foo":[1,2,3]}`,
			},
			TaskData: map[string]any{
				"type":   "update",
				"dst":    "example.json",
				"match":  "$.foo",
				"action": "append",
				"src":    []any{4, 5},
			},
			EndFiles: map[string]any{
				"example.json": `{
    "foo": [
        1,
        2,
        3,
        4,
        5
    ]
}`,
			},
		},
		{
			Desc: "replaces JSON data in dst",
			StartFiles: map[string]any{
				"example.json": `{"foo":[1,2,3]}`,
			},
			TaskData: map[string]any{
				"type":   "update",
				"dst":    "example.json",
				"match":  "$.foo",
				"action": "replace",
				"src":    []any{4, 5},
			},
			EndFiles: map[string]any{
				"example.json": `{
    "foo": [
        4,
        5
    ]
}`,
			},
		},
		{
			Desc: "deletes JSON data in dst",
			StartFiles: map[string]any{
				"example.json": `{"foo":[1,2,3]}`,
			},
			TaskData: map[string]any{
				"type":   "update",
				"dst":    "example.json",
				"match":  "$.foo",
				"action": "delete",
			},
			EndFiles: map[string]any{
				"example.json": `{}`,
			},
		},
		{
			Desc: "matches root element if pattern evaluates to empty string",
			StartFiles: map[string]any{
				"example.json": `{"foo":[1,2,3]}`,
			},
			TaskData: map[string]any{
				"type":   "update",
				"dst":    "example.json",
				"action": "replace",
				"src": map[string]any{
					"bar": true,
				},
			},
			EndFiles: map[string]any{
				"example.json": `{
    "bar": true
}`,
			},
		},
		{
			Desc: "parses src path content before updating",
			StartFiles: map[string]any{
				"example.json": `{"foo":[1,2,3]}`,
			},
			TaskData: map[string]any{
				"type":   "update",
				"dst":    "example.json",
				"action": "append",
				"match":  "$.foo",
				"src":    "valid.json",
			},
			Values: map[string]any{
				"SrcPath": srcPath,
			},
			EndFiles: map[string]any{
				"example.json": `{
    "foo": [
        1,
        2,
        3,
        "aaa",
        "bbb",
        "ccc"
    ]
}`,
			},
		},
		{
			Desc: "ensures pattern default before updating",
			StartFiles: map[string]any{
				"example.json": `{}`,
			},
			TaskData: map[string]any{
				"type":   "update",
				"dst":    "example.json",
				"action": "append",
				"src":    "valid.json",
				"match": map[string]any{
					"pattern": "$.foo",
					"default": []any{},
				},
			},
			Values: map[string]any{
				"SrcPath": srcPath,
			},
			EndFiles: map[string]any{
				"example.json": `{
    "foo": [
        "aaa",
        "bbb",
        "ccc"
    ]
}`,
			},
		},
		{
			Desc: "does not append duplicate items when array_mode is upsert",
			StartFiles: map[string]any{
				"example.json": `{"foo": ["aaa"]}`,
			},
			TaskData: map[string]any{
				"type":  "update",
				"dst":   "example.json",
				"match": "$.foo",
				"action": map[string]any{
					"type":  "append",
					"merge": "upsert",
				},
				"src": []any{
					"aaa",
					"bbb",
					"ccc",
				},
			},
			EndFiles: map[string]any{
				"example.json": `{
    "foo": [
        "aaa",
        "bbb",
        "ccc"
    ]
}`,
			},
		},

		{
			Desc: "prepends YAML data in dst",
			StartFiles: map[string]any{
				"example.yml": "foo: [1,2,3]\n",
			},
			TaskData: map[string]any{
				"type":   "update",
				"dst":    "example.yml",
				"match":  "$.foo",
				"action": "prepend",
				"src":    []any{4, 5},
			},
			EndFiles: map[string]any{
				"example.yml": "foo:\n" +
					"  - 4\n" +
					"  - 5\n" +
					"  - 1\n" +
					"  - 2\n" +
					"  - 3\n",
			},
		},
		{
			Desc: "appends YAML data in dst",
			StartFiles: map[string]any{
				"example.yml": "foo: [1,2,3]\n",
			},
			TaskData: map[string]any{
				"type":   "update",
				"dst":    "example.yml",
				"match":  "$.foo",
				"action": "append",
				"src":    []any{4, 5},
			},
			EndFiles: map[string]any{
				"example.yml": "foo:\n" +
					"  - 1\n" +
					"  - 2\n" +
					"  - 3\n" +
					"  - 4\n" +
					"  - 5\n",
			},
		},
		{
			Desc: "replaces YAML data in dst",
			StartFiles: map[string]any{
				"example.yml": "foo: [1,2,3]\n",
			},
			TaskData: map[string]any{
				"type":   "update",
				"dst":    "example.yml",
				"match":  "$.foo",
				"action": "replace",
				"src":    []any{4, 5},
			},
			EndFiles: map[string]any{
				"example.yml": "foo:\n" +
					"  - 4\n" +
					"  - 5\n",
			},
		},
		{
			Desc: "deletes YAML data in dst",
			StartFiles: map[string]any{
				"example.yml": "foo: [1,2,3]\n",
			},
			TaskData: map[string]any{
				"type":   "update",
				"dst":    "example.yml",
				"match":  "$.foo",
				"action": "delete",
			},
			EndFiles: map[string]any{
				"example.yml": "{}\n",
			},
		},

		{
			Desc: "[missing:ignore] ignores missing paths",
			TaskData: map[string]any{
				"type":  "update",
				"dst":   "./README.md",
				"match": "foo",
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
				"match":   "foo",
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
