package stamp

import (
	"bytes"
	"fmt"
	"os"
	"regexp"

	"github.com/ohler55/ojg/jp"
	"github.com/swaggest/jsonschema-go"
	"github.com/twelvelabs/termite/render"

	"github.com/twelvelabs/stamp/internal/mdutil"
	"github.com/twelvelabs/stamp/internal/modify"
)

type UpdateTask struct {
	Common `mapstructure:",squash"`

	Action         UpdateAction    `mapstructure:"action"      title:"UpdateAction"`
	DescriptionTpl render.Template `mapstructure:"description" title:"Description" description:"An optional description of what is being updated."` //nolint: lll
	Dst            Destination     `mapstructure:"dst"         title:"Destination" required:"true"`
	Match          UpdateMatch     `mapstructure:"match"       title:"UpdateMatch"`
	Src            Source          `mapstructure:"src"         title:"Source" required:"true"`
	Type           string          `mapstructure:"type"        title:"Type"   required:"true" description:"Updates a file in the destination directory." const:"update" default:"update"` //nolint: lll
}

// PrepareJSONSchema implements the jsonschema.Preparer interface.
func (t *UpdateTask) PrepareJSONSchema(schema *jsonschema.Schema) error {
	schema.WithTitle("UpdateTask")
	schema.WithDescription(mdutil.ToMarkdown(`
		Updates a file in the destination directory.

		The default behavior is to replace the entire file with the
		source content, but you can optionally specify alternate
		[actions](#action) (prepend, append, or delete) or [target](#match)
		a subsection of the destination file.
		If the destination file is structured (JSON, YAML), then you
		may target a JSON path pattern, otherwise it will be treated
		as plain text and you can target via regular expression.

		Examples:

		__CODE_BLOCK__yaml
		tasks:
			- type: update
				# Render <./_src/COPYRIGHT.tpl> and append it
				# to the end of the README.
				# If the README does not exist in the destination dir,
				# then do nothing.
				src:
					path: "COPYRIGHT.tpl"
				action:
					type: "append"
				dst:
					path: "README.md"
		__CODE_BLOCK__

		__CODE_BLOCK__yaml
		tasks:
			- type: update
				# Update <./package.json> in the destination dir.
				# If the file is missing, create it.
				dst:
					path: "package.json"
					missing: "touch"
				# Don't update the entire file - just the dependencies section.
				# If the dependencies section is missing, initialize it to an empty object.
				match:
					pattern: "$.dependencies"
					default: {}
				# Append (i.e. merge) the source content to the dependencies section.
				# The default behavior is to fully replace the matched pattern
				# with the source content.
				action:
					type: "append"
				# Use this inline object as the source content.
				# We could alternately reference a source file
				# containing a JSON object.
				src:
					content:
						lodash: "4.17.21"
		__CODE_BLOCK__
	`))

	return nil
}

type UpdateAction struct {
	Type      modify.Action    `mapstructure:"type"  title:"Type"  default:"replace"`
	MergeType modify.MergeType `mapstructure:"merge" title:"Merge" default:"concat"`
}

// PrepareJSONSchema implements the jsonschema.Preparer interface.
func (UpdateAction) PrepareJSONSchema(schema *jsonschema.Schema) error {
	schema.WithTitle("UpdateAction")
	schema.WithDescription("The action to perform on the destination.")
	return nil
}

type UpdateMatch struct {
	PatternTpl render.Template `mapstructure:"pattern" title:"Pattern" default:"" description:"A regexp (content type: text) or JSON path expression (content type: json, yaml). When empty, will match everything."` //nolint: lll
	Default    any             `mapstructure:"default" title:"Default" description:"A default value to use if the JSON path expression is not found."`                                                                //nolint: lll
	Source     MatchSource     `mapstructure:"source"  title:"Source"  default:"line"`

	pattern string
}

// PrepareJSONSchema implements the jsonschema.Preparer interface.
func (UpdateMatch) PrepareJSONSchema(schema *jsonschema.Schema) error {
	schema.WithTitle("UpdateMatch")
	schema.WithDescription("Target a subset of the destination to update.")
	return nil
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

func (t *UpdateTask) TypeKey() string {
	return t.Type
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
	desc, err := t.DescriptionTpl.Render(values)
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
			if !ctx.DryRun {
				if err := os.WriteFile(t.Dst.Path(), []byte{}, DstFileMode); err != nil {
					ctx.Logger.Failure("fail", t.Dst.Path())
					return err
				}
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
