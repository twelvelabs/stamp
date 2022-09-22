package pkg

import (
	"path/filepath"
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
	val, _ := p.Metadata["Name"]
	if val, ok := val.(string); ok {
		return val
	}
	return ""
}

// SetName sets the package name.
func (p *Package) SetName(value string) {
	p.Metadata["Name"] = value
}

// Origin returns the path or URL used to install the package.
func (p *Package) Origin() string {
	val, _ := p.Metadata["Origin"]
	if val, ok := val.(string); ok {
		return val
	}
	return ""
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
