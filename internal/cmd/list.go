package cmd

import (
	"github.com/spf13/cobra"

	"github.com/twelvelabs/stamp/internal/stamp"
)

func NewListCmd(app *stamp.App) *cobra.Command {
	action := NewListAction(app)

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

func NewListAction(app *stamp.App) *ListAction {
	return &ListAction{
		App: app,
	}
}

type ListAction struct {
	*stamp.App
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

	a.UI.Out("Generators:\n")
	for _, p := range results {
		a.UI.Out(" - %s\n", p.Name())
	}
	return nil
}
