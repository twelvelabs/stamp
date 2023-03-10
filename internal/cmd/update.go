package cmd

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"

	"github.com/twelvelabs/stamp/internal/stamp"
)

func NewUpdateCmd(app *stamp.App) *cobra.Command {
	action := NewUpdateAction(app)

	cmd := &cobra.Command{
		Use:   "update [name]",
		Short: "Update generator to the latest version",
		Long:  "TODO",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := action.Setup(args)
			cobra.CheckErr(err)

			err = action.Validate()
			cobra.CheckErr(err)

			err = action.Run()
			cobra.CheckErr(err)
		},
	}
	return cmd
}

func NewUpdateAction(app *stamp.App) *UpdateAction {
	return &UpdateAction{
		App: app,
	}
}

type UpdateAction struct {
	*stamp.App

	Name string
}

func (a *UpdateAction) Setup(args []string) error {
	if len(args) >= 1 {
		a.Name = args[0]
	}
	return nil
}

func (a *UpdateAction) Validate() error {
	a.Name = strings.Trim(a.Name, " ")
	if a.Name == "" {
		return errors.New("name must not be blank")
	}
	return nil
}

func (a *UpdateAction) Run() error {
	a.UI.Out("Updating package: %s\n", a.Name)

	updated, err := a.Store.Update(a.Name)
	if err != nil {
		return err
	}

	children, err := updated.Children()
	if err != nil {
		return err
	}

	a.UI.Out(" - %s\n", updated.Name())
	for _, child := range children {
		a.UI.Out(" - %s\n", child.Name())
	}

	return nil
}
