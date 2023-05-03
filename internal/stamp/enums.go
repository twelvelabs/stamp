package stamp

import (
	"path/filepath"
	"strings"
)

//go:generate go-enum -f=$GOFILE --marshal --names

// ArrayMerge determines how to merge arrays when updating structured data.
// ENUM(concat, upsert, replace).
type ArrayMerge string

// ConflictConfig determines what to do when destination paths already exist.
// ENUM(keep, replace, prompt).
type ConflictConfig string

// ENUM(file, line).
type MatchSource string

// MissingConfig determines what to do when destination paths are missing.
// ENUM(ignore, error).
type MissingConfig string

// ENUM(json, yaml, text).
type FileType string

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

// IsStructured returns true if the receiver is JSON or YAML.
func (ft FileType) IsStructured() bool {
	switch ft {
	case FileTypeJson, FileTypeYaml:
		return true
	default:
		return false
	}
}
