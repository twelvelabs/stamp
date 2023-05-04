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

	Action      any           `mapstructure:"action"   default:"replace"`
	Description string        `mapstructure:"description"`
	Dst         string        `mapstructure:"dst"      validate:"required"`
	FileType    string        `mapstructure:"file_type"`
	Match       any           `mapstructure:"match"`
	Missing     MissingConfig `mapstructure:"missing"  validate:"required" default:"ignore"`
	Mode        string        `mapstructure:"mode"     validate:"omitempty"`
	Src         any           `mapstructure:"src"`

	action      actionConfig
	description string
	dstBytes    []byte
	dstPath     string
	fileType    FileType
	match       matchConfig
	mode        os.FileMode
	src         any
}

type actionConfig struct {
	Type      modify.Action    `mapstructure:"type"`
	MergeType modify.MergeType `mapstructure:"array_merge"`
}

type matchConfig struct {
	Pattern string      `mapstructure:"pattern"`
	Default any         `mapstructure:"default"`
	Source  MatchSource `mapstructure:"source"`
}

func (t *UpdateTask) Execute(ctx *TaskContext, values map[string]any) error {
	err := t.prepare(ctx, values)
	if err != nil {
		ctx.Logger.Failure("fail", t.dstPath)
		return err
	}

	if fsutil.PathExists(t.dstPath) {
		if err := t.updateDst(ctx); err != nil {
			ctx.Logger.Failure("fail", t.dstPath)
			return err
		}
		updateMsg := t.dstPath
		if t.description != "" {
			updateMsg = fmt.Sprintf("%s (%s)", t.dstPath, t.description)
		} else if t.fileType.IsStructured() {
			updateMsg = fmt.Sprintf("%s (%s)", t.dstPath, t.match.Pattern)
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
			t.src = srcContent
			// TODO: move this
			if t.fileType.IsStructured() {
				t.src, err = t.parse([]byte(srcContent))
				if err != nil {
					return fmt.Errorf("parse src path: %w", err)
				}
			}
		} else {
			// Path didn't exist, just render as content.
			t.src, err = render.String(src, values)
			if err != nil {
				return fmt.Errorf("render src: %w", err)
			}
		}
	} else {
		// Not a string, so must be structured data.
		t.src, err = render.Any(t.Src, values)
		if err != nil {
			return fmt.Errorf("render src: %w", err)
		}
	}

	if t.Mode != "" {
		t.mode, err = t.RenderMode(t.Mode, values)
		if err != nil {
			return fmt.Errorf("parse mode: %w", err)
		}
	}

	// Action config can be provided as either a string or an object in YAML.
	t.action = actionConfig{
		Type:      modify.ActionReplace,
		MergeType: modify.MergeTypeConcat,
	}
	if obj, ok := t.Action.(map[string]any); ok {
		// action:
		//   type: append
		//   array_merge: upsert
		err = mapstructure.Decode(obj, &t.action)
	} else if str, ok := t.Action.(string); ok {
		// action: append
		t.action.Type, err = modify.ParseAction(str)
	}
	if err != nil {
		return fmt.Errorf("parse action: %w", err)
	}

	// Match config can be provided as either a string or an object in YAML.
	t.match = matchConfig{}
	if obj, ok := t.Match.(map[string]any); ok {
		// match:
		//   pattern: $.items
		//   default: []
		err = mapstructure.Decode(obj, &t.match)
		if err != nil {
			return fmt.Errorf("parse match: %w", err)
		}
	} else {
		// match: $.items
		// match: ^foo(\w+)$
		// match: 123
		t.match.Pattern = cast.ToString(t.Match)
	}

	// Render the match path (default to matching everything if none provided).
	// Depends on: FileType
	t.match.Pattern = t.Render(t.match.Pattern, values)
	if t.match.Pattern == "" {
		if t.fileType.IsStructured() {
			t.match.Pattern = "$" // root node
		} else {
			t.match.Source = MatchSourceFile
			t.match.Pattern = "(?s)^(.*)$" // `?s` causes . to match newlines
		}
	}

	t.description = t.Render(t.Description, values)

	return nil
}

func (t *UpdateTask) updateDst(ctx *TaskContext) error {
	if ctx.DryRun {
		return nil
	}

	// Update the content (using the pattern and replacement values)
	var err error
	if t.fileType.IsStructured() {
		t.dstBytes, err = t.replaceStructured(t.dstBytes, t.match.Pattern, t.src)
	} else {
		t.dstBytes, err = t.replaceText(t.dstBytes, t.match.Pattern, t.src)
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
	if t.action.Type == modify.ActionDelete { //nolint:nestif
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
		modifierOpt := modify.WithMergeType(t.action.MergeType)
		modifier := modify.Modifier(t.action.Type, repl, modifierOpt)
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

	replStr, err := cast.ToStringE(repl)
	if err != nil {
		return nil, fmt.Errorf("replacement: %w", err)
	}

	srcBytes := []byte(replStr)
	replacerFunc := func(dst []byte) []byte {
		// The src content may contain capture group placeholders ("${1}", etc).
		// We need to expand those manually.
		matches := re.FindSubmatchIndex(dst)
		srcExpanded := re.Expand([]byte{}, srcBytes, dst, matches)

		// The modifier doesn't yet know how to handle byte arrays...
		srcStr := string(srcExpanded)
		dstStr := string(dst)

		// Use the expanded src content to create a modifier func,
		// then use it to modify the dst content.
		// This allows us to use the same modification logic
		// (and merge behavior) as when working with structured data.
		modifierOpt := modify.WithMergeType(t.action.MergeType)
		modifier := modify.Modifier(t.action.Type, srcStr, modifierOpt)
		modified, _ := modifier(dstStr)

		// Finally return the modified dst content back to the regexp
		// so that it can be used to replace the current match.
		dstStr = modified.(string)
		return []byte(dstStr)
	}

	if t.match.Source == MatchSourceFile {
		return re.ReplaceAllFunc(content, replacerFunc), nil
	}
	newline := []byte("\n")
	updatedLines := [][]byte{}
	for _, line := range bytes.Split(content, newline) {
		updatedLine := re.ReplaceAllFunc(line, replacerFunc)
		updatedLines = append(updatedLines, updatedLine)
	}
	return bytes.Join(updatedLines, newline), nil
}
