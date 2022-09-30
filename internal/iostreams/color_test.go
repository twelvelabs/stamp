package iostreams

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsColorEnabled(t *testing.T) {
	// reset env vars on teardown
	old_NO_COLOR := os.Getenv("NO_COLOR")
	old_CLICOLOR := os.Getenv("CLICOLOR")
	old_CLICOLOR_FORCE := os.Getenv("CLICOLOR_FORCE")
	defer func() {
		os.Setenv("NO_COLOR", old_NO_COLOR)
		os.Setenv("CLICOLOR", old_CLICOLOR)
		os.Setenv("CLICOLOR_FORCE", old_CLICOLOR_FORCE)
	}()

	tests := []struct {
		desc           string
		NO_COLOR       string
		CLICOLOR       string
		CLICOLOR_FORCE string
		disabled       bool
		forced         bool
		enabled        bool
	}{
		{
			desc:           "colors are enabled when all env vars are unset",
			NO_COLOR:       "",
			CLICOLOR:       "",
			CLICOLOR_FORCE: "",
			disabled:       false,
			forced:         false,
			enabled:        true,
		},

		{
			desc:           "colors are disabled whenever NO_COLOR has a value",
			NO_COLOR:       "something",
			CLICOLOR:       "",
			CLICOLOR_FORCE: "",
			disabled:       true,
			forced:         false,
			enabled:        false,
		},

		{
			desc:           "colors are disabled whenever CLICOLOR is set to 0",
			NO_COLOR:       "",
			CLICOLOR:       "0",
			CLICOLOR_FORCE: "",
			disabled:       true,
			forced:         false,
			enabled:        false,
		},

		{
			desc:           "colors are not disabled if CLICOLOR is set to anything other than 0",
			NO_COLOR:       "",
			CLICOLOR:       "something",
			CLICOLOR_FORCE: "",
			disabled:       false,
			forced:         false,
			enabled:        true,
		},

		{
			desc:           "colors are forced when CLICOLOR_FORCE is non-zero",
			NO_COLOR:       "",
			CLICOLOR:       "",
			CLICOLOR_FORCE: "1",
			disabled:       false,
			forced:         true,
			enabled:        true,
		},
		{
			desc:           "forcing colors takes priority over disabling",
			NO_COLOR:       "1",
			CLICOLOR:       "",
			CLICOLOR_FORCE: "1",
			disabled:       true,
			forced:         true,
			enabled:        true,
		},
		{
			desc:           "colors are not forced when CLICOLOR_FORCE is set to 0",
			NO_COLOR:       "",
			CLICOLOR:       "",
			CLICOLOR_FORCE: "0",
			disabled:       false,
			forced:         false,
			enabled:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			os.Setenv("NO_COLOR", tt.NO_COLOR)
			os.Setenv("CLICOLOR", tt.CLICOLOR)
			os.Setenv("CLICOLOR_FORCE", tt.CLICOLOR_FORCE)
			assert.Equal(t, tt.disabled, EnvColorDisabled())
			assert.Equal(t, tt.forced, EnvColorForced())
			assert.Equal(t, tt.enabled, IsColorEnabled())
		})
	}
}
