package generate_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twelvelabs/stamp/internal/iostreams"
	"github.com/twelvelabs/stamp/internal/task"
	"github.com/twelvelabs/stamp/internal/task/common"
	"github.com/twelvelabs/stamp/internal/task/generate"
	"github.com/twelvelabs/stamp/internal/test_util"
	"github.com/twelvelabs/stamp/internal/value"
)

func TestGenerateTask(t *testing.T) {
	tests := []struct {
		Name   string
		Input  map[string]any
		Output interface{}
		Err    string
	}{
		{
			Name: "empty",
			Input: map[string]any{
				"type": "generate",
			},
			Output: nil,
			Err:    "Src is a required field, Dst is a required field",
		},
		{
			Name: "only src",
			Input: map[string]any{
				"type": "generate",
				"src":  "example.tpl",
			},
			Output: nil,
			Err:    "Dst is a required field",
		},
		{
			Name: "only dst",
			Input: map[string]any{
				"type": "generate",
				"dst":  "example.txt",
			},
			Output: nil,
			Err:    "Src is a required field",
		},
		{
			Name: "valid",
			Input: map[string]any{
				"type": "generate",
				"src":  "example.tpl",
				"dst":  "example.txt",
			},
			Output: &generate.Task{
				Common: common.Common{
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
			actual, err := task.NewTask(test.Input)

			assert.Equal(t, test.Output, actual)
			if test.Err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, test.Err)
			}
		})
	}
}

func TestExecute(t *testing.T) {
	templatesDir := filepath.Join("..", "..", "..", "testdata", "templates")
	tests := []struct {
		Desc       string
		TaskData   map[string]any
		Values     map[string]any
		Prompter   *value.PrompterMock
		StartFiles map[string]string
		EndFiles   map[string]string
		Err        string
	}{
		{
			Desc: "generates individual files",
			TaskData: map[string]any{
				"type": "generate",
				"src":  "{{ .SrcPath }}/README.md",
				"dst":  "{{ .DstPath }}/README.md",
			},
			Values: map[string]any{
				"ProjectName": "My Project",
				"SrcPath":     templatesDir,
				"DstPath":     "TBD",
			},
			EndFiles: map[string]string{
				"README.md": "# My Project\n",
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
				"DstPath":     "TBD",
			},
			EndFiles: map[string]string{
				"nested/README.md":  "# My Project\n",
				"nested/bin/aaa.sh": "#!/bin/bash\necho \"Hi from aaa in My Project\"\n",
				"nested/bin/bbb.sh": "#!/bin/bash\necho \"Hi from bbb in My Project\"\n",
				"nested/docs/":      "",
			},
			Err: "",
		},

		{
			Desc: "[conflict:keep] will keep without prompting",
			StartFiles: map[string]string{
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
				"DstPath":     "TBD",
			},
			Prompter: &value.PrompterMock{},
			EndFiles: map[string]string{
				"README.md": "Pre-existing content",
			},
			Err: "",
		},

		{
			Desc: "[conflict:replace] will replace without prompting",
			StartFiles: map[string]string{
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
				"DstPath":     "TBD",
			},
			Prompter: &value.PrompterMock{},
			EndFiles: map[string]string{
				"README.md": "# My Project\n",
			},
			Err: "",
		},

		{
			Desc: "[conflict:prompt] will prompt to overwrite and replace the file if user confirms",
			StartFiles: map[string]string{
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
				"DstPath":     "TBD",
			},
			Prompter: &value.PrompterMock{
				ConfirmFunc: value.NewConfirmFunc(true, nil),
			},
			EndFiles: map[string]string{
				"README.md": "# My Project\n",
			},
			Err: "",
		},
		{
			Desc: "[conflict:prompt] will prompt to overwrite and keep the file if user responds no",
			StartFiles: map[string]string{
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
				"DstPath":     "TBD",
			},
			Prompter: &value.PrompterMock{
				ConfirmFunc: value.NewConfirmFunc(false, nil),
			},
			EndFiles: map[string]string{
				"README.md": "Pre-existing content",
			},
			Err: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Desc, func(t *testing.T) {
			defer test_util.Cleanup()

			// Create a temp dir
			tmpDir := test_util.MkdirTemp(t)
			tt.Values["DstPath"] = tmpDir

			// Populate the temp dir w/ any pre-existing files needed for the run
			if tt.StartFiles != nil {
				for path, content := range tt.StartFiles {
					path = tmpDir + "/" + path
					_ = os.WriteFile(path, []byte(content), 0666)
				}
			}

			// Create a new task and execute it
			task, err := task.NewTask(tt.TaskData)
			assert.NoError(t, err)
			ios, _, _, _ := iostreams.Test()
			err = task.Execute(tt.Values, ios, tt.Prompter, false)

			// Ensure the expected files were generated
			if tt.EndFiles != nil {
				for path, content := range tt.EndFiles {
					path = tmpDir + "/" + path
					if strings.HasSuffix(path, "/") {
						assert.DirExists(t, path)
					} else {
						assert.FileExists(t, path)
						buf, _ := os.ReadFile(path)
						assert.Equal(t, content, string(buf))
					}
				}
			}

			if tt.Err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.Err)
			}
		})
	}
}
