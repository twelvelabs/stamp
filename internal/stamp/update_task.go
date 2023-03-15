package stamp

type UpdateTask struct {
	Common `mapstructure:",squash"`

	Dst     string        `validate:"required"`
	Missing MissingConfig `validate:"required" default:"ignore"`
	Pattern string        `validate:"required"`
	Value   string        ``
}

func (t *UpdateTask) Execute(ctx *TaskContext, values map[string]any) error {
	_, err := t.RenderRequired("dst", t.Dst, values)
	if err != nil {
		return err
	}

	return nil
}
