package gen

import (
	"errors"
	"fmt"
	"path/filepath"
	"reflect"

	// cspell:disable
	"github.com/gobuffalo/flect"

	// cspell:enable

	"github.com/twelvelabs/stamp/internal/pkg"
	"github.com/twelvelabs/stamp/internal/task"
	"github.com/twelvelabs/stamp/internal/value"
)

var (
	ErrNilPackage = errors.New("nil package")
)

func NewGenerator(pkg *pkg.Package) (*Generator, error) {
	if pkg == nil {
		return nil, ErrNilPackage
	}

	gen := &Generator{
		Package: pkg,
		Values:  value.NewValueSet(),
		Tasks:   task.NewTaskSet(),
	}

	dstPath, err := filepath.Abs(".")
	if err != nil {
		return nil, err
	}
	gen.Values.Set("SrcPath", gen.Path())
	gen.Values.Set("DstPath", dstPath)

	for _, vm := range gen.valueMetadata() {
		v, err := value.NewValue(vm)
		if err != nil {
			return nil, fmt.Errorf("generator metadata invalid: %w", err)
		}
		gen.Values.Add(v)
	}

	for _, tm := range gen.taskMetadata() {
		t, err := task.NewTask(tm)
		if err != nil {
			return nil, fmt.Errorf("generator metadata invalid: %w", err)
		}
		gen.Tasks.Add(t)
	}

	return gen, nil
}

func NewGenerators(packages []*pkg.Package) ([]*Generator, error) {
	generators := []*Generator{}
	for _, p := range packages {
		generator, err := NewGenerator(p)
		if err != nil {
			return nil, err
		}
		generators = append(generators, generator)
	}
	return generators, nil
}

type Generator struct {
	*pkg.Package

	Values *value.ValueSet
	Tasks  *task.TaskSet
}

func (g *Generator) taskMetadata() []map[string]any {
	return g.metadataSliceOfMaps("tasks")
}

func (g *Generator) valueMetadata() []map[string]any {
	return g.metadataSliceOfMaps("values")
}

// Returns a slice of maps for the given metadata key.
func (g *Generator) metadataSliceOfMaps(key string) []map[string]any {
	items := []map[string]any{}
	for i, m := range g.metadataSlice(key) {
		item, ok := m.(map[string]any)
		if !ok {
			panic(NewMetadataTypeCastError(
				fmt.Sprintf("%s[%d]", key, i),
				m,
				"map[string]any",
			))
		}
		items = append(items, item)
	}
	return items
}

// returns a slice value for the given metadata key.
func (g *Generator) metadataSlice(key string) []any {
	val := g.metadataLookup(key)
	switch val.(type) {
	case []any:
		return val.([]any)
	case nil:
		return []any{}
	default:
		panic(NewMetadataTypeCastError(key, val, "[]any"))
	}
}

// Returns the value for key in metadata.
// If not found, then tries the camel-cased version of key.
func (g *Generator) metadataLookup(key string) any {
	if val, ok := g.Metadata[key]; ok {
		return val
	}

	if val, ok := g.Metadata[flect.Pascalize(key)]; ok {
		return val
	}
	return nil
}

type MetadataTypeCastError struct {
	key          string
	value        any
	expectedType string
	actualType   string
}

func NewMetadataTypeCastError(key string, value any, expectedType string) MetadataTypeCastError {
	actualType := reflect.TypeOf(value).String()
	return MetadataTypeCastError{
		key:          key,
		value:        value,
		expectedType: expectedType,
		actualType:   actualType,
	}
}

func (e MetadataTypeCastError) Error() string {
	return fmt.Sprintf(
		"generator metadata invalid: '%s' should be '%s', is '%s'",
		e.key,
		e.expectedType,
		e.actualType,
	)
}
