package prompt

import (
	"fmt"
	"io"
	"strings"

	"github.com/AlecAivazis/survey/v2" //cspell: disable-line

	"github.com/twelvelabs/stamp/internal/value"
)

type FileReader interface {
	io.Reader
	Fd() uintptr
}

type FileWriter interface {
	io.Writer
	Fd() uintptr
}

func NewSurveyPrompter(in FileReader, out FileWriter, err FileWriter) value.Prompter {
	return &surveyPrompter{
		stdin:  in,
		stdout: out,
		stderr: err,
	}
}

// The default prompter that delegates to the `survey` package.
type surveyPrompter struct {
	stdin  FileReader
	stdout FileWriter
	stderr FileWriter
}

func (p *surveyPrompter) Confirm(prompt string, defaultValue bool, help string, validationRules string) (bool, error) {
	var result bool

	opts := []survey.AskOpt{
		survey.WithValidator(Validate(prompt, validationRules)),
	}

	err := p.ask(&survey.Confirm{
		Message: prompt,
		Help:    help,
		Default: defaultValue,
	}, &result, opts...)

	return result, err
}

func (p *surveyPrompter) Input(prompt string, defaultValue string, help string, validationRules string) (string, error) {
	var result string

	opts := []survey.AskOpt{
		survey.WithValidator(Validate(prompt, validationRules)),
	}

	err := p.ask(&survey.Input{
		Message: prompt,
		Help:    help,
		Default: defaultValue,
	}, &result, opts...)

	return result, err
}

func (p *surveyPrompter) ask(q survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
	opts = append(opts, survey.WithStdio(p.stdin, p.stdout, p.stderr))
	// survey.AskOne() doesn't allow passing in a transform func,
	// so we need to call survey.Ask().
	qs := []*survey.Question{
		{
			Prompt:    q,
			Transform: TrimSpace,
		},
	}
	err := survey.Ask(qs, response, opts...)
	if err == nil {
		return nil
	}
	return fmt.Errorf("could not prompt: %w", err)
}

var (
	_ survey.Transformer = TrimSpace
)

// Custom survey.Transformer that removes leading and trailing whitespace
// from string values. Non-string values are a no-op.
func TrimSpace(val any) any {
	if str, ok := val.(string); ok {
		return strings.TrimSpace(str)
	}
	return val
}

// Returns a survey.Validator that delegates to value.Validate.
func Validate(key string, rules string) survey.Validator {
	return func(val any) error {
		return value.ValidateKeyVal(key, val, rules)
	}
}
