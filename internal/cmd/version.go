package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/twelvelabs/stamp/internal/stamp"
)

func NewVersionCmd(app *stamp.App) *cobra.Command {
	action := NewVersionAction(app)

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show full version info",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := action.Validate(args); err != nil {
				return err
			}
			return action.Run(cmd.Context())
		},
	}

	return cmd
}

func NewVersionAction(app *stamp.App) *VersionAction {
	return &VersionAction{
		App: app,
	}
}

type VersionAction struct {
	*stamp.App
}

func (a *VersionAction) Validate(_ []string) error {
	return nil
}

func (a *VersionAction) Run(_ context.Context) error {
	fmt.Fprintln(a.IO.Out, "Version:", a.Meta.Version)
	fmt.Fprintln(a.IO.Out, "GOOS:", a.Meta.GOOS)
	fmt.Fprintln(a.IO.Out, "GOARCH:", a.Meta.GOARCH)
	fmt.Fprintln(a.IO.Out, "")
	fmt.Fprintln(a.IO.Out, "Build Time:", a.Meta.BuildTime.Format(time.RFC3339))
	fmt.Fprintln(a.IO.Out, "Build Commit:", a.Meta.BuildCommit)
	fmt.Fprintln(a.IO.Out, "Build Version:", a.Meta.BuildVersion)
	fmt.Fprintln(a.IO.Out, "Build Checksum:", a.Meta.BuildChecksum)
	fmt.Fprintln(a.IO.Out, "Build Go Version:", a.Meta.BuildGoVersion)
	return nil
}
