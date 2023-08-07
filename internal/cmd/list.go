package cmd

import (
	"github.com/rodaine/table"
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
			return action.Run()
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

func (a *ListAction) Setup(_ *cobra.Command, _ []string) error {
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

	tbl := table.New("Name", "Description", "Origin").WithWriter(a.IO.Out)
	for _, p := range results {
		tbl.AddRow(p.Name(), p.Description(), p.Origin())
	}
	tbl.Print()

	return nil
}
