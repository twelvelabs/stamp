package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/twelvelabs/stamp/internal/core"
	"github.com/twelvelabs/stamp/internal/gen"
	"github.com/twelvelabs/stamp/internal/iostreams"
)

func NewListCmd(app *core.App) *cobra.Command {
	action := &ListAction{
		IO:    app.IO,
		Store: app.Store,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List installed generators",
		Long:  "TODO",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := action.Setup(cmd, args); err != nil {
				return err
			}
			if err := action.Validate(); err != nil {
				return err
			}
			if err := action.Run(); err != nil {
				return err
			}
			return nil
		},
	}

	return cmd
}

type ListAction struct {
	IO    *iostreams.IOStreams
	Store *gen.Store
}

func (a *ListAction) Setup(cmd *cobra.Command, args []string) error {
	return nil
}
func (a *ListAction) Validate() error {
	return nil
}
func (a *ListAction) Run() error {
	results, err := a.Store.LoadAll()
	if err != nil {
		return err
	}

	fmt.Fprintln(a.IO.Err, "Generators:")
	for _, p := range results {
		fmt.Fprintf(a.IO.Err, " - %s\n", p.Name())
	}
	return nil
}
