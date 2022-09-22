package generate

//cspell:words oneof

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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

	// validate both paths now so we don't have to everywhere else
	if _, err := t.renderPath(values, t.Src); err != nil {
		return err
	}
	dst, err := t.renderPath(values, t.Dst)
	if err != nil {
		return err
	}

	if _, err := os.Stat(dst); os.IsNotExist(err) {
		return t.generate(values, ios, prompter)
	} else {
		switch t.Conflict {
		case ConflictPrompt:
			return t.prompt(values, ios, prompter)
		case ConflictKeep:
			return t.keep(values, ios, prompter)
		case ConflictReplace:
			return t.replace(values, ios, prompter)
		default:
			return fmt.Errorf("unknown conflict type: %v", t.Conflict)
		}
	}
}

func (t *Task) generate(values map[string]any, ios *iostreams.IOStreams, prompter value.Prompter) error {
	dst, _ := t.renderPath(values, t.Dst)
	if err := t.createDst(values, ios); err != nil {
		t.LogFailure(ios, "fail", dst)
		return err
	}
	t.LogSuccess(ios, "generate", dst)
	return nil
}

func (t *Task) keep(values map[string]any, ios *iostreams.IOStreams, prompter value.Prompter) error {
	dst, _ := t.renderPath(values, t.Dst)
	t.LogSuccess(ios, "keep", dst)
	return nil
}

func (t *Task) replace(values map[string]any, ios *iostreams.IOStreams, prompter value.Prompter) error {
	dst, _ := t.renderPath(values, t.Dst)
	if err := t.deleteDst(values, ios); err != nil {
		t.LogFailure(ios, "fail", dst)
		return err
	}
	if err := t.createDst(values, ios); err != nil {
		t.LogFailure(ios, "fail", dst)
		return err
	}
	t.LogSuccess(ios, "replace", dst)
	return nil
}

func (t *Task) prompt(values map[string]any, ios *iostreams.IOStreams, prompter value.Prompter) error {
	dst, _ := t.renderPath(values, t.Dst)
	t.LogWarning(ios, "conflict", fmt.Sprintf("%s already exists", dst))

	overwrite, err := prompter.Confirm("Overwrite", false, "", "")
	if err != nil {
		return err
	}

	if overwrite {
		return t.replace(values, ios, prompter)
	} else {
		return t.keep(values, ios, prompter)
	}
}

func (t *Task) createDst(values map[string]any, ios *iostreams.IOStreams) error {
	if t.DryRun {
		return nil
	}
	src, _ := t.renderPath(values, t.Src)
	dst, _ := t.renderPath(values, t.Dst)
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
func (t *Task) deleteDst(values map[string]any, ios *iostreams.IOStreams) error {
	if t.DryRun {
		return nil
	}
	dst, _ := t.renderPath(values, t.Dst)
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
