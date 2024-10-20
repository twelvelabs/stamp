package stamp

import (
	"encoding/json"
	"fmt"

	"github.com/twelvelabs/termite/render"
)

type PluginTask struct {
	Common `mapstructure:",squash"`

	DescriptionTpl render.Template `mapstructure:"description" title:"Description"                 description:"Optional description of the task." validate:"required"`              //nolint: lll
	FunctionTpl    render.Template `mapstructure:"function"    title:"Function"    required:"true" description:"The name of the plugin function."  validate:"required"`              //nolint: lll
	NameTpl        render.Template `mapstructure:"name"        title:"Name"        required:"true" description:"The name of the plugin."           validate:"required"`              //nolint: lll
	Values         map[string]any  `mapstructure:"values"      title:"Values"                      description:"Additional key/value pairs to pass to the plugin." default:"{}"`     //nolint: lll
	Type           string          `mapstructure:"type"        title:"Type"        required:"true" description:"Executes a plugin." const:"plugin"                 default:"plugin"` //nolint: lll

	Generator *Generator
}

func (t *PluginTask) TypeKey() string {
	return t.Type
}

func (t *PluginTask) Execute(ctx *TaskContext, values map[string]any) error {
	// Render plugin name.
	pluginName, err := t.NameTpl.Render(values)
	if err != nil {
		ctx.Logger.Failure("fail", pluginName)
		return err
	}

	// Render function name.
	funcName, err := t.FunctionTpl.Render(values)
	if err != nil {
		ctx.Logger.Failure("fail", pluginName)
		return err
	}

	// Render description.
	desc, err := t.DescriptionTpl.Render(values)
	if err != nil {
		ctx.Logger.Failure("fail", pluginName)
		return err
	}

	message := fmt.Sprintf("%s (%s)", pluginName, funcName)
	if desc != "" {
		message = desc
	}

	// Load the plugin.
	plugin, err := t.Generator.LoadPlugin(pluginName)
	if err != nil {
		ctx.Logger.Failure("fail", message)
		return err
	}

	// Render (and copy in) any additional values provided by the task.
	inVals, err := render.Map(t.Values, values)
	if err != nil {
		ctx.Logger.Failure("fail", message)
		return fmt.Errorf("plugin input: %w", err)
	}
	for k, v := range inVals {
		values[k] = v
	}

	// JSON encode the values.
	input, err := json.Marshal(values)
	if err != nil {
		ctx.Logger.Failure("fail", message)
		return fmt.Errorf("plugin input: %w", err)
	}

	// Call the plugin (passing in the values JSON).
	_, output, err := plugin.Call(funcName, input)
	if err != nil {
		ctx.Logger.Failure("fail", message)
		return fmt.Errorf("plugin call: %w", err)
	}

	// JSON decode the plugin output.
	var outVals map[string]any
	if err := json.Unmarshal(output, &outVals); err != nil {
		ctx.Logger.Failure("fail", message)
		return fmt.Errorf("plugin output: %w", err)
	}

	// Copy the plugin output into the values map,
	// so those values can be consumed by downstream tasks.
	for k, v := range outVals {
		if k == "SrcPath" || k == "DstPath" {
			// Prevent plugins from altering sensitive values
			// used to sandbox generators.
			ctx.Logger.Failure("fail", message)
			return fmt.Errorf("plugin output key invalid: %s", k)
		}
		values[k] = v
	}

	ctx.Logger.Success("exec", message)
	return nil
}
