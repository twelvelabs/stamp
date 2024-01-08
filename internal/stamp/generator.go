package stamp

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/twelvelabs/termite/render"

	"github.com/twelvelabs/stamp/internal/fsutil"
	"github.com/twelvelabs/stamp/internal/pkg"
	"github.com/twelvelabs/stamp/internal/value"
)

var (
	ErrNilPackage = errors.New("nil package")
	ErrNilStore   = errors.New("nil store")
	ErrNotFound   = errors.New("generator not found")

	metaFileName = "generator.yaml"
)

func init() {
	funcMap := render.DefaultFuncMap()

	// Allow generators to access to their name in templates.
	// Primary use case is for the "generator generator" to auto-populate the name on create.
	funcMap["generatorName"] = GeneratorName
	funcMap["generatorNameForCreate"] = GeneratorNameForCreate

	render.FuncMap = funcMap
}

// GeneratorName returns the name (i.e. "foo:bar:baz") for the generator at path.
// Returns an empty string if path is not a generator.
func GeneratorName(path string) string {
	p, err := pkg.LoadPackage(path, metaFileName)
	if err != nil {
		return ""
	}
	return p.Name()
}

// GeneratorNameForCreate returns the correct name for a new generator at the given path.
//
// It searches for the furthest ancestor directory containing a generator.
// If found, the name will be prefixed with that generator name.
// If no ancestor is found then the name will be the segment following
// the final path separator.
func GeneratorNameForCreate(path string) string {
	path, _ = filepath.Abs(path)
	separator := string(filepath.Separator)
	segments := strings.Split(path, separator)

	// Find the furthest ancestor that is a generator and start the name from there.
	// If no ancestor generator, default the immediate directory name.
	var idx int
	for idx = range segments {
		// Reconstruct the path for this segment range.
		// Note: `filepath.Join` ignores empty segments, so add the leading "/" back in.
		subpath := filepath.Join(segments[0 : idx+1]...)
		if !strings.HasPrefix(subpath, separator) {
			subpath = separator + subpath
		}
		// Try loading a package for the subpath.
		p, err := pkg.LoadPackage(subpath, metaFileName)
		if err != nil {
			continue // NOT a valid package
		}
		// Found a valid package. Ensure we're using the root package name
		// (it may be different from the path segment).
		segments[idx] = p.Name()
		break
	}

	return strings.Join(segments[idx:], ":")
}

// NewGeneratorFromPath returns the generator located at path (or ErrNotFound).
func NewGeneratorFromPath(store *Store, path string) (*Generator, error) {
	if !fsutil.PathIsDir(path) {
		return nil, ErrNotFound
	}
	if fsutil.NoPathExists(filepath.Join(path, store.MetaFile)) {
		return nil, ErrNotFound
	}

	p, err := pkg.LoadPackage(path, store.MetaFile)
	if err != nil {
		return nil, err
	}
	return NewGenerator(store, p)
}

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

// ShortDescription returns the first line of the description.
func (g *Generator) ShortDescription() string {
	lines := strings.Split(g.Description(), "\n")
	return lines[0]
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
