package encode

// Encoder is an interface type that can convert data structures
// to/from encoded byte arrays.
type Encoder interface {
	// Decode deserializes the given byte array into a data structure.
	Decode(encoded []byte) (any, error)

	// Encode serializes the given data structure into a byte array.
	Encode(data any) ([]byte, error)
}
