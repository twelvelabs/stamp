package pkg

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
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

func NewNotFoundError(metaFile string) NotFoundError {
	return NotFoundError{
		metaFile: metaFile,
	}
}

type NotFoundError struct {
	metaFile string
}

func (e NotFoundError) Error() string {
	// This assumes that the meta filename describes what the package is.
	// i.e. if your package is a "widget", then name the file widget.yaml
	kind := strings.TrimSuffix(e.metaFile, filepath.Ext(e.metaFile))
	return fmt.Sprintf("%s not found", kind)
}
