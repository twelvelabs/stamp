package generate

//cspell:words oneof

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/twelvelabs/stamp/internal/iostreams"
	"github.com/twelvelabs/stamp/internal/task/common"
	"github.com/twelvelabs/stamp/internal/value"
)

//go:generate go-enum -f=$GOFILE --marshal --names

// Conflict determines what to do when Dst path already exists.
// ENUM(keep, replace, prompt)
type Conflict string

const (
	DST_DIR_MODE os.FileMode = 0755
)

type Task struct {
	common.Common `mapstructure:",squash"`

	Src      string   `validate:"required"`
	Dst      string   `validate:"required"`
	Mode     string   `validate:"required,posix-mode" default:"0666"`
	Conflict Conflict `validate:"required,oneof=keep replace prompt" default:"prompt"`
}

func (t *Task) Execute(values map[string]any, ios *iostreams.IOStreams, prompter value.Prompter, dryRun bool) error {
	t.DryRun = dryRun

	src, err := t.renderPath(values, t.Src)
	if err != nil {
		return err
	}
	dst, err := t.renderPath(values, t.Dst)
	if err != nil {
		return err
	}

	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	if info.IsDir() {
		src = strings.TrimSuffix(src, "/")
		// src is a dir; walk and call dispatch on each file
		return filepath.Walk(src, func(srcPath string, srcPathInfo fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			dstPath := filepath.Join(dst, strings.TrimPrefix(srcPath, src))
			if srcPathInfo.IsDir() {
				return t.createDstDir(dstPath)
			} else {
				return t.dispatch(values, ios, prompter, srcPath, dstPath)
			}
		})
	} else {
		// src is a single file
		return t.dispatch(values, ios, prompter, src, dst)
	}
}

// dispatch looks for conflicts and delegates to the correct generation method.
func (t *Task) dispatch(values map[string]any, ios *iostreams.IOStreams, prompter value.Prompter, src string, dst string) error {
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		return t.generate(values, ios, prompter, src, dst)
	} else {
		switch t.Conflict {
		case ConflictPrompt:
			return t.prompt(values, ios, prompter, src, dst)
		case ConflictKeep:
			return t.keep(values, ios, prompter, src, dst)
		case ConflictReplace:
			return t.replace(values, ios, prompter, src, dst)
		default:
			return fmt.Errorf("unknown conflict type: %v", t.Conflict)
		}
	}
}

// generate is called to generate a non-existing dst file.
func (t *Task) generate(values map[string]any, ios *iostreams.IOStreams, prompter value.Prompter, src string, dst string) error {
	if err := t.createDst(values, src, dst); err != nil {
		t.LogFailure(ios, "fail", dst)
		return err
	}
	t.LogSuccess(ios, "generate", dst)
	return nil
}

// keep is called when keeping an existing dst file.
func (t *Task) keep(values map[string]any, ios *iostreams.IOStreams, prompter value.Prompter, src string, dst string) error {
	t.LogSuccess(ios, "keep", dst)
	return nil
}

// replace is called when replacing an existing dst file.
func (t *Task) replace(values map[string]any, ios *iostreams.IOStreams, prompter value.Prompter, src string, dst string) error {
	if err := t.deleteDst(dst); err != nil {
		t.LogFailure(ios, "fail", dst)
		return err
	}
	if err := t.createDst(values, src, dst); err != nil {
		t.LogFailure(ios, "fail", dst)
		return err
	}
	t.LogSuccess(ios, "replace", dst)
	return nil
}

// prompt is called to prompt the user for how to resolve a dst file conflict.
// delegates to keep or replace depending on their response.
func (t *Task) prompt(values map[string]any, ios *iostreams.IOStreams, prompter value.Prompter, src string, dst string) error {
	t.LogWarning(ios, "conflict", fmt.Sprintf("%s already exists", dst))
	overwrite, err := prompter.Confirm("Overwrite", false, "", "")
	if err != nil {
		return err
	}
	if overwrite {
		return t.replace(values, ios, prompter, src, dst)
	} else {
		return t.keep(values, ios, prompter, src, dst)
	}
}

func (t *Task) createDstDir(dst string) error {
	if t.DryRun {
		return nil
	}
	return os.MkdirAll(dst, DST_DIR_MODE)
}

func (t *Task) createDst(values map[string]any, src string, dst string) error {
	if t.DryRun {
		return nil
	}
	mode, err := t.renderMode(values, t.Mode)
	if err != nil {
		return err
	}

	// parse the template
	template, err := template.ParseFiles(src)
	if err != nil {
		return err
	}

	// create base dst dirs
	if err := os.MkdirAll(filepath.Dir(dst), DST_DIR_MODE); err != nil {
		return err
	}

	// create dst
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()

	// set perms
	err = os.Chmod(dst, mode)
	if err != nil {
		return err
	}

	// render template
	err = template.Execute(f, values)
	if err != nil {
		return err
	}

	return nil
}
func (t *Task) deleteDst(dst string) error {
	if t.DryRun {
		return nil
	}
	if err := os.Remove(dst); err != nil {
		return err
	}
	return nil
}

func (t *Task) renderPath(values map[string]any, path string) (string, error) {
	rendered := t.Common.Render(path, values)
	if rendered == "" {
		return "", fmt.Errorf("src or dst path '%s' evaluated to an empty string", path)
	}
	return rendered, nil
}

func (t *Task) renderMode(values map[string]any, mode string) (os.FileMode, error) {
	rendered := t.Common.Render(mode, values)
	if rendered == "" {
		return 0, fmt.Errorf("mode '%s' evaluated to an empty string", mode)
	}
	parsed, err := strconv.ParseInt(rendered, 8, 64)
	if err != nil {
		return 0, err
	}
	return os.FileMode(parsed), nil
}
