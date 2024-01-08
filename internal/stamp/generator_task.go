package stamp

import (
	"github.com/mitchellh/copystructure"
	"github.com/swaggest/jsonschema-go"

	"github.com/twelvelabs/stamp/internal/mdutil"
)

type GeneratorTask struct {
	Common `mapstructure:",squash"`

	Name   string         `mapstructure:"name"   title:"Name" required:"true" description:"The name of the generator to execute." validate:"required"`         //nolint: lll
	Values map[string]any `mapstructure:"values" title:"Values"               description:"Optional key/value pairs to pass to the generator." default:"{}"`   //nolint: lll
	Type   string         `mapstructure:"type"   title:"Type" required:"true" description:"Executes another generator." const:"generator" default:"generator"` //nolint: lll
}

// PrepareJSONSchema implements the jsonschema.Preparer interface.
func (t *GeneratorTask) PrepareJSONSchema(schema *jsonschema.Schema) error {
	schema.WithTitle("GeneratorTask")
	schema.WithDescription(mdutil.ToMarkdown(`
		Executes another generator.

		All values defined by the included generator are prepended
		to the including generator's values list. This allows the
		including generator to redefine values if needed.

		You can also optionally define values to set in the included generator.
		This can be useful to prevent the user from being prompted for
		values that do not make sense in your use case (see example).

		Example:

		__CODE_BLOCK__yaml
		name: "python-api"

		tasks:
			- type: "generator"
				# Executes the gitignore generator, pre-setting
				# the "Language" value so the user isn't prompted.
				name: "gitignore"
				values:
					Language: "python"

			# ... other tasks ...
		__CODE_BLOCK__
	`))

	schema.Properties["values"].TypeObject.
		WithExamples(
			map[string]any{
				"ValueKeyOne": "foo",
				"ValueKeyTwo": "bar",
			},
		)

	return nil
}

func (t *GeneratorTask) TypeKey() string {
	return t.Type
}

func (t *GeneratorTask) Execute(ctx *TaskContext, values map[string]any) error {
	gen, err := t.GetGenerator(ctx.Store)
	if err != nil {
		return err
	}

	// Deep copy so that any value mutation done by this generator
	// doesn't leak up to the caller.
	copied, err := copystructure.Copy(values)
	if err != nil {
		return err
	}
	data := copied.(map[string]any)
	// Doing this here (in addition to the setting in GetGenerator) to cover
	// the case where the same generator name is used in multiple tasks.
	// In that scenario, `values` would contain those from the last task added.
	for k, v := range t.Values {
		data[k] = v
	}

	return gen.Tasks.Execute(ctx, data)
}

func (t *GeneratorTask) GetGenerator(store *Store) (*Generator, error) {
	gen, err := store.Load(t.Name)
	if err != nil {
		return nil, err
	}
	// Setting the value overrides here so that when the values are added to the
	// delegating generator (See NewGenerator()), they have the correct data.
	for k, v := range t.Values {
		err = gen.Values.Set(k, v)
		if err != nil {
			return nil, err
		}
	}
	return gen, nil
}
