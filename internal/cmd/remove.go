package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/twelvelabs/stamp/internal/core"
	"github.com/twelvelabs/stamp/internal/gen"
	"github.com/twelvelabs/stamp/internal/iostreams"
	"github.com/twelvelabs/stamp/internal/value"
)

func NewRemoveCmd(app *core.App) *cobra.Command {
	action := &RemoveAction{
		IO:       app.IO,
		Prompter: app.Prompter,
		Store:    app.Store,
	}

	cmd := &cobra.Command{
		Use:   "remove [name]",
		Short: "Remove generator",
		Long:  "TODO",
		Args:  cobra.ExactArgs(1),
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

type RemoveAction struct {
	Store    *gen.Store
	IO       *iostreams.IOStreams
	Prompter value.Prompter

	Name string
}

func (a *RemoveAction) Setup(cmd *cobra.Command, args []string) error {
	if len(args) >= 1 {
		a.Name = args[0]
	}
	return nil
}

func (a *RemoveAction) Validate() error {
	a.Name = strings.Trim(a.Name, " ")
	if a.Name == "" {
		return errors.New("name must not be blank")
	}
	return nil
}

func (a *RemoveAction) Run() error {
	generator, err := a.Store.Load(a.Name)
	if err != nil {
		return err
	}

	children, err := generator.Children()
	if err != nil {
		return err
	}

	fmt.Fprintln(a.IO.Out, "Removing the following packages:")
	fmt.Fprintln(a.IO.Out, " -", generator.Name())
	for _, child := range children {
		fmt.Fprintln(a.IO.Out, " -", child.Name())
	}

	confirmed, err := a.Prompter.Confirm("Are you sure?", false, "", "")
	if err != nil {
		return err
	}

	if confirmed {
		_, err = a.Store.Uninstall(a.Name)
		return err
	}
	return nil
}
