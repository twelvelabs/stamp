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
		Err      string
	}{
		{
			Name: "returns an error when missing field is invalid",
			TaskData: map[string]any{
				"type": "update",
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
				"type": "update",
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
			Desc: "returns an error if dst evaluates to empty string",
			TaskData: map[string]any{
				"type": "update",
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
			Desc: "returns an error if src fails to render",
			TaskData: map[string]any{
				"type": "update",
				"dst": map[string]any{
					"path": "README.md",
				},
				"src": map[string]any{
					"path": `{{ fail "boom" }}`,
				},
			},
			Err: "boom",
		},
		{
			Desc: "returns an error if description fails to render",
			TaskData: map[string]any{
				"type": "update",
				"dst": map[string]any{
					"path": "README.md",
				},
				"description": `{{ fail "boom" }}`,
			},
			Err: "boom",
		},
		{
			Desc: "returns an error if pattern can not be compiled to regexp",
			StartFiles: map[string]any{
				"README.md": "Hello World\n",
			},
			TaskData: map[string]any{
				"type": "update",
				"dst": map[string]any{
					"path": "README.md",
				},
				"match": map[string]any{
					"pattern": "(.}",
				},
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
				"type": "update",
				"dst": map[string]any{
					"path": "README.md",
				},
				"match": map[string]any{
					"pattern": "World",
				},
				"src": map[string]any{
					"content": struct{}{},
				},
			},
			EndFiles: map[string]any{
				"README.md": "Hello World\n",
			},
			Err: "unable to cast",
		},
		{
			Desc: "returns an error if dst can not be touched",
			TaskData: map[string]any{
				"type": "update",
				"dst": map[string]any{
					"path":    "invalid\x00file.txt",
					"missing": "touch",
				},
			},
			Err: "invalid argument",
		},

		{
			Desc: "prepends a string in dst",
			StartFiles: map[string]any{
				"README.md": "aaa bbb ccc\n",
			},
			TaskData: map[string]any{
				"type": "update",
				"dst": map[string]any{
					"path": "README.md",
				},
				"match": map[string]any{
					"pattern": "aaa",
				},
				"action": map[string]any{
					"type": "prepend",
				},
				"src": map[string]any{
					"content": "000 ",
				},
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
				"type": "update",
				"dst": map[string]any{
					"path": "README.md",
				},
				"match": map[string]any{
					"pattern": "ccc",
				},
				"action": map[string]any{
					"type": "append",
				},
				"src": map[string]any{
					"content": " ddd",
				},
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
				"type": "update",
				"dst": map[string]any{
					"path": "README.md",
				},
				"match": map[string]any{
					"pattern": "Hello (\\w+)",
				},
				"action": map[string]any{
					"type": "replace",
				},
				"src": map[string]any{
					"content": "Goodbye $1",
				},
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
				"type": "update",
				"dst": map[string]any{
					"path": "README.md",
				},
				"match": map[string]any{
					"pattern": "(?m)\\s*(\\w+)$",
				},
				"action": map[string]any{
					"type": "delete",
				},
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
				"type": "update",
				"dst": map[string]any{
					"path": "README.md",
				},
				"action": map[string]any{
					"type": "append",
				},
				"src": map[string]any{
					"content": "Goodbye\n",
				},
			},
			EndFiles: map[string]any{
				"README.md": "Hello World\nGoodbye\n",
			},
		},
		{
			Desc: "creates and updates a path when missing set to touch",
			TaskData: map[string]any{
				"type": "update",
				"dst": map[string]any{
					"path":    "README.md",
					"missing": "touch",
				},
				"src": map[string]any{
					"content": "Howdy",
				},
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
				"type": "update",
				"dst": map[string]any{
					"path": "README.md",
					"mode": "0755",
				},
				"match": map[string]any{
					"pattern": "Hello (\\w+)",
				},
				"action": map[string]any{
					"type": "replace",
				},
				"src": map[string]any{
					"content": "Goodbye $1",
				},
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
				"type": "update",
				"dst": map[string]any{
					"path": "README.md",
				},
				"match": map[string]any{
					"pattern": "Hello (\\w+)",
				},
				"action": map[string]any{
					"type": "replace",
				},
				"src": map[string]any{
					"content": "Goodbye $1",
				},
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
				"type": "update",
				"dst": map[string]any{
					"path": "README.md",
				},
				"action": map[string]any{
					"type": "replace",
				},
				"src": map[string]any{
					"path": "valid.txt",
				},
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
				"dst": map[string]any{
					"path": "README.md",
				},
				"match": map[string]any{
					"pattern": `^the (\w+) line might be (\w+)$`,
				},
				"src": map[string]any{
					"content": "the ${1} line IS ${2}",
				},
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
				"dst": map[string]any{
					"path": "README.md",
				},
				"match": map[string]any{
					"pattern": `(?ms) second(.+)baz$`,
					"source":  "file",
				},
				"action": map[string]any{
					"type": "delete",
				},
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
				"type": "update",
				"dst": map[string]any{
					"path": "example.json",
				},
				"match": map[string]any{
					"pattern": "$.foo",
				},
				"action": map[string]any{
					"type": "prepend",
				},
				"src": map[string]any{
					"content": []any{4, 5},
				},
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
				"type": "update",
				"dst": map[string]any{
					"path": "example.json",
				},
				"match": map[string]any{
					"pattern": "$.foo",
				},
				"action": map[string]any{
					"type": "append",
				},
				"src": map[string]any{
					"content": []any{4, 5},
				},
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
				"type": "update",
				"dst": map[string]any{
					"path": "example.json",
				},
				"match": map[string]any{
					"pattern": "$.foo",
				},
				"action": map[string]any{
					"type": "replace",
				},
				"src": map[string]any{
					"content": []any{4, 5},
				},
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
				"type": "update",
				"dst": map[string]any{
					"path": "example.json",
				},
				"match": map[string]any{
					"pattern": "$.foo",
				},
				"action": map[string]any{
					"type": "delete",
				},
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
				"type": "update",
				"dst": map[string]any{
					"path": "example.json",
				},
				"action": map[string]any{
					"type": "replace",
				},
				"src": map[string]any{
					"content": map[string]any{
						"bar": true,
					},
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
				"type": "update",
				"dst": map[string]any{
					"path": "example.json",
				},
				"action": map[string]any{
					"type": "append",
				},
				"match": map[string]any{
					"pattern": "$.foo",
				},
				"src": map[string]any{
					"path": "valid.json",
				},
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
				"type": "update",
				"dst": map[string]any{
					"path": "example.json",
				},
				"action": map[string]any{
					"type": "append",
				},
				"match": map[string]any{
					"pattern": "$.foo",
					"default": []any{},
				},
				"src": map[string]any{
					"path": "valid.json",
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
				"type": "update",
				"dst": map[string]any{
					"path": "example.json",
				},
				"match": map[string]any{
					"pattern": "$.foo",
				},
				"action": map[string]any{
					"type":  "append",
					"merge": "upsert",
				},
				"src": map[string]any{
					"content": []any{
						"aaa",
						"bbb",
						"ccc",
					},
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
				"type": "update",
				"dst": map[string]any{
					"path": "example.yml",
				},
				"match": map[string]any{
					"pattern": "$.foo",
				},
				"action": map[string]any{
					"type": "prepend",
				},
				"src": map[string]any{
					"content": []any{4, 5},
				},
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
				"type": "update",
				"dst": map[string]any{
					"path": "example.yml",
				},
				"match": map[string]any{
					"pattern": "$.foo",
				},
				"action": map[string]any{
					"type": "append",
				},
				"src": map[string]any{
					"content": []any{4, 5},
				},
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
				"type": "update",
				"dst": map[string]any{
					"path": "example.yml",
				},
				"match": map[string]any{
					"pattern": "$.foo",
				},
				"action": map[string]any{
					"type": "replace",
				},
				"src": map[string]any{
					"content": []any{4, 5},
				},
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
				"type": "update",
				"dst": map[string]any{
					"path": "example.yml",
				},
				"match": map[string]any{
					"pattern": "$.foo",
				},
				"action": map[string]any{
					"type": "delete",
				},
			},
			EndFiles: map[string]any{
				"example.yml": "{}\n",
			},
		},

		{
			Desc: "[missing:ignore] ignores missing paths",
			TaskData: map[string]any{
				"type": "update",
				"dst": map[string]any{
					"path": "README.md",
				},
				"match": map[string]any{
					"pattern": "foo",
				},
			},
			EndFiles: map[string]any{
				"README.md": false,
			},
		},
		{
			Desc: "[missing:error] returns an error when path is missing",
			TaskData: map[string]any{
				"type": "update",
				"dst": map[string]any{
					"path":    "README.md",
					"missing": "error",
				},
				"match": map[string]any{
					"pattern": "foo",
				},
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
