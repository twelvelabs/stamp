package cmd

import (
	"fmt"

	"github.com/alexeyco/simpletable"
	"github.com/fatih/color"
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
	table := simpletable.New()

	formatHeader := color.New(color.FgYellow, color.Underline).SprintfFunc()
	formatPublic := color.New(color.FgCyan).SprintfFunc()
	formatHidden := color.New(color.FgCyan, color.Faint).SprintfFunc()

	// Setup the header.
	table.Header.Cells = []*simpletable.Cell{
		{Text: formatHeader("Name")},
		{Text: formatHeader("Description")},
	}
	if showAll {
		cell := &simpletable.Cell{
			Text: formatHeader("Origin"),
		}
		table.Header.Cells = append(table.Header.Cells, cell)
	}

	// Setup the body.
	for _, g := range generators {
		if g.IsPrivate() {
			continue
		}
		if g.IsHidden() && !showAll {
			continue
		}

		name := formatPublic(g.Name())
		if g.IsHidden() {
			name = formatHidden(g.Name())
		}
		row := []*simpletable.Cell{
			{Text: name},
			{Text: g.ShortDescription()},
		}
		if showAll {
			row = append(row, &simpletable.Cell{
				Text: g.Origin(),
			})
		}

		table.Body.Cells = append(table.Body.Cells, row)
	}

	table.SetStyle(simpletable.StyleCompactClassic)
	fmt.Fprintln(out, table.String())
}
