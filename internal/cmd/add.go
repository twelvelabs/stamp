package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/twelvelabs/stamp/internal/core"
	"github.com/twelvelabs/stamp/internal/gen"
	"github.com/twelvelabs/stamp/internal/iostreams"
	"github.com/twelvelabs/stamp/internal/pkg"
)

func NewAddCmd(app *core.App) *cobra.Command {
	action := &AddAction{
		IO:    app.IO,
		Store: app.Store,
	}

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
			if err := action.Run(); err != nil {
				return err
			}
			return nil
		},
	}

	return cmd
}

type AddAction struct {
	IO    *iostreams.IOStreams
	Store *gen.Store

	Origin string
}

func (a *AddAction) Setup(cmd *cobra.Command, args []string) error {
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
	fmt.Fprintf(a.IO.Err, "Adding package from: %s\n", a.Origin)

	installed, err := a.Store.Install(a.Origin)
	if err != nil {
		return err
	}

	results, err := installed.Children()
	if err != nil {
		return err
	}
	results = append([]*pkg.Package{installed}, results...)

	for _, p := range results {
		fmt.Fprintf(a.IO.Err, " - %s\n", p.Name())
	}
	return nil
}
