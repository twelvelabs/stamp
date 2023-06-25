package stamp

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratorMetadata_ReflectSchema(t *testing.T) {
	metadata := &GeneratorMetadata{}

	// Should be able to reflect the JSON schema document.
	schema, err := metadata.ReflectSchema()
	assert.NotNil(t, schema)
	assert.NoError(t, err)

	// Should be able to convert that document into JSON.
	buf, err := json.Marshal(schema)
	assert.NoError(t, err)
	assert.NotEmpty(t, buf)
}
