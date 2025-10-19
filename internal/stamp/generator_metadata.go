package stamp

import (
	"github.com/swaggest/jsonschema-go"

	"github.com/twelvelabs/stamp/internal/mdutil"
	"github.com/twelvelabs/stamp/internal/value"
)

// GeneratorMetadata represents the structure of the generator.yaml file.
type GeneratorMetadata struct {
	Name        string         `mapstructure:"name" required:"true"`
	Description string         `mapstructure:"description"`
	Visibility  VisibilityType `mapstructure:"visibility" default:"public"`
	Values      []value.Value  `mapstructure:"values"`
	Tasks       []TaskSchema   `mapstructure:"tasks"`
}

var _ jsonschema.Preparer = &GeneratorMetadata{}

// PrepareJSONSchema implements the jsonschema.Preparer interface.
func (m *GeneratorMetadata) PrepareJSONSchema(schema *jsonschema.Schema) error {
	schema.WithTitle("Generator")
	schema.WithDescription("Stamp generator metadata.")
	schema.WithExtraPropertiesItem("markdownDescription", mdutil.ToMarkdown(`
		Stamp generator metadata.

		Example:

		__CODE_BLOCK__yaml
		name: "greet"
		description: "Generates a text file with a greeting message."

		# When run, the generator will prompt the user for these values.
		# Alternately, the user can pass them in via flags.
		values:
			# Values are prompted in the order they are defined.
			- key: "Name"
				default: "Some Name"

			# Subsequent values can reference those defined prior.
			# This allows for sensible, derived defaults.
			- key: "Greeting"
				default: "Hello, {{ .Name }}."

		# Next, the generator executes a series of tasks.
		# Tasks have access to the values defined above.
		tasks:
			# Render the inline content as a template string.
			# Write it to <./some_name.txt> in the destination directory.
			- type: create
				src:
					content: "{{ .Greeting }}"
				dst:
					path: "{{ .Name | underscore }}.txt"
		__CODE_BLOCK__


		__CODE_BLOCK__shell
		# Save the above to <./greeting/generator.yaml>.
		# The following will prompt for values, then write <./some_name.txt>.
		stamp new ./greeting

		# Pass an alternate destination dir as the second argument.
		# The following creates </some/other/dir/some_name.txt>.
		stamp new ./greeting /some/other/dir

		# Install the generator so you can refer to it by name
		# rather than filesystem path.
		stamp add ./greeting
		stamp new greet

		# You can also publish it to a git repo or upload it as an archive
		# and share it with others:
		stamp add git@github.com:username/my-generator.git
		stamp add github.com/username/my-generator
		stamp add https://example.com/my-generator.tar.gz
		__CODE_BLOCK__
	`))

	schema.WithSchema("https://json-schema.org/draft-07/schema")
	schema.WithID(
		"https://raw.githubusercontent.com/twelvelabs/stamp/main/docs/stamp.schema.json",
	)

	schema.Properties["name"].TypeObject.
		WithTitle("Name").
		WithDescription("The generator name.").
		WithPattern(`^[\w:_-]+$`).
		WithMinLength(1)

	schema.Properties["description"].TypeObject.
		WithTitle("Description").
		WithDescription(
			"The generator description. " +
				"The first line is shown when listing all generators. " +
				"The full description is used when viewing generator help/usage text.",
		)

	schema.Properties["visibility"].TypeObject.
		WithTitle("Visibility").
		WithDescription("How the generator may be viewed or invoked.")

	schema.Properties["values"].TypeObject.
		WithTitle("Values").
		WithDescription(
			"A list of generator input [values](https://github.com/twelvelabs/stamp/tree/main/docs/value.md).",
		)

	schema.Properties["tasks"].TypeObject.
		WithTitle("Tasks").
		WithDescription(
			"A list of generator [tasks](https://github.com/twelvelabs/stamp/tree/main/docs/task.md).",
		)

	return nil
}

// ReflectSchema generates a JSON schema document for the metadata file.
func (m *GeneratorMetadata) ReflectSchema() (jsonschema.Schema, error) {
	// Generate the schema.
	reflector := &jsonschema.Reflector{}
	schema, err := reflector.Reflect(m,
		jsonschema.StripDefinitionNamePrefix("Modify", "Stamp", "Value"),
		jsonschema.PropertyNameTag("mapstructure"),
	)
	if err != nil {
		return jsonschema.Schema{}, err
	}

	addTransformRules(&schema)
	setAdditionalProperties(&schema)
	setMarkdownDescription(&schema)

	return schema, nil
}

// Recursively sets the `markdownDescription` property to
// the content of `description`.
//
// This non-standard property is used by some LSP language servers
// when hovering over JSON or YAML fields.
// For example, VS Code uses either `description` or `markdownDescription`
// when rendering hover overlays, but escapes all markdown in `description`.
func setMarkdownDescription(schema *jsonschema.Schema) {
	if schema.Description != nil {
		ep := schema.ExtraProperties
		if ep == nil {
			ep = map[string]any{}
		}
		// Only write markdownDescription if not already set.
		if _, ok := ep["markdownDescription"]; !ok {
			schema.WithExtraPropertiesItem(
				"markdownDescription", *schema.Description,
			)
		}
	}
	for _, s := range schema.Properties {
		if s.TypeObject != nil {
			setMarkdownDescription(s.TypeObject)
		}
	}
	for _, s := range schema.Definitions {
		if s.TypeObject != nil {
			setMarkdownDescription(s.TypeObject)
		}
	}
}

func addTransformRules(schema *jsonschema.Schema) {
	enums := []any{}
	enumDescriptions := []any{}

	for _, t := range value.RegisteredTransformers() {
		enums = append(enums, t.Name)
		enumDescriptions = append(enumDescriptions, t.Description)
	}

	transform := &jsonschema.SchemaOrBool{}
	transform.TypeObjectEns().
		WithType(jsonschema.String.Type()).
		WithTitle("Transform").
		WithDescription("A transformation function used to process a value.").
		WithEnum(enums...).
		WithExtraPropertiesItem("enumDescriptions", enumDescriptions)

	schema.WithDefinitionsItem("Transform", *transform)
}

// Helper that squashes embedded structs up into the parent
// (similar to how mapstructure behaves when using the `squash` struct tag).
func setAdditionalProperties(schema *jsonschema.Schema) {
	// Create a new set of definitions that will replace the current set.
	newDefs := map[string]jsonschema.SchemaOrBool{}
	for k, def := range schema.Definitions {
		// Set additionalProperties to false for all `object` defs.
		// (though not if it's a (all|any|one)Of type, because that breaks validation).
		s := def.TypeObjectEns()
		sIsObj := schemaIsType(s, jsonschema.Object)
		if sIsObj && len(s.AllOf) == 0 && len(s.AnyOf) == 0 && len(s.OneOf) == 0 {
			s.AdditionalPropertiesEns().WithTypeBoolean(false)
		}
		// Add the updated definition to the new set.
		newDefs[k] = def
	}
	// Finally, update the schema w/ the new defs.
	schema.WithDefinitions(newDefs)
}

func schemaIsType(schema *jsonschema.Schema, st jsonschema.SimpleType) bool {
	// The jsonschema lib has two different ways of tracking type.
	// Single type.
	if schema.TypeEns().SimpleTypes != nil && *schema.TypeEns().SimpleTypes == st {
		return true
	}
	// Or a slice of possible types.
	for _, sst := range schema.TypeEns().SliceOfSimpleTypeValues {
		if st == sst {
			return true
		}
	}
	return false
}
