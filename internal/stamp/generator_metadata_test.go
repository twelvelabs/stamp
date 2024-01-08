package stamp

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/swaggest/jsonschema-go"
)

func TestGeneratorMetadata_ReflectSchema(t *testing.T) {
	require := require.New(t)

	metadata := &GeneratorMetadata{}

	// Should be able to reflect the JSON schema document.
	schema, err := metadata.ReflectSchema()
	require.NotNil(schema)
	require.NoError(err)

	// Should be able to convert that document into JSON.
	buf, err := json.Marshal(schema)
	require.NoError(err)
	require.NotEmpty(buf)
}

func TestGeneratorMetadata_ReflectSchema_SetsMarkdownDescription(t *testing.T) {
	require := require.New(t)

	metadata := &GeneratorMetadata{}
	schema, err := metadata.ReflectSchema()
	require.NoError(err)

	// Each prop should have markdownDescription attribute.
	for key, p := range schema.Properties {
		prop := p.TypeObjectEns()
		name := fmt.Sprintf("prop: %s", key)
		requireMarkdownDescription(t, name, prop)
	}
}

func requireMarkdownDescription(t *testing.T, name string, schema *jsonschema.Schema) {
	t.Helper()
	require := require.New(t)

	require.NotNil(schema.Description, "Description should not be nil for %s", name)
	require.NotNil(schema.ExtraProperties, "ExtraProperties should not be nil for %s", name)

	desc := *schema.Description
	mdDesc := schema.ExtraProperties["markdownDescription"]
	require.NotEmpty(desc, "Description should not be empty for %s", name)
	require.NotEmpty(mdDesc, "MarkdownDescription should not be empty for %s", name)
	require.Equal(desc, mdDesc, "descriptions should be equal for %s", name)
}
