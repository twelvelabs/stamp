package encode

import (
	"bytes"
	"fmt"

	"github.com/spf13/cast"
)

var _ Encoder = &TextEncoder{}

type TextEncoder struct {
}

// Decode returns a copy of the given byte array.
func (e *TextEncoder) Decode(encoded []byte) (any, error) {
	return bytes.Clone(encoded), nil
}

// Encode serializes the given data structure into a byte array.
func (e *TextEncoder) Encode(data any) ([]byte, error) {
	casted, err := cast.ToStringE(data)
	if err != nil {
		return nil, fmt.Errorf("text encode: unable to cast: %#v", data)
	}
	return []byte(casted), nil
}
