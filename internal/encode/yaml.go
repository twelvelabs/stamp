package encode

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v3"
)

var _ Encoder = &YAMLEncoder{}

type YAMLEncoder struct {
}

// Decode deserializes the given YAML encoded byte array into a data structure.
func (e *YAMLEncoder) Decode(encoded []byte) (any, error) {
	var data any
	err := yaml.Unmarshal(encoded, &data)
	if err != nil {
		return nil, fmt.Errorf("yaml decode: %w", err)
	}
	return data, nil
}

// Encode serializes the given data structure into a YAML encoded byte array.
func (e *YAMLEncoder) Encode(data any) ([]byte, error) {
	b := &bytes.Buffer{}
	encoder := yaml.NewEncoder(b)
	encoder.SetIndent(2)
	err := encoder.Encode(data)
	if err != nil {
		return nil, fmt.Errorf("yaml encode: %w", err)
	}
	return b.Bytes(), nil
}
