package gen

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/twelvelabs/termite/render"
)

const (
	DstDirMode os.FileMode = 0755
)

type GenerateTask struct {
	Common `mapstructure:",squash"`

	Src      string   `validate:"required"`
	Dst      string   `validate:"required"`
	Mode     string   `validate:"required,posix-mode" default:"0666"`
	Conflict Conflict `validate:"required" default:"prompt"`
}

func (t *GenerateTask) Execute(ctx *TaskContext, values map[string]any) error {
	t.DryRun = ctx.DryRun

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
			}
			return t.dispatch(ctx, values, srcPath, dstPath)
		})
	}
	// src is a single file
	return t.dispatch(ctx, values, src, dst)
}

// dispatch looks for conflicts and delegates to the correct generation method.
func (t *GenerateTask) dispatch(ctx *TaskContext, values map[string]any, src string, dst string) error {
	if _, err := os.Stat(dst); os.IsNotExist(err) && dst != "" {
		return t.generate(ctx, values, src, dst)
	}
	switch t.Conflict {
	case ConflictPrompt:
		return t.prompt(ctx, values, src, dst)
	case ConflictKeep:
		return t.keep(ctx, values, src, dst)
	case ConflictReplace:
		return t.replace(ctx, values, src, dst)
	default:
		return fmt.Errorf("unknown conflict type: %v", t.Conflict)
	}
}

// generate is called to generate a non-existing dst file.
func (t *GenerateTask) generate(ctx *TaskContext, values map[string]any, src string, dst string) error {
	if err := t.createDst(values, src, dst); err != nil {
		ctx.Logger.Failure("fail", dst)
		return err
	}
	ctx.Logger.Success("generate", dst)
	return nil
}

// keep is called when keeping an existing dst file.
func (t *GenerateTask) keep(ctx *TaskContext, _ map[string]any, _ string, dst string) error {
	ctx.Logger.Success("keep", dst)
	return nil
}

// replace is called when replacing an existing dst file.
func (t *GenerateTask) replace(ctx *TaskContext, values map[string]any, src string, dst string) error {
	if err := t.deleteDst(dst); err != nil {
		ctx.Logger.Failure("fail", dst)
		return err
	}
	if err := t.createDst(values, src, dst); err != nil {
		ctx.Logger.Failure("fail", dst)
		return err
	}
	ctx.Logger.Success("replace", dst)
	return nil
}

// prompt is called to prompt the user for how to resolve a dst file conflict.
// delegates to keep or replace depending on their response.
func (t *GenerateTask) prompt(ctx *TaskContext, values map[string]any, src string, dst string) error {
	ctx.Logger.Warning("conflict", "%s already exists", dst)
	overwrite, err := ctx.Prompter.Confirm("Overwrite", false, "", "")
	if err != nil {
		return err
	}
	if overwrite {
		return t.replace(ctx, values, src, dst)
	}
	return t.keep(ctx, values, src, dst)
}

func (t *GenerateTask) createDstDir(dst string) error {
	if t.DryRun {
		return nil
	}
	return os.MkdirAll(dst, DstDirMode)
}

func (t *GenerateTask) createDst(values map[string]any, src string, dst string) error {
	if t.DryRun {
		return nil
	}
	mode, err := t.parseMode(values, t.Mode)
	if err != nil {
		return err
	}

	// render the src template
	rendered, err := render.File(src, values)
	if err != nil {
		return err
	}

	// create base dst dirs
	if err := os.MkdirAll(filepath.Dir(dst), DstDirMode); err != nil {
		return err
	}

	// create dst
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(rendered)
	if err != nil {
		return err
	}

	// set perms
	err = os.Chmod(dst, mode)
	if err != nil {
		return err
	}

	return nil
}

func (t *GenerateTask) deleteDst(dst string) error {
	if t.DryRun {
		return nil
	}
	if err := os.Remove(dst); err != nil {
		return err
	}
	return nil
}

func (t *GenerateTask) renderPath(values map[string]any, path string) (string, error) {
	rendered := t.Common.Render(path, values)
	if rendered == "" {
		return "", fmt.Errorf("path '%s' evaluated to an empty string", path)
	}
	return rendered, nil
}

func (t *GenerateTask) parseMode(_ map[string]any, mode string) (os.FileMode, error) {
	parsed, err := strconv.ParseInt(mode, 8, 64)
	if err != nil {
		return 0, err
	}
	return os.FileMode(parsed), nil
}
