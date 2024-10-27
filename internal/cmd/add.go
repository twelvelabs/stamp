package cmd

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"

	"github.com/twelvelabs/stamp/internal/stamp"
)

func NewAddCmd(app *stamp.App) *cobra.Command {
	action := NewAddAction(app)

	cmd := &cobra.Command{
		Use:   "add [origin]",
		Short: "Add a new generator",
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

func NewAddAction(app *stamp.App) *AddAction {
	return &AddAction{
		App: app,
	}
}

type AddAction struct {
	*stamp.App

	Origin string
}

func (a *AddAction) Setup(_ *cobra.Command, args []string) error {
	if len(args) >= 1 {
		a.Origin = args[0]
	}
	return nil
}

func (a *AddAction) Validate() error {
	a.Origin = strings.Trim(a.Origin, " ")
	if a.Origin == "" {
		return errors.New("origin must not be blank")
	}
	return nil
}

func (a *AddAction) Run() error {
	a.UI.ProgressIndicator.StartWithLabel("Installing")

	installed, err := a.Store.Install(a.Origin)
	if err != nil {
		a.UI.ProgressIndicator.Stop()
		a.UI.Out(a.UI.FailureIcon() + " Install failed\n")
		return err
	}

	a.UI.ProgressIndicator.Stop()
	a.UI.Out(a.UI.SuccessIcon()+" Installed package: %s\n", installed.Name())
	a.UI.Out("\n")

	generators, err := a.Store.AsGenerators(installed.All())
	if err != nil {
		return err
	}
	renderGeneratorList(generators, false, a.IO.Out)

	return nil
}
