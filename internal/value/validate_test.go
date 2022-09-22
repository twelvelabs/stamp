package value

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	var tests = []struct {
		Rule  string
		Value any
		Err   string
	}{
		{
			Rule:  "",
			Value: "",
			Err:   "",
		},
		{
			Rule:  "unknown",
			Value: "",
			Err:   "undefined rule [my-val: unknown]",
		},
		{
			Rule:  "lowercase",
			Value: false,
			Err:   "invalid rule for bool [my-val: lowercase]",
		},
		{
			Rule:  "kebabcase",
			Value: "not_kebab_case",
			Err:   "my-val must be kebabcase",
		},
		{
			Rule:  "kebabcase",
			Value: "is-kebab-case",
			Err:   "",
		},
		{
			Rule:  "not-blank",
			Value: " ",
			Err:   "my-val must not be blank",
		},
		{
			Rule:  "not-blank",
			Value: "is-not-blank",
			Err:   "",
		},
		{
			Rule:  "posix-mode",
			Value: "12345",
			Err:   "my-val must be a valid posix file mode",
		},
		{
			Rule:  "posix-mode",
			Value: "0755",
			Err:   "",
		},
		{
			Rule:  "posix-mode",
			Value: "755",
			Err:   "",
		},
		{
			Rule:  "gt=0,lte=10",
			Value: 12,
			Err:   "my-val must be 10 or less",
		},
		{
			Rule:  "gt=0,lte=10",
			Value: 5,
			Err:   "",
		},
	}

	for _, test := range tests {
		t.Run(test.Rule, func(t *testing.T) {
			err := ValidateKeyVal("my-val", test.Value, test.Rule)

			if test.Err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, test.Err)
			}
		})
	}
}
