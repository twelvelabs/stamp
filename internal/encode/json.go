package encode

import (
	"encoding/json"
	"fmt"

	"github.com/ohler55/ojg/oj"
)

var _ Encoder = &JSONEncoder{}

type JSONEncoder struct {
}

// Decode deserializes the given JSON encoded byte array into a data structure.
func (e *JSONEncoder) Decode(encoded []byte) (any, error) {
	data, err := oj.Parse(encoded)
	if err != nil {
		return nil, fmt.Errorf("json decode: %w", err)
	}
	return data, nil
}

// Encode serializes the given data structure into a JSON encoded byte array.
func (e *JSONEncoder) Encode(data any) ([]byte, error) {
	// Note: using standard lib to marshal because it sorts JSON object keys
	// (oj does not and it looks ugly when adding new keys).
	content, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("json encode: %w", err)
	}
	return content, nil
}
