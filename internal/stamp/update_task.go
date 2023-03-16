package stamp

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/twelvelabs/stamp/internal/fsutil"
)

type UpdateTask struct {
	Common `mapstructure:",squash"`

	Dst     string        `validate:"required"`
	Missing MissingConfig `validate:"required" default:"ignore"`
	Mode    string        `validate:"omitempty,posix-mode"`
	Parse   string        ``
	Pattern string        `validate:"required"`
	Action  UpdateAction  `validate:"required" default:"replace"`
	Content string        ``
}

func (t *UpdateTask) Execute(ctx *TaskContext, values map[string]any) error {
	dst, err := t.RenderRequired("dst", t.Dst, values)
	if err != nil {
		ctx.Logger.Failure("fail", dst)
		return err
	}

	if fsutil.PathExists(dst) {
		if err := t.updateDst(ctx, values, dst); err != nil {
			ctx.Logger.Failure("fail", dst)
			return err
		}
		ctx.Logger.Success("update", dst)
		return nil
	} else if t.Missing == MissingConfigError {
		ctx.Logger.Failure("fail", dst)
		return ErrPathNotFound
	}

	return nil
}

func (t *UpdateTask) updateDst(ctx *TaskContext, values map[string]any, dst string) error {
	if ctx.DryRun {
		return nil
	}

	// resolve dst content
	content, err := os.ReadFile(dst)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}
	// resolve pattern
	pattern, err := t.RenderRequired("pattern", t.Pattern, values)
	if err != nil {
		return fmt.Errorf("resolve pattern: %w", err)
	}
	// resolve the replacement value
	replacement := t.Render(t.Content, values)
	// resolve the parse field
	parse := t.Render(t.Parse, values)
	if parse == "true" {
		// "parse: true" is shorthand for "figure out file type from the extension".
		parse = strings.TrimPrefix(filepath.Ext(dst), ".")
	}

	// Update the content (using the pattern and replacement values)
	switch parse {
	case "json":
		content, err = t.replaceJSON(content, pattern, replacement)
	case "yaml", "yml":
		content, err = t.replaceYAML(content, pattern, replacement)
	default:
		content, err = t.replaceText(content, pattern, replacement)
	}
	if err != nil {
		return fmt.Errorf("update content: %w", err)
	}

	// update dst
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(content)
	if err != nil {
		return err
	}

	// set permissions (if configured)
	if t.Mode != "" {
		mode, err := t.RenderMode(t.Mode, values)
		if err != nil {
			return fmt.Errorf("resolve mode: %w", err)
		}
		err = os.Chmod(dst, mode)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *UpdateTask) replaceJSON(content []byte, pattern, replacement string) ([]byte, error) {
	return nil, nil
}

func (t *UpdateTask) replaceYAML(content []byte, pattern, replacement string) ([]byte, error) {
	return nil, nil
}

func (t *UpdateTask) replaceText(content []byte, pattern, repl string) ([]byte, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("compile pattern: %w", err)
	}

	var replacement []byte
	switch t.Action {
	case UpdateActionAppend:
		replacement = append([]byte("$0"), []byte(repl)...)
	case UpdateActionPrepend:
		replacement = append([]byte(repl), []byte("$0")...)
	case UpdateActionReplace:
		replacement = []byte(repl)
	case UpdateActionDelete:
		replacement = []byte{}
	}

	return re.ReplaceAll(content, replacement), nil
}
