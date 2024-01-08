package stamp

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/swaggest/jsonschema-go"

	"github.com/twelvelabs/stamp/internal/mdutil"
)

const (
	DstDirMode  os.FileMode = 0755
	DstFileMode os.FileMode = 0666
)

type CreateTask struct {
	Common `mapstructure:",squash"`

	Dst  Destination `mapstructure:"dst"  required:"true"`
	Src  Source      `mapstructure:"src"  required:"true"`
	Type string      `mapstructure:"type" required:"true" const:"create" default:"create"`
}

// PrepareJSONSchema implements the jsonschema.Preparer interface.
func (t *CreateTask) PrepareJSONSchema(schema *jsonschema.Schema) error {
	schema.WithTitle("CreateTask")
	schema.WithDescription(mdutil.ToMarkdown(`
		Creates a new path in the destination directory.

		When using source templates, the [src.path](source_path.md#path)
		attribute may be a file or a directory path. When the latter,
		the source directory will be copied to the destination path recursively.

		Examples:

		__CODE_BLOCK__yaml
		tasks:
			- type: create
				# Render <./_src/README.tpl> (using the values defined in the generator)
				# and write it to <./README.md> in the destination directory.
				# If the README file already exists in the destination dir,
				# keep the existing file and do not bother prompting the user.
				src:
					path: "README.tpl"
				dst:
					path: "README.md"
					conflict: keep
		__CODE_BLOCK__

		__CODE_BLOCK__yaml
		values:
			- key: "FirstName"
				default: "Some Name"

		tasks:
			- type: create
				# Render the inline content as a template and write it to
				# <./some_name/greeting.txt> in the destination directory.
				src:
					content: "Hello, {{ .FirstName }}!"
				dst:
					path: "{{ .FirstName | underscore }}/greeting.txt"
		__CODE_BLOCK__

		__CODE_BLOCK__yaml
		tasks:
			- type: create
				# Render all the files in <./_src/scripts/> (using the values defined in the generator),
				# copy them to <./scripts/> in the destination directory, then make them executable.
				src:
					path: "scripts/"
				dst:
					path: "scripts/"
					mode: "0755"
		__CODE_BLOCK__
	`))

	schema.Properties["dst"].TypeObject.
		WithTitle("Destination").
		WithExamples(
			map[string]any{
				"path": "README.md",
			},
			map[string]any{
				"path": "bin/build.sh",
				"mode": "0755",
			},
		)

	schema.Properties["src"].TypeObject.
		WithTitle("Source").
		WithExamples(
			map[string]any{
				"path": "README.tpl",
			},
			map[string]any{
				"content": "Hello, {{ .FirstName }}!",
			},
		)

	schema.Properties["type"].TypeObject.
		WithTitle("Type").
		WithDescription("Creates a new path in the destination directory.").
		WithExamples("create")

	return nil
}

func (t *CreateTask) TypeKey() string {
	return t.Type
}

func (t *CreateTask) Execute(ctx *TaskContext, values map[string]any) error {
	if err := t.Dst.SetValues(values); err != nil {
		ctx.Logger.Failure("fail", t.Dst.Path())
		return err
	}
	if err := t.Src.SetValues(values); err != nil {
		ctx.Logger.Failure("fail", t.Dst.Path())
		return err
	}

	if t.Src.IsDir() {
		// src is a dir; walk and call dispatch on each file
		srcRoot := strings.TrimSuffix(t.Src.Path(), "/")
		dstRoot := strings.TrimSuffix(t.Dst.Path(), "/")
		return filepath.Walk(srcRoot, func(srcPath string, srcPathInfo fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// Construct the dst path by replacing `srcRoot` with `dstRoot`.
			dstPath := filepath.Join(dstRoot, strings.TrimPrefix(srcPath, srcRoot))

			// If the src path is a dir, create the dst dir and move on.
			if srcPathInfo.IsDir() {
				return t.createDstDir(ctx, dstPath)
			}

			// Otherwise create new Source and Destination structs and dispatch.
			src, err := NewSourceWithValues(srcPath, values)
			if err != nil {
				return err
			}
			mode, _ := t.Dst.ModeTpl.Render(values)
			dst, err := NewDestinationWithValues(dstPath, mode, values)
			if err != nil {
				return err
			}
			return t.dispatch(ctx, src, dst)
		})
	}

	// src is a single file (or inline content)
	return t.dispatch(ctx, t.Src, t.Dst)
}

// dispatch looks for conflicts and delegates to the correct generation method.
func (t *CreateTask) dispatch(ctx *TaskContext, src Source, dst Destination) error {
	if !dst.Exists() {
		return t.create(ctx, src, dst)
	}
	switch t.Dst.Conflict {
	case ConflictConfigKeep:
		return t.keep(ctx, src, dst)
	case ConflictConfigReplace:
		return t.replace(ctx, src, dst)
	default: // ConflictConfigPrompt
		return t.prompt(ctx, src, dst)
	}
}

// create is called to create a non-existing dst file.
func (t *CreateTask) create(ctx *TaskContext, src Source, dst Destination) error {
	if err := t.createDst(ctx, src, dst); err != nil {
		ctx.Logger.Failure("fail", dst.Path())
		return err
	}
	ctx.Logger.Success("create", dst.Path())
	return nil
}

// keep is called when keeping an existing dst file.
func (t *CreateTask) keep(ctx *TaskContext, _ Source, dst Destination) error {
	ctx.Logger.Success("keep", dst.Path())
	return nil
}

// replace is called when replacing an existing dst file.
func (t *CreateTask) replace(ctx *TaskContext, src Source, dst Destination) error {
	if err := t.deleteDst(ctx, src, dst); err != nil {
		ctx.Logger.Failure("fail", dst.Path())
		return err
	}
	if err := t.createDst(ctx, src, dst); err != nil {
		ctx.Logger.Failure("fail", dst.Path())
		return err
	}
	ctx.Logger.Success("replace", dst.Path())
	return nil
}

// prompt is called to prompt the user for how to resolve a dst file conflict.
// delegates to keep or replace depending on their response.
func (t *CreateTask) prompt(ctx *TaskContext, src Source, dst Destination) error {
	ctx.Logger.Warning("conflict", "%s already exists", dst.Path())
	overwrite, err := ctx.UI.Confirm("Overwrite", false)
	if err != nil {
		return err
	}
	if overwrite {
		return t.replace(ctx, src, dst)
	}
	return t.keep(ctx, src, dst)
}

func (t *CreateTask) createDstDir(ctx *TaskContext, path string) error {
	if ctx.DryRun {
		return nil
	}
	return os.MkdirAll(path, DstDirMode)
}

func (t *CreateTask) createDst(ctx *TaskContext, src Source, dst Destination) error {
	if ctx.DryRun {
		return nil
	}
	return dst.Write(src.Content())
}

func (t *CreateTask) deleteDst(ctx *TaskContext, _ Source, dst Destination) error {
	if ctx.DryRun {
		return nil
	}
	return dst.Delete()
}
