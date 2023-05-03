package stamp

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cast"
	"github.com/twelvelabs/termite/render"

	"github.com/twelvelabs/stamp/internal/fsutil"
)

const (
	DstDirMode os.FileMode = 0755
)

type CreateTask struct {
	Common `mapstructure:",squash"`

	Src      string         `validate:"required"`
	Dst      string         `validate:"required"`
	Mode     string         `validate:"required,posix-mode" default:"0666"`
	Conflict ConflictConfig `validate:"required" default:"prompt"`
}

func (t *CreateTask) Execute(ctx *TaskContext, values map[string]any) error {
	t.DryRun = ctx.DryRun

	srcRoot := cast.ToString(values["SrcPath"])
	src, err := t.RenderPath("src", t.Src, srcRoot, values)
	if err != nil {
		return err
	}
	dstRoot := cast.ToString(values["DstPath"])
	dst, err := t.RenderPath("dst", t.Dst, dstRoot, values)
	if err != nil {
		return err
	}

	info, _ := os.Stat(src)
	if info != nil && info.IsDir() {
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

	// src is a single file (or inline content)
	return t.dispatch(ctx, values, src, dst)
}

// dispatch looks for conflicts and delegates to the correct generation method.
func (t *CreateTask) dispatch(ctx *TaskContext, values map[string]any, src string, dst string) error {
	if _, err := os.Stat(dst); os.IsNotExist(err) && dst != "" {
		return t.create(ctx, values, src, dst)
	}
	switch t.Conflict {
	case ConflictConfigPrompt:
		return t.prompt(ctx, values, src, dst)
	case ConflictConfigKeep:
		return t.keep(ctx, values, src, dst)
	case ConflictConfigReplace:
		return t.replace(ctx, values, src, dst)
	default:
		return fmt.Errorf("unknown conflict type: %v", t.Conflict)
	}
}

// create is called to create a non-existing dst file.
func (t *CreateTask) create(ctx *TaskContext, values map[string]any, src string, dst string) error {
	if err := t.createDst(values, src, dst); err != nil {
		ctx.Logger.Failure("fail", dst)
		return err
	}
	ctx.Logger.Success("create", dst)
	return nil
}

// keep is called when keeping an existing dst file.
func (t *CreateTask) keep(ctx *TaskContext, _ map[string]any, _ string, dst string) error {
	ctx.Logger.Success("keep", dst)
	return nil
}

// replace is called when replacing an existing dst file.
func (t *CreateTask) replace(ctx *TaskContext, values map[string]any, src string, dst string) error {
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
func (t *CreateTask) prompt(ctx *TaskContext, values map[string]any, src string, dst string) error {
	ctx.Logger.Warning("conflict", "%s already exists", dst)
	overwrite, err := ctx.UI.Confirm("Overwrite", false)
	if err != nil {
		return err
	}
	if overwrite {
		return t.replace(ctx, values, src, dst)
	}
	return t.keep(ctx, values, src, dst)
}

func (t *CreateTask) createDstDir(dst string) error {
	if t.DryRun {
		return nil
	}
	return os.MkdirAll(dst, DstDirMode)
}

func (t *CreateTask) createDst(values map[string]any, src string, dst string) error {
	if t.DryRun {
		return nil
	}

	// render and parse mode
	mode, err := t.RenderMode(t.Mode, values)
	if err != nil {
		return err
	}

	var rendered string
	// Src can be:
	// - A path to a template file containing the source content.
	// - An inline string literal to render and use as source content.
	if fsutil.PathExists(src) {
		rendered, err = render.File(src, values)
		if err != nil {
			return err
		}
	} else {
		rendered = t.Render(t.Src, values)
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

func (t *CreateTask) deleteDst(dst string) error {
	if t.DryRun {
		return nil
	}
	return os.Remove(dst)
}
