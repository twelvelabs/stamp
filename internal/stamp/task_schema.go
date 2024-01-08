package stamp

import (
	"github.com/swaggest/jsonschema-go"
)

// TaskSchema represents the Task interface in the JSON Schema definition.
type TaskSchema struct {
	Type string `mapstructure:"type" title:"Type" description:"The task type." required:"true"`
}

var _ jsonschema.Preparer = TaskSchema{}

// PrepareJSONSchema implements the jsonschema.Preparer interface.
func (s TaskSchema) PrepareJSONSchema(schema *jsonschema.Schema) error {
	schema.WithTitle("Task")
	schema.WithDescription(
		"A task to execute in the destination directory.",
	)

	typeKeys := []any{}
	for _, t := range AllTasks() {
		typeKeys = append(typeKeys, t.TypeKey())
	}
	schema.Properties["type"].TypeObject.WithEnum(typeKeys...)

	return nil
}

var _ jsonschema.OneOfExposer = TaskSchema{}

// PrepareJSONSchema implements the jsonschema.OneOfExposer interface.
func (s TaskSchema) JSONSchemaOneOf() []any {
	tasks := []any{}
	for _, t := range AllTasks() {
		tasks = append(tasks, t)
	}
	return tasks
}
