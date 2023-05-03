package stamp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/mitchellh/mapstructure"
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

	Src         any           `mapstructure:"src"`
	Dst         string        `mapstructure:"dst"      validate:"required"`
	Match       any           `mapstructure:"match"`
	Missing     MissingConfig `mapstructure:"missing"  validate:"required" default:"ignore"`
	Mode        string        `mapstructure:"mode"     validate:"omitempty"`
	Action      modify.Action `mapstructure:"action"   validate:"required" default:"replace"`
	FileType    string        `mapstructure:"file_type"`
	Description string        `mapstructure:"description"`
	Upsert      bool          `mapstructure:"upsert"`

	dstPath     string
	dstBytes    []byte
	description string
	fileType    FileType
	match       matchConfig
	mode        os.FileMode
	replacement any
}

type matchConfig struct {
	Path    string `mapstructure:"path"`
	Default any    `mapstructure:"default"`
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
		if t.description != "" {
			updateMsg = fmt.Sprintf("%s (%s)", t.dstPath, t.description)
		} else if t.fileType.IsStructured() {
			updateMsg = fmt.Sprintf("%s (%s)", t.dstPath, t.match.Path)
		}
		ctx.Logger.Success("update", updateMsg)
		return nil
	} else if t.Missing == MissingConfigError {
		ctx.Logger.Failure("fail", t.dstPath)
		return ErrPathNotFound
	}

	return nil
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
	if t.FileType != "" {
		t.fileType, err = ParseFileType(t.FileType)
	} else {
		t.fileType, err = ParseFileTypeFromPath(t.dstPath)
	}
	if err != nil {
		return fmt.Errorf("parse file_type: %w", err)
	}

	// Src can be a few different things:
	// - A path to a template file containing the source content.
	// - A string literal to render and use as source content.
	// - Structured data.
	// Depends on: FileType
	if src, ok := t.Src.(string); ok { //nolint:nestif
		// Try rendering as a file path.
		srcRoot := cast.ToString(values["SrcPath"])
		srcPath, err := t.RenderPath("src", src, srcRoot, values)
		if err != nil {
			return fmt.Errorf("resolve src path: %w", err)
		}

		if fsutil.PathExists(srcPath) {
			// Render the content located at the path.
			srcContent, err := render.File(srcPath, values)
			if err != nil {
				return fmt.Errorf("render src path: %w", err)
			}
			// Parse the src content if we need structured data,
			// otherwise the content itself is the replacement.
			t.replacement = srcContent
			// TODO: move this
			if t.fileType.IsStructured() {
				t.replacement, err = t.parse([]byte(srcContent))
				if err != nil {
					return fmt.Errorf("parse src path: %w", err)
				}
			}
		} else {
			// Path didn't exist, just render as content.
			t.replacement = t.Render(src, values)
		}
	} else {
		// Not a string, so must be structured data.
		t.replacement = t.Src
	}

	if t.Mode != "" {
		t.mode, err = t.RenderMode(t.Mode, values)
		if err != nil {
			return fmt.Errorf("parse mode: %w", err)
		}
	}

	// Match config can be provided as either a string or an object in YAML.
	t.match = matchConfig{}
	if m, ok := t.Match.(map[string]any); ok {
		// match:
		//   path: $.items
		//   default: []
		err := mapstructure.Decode(m, &t.match)
		if err != nil {
			return fmt.Errorf("parse match: %w", err)
		}
	} else {
		// match: $.items
		// match: ^foo(\w+)$
		// match: 123
		t.match.Path = cast.ToString(t.Match)
	}

	// Render the match path (default to matching everything if none provided).
	// Depends on: FileType
	t.match.Path = t.Render(t.match.Path, values)
	if t.match.Path == "" {
		if t.fileType.IsStructured() {
			t.match.Path = "$" // root node
		} else {
			t.match.Path = "(?s)^(.*)$" // `?s` causes . to match newlines
		}
	}

	t.description = t.Render(t.Description, values)

	return nil
}

func (t *UpdateTask) updateDst(ctx *TaskContext, _ map[string]any, _ string) error {
	if ctx.DryRun {
		return nil
	}

	// Update the content (using the pattern and replacement values)
	var err error
	if t.fileType.IsStructured() {
		t.dstBytes, err = t.replaceStructured(t.dstBytes, t.match.Path, t.replacement)
	} else {
		t.dstBytes, err = t.replaceText(t.dstBytes, t.match.Path, t.replacement)
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

	switch t.fileType {
	case FileTypeJson:
		data, err = oj.Parse(content)
		if err != nil {
			return nil, fmt.Errorf("json parse: %w", err)
		}
	case FileTypeYaml:
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

	switch t.fileType {
	case FileTypeJson:
		// Note: using standard lib to marshal because it sorts JSON object keys
		// (oj does not and it looks ugly when adding new keys).
		content, err = json.MarshalIndent(data, "", "    ")
		if err != nil {
			return nil, fmt.Errorf("json marshal: %w", err)
		}
	case FileTypeYaml:
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
	if t.Action == modify.ActionDelete { //nolint:nestif
		data, err = exp.Remove(data)
		if err != nil {
			return nil, fmt.Errorf("json path remove: %w", err)
		}
	} else {
		if !exp.Has(data) && t.match.Default != nil {
			err = exp.Set(data, t.match.Default)
			if err != nil {
				return nil, fmt.Errorf("json path set default: %w", err)
			}
		}
		modifier := modify.Modifier(t.Action, repl, modify.WithUpsert(t.Upsert))
		data, err = exp.Modify(data, modifier)
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
	reMatchBytes := []byte("${0}")

	replStr, err := cast.ToStringE(repl)
	if err != nil {
		return nil, fmt.Errorf("replacement: %w", err)
	}
	replBytes := []byte(replStr)

	shouldPerformAction := !t.Upsert || (t.Upsert && !bytes.Contains(content, replBytes))

	var replacement []byte
	switch t.Action {
	case modify.ActionAppend:
		replacement = append(replacement, reMatchBytes...)
		if shouldPerformAction {
			replacement = append(replacement, replBytes...)
		}
	case modify.ActionPrepend:
		if shouldPerformAction {
			replacement = append(replacement, replBytes...)
		}
		replacement = append(replacement, reMatchBytes...)
	case modify.ActionReplace:
		replacement = append(replacement, replBytes...)
	case modify.ActionDelete:
		replacement = []byte{}
	}

	return re.ReplaceAll(content, replacement), nil
}
