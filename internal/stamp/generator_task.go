package stamp

import (
	"github.com/mitchellh/copystructure"
)

type GeneratorTask struct {
	Common `mapstructure:",squash"`

	Name   string         `mapstructure:"name" validate:"required"`
	Values map[string]any `mapstructure:"values" default:"{}"`
	Type   string         `mapstructure:"type" const:"generator" description:"Executes another generator."`
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
