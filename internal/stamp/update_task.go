package stamp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
	"github.com/spf13/cast"
	"github.com/twelvelabs/termite/render"
	"gopkg.in/yaml.v3"

	"github.com/twelvelabs/stamp/internal/fsutil"
	"github.com/twelvelabs/stamp/internal/modify"
)

const (
	fileTypeJSON = "json"
	fileTypeYAML = "yaml"
	fileTypeYML  = "yml"
)

type UpdateTask struct {
	Common `mapstructure:",squash"`

	Src        string        `mapstructure:"src"`
	SrcContent any           `mapstructure:"src_content"`
	Dst        string        `mapstructure:"dst"      validate:"required"`
	Missing    MissingConfig `mapstructure:"missing"  validate:"required" default:"ignore"`
	Mode       string        `mapstructure:"mode"     validate:"omitempty,posix-mode"`
	Pattern    string        `mapstructure:"pattern"`
	Action     modify.Action `mapstructure:"action"   validate:"required" default:"replace"`
	FileType   string        `mapstructure:"file_type"`

	dstPath     string
	dstBytes    []byte
	mode        os.FileMode
	pattern     string
	replacement any
}

func (t *UpdateTask) Execute(ctx *TaskContext, values map[string]any) error {
	err := t.prepare(ctx, values)
	if err != nil {
		ctx.Logger.Failure("fail", t.dstPath)
		return err
	}

	if fsutil.PathExists(t.dstPath) {
		if err := t.updateDst(ctx, values, t.dstPath); err != nil {
			ctx.Logger.Failure("fail", t.dstPath)
			return err
		}
		updateMsg := t.dstPath
		if t.isStructured(t.FileType) {
			updateMsg = fmt.Sprintf("%s (%s)", t.dstPath, t.pattern)
		}
		ctx.Logger.Success("update", updateMsg)
		return nil
	} else if t.Missing == MissingConfigError {
		ctx.Logger.Failure("fail", t.dstPath)
		return ErrPathNotFound
	}

	return nil
}

func (t *UpdateTask) isStructured(fileType string) bool {
	switch fileType {
	case fileTypeJSON, fileTypeYAML, fileTypeYML:
		return true
	default:
		return false
	}
}

// prepare post-processes and validates the task YAML fields.
func (t *UpdateTask) prepare(_ *TaskContext, values map[string]any) error {
	var err error

	dstRoot := cast.ToString(values["DstPath"])
	t.dstPath, err = t.RenderPath("dst", t.Dst, dstRoot, values)
	if err != nil {
		return fmt.Errorf("resolve dst path: %w", err)
	}
	if fsutil.PathExists(t.dstPath) {
		t.dstBytes, err = os.ReadFile(t.dstPath)
		if err != nil {
			return fmt.Errorf("read dst path: %w", err)
		}
	}

	// Depends on: DstPath
	if t.FileType == "" {
		// No explicit file type provided, infer from file extension.
		t.FileType = strings.ToLower(strings.TrimPrefix(filepath.Ext(t.dstPath), "."))
	}

	// Depends on: FileType
	if t.Src != "" && t.SrcContent != nil { //nolint:nestif
		return fmt.Errorf("src and src_content fields are mutually exclusive")
	} else if t.Src != "" {
		srcRoot := cast.ToString(values["SrcPath"])
		srcPath, err := t.RenderPath("src", t.Src, srcRoot, values)
		if err != nil {
			return fmt.Errorf("resolve src path: %w", err)
		}
		srcContent, err := render.File(srcPath, values)
		if err != nil {
			return fmt.Errorf("render src path: %w", err)
		}
		t.replacement = srcContent
		if t.isStructured(t.FileType) {
			t.replacement, err = t.parse([]byte(srcContent))
			if err != nil {
				return fmt.Errorf("parse src path: %w", err)
			}
		}
	} else if s, ok := t.SrcContent.(string); ok {
		t.replacement = t.Render(s, values)
	} else {
		t.replacement = t.SrcContent
	}

	if t.Mode != "" {
		t.mode, err = t.RenderMode(t.Mode, values)
		if err != nil {
			return fmt.Errorf("resolve dst mode: %w", err)
		}
	}

	// Depends on: FileType
	t.pattern = t.Render(t.Pattern, values)
	if t.pattern == "" {
		// Match the entire file if pattern is empty.
		if t.isStructured(t.FileType) {
			t.pattern = "$"
		} else {
			t.pattern = "(?s)^(.*)$"
		}
	}

	return nil
}

func (t *UpdateTask) updateDst(ctx *TaskContext, _ map[string]any, _ string) error {
	if ctx.DryRun {
		return nil
	}

	// Update the content (using the pattern and replacement values)
	var err error
	if t.isStructured(t.FileType) {
		t.dstBytes, err = t.replaceStructured(t.dstBytes, t.pattern, t.replacement)
	} else {
		t.dstBytes, err = t.replaceText(t.dstBytes, t.pattern, t.replacement)
	}
	if err != nil {
		return fmt.Errorf("update content: %w", err)
	}

	// update dst
	f, err := os.Create(t.dstPath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(t.dstBytes)
	if err != nil {
		return err
	}

	// set permissions (if configured)
	if t.mode != 0 {
		err = os.Chmod(t.dstPath, t.mode)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *UpdateTask) parse(content []byte) (any, error) {
	var data any
	var err error

	switch t.FileType {
	case fileTypeJSON:
		data, err = oj.Parse(content)
		if err != nil {
			return nil, fmt.Errorf("json parse: %w", err)
		}
	case fileTypeYAML, fileTypeYML:
		err := yaml.Unmarshal(content, &data)
		if err != nil {
			return nil, fmt.Errorf("yaml parse: %w", err)
		}
	default:
		data = content
	}

	return data, nil
}

func (t *UpdateTask) marshal(data any) ([]byte, error) {
	var content []byte
	var err error

	switch t.FileType {
	case fileTypeJSON:
		// Note: using standard lib to marshal because it sorts JSON object keys
		// (oj does not and it looks ugly when adding new keys).
		content, err = json.MarshalIndent(data, "", "    ")
		if err != nil {
			return nil, fmt.Errorf("json marshal: %w", err)
		}
	case fileTypeYAML, fileTypeYML:
		b := &bytes.Buffer{}
		encoder := yaml.NewEncoder(b)
		encoder.SetIndent(2)
		err = encoder.Encode(data)
		if err != nil {
			return nil, fmt.Errorf("yaml marshal: %w", err)
		}
		content = b.Bytes()
	default:
		content = data.([]byte)
	}

	return content, nil
}

func (t *UpdateTask) replaceStructured(content []byte, pattern string, repl any) ([]byte, error) {
	data, err := t.parse(content)
	if err != nil {
		return nil, err
	}

	// parse pattern as a JSON path expression
	exp, err := jp.ParseString(pattern)
	if err != nil {
		return nil, fmt.Errorf("json path parse: %w", err)
	}
	// use the expression to modify the data structure
	if t.Action == modify.ActionDelete {
		data, err = exp.Remove(data)
		if err != nil {
			return nil, fmt.Errorf("json path remove: %w", err)
		}
	} else {
		data, err = exp.Modify(data, modify.Modifier(t.Action, repl))
		if err != nil {
			return nil, fmt.Errorf("json path modify: %w", err)
		}
	}

	marshalled, err := t.marshal(data)
	if err != nil {
		return nil, err
	}

	return marshalled, nil
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
		replacement = append([]byte("${0}"), []byte(replStr)...)
	case modify.ActionPrepend:
		replacement = append([]byte(replStr), []byte("${0}")...)
	case modify.ActionReplace:
		replacement = []byte(replStr)
	case modify.ActionDelete:
		replacement = []byte{}
	}

	return re.ReplaceAll(content, replacement), nil
}
