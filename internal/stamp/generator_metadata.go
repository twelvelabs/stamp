package stamp

import (
	"github.com/swaggest/jsonschema-go"

	"github.com/twelvelabs/stamp/internal/value"
)

// GeneratorMetadata represents the structure of the generator.yaml file.
type GeneratorMetadata struct {
	Name        string        `mapstructure:"name"`
	Description string        `mapstructure:"description"`
	Values      []value.Value `mapstructure:"values"`
	Tasks       []TaskSchema  `mapstructure:"tasks"`
}

var _ jsonschema.Preparer = &GeneratorMetadata{}

// PrepareJSONSchema implements the jsonschema.Preparer interface.
func (m *GeneratorMetadata) PrepareJSONSchema(schema *jsonschema.Schema) error {
	schema.WithTitle("Generator")
	schema.WithDescription("Stamp generator metadata.")

	schema.WithSchema("https://json-schema.org/draft-07/schema")
	schema.WithID(
		"https://raw.githubusercontent.com/twelvelabs/stamp/main/docs/stamp.schema.json",
	)

	schema.Properties["name"].TypeObject.
		WithDescription("The generator name.").
		WithPattern(`^[\w:_-]+$`).
		WithMinLength(1)

	schema.Properties["description"].TypeObject.
		WithDescription(
			"The generator description. " +
				"The first line is shown when listing all generators. " +
				"The full description is used when viewing generator help/usage text.",
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
		jsonschema.InterceptProp(func(params jsonschema.InterceptPropParams) error {
			// switch params.Name {
			// case "Common":
			// }
			return nil
		}),
	)
	if err != nil {
		return jsonschema.Schema{}, err
	}

	schema = squashEmbeddedStruct(schema, "Common")

	return schema, nil
}

// Helper that squashes embedded structs up into the parent
// (similar to how mapstructure behaves when using the `squash` struct tag).
func squashEmbeddedStruct(schema jsonschema.Schema, embeddedName string) jsonschema.Schema {
	embeddedDef, ok := schema.Definitions[embeddedName]
	if !ok {
		return schema // not found; nothing to do.
	}

	// Create a new set of definitions that will replace the current set.
	newDefs := map[string]jsonschema.SchemaOrBool{}
	for k, def := range schema.Definitions {
		// If this definition is the embedded one, then remove it from the new set.
		if k == embeddedName {
			continue
		}

		// Set additionalProperties to false for all `object` defs.
		// (though not if it's a (all|any|one)Of type, because that breaks validation).
		s := def.TypeObjectEns()
		sIsObj := schemaIsType(s, jsonschema.Object)
		if sIsObj && len(s.AllOf) == 0 && len(s.AnyOf) == 0 && len(s.OneOf) == 0 {
			s.AdditionalPropertiesEns().WithTypeBoolean(false)
		}

		if _, ok := def.TypeObjectEns().Properties[embeddedName]; !ok {
			// This def doesn't contain the embedded struct.
			// Just need to add it to the new set of defs and move on.
			newDefs[k] = def
			continue
		}

		// Otherwise, create a new set of properties for the definition.
		// that flatten the embedded struct's props down into the definition.
		newProps := map[string]jsonschema.SchemaOrBool{}
		for k, v := range embeddedDef.TypeObjectEns().Properties {
			newProps[k] = v
		}
		for k, v := range def.TypeObjectEns().Properties {
			if k == embeddedName {
				continue
			}
			newProps[k] = v
		}
		// Replace the definition's props with the new, flattened set.
		def.TypeObjectEns().WithProperties(newProps)

		// Add the updated definition to the new set.
		newDefs[k] = def
	}

	// Finally, update the schema w/ the new defs.
	schema.WithDefinitions(newDefs)

	return schema
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
