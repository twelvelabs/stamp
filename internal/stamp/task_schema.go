package stamp

import (
	"github.com/swaggest/jsonschema-go"
)

// TaskSchema represents the Task interface in the JSON Schema definition.
type TaskSchema struct {
}

var _ jsonschema.Preparer = TaskSchema{}

// PrepareJSONSchema implements the jsonschema.Preparer interface.
func (s TaskSchema) PrepareJSONSchema(schema *jsonschema.Schema) error {
	schema.
		WithDescription("A generator task.")
	return nil
}

var _ jsonschema.OneOfExposer = TaskSchema{}

// PrepareJSONSchema implements the jsonschema.OneOfExposer interface.
func (s TaskSchema) JSONSchemaOneOf() []any {
	return AllTasks()
}
