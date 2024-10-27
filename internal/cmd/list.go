package cmd

import (
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
	"github.com/twelvelabs/termite/ui"

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

	cmd.Flags().BoolVarP(&action.ShowAll, "all", "a", action.ShowAll, "Show all columns.")

	return cmd
}

func NewListAction(app *stamp.App) *ListAction {
	return &ListAction{
		App: app,
	}
}

type ListAction struct {
	*stamp.App
	ShowAll bool
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

	renderGeneratorList(results, a.ShowAll, a.IO.Out)

	return nil
}

func renderGeneratorList(
	generators []*stamp.Generator,
	showAll bool,
	out ui.IOStream,
) {
	columns := []any{"Name", "Description"}
	if showAll {
		columns = append(columns, "Origin")
	}
	tbl := table.New(columns...).WithWriter(out)

	for _, g := range generators {
		if g.IsPrivate() {
			continue
		}
		if g.IsHidden() && !showAll {
			continue
		}
		row := []any{g.Name(), g.ShortDescription()}
		if showAll {
			row = append(row, g.Origin())
		}
		tbl.AddRow(row...)
	}

	tbl.Print()
}
