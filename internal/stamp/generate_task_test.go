package stamp

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twelvelabs/termite/testutil"

	"github.com/twelvelabs/stamp/internal/value"
)

func TestNewTask_WhenTypeIsGenerate(t *testing.T) {
	tests := []struct {
		Name     string
		TaskData map[string]any
		Task     interface{}
		Err      string
	}{
		{
			Name: "returns an error when both src and dst are missing",
			TaskData: map[string]any{
				"type": "generate",
			},
			Task: nil,
			Err:  "Src is a required field, Dst is a required field",
		},
		{
			Name: "returns an error when dst is missing",
			TaskData: map[string]any{
				"type": "generate",
				"src":  "example.tpl",
			},
			Task: nil,
			Err:  "Dst is a required field",
		},
		{
			Name: "returns an error when src is missing",
			TaskData: map[string]any{
				"type": "generate",
				"dst":  "example.txt",
			},
			Task: nil,
			Err:  "Src is a required field",
		},
		{
			Name: "returns an error when mode is invalid",
			TaskData: map[string]any{
				"type": "generate",
				"src":  "example.tpl",
				"dst":  "example.txt",
				"mode": "not a posix-mode",
			},
			Task: nil,
			Err:  "Mode must be a valid posix file mode",
		},
		{
			Name: "returns an error when conflict is invalid",
			TaskData: map[string]any{
				"type":     "generate",
				"src":      "example.tpl",
				"dst":      "example.txt",
				"conflict": "unknown",
			},
			Task: nil,
			Err:  "unknown is not a valid Conflict",
		},
		{
			Name: "returns the task when all fields are valid",
			TaskData: map[string]any{
				"type": "generate",
				"src":  "example.tpl",
				"dst":  "example.txt",
			},
			Task: &GenerateTask{
				Common: Common{
					If:   "true",
					Each: "",
				},
				Src:      "example.tpl",
				Dst:      "example.txt",
				Conflict: "prompt",
				Mode:     "0666",
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

func TestGenerateTask_Execute(t *testing.T) { //nolint:maintidx
	templatesDir, _ := filepath.Abs(filepath.Join("testdata", "templates"))
	tests := []struct {
		Desc       string
		DryRun     bool
		TaskData   map[string]any
		Values     map[string]any
		Prompter   *value.PrompterMock
		StartFiles map[string]any
		EndFiles   map[string]any
		Err        string
	}{
		{
			Desc: "returns an error if src evaluates to empty string",
			TaskData: map[string]any{
				"type": "generate",
				"src":  "{{ .Empty }}",
				"dst":  "{{ .DstPath }}/README.md",
			},
			Values: map[string]any{
				"DstPath": ".",
				"Empty":   "",
			},
			Err: "path '{{ .Empty }}' evaluated to an empty string",
		},
		{
			Desc: "returns an error if dst evaluates to empty string",
			TaskData: map[string]any{
				"type": "generate",
				"src":  "{{ .SrcPath }}/README.md",
				"dst":  "{{ .Empty }}",
			},
			Values: map[string]any{
				"SrcPath": templatesDir,
				"Empty":   "",
			},
			Err: "path '{{ .Empty }}' evaluated to an empty string",
		},
		{
			Desc: "returns an error if src does not exist",
			TaskData: map[string]any{
				"type": "generate",
				"src":  "{{ .SrcPath }}/missing.md",
				"dst":  "{{ .DstPath }}/README.md",
			},
			Values: map[string]any{
				"SrcPath": templatesDir,
				"DstPath": ".",
			},
			Err: "missing.md: no such file or directory",
		},

		{
			Desc: "generates a single file",
			TaskData: map[string]any{
				"type": "generate",
				"src":  "{{ .SrcPath }}/README.md",
				"dst":  "{{ .DstPath }}/README.md",
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
				"type": "generate",
				"src":  "{{ .SrcPath }}/README.md",
				"dst":  "{{ .DstPath }}/README.md",
				"mode": "0755",
			},
			Values: map[string]any{
				"ProjectName": "My Project",
				"SrcPath":     templatesDir,
				"DstPath":     ".",
			},
			EndFiles: map[string]any{
				"README.md": 0o755,
			},
			Err: "",
		},
		{
			Desc:   "does not generate a file during a dry run",
			DryRun: true,
			TaskData: map[string]any{
				"type": "generate",
				"src":  "{{ .SrcPath }}/README.md",
				"dst":  "{{ .DstPath }}/README.md",
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
				"type": "generate",
				"src":  "{{ .SrcPath }}/nested/",
				"dst":  "{{ .DstPath }}/nested/",
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
			Desc:   "does not generate files or directories during a dry run",
			DryRun: true,
			TaskData: map[string]any{
				"type": "generate",
				"src":  "{{ .SrcPath }}/nested/",
				"dst":  "{{ .DstPath }}/nested/",
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
				"type":     "generate",
				"src":      "{{ .SrcPath }}/README.md",
				"dst":      "{{ .DstPath }}/README.md",
				"conflict": "keep",
			},
			Values: map[string]any{
				"ProjectName": "My Project",
				"SrcPath":     templatesDir,
				"DstPath":     ".",
			},
			Prompter: &value.PrompterMock{},
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
				"type":     "generate",
				"src":      "{{ .SrcPath }}/README.md",
				"dst":      "{{ .DstPath }}/README.md",
				"conflict": "replace",
			},
			Values: map[string]any{
				"ProjectName": "My Project",
				"SrcPath":     templatesDir,
				"DstPath":     ".",
			},
			Prompter: &value.PrompterMock{},
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
				"type":     "generate",
				"src":      "{{ .SrcPath }}/README.md",
				"dst":      "{{ .DstPath }}/README.md",
				"conflict": "prompt",
			},
			Values: map[string]any{
				"ProjectName": "My Project",
				"SrcPath":     templatesDir,
				"DstPath":     ".",
			},
			Prompter: &value.PrompterMock{
				ConfirmFunc: value.NewConfirmFunc(true, nil),
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
				"type":     "generate",
				"src":      "{{ .SrcPath }}/README.md",
				"dst":      "{{ .DstPath }}/README.md",
				"conflict": "prompt",
			},
			Values: map[string]any{
				"ProjectName": "My Project",
				"SrcPath":     templatesDir,
				"DstPath":     ".",
			},
			Prompter: &value.PrompterMock{
				ConfirmFunc: value.NewConfirmFunc(false, nil),
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
				"type":     "generate",
				"src":      "{{ .SrcPath }}/README.md",
				"dst":      "{{ .DstPath }}/README.md",
				"conflict": "prompt",
			},
			Values: map[string]any{
				"ProjectName": "My Project",
				"SrcPath":     templatesDir,
				"DstPath":     ".",
			},
			Prompter: &value.PrompterMock{
				ConfirmFunc: value.NewConfirmFunc(false, errors.New("boom")),
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

				// Create a new task and execute it
				task, err := NewTask(tt.TaskData)
				assert.NoError(t, err)

				app := NewTestApp()
				app.Prompter = tt.Prompter

				ctx := NewTaskContext(app, tt.DryRun)
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

func TestGenerateTask_DispatchErrorsOnInvalidConflict(t *testing.T) {
	// invalid conflicts should always be caught by `NewTask`,
	// but testing here for full coverage.
	task := &GenerateTask{
		Conflict: "unknown",
	}
	app := NewTestApp()
	ctx := NewTaskContext(app, false)
	values := map[string]any{}

	err := task.dispatch(ctx, values, "", "")
	assert.ErrorContains(t, err, "unknown conflict type")
}

func TestGenerateTask_DispatchErrorsOnInvalidMode(t *testing.T) {
	// invalid modes should always be caught by `NewTask`,
	// but testing here for full coverage.
	task := &GenerateTask{
		Mode: "unknown",
	}
	app := NewTestApp()
	ctx := NewTaskContext(app, false)
	values := map[string]any{}

	err := task.dispatch(ctx, values, "", "/do-not-generate")
	assert.NoFileExists(t, "/do-not-generate")
	assert.ErrorContains(t, err, "invalid syntax")
}
