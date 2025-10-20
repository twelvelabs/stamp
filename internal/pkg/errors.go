package pkg

import (
	"fmt"
	"reflect"
)

func NewMetadataTypeCastError(key string, value any, expectedType string) MetadataTypeCastError {
	actualType := reflect.TypeOf(value).String()
	return MetadataTypeCastError{
		key:          key,
		value:        value,
		expectedType: expectedType,
		actualType:   actualType,
	}
}

type MetadataTypeCastError struct {
	key          string
	value        any
	expectedType string
	actualType   string
}

func (e MetadataTypeCastError) Error() string {
	return fmt.Sprintf(
		"metadata invalid: '%s' should be '%s', is '%s'",
		e.key,
		e.expectedType,
		e.actualType,
	)
}
