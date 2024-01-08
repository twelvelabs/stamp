package stamp

import (
	"github.com/swaggest/jsonschema-go"
)

type DeleteTask struct {
	Common `mapstructure:",squash"`

	Dst  Destination `mapstructure:"dst"  required:"true" title:"Destination"`
	Type string      `mapstructure:"type" required:"true" title:"Type" description:"Deletes a path in the destination directory." const:"delete" default:"delete"` //nolint: lll
}

// PrepareJSONSchema implements the jsonschema.Preparer interface.
func (t *DeleteTask) PrepareJSONSchema(schema *jsonschema.Schema) error {
	schema.WithTitle("DeleteTask")
	schema.WithDescription("Deletes a path in the destination directory.")
	return nil
}

func (t *DeleteTask) TypeKey() string {
	return t.Type
}

func (t *DeleteTask) Execute(ctx *TaskContext, values map[string]any) error {
	err := t.Dst.SetValues(values)
	if err != nil {
		return err
	}

	if t.Dst.Exists() {
		if err := t.deleteDst(ctx); err != nil {
			ctx.Logger.Failure("fail", t.Dst.Path())
			return err
		}
		ctx.Logger.Success("delete", t.Dst.Path())
		return nil
	} else if t.Dst.Missing == MissingConfigError {
		ctx.Logger.Failure("fail", t.Dst.Path())
		return ErrPathNotFound
	}

	return nil
}

func (t *DeleteTask) deleteDst(ctx *TaskContext) error {
	if ctx.DryRun {
		return nil
	}
	return t.Dst.Delete()
}
