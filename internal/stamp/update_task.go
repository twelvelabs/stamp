package stamp

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
	"github.com/spf13/cast"

	"github.com/twelvelabs/stamp/internal/fsutil"
	"github.com/twelvelabs/stamp/internal/modify"
)

type UpdateTask struct {
	Common `mapstructure:",squash"`

	Dst     string        `validate:"required"`
	Missing MissingConfig `validate:"required" default:"ignore"`
	Mode    string        `validate:"omitempty,posix-mode"`
	Parse   string        ``
	Pattern string        `validate:"required"`
	Action  modify.Action `validate:"required" default:"replace"`
	Content any           ``
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
	var replacement any
	if s, ok := t.Content.(string); ok {
		replacement = t.Render(s, values)
	} else {
		replacement = t.Content
	}
	// resolve the parse field
	parse := t.Render(t.Parse, values)
	if parse == "" {
		// An unspecified parse value implies plain text.
		parse = "text"
	} else if parse == "true" {
		// "parse: true" is shorthand for "figure out file type from the extension".
		parse = strings.TrimPrefix(filepath.Ext(dst), ".")
	}

	// Update the content (using the pattern and replacement values)
	switch parse {
	case "json":
		content, err = t.replaceJSON(content, pattern, replacement)
	case "yaml", "yml":
		content, err = t.replaceYAML(content, pattern, replacement)
	case "text", "txt":
		content, err = t.replaceText(content, pattern, replacement)
	default:
		return fmt.Errorf("unable to parse file type: %s", parse)
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

func (t *UpdateTask) replaceJSON(content []byte, pattern string, repl any) ([]byte, error) {
	// parse the JSON into a data structure
	data, err := oj.Parse(content)
	if err != nil {
		return nil, err
	}

	// parse the JSON path expression
	exp, err := jp.ParseString(pattern)
	if err != nil {
		return nil, err
	}

	// modify the data structure
	if t.Action == modify.ActionDelete {
		_, err := exp.Remove(data)
		if err != nil {
			return nil, err
		}
	} else {
		_, err := exp.Modify(data, modify.Modifier(t.Action, repl))
		if err != nil {
			return nil, err
		}
	}

	// convert back to JSON
	marshalled, err := oj.Marshal(data)
	if err != nil {
		return nil, err
	}
	return marshalled, nil
}

func (t *UpdateTask) replaceYAML(content []byte, pattern string, repl any) ([]byte, error) {
	return nil, nil
}

func (t *UpdateTask) replaceText(content []byte, pattern string, repl any) ([]byte, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("pattern: %w", err)
	}
	replStr, err := cast.ToStringE(repl)
	if err != nil {
		return nil, fmt.Errorf("replacement: %w", err)
	}

	var replacement []byte
	switch t.Action {
	case modify.ActionAppend:
		replacement = append([]byte("$0"), []byte(replStr)...)
	case modify.ActionPrepend:
		replacement = append([]byte(replStr), []byte("$0")...)
	case modify.ActionReplace:
		replacement = []byte(replStr)
	case modify.ActionDelete:
		replacement = []byte{}
	}

	return re.ReplaceAll(content, replacement), nil
}
