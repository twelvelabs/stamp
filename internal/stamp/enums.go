package stamp

import (
	"path/filepath"
	"strings"

	"github.com/twelvelabs/stamp/internal/encode"
)

// cspell: words: createtask updatetask
//go:generate go-enum -f=$GOFILE -t ../enums.tmpl --marshal --names --nocomments

// Determines what to do when creating a new file and
// the destination path already exists.
//
// > [!IMPORTANT]
// > Only used in [create] tasks.
//
// [create]: https://github.com/twelvelabs/stamp/tree/main/docs/create_task.md
/*
	ENUM(
		keep     // Keep the existing path. The task becomes a noop.
		replace  // Replace the existing path.
		prompt   // Prompt the user.
	).
*/
type ConflictConfig string

// Determines how regexp patterns should be applied.
/*
	ENUM(
		file  // Match the entire file.
		line  // Match each line.
	).
*/
type MatchSource string

// Determines what to do when updating an existing file and
// the destination path is missing.
//
// > [!IMPORTANT]
// > Only used in [update] and [delete] tasks.
//
// [update]: https://github.com/twelvelabs/stamp/tree/main/docs/update_task.md
// [delete]: https://github.com/twelvelabs/stamp/tree/main/docs/delete_task.md
/*
	ENUM(
		ignore  // Do nothing. The task becomes a noop.
		touch   // Create an empty file.
		error   // Raise an error.
	).
*/
type MissingConfig string

// Determines the visibility of the generator.
/*
	ENUM(
		public   // Callable anywhere.
		hidden   // Public, but hidden in the generator list.
		private  // Only callable as a sub-generator. Never displayed.
	).
*/
type VisibilityType string

// Specifies the content type of the file.
// Inferred from the file extension by default.
//
// When the content type is JSON or YAML, the file will be
// parsed into a data structure before use.
// When updating files, the content type determines
// the behavior of the [match.pattern] attribute.
//
// [match.pattern]: https://github.com/twelvelabs/stamp/tree/main/docs/match.md#pattern
/*
	ENUM(
		json
		yaml
		text
	).
*/
type FileType string

// Encoder returns the encoder for this content type.
func (ft FileType) Encoder() encode.Encoder {
	switch ft {
	case FileTypeJson:
		return &encode.JSONEncoder{}
	case FileTypeYaml:
		return &encode.YAMLEncoder{}
	default: // FileTypeText
		return &encode.TextEncoder{}
	}
}

// IsStructured returns true if the receiver is JSON or YAML.
func (ft FileType) IsStructured() bool {
	switch ft {
	case FileTypeJson, FileTypeYaml:
		return true
	default:
		return false
	}
}

// ParseFileTypeWithFallback parses value into a FileType or,
// if value is empty, attempts to infer the FileType from the given path.
func ParseFileTypeWithFallback(value string, path string) (FileType, error) {
	if value != "" {
		return ParseFileType(value)
	}
	return ParseFileTypeFromPath(path)
}

// ParseFileTypeFromPath returns the correct file type for the given path.
func ParseFileTypeFromPath(path string) (FileType, error) {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(path), "."))
	switch ext {
	case "json":
		return FileTypeJson, nil
	case "yaml", "yml":
		return FileTypeYaml, nil
	default:
		return FileTypeText, nil
	}
}
