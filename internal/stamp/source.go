package stamp

import (
	"fmt"

	"github.com/spf13/cast"
	"github.com/swaggest/jsonschema-go"
	"github.com/twelvelabs/termite/render"

	"github.com/twelvelabs/stamp/internal/fsutil"
)

// NewSourceWithValues returns a new source set with the given values.
func NewSourceWithValues(path string, values map[string]any) (Source, error) {
	src := Source{}

	pathTpl, err := render.Compile(path)
	if err != nil {
		return src, fmt.Errorf("new source: %w", err)
	}
	src.PathTpl = *pathTpl

	if err := src.SetValues(values); err != nil {
		return src, fmt.Errorf("new source: %w", err)
	}

	return src, nil
}

type Source struct {
	ContentTypeTpl render.Template `mapstructure:"content_type"`
	InlineContent  any             `mapstructure:"content"`
	PathTpl        render.Template `mapstructure:"path"`

	content     any
	contentType FileType
	path        string
}

// SourceWithContent represents one version of Source in the JSON schema.
type SourceWithContent struct {
	ContentType FileType `mapstructure:"content_type"`
	Content     any      `mapstructure:"content" required:"true" description:"Inline content. Can be any type. String keys and/or values will be rendered as templates."` //nolint: lll
}

// SourceWithPath represents one version of Source in the JSON schema.
type SourceWithPath struct {
	ContentType FileType `mapstructure:"content_type"`
	Path        string   `mapstructure:"path" required:"true" description:"The file path relative to the source directory. Attempts to traverse outside the source directory will raise a runtime error."` //nolint: lll
}

var _ jsonschema.Preparer = Source{}

// PrepareJSONSchema implements the jsonschema.Preparer interface.
func (s Source) PrepareJSONSchema(schema *jsonschema.Schema) error {
	schema.WithDescription("The source path or inline content.")
	// Reset properties and just rely on `oneOf`.
	schema.WithProperties(map[string]jsonschema.SchemaOrBool{})
	return nil
}

var _ jsonschema.OneOfExposer = Source{}

// PrepareJSONSchema implements the jsonschema.OneOfExposer interface.
func (s Source) JSONSchemaOneOf() []any {
	return []any{
		SourceWithContent{},
		SourceWithPath{},
	}
}

// Content returns the parsed file content.
func (s *Source) Content() any {
	return s.content
}

// ContentBytes returns an encoded byte array of the content.
func (s *Source) ContentBytes() ([]byte, error) {
	return s.contentType.Encoder().Encode(s.content)
}

// ContentType returns the content type of the file.
// If no content type is supplied, one will be inferred from the path.
func (s *Source) ContentType() FileType {
	return s.contentType
}

// Path returns the absolute path to the file.
func (s *Source) Path() string {
	return s.path
}

// Exists returns true if the file path exists.
func (s *Source) Exists() bool {
	return fsutil.PathExists(s.path)
}

// IsDir returns true if the file path is a directory.
func (s *Source) IsDir() bool {
	return fsutil.PathIsDir(s.path)
}

// SetValues calculates source properties using the given values.
func (s *Source) SetValues(values map[string]any) error {
	var err error

	// Render path.
	s.path, err = s.PathTpl.Render(values)
	if err != nil {
		return fmt.Errorf("src path render: %w", err)
	}
	// Render inline content.
	s.content, err = render.Any(s.InlineContent, values)
	if err != nil {
		return fmt.Errorf("src inline content render: %w", err)
	}
	// Render and parse content type.
	ct, err := s.ContentTypeTpl.Render(values)
	if err != nil {
		return fmt.Errorf("src content_type render: %w", err)
	}
	s.contentType, err = ParseFileTypeWithFallback(ct, s.path)
	if err != nil {
		return fmt.Errorf("src content_type parse: %w", err)
	}

	// Ensure mutually exclusive fields.
	if s.content != nil && s.path != "" {
		return fmt.Errorf("src: path and content are mutually exclusive")
	}

	// Bail early if path is empty (everything below requires one).
	if s.path == "" {
		return nil
	}

	// Validate path.
	s.path, err = fsutil.EnsurePathRelativeToRoot(
		s.path, cast.ToString(values["SrcPath"]),
	)
	if err != nil {
		return fmt.Errorf("src path validate: %w", err)
	}

	if s.Exists() && !s.IsDir() {
		// Render the content located at the path.
		rendered, err := render.File(s.path, values)
		if err != nil {
			return fmt.Errorf("src path content render: %w", err)
		}

		// Decode the rendered content.
		encoder := s.contentType.Encoder()
		s.content, err = encoder.Decode([]byte(rendered))
		if err != nil {
			return fmt.Errorf("src path content decode: %w", err)
		}
	}

	return nil
}
