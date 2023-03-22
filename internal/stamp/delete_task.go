package stamp

import (
	"os"

	"github.com/spf13/cast"

	"github.com/twelvelabs/stamp/internal/fsutil"
)

type DeleteTask struct {
	Common `mapstructure:",squash"`

	Dst     string        `validate:"required"`
	Missing MissingConfig `validate:"required" default:"ignore"`
}

func (t *DeleteTask) Execute(ctx *TaskContext, values map[string]any) error {
	dstRoot := cast.ToString(values["DstPath"])
	dst, err := t.RenderPath("dst", t.Dst, dstRoot, values)
	if err != nil {
		return err
	}

	if fsutil.PathExists(dst) {
		if err := t.deleteDst(ctx, dst); err != nil {
			ctx.Logger.Failure("fail", dst)
			return err
		}
		ctx.Logger.Success("delete", dst)
		return nil
	} else if t.Missing == MissingConfigError {
		ctx.Logger.Failure("fail", dst)
		return ErrPathNotFound
	}

	return nil
}

func (t *DeleteTask) deleteDst(ctx *TaskContext, dst string) error {
	if ctx.DryRun {
		return nil
	}
	return os.RemoveAll(dst)
}
