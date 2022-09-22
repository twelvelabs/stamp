package value

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransform(t *testing.T) {
	var tests = []struct {
		Rule   string
		Input  any
		Output any
		Err    string
	}{
		{
			Rule:   "", // no rule should be a no-op
			Input:  "foo",
			Output: "foo",
			Err:    "",
		},
		{
			Rule:   " unknown ", // rule name should be trimmed in error msg
			Input:  "foo",
			Output: nil,
			Err:    "undefined transform [my-val: unknown]",
		},
		{
			Rule:   "trim",
			Input:  "  foo  ",
			Output: "foo",
			Err:    "",
		},
		{
			Rule:   "uppercase",
			Input:  "foo",
			Output: "FOO",
			Err:    "",
		},
		{
			Rule:   "lowercase",
			Input:  "FOO",
			Output: "foo",
			Err:    "",
		},
		{
			Rule:   "kebabcase",
			Input:  "FOO_BAR",
			Output: "foo-bar",
			Err:    "",
		},
		{
			Rule:   "camelcase",
			Input:  "foo bar",
			Output: "FooBar",
			Err:    "",
		},
		{
			Rule:   "expand-path",
			Input:  "~/../../../../../../home/${USER}",
			Output: os.ExpandEnv("/home/${USER}"),
			Err:    "",
		},
		{
			Rule:   "trim, kebabcase, uppercase", // should be able to combine rules
			Input:  "  foo bar  ",
			Output: "FOO-BAR",
			Err:    "",
		},
	}

	for _, test := range tests {
		t.Run(test.Rule, func(t *testing.T) {
			result, err := Transform("my-val", test.Input, test.Rule)

			assert.Equal(t, test.Output, result)

			if test.Err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, test.Err)
			}
		})
	}
}
