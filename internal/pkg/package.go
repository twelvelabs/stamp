package pkg

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/flect"
)

type Package struct {
	Metadata map[string]any

	path     string
	metaFile string
}

// MetaFile returns the name of the metadata file.
func (p *Package) MetaFile() string {
	return p.metaFile
}

// MetaPath returns the full path to the metadata file.
func (p *Package) MetaPath() string {
	return filepath.Join(p.Path(), p.MetaFile())
}

// Name of the package.
// Must match the filesystem path (relative to the store) of the package
// with separators replaced with `:`. For example:
//
//	"/package/foo/bar/baz" => "foo:bar:baz"
func (p *Package) Name() string {
	return p.MetadataString("name")
}

// SetName sets the package name.
func (p *Package) SetName(value string) {
	p.Metadata["Name"] = value
}

// Description returns the description of the package.
func (p *Package) Description() string {
	return p.MetadataString("description")
}

// SetName sets the package description.
func (p *Package) SetDescription(value string) {
	p.Metadata["Description"] = value
}

// ShortDescription returns the first line of the description.
func (p *Package) ShortDescription() string {
	lines := strings.Split(p.Description(), "\n")
	return lines[0]
}

// Origin returns the path or URL used to install the package.
func (p *Package) Origin() string {
	return p.MetadataString("origin")
}

// SetOrigin sets the package origin.
func (p *Package) SetOrigin(value string) {
	p.Metadata["Origin"] = value
}

// Path returns the filesystem path of the package.
func (p *Package) Path() string {
	return p.path
}

// SetPath sets the package origin.
func (p *Package) SetPath(value string) {
	p.path = value
}

// Returns all nested sub-packages ordered by name.
func (p *Package) Children() ([]*Package, error) {
	return LoadPackages(p.Path(), p.MetaFile())
}

// Returns the receiver and all children.
func (p *Package) All() ([]*Package, error) {
	children, err := LoadPackages(p.Path(), p.MetaFile())
	if err != nil {
		return nil, err
	}
	packages := []*Package{}
	packages = append(packages, p)
	packages = append(packages, children...)
	return packages, nil
}

func (p *Package) Parent() *Package {
	parentPath := filepath.Dir(p.Path())
	pkg, _ := LoadPackage(parentPath, p.MetaFile())
	return pkg
}

func (p *Package) Root() *Package {
	node := p
	for {
		n := node.Parent()
		if n == nil {
			break
		}
		node = n
	}
	return node
}

// MetadataMapSlice returns a slice of maps for the given metadata key.
func (p *Package) MetadataMapSlice(key string) []map[string]any {
	items := []map[string]any{}
	for i, m := range p.MetadataSlice(key) {
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

// MetadataSlice returns a slice value for the given metadata key.
func (p *Package) MetadataSlice(key string) []any {
	val := p.MetadataLookup(key)
	switch v := val.(type) {
	case []any:
		return v
	case nil:
		return []any{}
	default:
		panic(NewMetadataTypeCastError(key, val, "[]any"))
	}
}

// MetadataString returns a string value for the given metadata key.
func (p *Package) MetadataString(key string) string {
	val := p.MetadataLookup(key)
	switch v := val.(type) {
	case string:
		return v
	case nil:
		return ""
	default:
		panic(NewMetadataTypeCastError(key, val, "string"))
	}
}

// MetadataLookup returns the value for key in Metadata.
// If not found, then tries the pascalized version of key.
func (p *Package) MetadataLookup(key string) any {
	if val, ok := p.Metadata[key]; ok {
		return val
	}

	if val, ok := p.Metadata[flect.Pascalize(key)]; ok {
		return val
	}
	return nil
}
