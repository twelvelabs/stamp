package cmd

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"

	"github.com/twelvelabs/stamp/internal/core"
	"github.com/twelvelabs/stamp/internal/gen"
)

func NewUpdateCmd(app *core.App) *cobra.Command {
	action := &UpdateAction{
		Store: app.Store,
	}

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

			err = action.Run(cmd.InOrStdin(), cmd.OutOrStdout())
			cobra.CheckErr(err)
		},
	}
	return cmd
}

type UpdateAction struct {
	Store *gen.Store
	Name  string
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

func (a *UpdateAction) Run(in io.Reader, out io.Writer) error {
	fmt.Fprintln(out, "Updating package:", a.Name)

	updated, err := a.Store.Update(a.Name)
	if err != nil {
		return err
	}

	children, err := updated.Children()
	if err != nil {
		return err
	}

	fmt.Fprintln(out, " -", updated.Name())
	for _, child := range children {
		fmt.Fprintln(out, " -", child.Name())
	}

	return nil
}
