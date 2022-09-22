package value

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	uni             *ut.UniversalTranslator
	translator      ut.Translator
	validate        *validator.Validate
	validUndefRegex = regexp.MustCompile(
		`^Undefined validation function '(?P<tag>[\w\.\-_]+)'`,
	)
	validBadTypeRegex = regexp.MustCompile(
		`^Bad field type (?P<field_type>.+)`,
	)
	validKebabRegex     = regexp.MustCompile(`^[a-z\-\d]+$`)
	validPosixModeRegex = regexp.MustCompile(`^[0-7]?[0-7]{3}$`) // https://rubular.com/r/WY16zVRPA90l2K
)

func init() {
	en := en.New()
	uni = ut.New(en, en)

	translator, _ = uni.GetTranslator("en")

	validate = validator.New()
	validate.RegisterValidation("kebabcase", IsKebabCase)
	validate.RegisterValidation("not-blank", validators.NotBlank)
	validate.RegisterValidation("posix-mode", IsPosixMode)

	en_translations.RegisterDefaultTranslations(validate, translator)

	validate.RegisterTranslation(
		"kebabcase",
		translator,
		func(ut ut.Translator) error {
			return ut.Add("kebabcase", "{0} must be kebabcase", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("kebabcase", fe.Field())
			return t
		},
	)
	validate.RegisterTranslation(
		"not-blank",
		translator,
		func(ut ut.Translator) error {
			return ut.Add("not-blank", "{0} must not be blank", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("not-blank", fe.Field())
			return t
		},
	)
	validate.RegisterTranslation(
		"posix-mode",
		translator,
		func(ut ut.Translator) error {
			return ut.Add("posix-mode", "{0} must be a valid posix file mode", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("posix-mode", fe.Field())
			return t
		},
	)
}

// Custom validator that ensures a POSIX permission number
// (either octal or int form).
func IsPosixMode(fl validator.FieldLevel) bool {
	field := fl.Field()
	return validPosixModeRegex.MatchString(strings.TrimSpace(field.String()))
}

// Custom validator that ensures kebab-case.
func IsKebabCase(fl validator.FieldLevel) bool {
	field := fl.Field()
	return validKebabRegex.MatchString(strings.TrimSpace(field.String()))
}

// ValidateKeyVal validates the key/value pair using rules.
func ValidateKeyVal(key string, value any, rules string) (err error) {
	defer func() {
		if panicVal := recover(); panicVal != nil {
			err = validatorPanicToErr(key, rules, panicVal)
		}
	}()

	err = validate.Var(value, rules)
	if err == nil {
		return
	}

	if errs, ok := err.(validator.ValidationErrors); ok {
		err = validatorErrsToErr(key, errs)
	}

	return
}

// ValidateStruct validates the struct using "validate" field tags.
func ValidateStruct(data any) error {
	err := validate.Struct(data)
	if err == nil {
		return nil
	}

	if errs, ok := err.(validator.ValidationErrors); ok {
		err = validatorErrsToErr("", errs)
	}
	return err
}

// Converts a [validator.ValidationErrors] into a standard error.
func validatorErrsToErr(key string, errs validator.ValidationErrors) error {
	translations := []string{}
	for _, fe := range errs {
		translation := fe.Translate(translator)
		if strings.HasPrefix(translation, " ") {
			// All the translations start w/ "{0}" (the field name),
			// but validate.Var() doesn't have access to the field name
			// and thus they end up starting w/ an empty space.
			// Add the field name back in.
			translation = fmt.Sprintf("%s %s", key, strings.TrimSpace(translation))
		} else {
			// If the translation is missing, `fe.Translate()` falls back
			// to the default error message - which has the same field name issue.
			translation = strings.ReplaceAll(translation, "''", fmt.Sprintf("'%s'", key))
		}
		translations = append(translations, translation)
	}
	s := strings.Join(translations, ", ")
	return errors.New(s)
}

// Converts a recovered panic value into an error.
func validatorPanicToErr(key string, rules string, panicVal any) error {
	str, ok := panicVal.(string)
	if !ok {
		return fmt.Errorf("%v", panicVal)
	}

	// reword to match the other validation errors
	if validUndefRegex.MatchString(str) {
		matches := validUndefRegex.FindStringSubmatch(str)
		tag := matches[validUndefRegex.SubexpIndex("tag")]
		str = fmt.Sprintf("undefined rule [%s: %s]", key, tag)
	}
	if validBadTypeRegex.MatchString(str) {
		matches := validBadTypeRegex.FindStringSubmatch(str)
		fieldType := matches[validBadTypeRegex.SubexpIndex("field_type")]
		str = fmt.Sprintf("invalid rule for %s [%s: %s]", fieldType, key, rules)
	}

	return errors.New(str)
}
