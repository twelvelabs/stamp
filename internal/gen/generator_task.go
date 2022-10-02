package gen

import "github.com/mitchellh/copystructure" //cspell: disable-line

type GeneratorTask struct {
	Common `mapstructure:",squash"`

	Name  string         `validate:"required"`
	Extra map[string]any `default:"{}"`
}

func (t *GeneratorTask) Execute(ctx *TaskContext, values map[string]any) error {
	t.DryRun = ctx.DryRun

	gen, err := t.GetGenerator(ctx.Store)
	if err != nil {
		return err
	}

	copied, err := copystructure.Copy(values)
	if err != nil {
		return err
	}
	data := copied.(map[string]any)
	for k, v := range t.Extra {
		data[k] = v
	}

	return gen.Tasks.Execute(ctx, data)
}

func (t *GeneratorTask) GetGenerator(store *Store) (*Generator, error) {
	gen, err := store.Load(t.Name)
	if err != nil {
		return nil, err
	}
	// TODO: set values using Extra
	return gen, nil
}
