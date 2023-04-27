package stamp

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/twelvelabs/stamp/internal/pkg"
	"github.com/twelvelabs/stamp/internal/value"
)

var (
	ErrNilPackage = errors.New("nil package")
	ErrNilStore   = errors.New("nil store")
)

func NewGenerator(store *Store, pkg *pkg.Package) (*Generator, error) {
	if store == nil {
		return nil, ErrNilStore
	}
	if pkg == nil {
		return nil, ErrNilPackage
	}

	gen := &Generator{
		Package: pkg,
		Values:  value.NewValueSet(),
		Tasks:   NewTaskSet(),
	}

	for _, tm := range gen.taskMetadata() {
		t, err := NewTask(tm)
		if err != nil {
			return nil, fmt.Errorf("generator metadata invalid: %w", err)
		}
		gen.Tasks.Add(t)

		if gt, ok := t.(*GeneratorTask); ok {
			// This is a task that calls out to a sub-generator,
			// merge in the sub-generator's values.
			subGen, err := gt.GetGenerator(store)
			if err != nil {
				return nil, fmt.Errorf("unable to load sub-generator '%s': %w", gt.Name, err)
			}
			for _, val := range subGen.Values.All() {
				gen.Values.Add(val)
			}
		}
	}

	for _, vm := range gen.valueMetadata() {
		v, err := value.NewValue(vm)
		if err != nil {
			return nil, fmt.Errorf("generator metadata invalid: %w", err)
		}
		gen.Values.Add(v)
	}

	if len(gen.Values.Args()) == 0 {
		gen.Values.Prepend(&value.Value{
			Key:             "DstPath",
			Name:            "Destination Path",
			Help:            "The path to generate files to.",
			DataType:        value.DataTypeString,
			Default:         ".",
			InputMode:       value.InputModeArg,
			PromptConfig:    value.PromptConfigOnEmpty,
			TransformRules:  "trim,expand-path",
			ValidationRules: "required",
		})
	}

	dstPath, err := filepath.Abs(".")
	if err != nil {
		return nil, err
	}
	gen.Tasks.SrcPath = gen.SrcPath()
	gen.Tasks.DstPath = dstPath

	return gen, nil
}

func NewGenerators(store *Store, packages []*pkg.Package) ([]*Generator, error) {
	if store == nil {
		return nil, ErrNilStore
	}
	generators := []*Generator{}
	for _, p := range packages {
		generator, err := NewGenerator(store, p)
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
	Tasks  *TaskSet
}

func (g *Generator) Description() string {
	return g.MetadataString("description")
}

func (g *Generator) SrcPath() string {
	return filepath.Join(g.Path(), "_src")
}

func (g *Generator) taskMetadata() []map[string]any {
	return g.MetadataMapSlice("tasks")
}

func (g *Generator) valueMetadata() []map[string]any {
	return g.MetadataMapSlice("values")
}
