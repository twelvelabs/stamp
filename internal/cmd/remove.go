package cmd

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"

	"github.com/twelvelabs/stamp/internal/core"
)

func NewRemoveCmd(app *core.App) *cobra.Command {
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
			if err := action.Run(); err != nil {
				return err
			}
			return nil
		},
	}

	return cmd
}

func NewRemoveAction(app *core.App) *RemoveAction {
	return &RemoveAction{
		App: app,
	}
}

type RemoveAction struct {
	*core.App

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

	a.UI.Out("Removing the following packages:\n")
	a.UI.Out(" - %s\n", generator.Name())
	for _, child := range children {
		a.UI.Out(" - %s\n", child.Name())
	}

	ok, err := a.UI.Confirm("Are you sure", false)
	if err != nil {
		return err
	}

	if ok {
		_, err = a.Store.Uninstall(a.Name)
		return err
	}
	return nil
}
