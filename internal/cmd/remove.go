package cmd

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"

	"github.com/twelvelabs/stamp/internal/stamp"
)

func NewRemoveCmd(app *stamp.App) *cobra.Command {
	action := NewRemoveAction(app)

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
			return action.Run()
		},
	}

	return cmd
}

func NewRemoveAction(app *stamp.App) *RemoveAction {
	return &RemoveAction{
		App: app,
	}
}

type RemoveAction struct {
	*stamp.App

	Name string
}

func (a *RemoveAction) Setup(_ *cobra.Command, args []string) error {
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

	generators, err := a.Store.AsGenerators(generator.All())
	if err != nil {
		return err
	}
	renderGeneratorList(generators, false, a.IO.Out)

	a.UI.Out("\n")
	ok, err := a.UI.Confirm("Remove these packages", false)
	if err != nil {
		return err
	}

	if ok {
		_, err = a.Store.Uninstall(a.Name)
		return err
	}
	return nil
}
