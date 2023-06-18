package stamp

import (
	"bytes"
	"fmt"
	"os"
	"regexp"

	"github.com/ohler55/ojg/jp"
	"github.com/twelvelabs/termite/render"

	"github.com/twelvelabs/stamp/internal/modify"
)

type UpdateTask struct {
	Common `mapstructure:",squash"`

	Action      UpdateAction    `mapstructure:"action"`
	Description render.Template `mapstructure:"description"`
	Dst         Destination     `mapstructure:"dst"`
	Match       UpdateMatch     `mapstructure:"match"`
	Src         Source          `mapstructure:"src"`
}

type UpdateAction struct {
	Type      modify.Action    `mapstructure:"type" default:"replace"`
	MergeType modify.MergeType `mapstructure:"merge" default:"concat"`
}

type UpdateMatch struct {
	PatternTpl render.Template `mapstructure:"pattern"`
	Default    any             `mapstructure:"default"`
	Source     MatchSource     `mapstructure:"source" default:"line"`

	pattern string
}

// Returns the rendered match pattern. SetPattern must be called first.
func (um *UpdateMatch) Pattern() string {
	return um.pattern
}

// SetPattern sets the given match pattern.
// Matches everything if the pattern is empty.
func (um *UpdateMatch) SetPattern(pat string, ct FileType) {
	if pat != "" {
		um.pattern = pat
	} else {
		if ct.IsStructured() {
			um.pattern = "$" // root node
		} else {
			um.Source = MatchSourceFile
			um.pattern = "(?s)^(.*)$" // `?s` causes . to match newlines
		}
	}
}

func (t *UpdateTask) Execute(ctx *TaskContext, values map[string]any) error {
	// Render the source and destination.
	if err := t.Dst.SetValues(values); err != nil {
		ctx.Logger.Failure("fail", t.Dst.Path())
		return err
	}
	if err := t.Src.SetValues(values); err != nil {
		ctx.Logger.Failure("fail", t.Dst.Path())
		return err
	}

	// Render description.
	desc, err := t.Description.Render(values)
	if err != nil {
		ctx.Logger.Failure("fail", t.Dst.Path())
		return err
	}

	// Render the match pattern.
	pattern, err := t.Match.PatternTpl.Render(values)
	if err != nil {
		ctx.Logger.Failure("fail", t.Dst.Path())
		return err
	}
	t.Match.SetPattern(pattern, t.Dst.ContentType())

	// Handle missing destination path.
	if !t.Dst.Exists() {
		switch t.Dst.Missing {
		case MissingConfigError:
			ctx.Logger.Failure("fail", t.Dst.Path())
			return ErrPathNotFound
		case MissingConfigTouch:
			if err := os.WriteFile(t.Dst.Path(), []byte{}, DstFileMode); err != nil {
				ctx.Logger.Failure("fail", t.Dst.Path())
				return err
			}
		default: // MissingConfigIgnore:
			return nil
		}
	}

	// Update the file.
	if err := t.updateDst(ctx); err != nil {
		ctx.Logger.Failure("fail", t.Dst.Path())
		return err
	}

	// Log success.
	updateMsg := t.Dst.Path()
	if desc != "" {
		// Include custom, generator supplied description.
		updateMsg = fmt.Sprintf("%s (%s)", t.Dst.Path(), desc)
	} else if t.Dst.ContentType().IsStructured() {
		// Or the JSON path expression.
		updateMsg = fmt.Sprintf("%s (%s)", t.Dst.Path(), t.Match.Pattern())
	}
	ctx.Logger.Success("update", updateMsg)

	return nil
}

func (t *UpdateTask) updateDst(ctx *TaskContext) error {
	if ctx.DryRun {
		return nil
	}

	var updated any
	var err error

	if t.Dst.ContentType().IsStructured() {
		updated, err = t.replaceStructured()
	} else {
		updated, err = t.replaceText()
	}
	if err != nil {
		return fmt.Errorf("update content: %w", err)
	}

	if err := t.Dst.Write(updated); err != nil {
		return fmt.Errorf("update content: %w", err)
	}

	return nil
}

func (t *UpdateTask) replaceStructured() (any, error) {
	data := t.Dst.Content()
	pattern := t.Match.Pattern()
	repl := t.Src.Content()

	// parse pattern as a JSON path expression
	exp, err := jp.ParseString(pattern)
	if err != nil {
		return nil, fmt.Errorf("json path parse: %w", err)
	}
	// use the expression to modify the data structure
	if t.Action.Type == modify.ActionDelete { //nolint:nestif
		data, err = exp.Remove(data)
		if err != nil {
			return nil, fmt.Errorf("json path remove: %w", err)
		}
	} else {
		if !exp.Has(data) && t.Match.Default != nil {
			err = exp.Set(data, t.Match.Default)
			if err != nil {
				return nil, fmt.Errorf("json path set default: %w", err)
			}
		}
		modifierOpt := modify.WithMergeType(t.Action.MergeType)
		modifier := modify.Modifier(t.Action.Type, repl, modifierOpt)
		data, err = exp.Modify(data, modifier)
		if err != nil {
			return nil, fmt.Errorf("json path modify: %w", err)
		}
	}

	return data, nil
}

func (t *UpdateTask) replaceText() (any, error) {
	dstBytes, err := t.Dst.ContentBytes()
	if err != nil {
		return nil, fmt.Errorf("dst bytes: %w", err)
	}

	re, err := regexp.Compile(t.Match.Pattern())
	if err != nil {
		return nil, fmt.Errorf("match pattern: %w", err)
	}

	srcBytes, err := t.Src.ContentBytes()
	if err != nil {
		return nil, fmt.Errorf("src bytes: %w", err)
	}

	// Using this replacement func (as opposed to a simple call to `re.ReplaceAll`)
	// because it allows us to use the `modify` package, and thus get the same
	// merge behavior as when working with structured data.
	replacerFunc := func(dst []byte) []byte {
		// The src content may contain capture group placeholders ("${1}", etc).
		// We need to expand those manually.
		matches := re.FindSubmatchIndex(dst)
		srcExpanded := re.Expand([]byte{}, srcBytes, dst, matches)

		// Modify the dst content with the src content.
		modifyFunc := modify.Modifier(
			t.Action.Type,
			srcExpanded,
			modify.WithMergeType(t.Action.MergeType),
		)
		modified, _ := modifyFunc(dst)

		// Finally return the modified dst content back to the regexp
		// so that it can be used to replace the current match.
		return modified.([]byte)
	}

	if t.Match.Source == MatchSourceFile {
		return re.ReplaceAllFunc(dstBytes, replacerFunc), nil
	}
	newline := []byte("\n")
	updatedLines := [][]byte{}
	for _, line := range bytes.Split(dstBytes, newline) {
		updatedLine := re.ReplaceAllFunc(line, replacerFunc)
		updatedLines = append(updatedLines, updatedLine)
	}
	return bytes.Join(updatedLines, newline), nil
}
