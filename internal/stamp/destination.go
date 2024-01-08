package stamp

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/cast"
	"github.com/swaggest/jsonschema-go"
	"github.com/twelvelabs/termite/render"

	"github.com/twelvelabs/stamp/internal/fsutil"
)

// NewDestinationWithValues returns a new destination set with the given values.
func NewDestinationWithValues(path string, mode string, values map[string]any) (Destination, error) {
	dst := Destination{}

	pathTpl, err := render.Compile(path)
	if err != nil {
		return dst, fmt.Errorf("new destination: %w", err)
	}
	dst.PathTpl = *pathTpl

	modeTpl, err := render.Compile(mode)
	if err != nil {
		return dst, fmt.Errorf("new destination: %w", err)
	}
	dst.ModeTpl = *modeTpl

	if err := dst.SetValues(values); err != nil {
		return dst, fmt.Errorf("new destination: %w", err)
	}

	return dst, nil
}

type Destination struct {
	ContentTypeTpl render.Template `mapstructure:"content_type" title:"Content Type"`
	Conflict       ConflictConfig  `mapstructure:"conflict"     title:"Conflict"  default:"prompt"`
	Missing        MissingConfig   `mapstructure:"missing"      title:"Missing"   default:"ignore"`
	ModeTpl        render.Template `mapstructure:"mode"         title:"Mode"`
	PathTpl        render.Template `mapstructure:"path"         title:"Path" required:"true"`

	content     any
	contentType FileType
	mode        os.FileMode
	path        string
}

// PrepareJSONSchema implements the jsonschema.Preparer interface.
func (d Destination) PrepareJSONSchema(schema *jsonschema.Schema) error {
	schema.WithTitle("Destination")
	schema.WithDescription("The destination path.")
	if prop, ok := schema.Properties["content_type"]; ok {
		prop.TypeObjectEns().
			WithDescription(d.contentType.Description()).
			WithEnum(d.contentType.Enum()...)
	}
	if prop, ok := schema.Properties["mode"]; ok {
		prop.TypeObjectEns().
			WithDefault("0666").
			WithDescription(
				"An optional [POSIX mode](https://en.wikipedia.org/wiki/File-system_permissions#Numeric_notation) "+
					"to set on the file path.",
			).
			WithExamples(
				"0755",
				"{{ .ModeValue }}",
			).
			WithPattern(`\{\{(.*)\}\}|\d{4}`) // https://rubular.com/r/t2lxVpKWQs5aeR
	}
	if prop, ok := schema.Properties["path"]; ok {
		prop.TypeObjectEns().
			WithDescription(
				"The file path relative to the destination directory. " +
					"Attempts to traverse outside the destination directory will raise a runtime error" +
					"\n\n" +
					"When creating new files, the [conflict](#conflict) attribute " +
					"will be used if the path already exists. " +
					"When updating or deleting files, the [missing](#missing) attribute " +
					"will be used if the path does not exist.",
			)
	}
	return nil
}

// Content returns the parsed file content.
func (d *Destination) Content() any {
	return d.content
}

// ContentBytes returns an encoded byte array of the content.
func (d *Destination) ContentBytes() ([]byte, error) {
	return d.contentType.Encoder().Encode(d.content)
}

// ContentType returns the content type of the file.
// If no content type is supplied, one will be inferred from the path.
func (d *Destination) ContentType() FileType {
	return d.contentType
}

// Mode returns the optional file mode to apply when writing
// the destination file.
func (d *Destination) Mode() os.FileMode {
	return d.mode
}

// Path returns the absolute path to the file.
func (d *Destination) Path() string {
	return d.path
}

// Exists returns true if the file path exists.
func (d *Destination) Exists() bool {
	return fsutil.PathExists(d.path)
}

// IsDir returns true if the file path is a directory.
func (d *Destination) IsDir() bool {
	return fsutil.PathIsDir(d.path)
}

// SetValues calculates destination properties using the given values.
func (d *Destination) SetValues(values map[string]any) error {
	var err error

	// Render and validate path.
	d.path, err = d.PathTpl.RenderRequired(values)
	if err != nil {
		return fmt.Errorf("dst path render: %w", err)
	}
	d.path, err = fsutil.EnsurePathRelativeToRoot(
		d.path, cast.ToString(values["DstPath"]),
	)
	if err != nil {
		return fmt.Errorf("dst path validate: %w", err)
	}

	// Render and parse content type.
	ct, err := d.ContentTypeTpl.Render(values)
	if err != nil {
		return fmt.Errorf("dst content_type render: %w", err)
	}
	d.contentType, err = ParseFileTypeWithFallback(ct, d.path)
	if err != nil {
		return fmt.Errorf("dst content_type parse: %w", err)
	}

	// Render and parse mode.
	mode, err := d.ModeTpl.Render(values)
	if err != nil {
		return fmt.Errorf("dst mode render: %w", err)
	}
	if mode != "" {
		parsed, err := strconv.ParseUint(mode, 8, 32)
		if err != nil {
			return fmt.Errorf("dst mode parse: %w", err)
		}
		d.mode = os.FileMode(parsed)
	} else if !d.Exists() {
		d.mode = DstFileMode
	}

	// Decode content.
	if d.Exists() && !d.IsDir() {
		content, err := os.ReadFile(d.path)
		if err != nil {
			return fmt.Errorf("dst path read: %w", err)
		}
		encoder := d.contentType.Encoder()
		d.content, err = encoder.Decode(content)
		if err != nil {
			return fmt.Errorf("dst path decode: %w", err)
		}
	}

	return nil
}

// Write encodes data and writes the resulting bytes
// to the destination file.
func (d *Destination) Write(data any) error {
	// Encode to byte array.
	buf, err := d.contentType.Encoder().Encode(data)
	if err != nil {
		return fmt.Errorf("dst encode: %w", err)
	}

	// Ensure base dirs.
	if err := os.MkdirAll(filepath.Dir(d.path), DstDirMode); err != nil {
		return err
	}

	// Write file
	f, err := os.Create(d.path)
	if err != nil {
		return fmt.Errorf("dst create: %w", err)
	}
	defer f.Close()
	_, err = f.Write(buf)
	if err != nil {
		return fmt.Errorf("dst write: %w", err)
	}

	// Set permissions (if configured).
	if d.mode != 0 {
		err = os.Chmod(d.path, d.mode)
		if err != nil {
			return fmt.Errorf("dst chmod: %w", err)
		}
	}

	// Set new content.
	d.content = data

	return nil
}

// Delete removes the destination file.
func (d *Destination) Delete() error {
	return os.RemoveAll(d.path)
}
