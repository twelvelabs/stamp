package stamp

import (
	"path/filepath"
	"strings"

	"github.com/twelvelabs/stamp/internal/encode"
)

//go:generate go-enum -f=$GOFILE -t ../enums.tmpl --marshal --names

// ConflictConfig determines what to do when destination paths already exist.
// ENUM(keep, replace, prompt).
type ConflictConfig string

// MatchSource determines whether match patterns should be applied per-line or to the entire file.
// ENUM(file, line).
type MatchSource string

// MissingConfig determines what to do when destination paths are missing.
// ENUM(ignore, touch, error).
type MissingConfig string

// FileType specifies the content type of the destination path.
// ENUM(json, yaml, text).
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
