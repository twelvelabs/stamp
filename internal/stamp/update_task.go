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
		if t.isStructuredFileType() {
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

func (t *UpdateTask) isStructuredFileType() bool {
	switch t.FileType {
	case "json", "yaml", "yml":
		return true
	default:
		return false
	}
}

// prepare post-processes and validates the task YAML fields.
func (t *UpdateTask) prepare(_ *TaskContext, values map[string]any) error {
	var err error

	if t.Src != "" && t.SrcContent != nil {
		return fmt.Errorf("src and src_content fields are mutually exclusive")
	}
	if t.Src != "" {
		srcRoot := cast.ToString(values["SrcPath"])
		srcPath, err := t.RenderPath("src", t.Src, srcRoot, values)
		if err != nil {
			return fmt.Errorf("resolve src path: %w", err)
		}
		// render the src template
		t.replacement, err = render.File(srcPath, values)
		if err != nil {
			return fmt.Errorf("render src path: %w", err)
		}
	} else if s, ok := t.SrcContent.(string); ok {
		t.replacement = t.Render(s, values)
	} else {
		t.replacement = t.SrcContent
	}

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

	if t.Mode != "" {
		t.mode, err = t.RenderMode(t.Mode, values)
		if err != nil {
			return fmt.Errorf("resolve dst mode: %w", err)
		}
	}

	if t.FileType == "" {
		// No explicit file type provided, infer from file extension.
		t.FileType = strings.TrimPrefix(filepath.Ext(t.dstPath), ".")
	}

	t.pattern = t.Render(t.Pattern, values)
	if t.pattern == "" {
		// Match the entire file if pattern is empty.
		if t.isStructuredFileType() {
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
	switch t.FileType {
	case "json":
		t.dstBytes, err = t.replaceJSON(t.dstBytes, t.pattern, t.replacement)
	case "yaml", "yml":
		t.dstBytes, err = t.replaceYAML(t.dstBytes, t.pattern, t.replacement)
	default:
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

func (t *UpdateTask) replaceJSON(content []byte, pattern string, repl any) ([]byte, error) {
	// parse the JSON into a data structure
	data, err := oj.Parse(content)
	if err != nil {
		return nil, fmt.Errorf("json parse: %w", err)
	}

	// modify the data structure
	data, err = t.modifyDataStructure(data, pattern, t.Action, repl)
	if err != nil {
		return nil, err
	}

	// convert back to JSON
	// Note: using standard lib to marshal because it sorts JSON object keys
	// (oj does not and it looks ugly when adding new keys).
	marshalled, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("json marshal: %w", err)
	}

	return marshalled, nil
}

func (t *UpdateTask) replaceYAML(content []byte, pattern string, repl any) ([]byte, error) {
	// parse the YAML into a data structure
	var data any
	err := yaml.Unmarshal(content, &data)
	if err != nil {
		return nil, err
	}

	// modify the data structure
	data, err = t.modifyDataStructure(data, pattern, t.Action, repl)
	if err != nil {
		return nil, err
	}

	// convert back to YAML
	buf := &bytes.Buffer{}
	encoder := yaml.NewEncoder(buf)
	encoder.SetIndent(2)
	err = encoder.Encode(data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (t *UpdateTask) modifyDataStructure(data any, pattern string, act modify.Action, repl any) (any, error) {
	// parse pattern as a JSON path expression
	exp, err := jp.ParseString(pattern)
	if err != nil {
		return nil, fmt.Errorf("json path parse: %w", err)
	}
	// use the expression to modify the data structure
	if act == modify.ActionDelete {
		data, err = exp.Remove(data)
		if err != nil {
			return nil, fmt.Errorf("json path remove: %w", err)
		}
	} else {
		data, err = exp.Modify(data, modify.Modifier(act, repl))
		if err != nil {
			return nil, fmt.Errorf("json path modify: %w", err)
		}
	}
	return data, nil
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
