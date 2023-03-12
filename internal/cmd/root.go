package cmd

import (
	"github.com/spf13/cobra"

	"github.com/twelvelabs/stamp/internal/stamp"
)

func NewRootCmd(app *stamp.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "stamp",
		Short:   "Stamp is project and file scaffolding tool.",
		Version: app.Meta.Version,
	}

	// cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.stamp.yaml)")
	cmd.AddCommand(NewAddCmd(app))
	cmd.AddCommand(NewListCmd(app))
	cmd.AddCommand(NewNewCmd(app))
	cmd.AddCommand(NewRemoveCmd(app))
	cmd.AddCommand(NewUpdateCmd(app))
	cmd.AddCommand(NewVersionCmd(app))

	return cmd
}
