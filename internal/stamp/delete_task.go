package stamp

import (
	"fmt"
	"os"

	"github.com/twelvelabs/stamp/internal/fsutil"
)

type DeleteTask struct {
	Common `mapstructure:",squash"`

	Dst     string        `validate:"required"`
	Missing MissingConfig `validate:"required" default:"ignore"`
}

func (t *DeleteTask) Execute(ctx *TaskContext, values map[string]any) error {
	dst, err := t.RenderRequired("dst", t.Dst, values)
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
		return fmt.Errorf("path does not exist")
	}

	return nil
}

func (t *DeleteTask) deleteDst(ctx *TaskContext, dst string) error {
	if ctx.DryRun {
		return nil
	}
	return os.RemoveAll(dst)
}
