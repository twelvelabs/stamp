package stamp

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twelvelabs/termite/testutil"
	"github.com/twelvelabs/termite/ui"
)

func TestNewTask_WhenTypeIsCreate(t *testing.T) {
	tests := []struct {
		Name     string
		TaskData map[string]any
		Err      string
	}{
		{
			Name: "returns an error when conflict is invalid",
			TaskData: map[string]any{
				"type": "create",
				"src": map[string]any{
					"path": "example.tpl",
				},
				"dst": map[string]any{
					"path":     "example.txt",
					"conflict": "unknown",
				},
			},
			Err: "unknown is not a valid Conflict",
		},
		{
			Name: "returns the task when all fields are valid",
			TaskData: map[string]any{
				"type": "create",
				"src": map[string]any{
					"path": "example.tpl",
				},
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

func TestCreateTask_Execute(t *testing.T) { //nolint:maintidx
	templatesDir, _ := filepath.Abs(filepath.Join("testdata", "templates"))
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
				"type": "create",
				"src": map[string]any{
					"path": "README.md",
				},
				"dst": map[string]any{
					"path": "{{ .Empty }}",
				},
			},
			Values: map[string]any{
				"SrcPath": templatesDir,
				"Empty":   "",
			},
			Err: "evaluated to an empty string",
		},

		{
			Desc: "generates a single file",
			TaskData: map[string]any{
				"type": "create",
				"src": map[string]any{
					"path": "README.md",
				},
				"dst": map[string]any{
					"path": "README.md",
				},
			},
			Values: map[string]any{
				"ProjectName": "My Project",
				"SrcPath":     templatesDir,
				"DstPath":     ".",
			},
			EndFiles: map[string]any{
				"README.md": "# My Project\n",
			},
			Err: "",
		},
		{
			Desc: "generates a single file with custom permissions",
			TaskData: map[string]any{
				"type": "create",
				"src": map[string]any{
					"path": "README.md",
				},
				"dst": map[string]any{
					"path": "README.md",
					"mode": "0755",
				},
			},
			Values: map[string]any{
				"ProjectName": "My Project",
				"SrcPath":     templatesDir,
				"DstPath":     ".",
			},
			EndFiles: map[string]any{
				"README.md": []any{
					"# My Project\n",
					0o755,
				},
			},
			Err: "",
		},
		{
			Desc: "generates a single file from inline content",
			TaskData: map[string]any{
				"type": "create",
				"src": map[string]any{
					"content": "Hello!",
				},
				"dst": map[string]any{
					"path": "README.md",
				},
			},
			Values: map[string]any{
				"SrcPath": templatesDir,
				"DstPath": ".",
			},
			EndFiles: map[string]any{
				"README.md": "Hello!",
			},
			Err: "",
		},
		{
			Desc:   "does not create a file during a dry run",
			DryRun: true,
			TaskData: map[string]any{
				"type": "create",
				"src": map[string]any{
					"path": "README.md",
				},
				"dst": map[string]any{
					"path": "README.md",
				},
			},
			Values: map[string]any{
				"ProjectName": "My Project",
				"SrcPath":     templatesDir,
				"DstPath":     ".",
			},
			EndFiles: map[string]any{
				"README.md": false,
			},
			Err: "",
		},

		{
			Desc: "generates entire directories of files",
			TaskData: map[string]any{
				"type": "create",
				"src": map[string]any{
					"path": "nested/",
				},
				"dst": map[string]any{
					"path": "nested/",
				},
			},
			Values: map[string]any{
				"ProjectName": "My Project",
				"SrcPath":     templatesDir,
				"DstPath":     ".",
			},
			EndFiles: map[string]any{
				"nested/README.md":  "# My Project\n",
				"nested/bin/aaa.sh": "#!/bin/bash\necho \"Hi from aaa in My Project\"\n",
				"nested/bin/bbb.sh": "#!/bin/bash\necho \"Hi from bbb in My Project\"\n",
				"nested/docs/":      "",
			},
			Err: "",
		},
		{
			Desc:   "does not create files or directories during a dry run",
			DryRun: true,
			TaskData: map[string]any{
				"type": "create",
				"src": map[string]any{
					"path": "nested/",
				},
				"dst": map[string]any{
					"path": "nested/",
				},
			},
			Values: map[string]any{
				"ProjectName": "My Project",
				"SrcPath":     templatesDir,
				"DstPath":     ".",
			},
			EndFiles: map[string]any{
				"nested/README.md":  false,
				"nested/bin/aaa.sh": false,
				"nested/bin/bbb.sh": false,
				"nested/docs/":      false,
			},
			Err: "",
		},

		{
			Desc: "[conflict:keep] will keep without prompting",
			StartFiles: map[string]any{
				"README.md": "Pre-existing content",
			},
			TaskData: map[string]any{
				"type": "create",
				"src": map[string]any{
					"path": "README.md",
				},
				"dst": map[string]any{
					"path":     "README.md",
					"conflict": "keep",
				},
			},
			Values: map[string]any{
				"ProjectName": "My Project",
				"SrcPath":     templatesDir,
				"DstPath":     ".",
			},
			EndFiles: map[string]any{
				"README.md": "Pre-existing content",
			},
			Err: "",
		},

		{
			Desc: "[conflict:replace] will replace without prompting",
			StartFiles: map[string]any{
				"README.md": "Pre-existing content",
			},
			TaskData: map[string]any{
				"type": "create",
				"src": map[string]any{
					"path": "README.md",
				},
				"dst": map[string]any{
					"path":     "README.md",
					"conflict": "replace",
				},
			},
			Values: map[string]any{
				"ProjectName": "My Project",
				"SrcPath":     templatesDir,
				"DstPath":     ".",
			},
			EndFiles: map[string]any{
				"README.md": "# My Project\n",
			},
			Err: "",
		},

		{
			Desc: "[conflict:prompt] will prompt to overwrite and replace the file if user confirms",
			StartFiles: map[string]any{
				"README.md": "Pre-existing content",
			},
			TaskData: map[string]any{
				"type": "create",
				"src": map[string]any{
					"path": "README.md",
				},
				"dst": map[string]any{
					"path":     "README.md",
					"conflict": "prompt",
				},
			},
			Values: map[string]any{
				"ProjectName": "My Project",
				"SrcPath":     templatesDir,
				"DstPath":     ".",
			},
			Setup: func(app *App) {
				app.UI.RegisterStub(
					ui.MatchConfirm("Overwrite"),
					ui.RespondBool(true),
				)
			},
			EndFiles: map[string]any{
				"README.md": "# My Project\n",
			},
			Err: "",
		},
		{
			Desc: "[conflict:prompt] will prompt to overwrite and keep the file if user responds no",
			StartFiles: map[string]any{
				"README.md": "Pre-existing content",
			},
			TaskData: map[string]any{
				"type": "create",
				"src": map[string]any{
					"path": "README.md",
				},
				"dst": map[string]any{
					"path":     "README.md",
					"conflict": "prompt",
				},
			},
			Values: map[string]any{
				"ProjectName": "My Project",
				"SrcPath":     templatesDir,
				"DstPath":     ".",
			},
			Setup: func(app *App) {
				app.UI.RegisterStub(
					ui.MatchConfirm("Overwrite"),
					ui.RespondBool(false),
				)
			},
			EndFiles: map[string]any{
				"README.md": "Pre-existing content",
			},
			Err: "",
		},
		{
			Desc: "[conflict:prompt] will return any prompter errors",
			StartFiles: map[string]any{
				"README.md": "Pre-existing content",
			},
			TaskData: map[string]any{
				"type": "create",
				"src": map[string]any{
					"path": "README.md",
				},
				"dst": map[string]any{
					"path":     "README.md",
					"conflict": "prompt",
				},
			},
			Values: map[string]any{
				"ProjectName": "My Project",
				"SrcPath":     templatesDir,
				"DstPath":     ".",
			},
			Setup: func(app *App) {
				app.UI.RegisterStub(
					ui.MatchConfirm("Overwrite"),
					ui.RespondError(errors.New("boom")),
				)
			},
			EndFiles: map[string]any{
				"README.md": "Pre-existing content",
			},
			Err: "boom",
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
