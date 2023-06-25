package cmd

import (
	"github.com/spf13/cobra"

	"github.com/twelvelabs/stamp/internal/stamp"
)

func NewRootCmd(app *stamp.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "stamp",
		Short:   "A project and file scaffolding tool.",
		Version: app.Meta.Version,
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
	}

	// cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.stamp.yaml)")
	cmd.AddCommand(NewAddCmd(app))
	cmd.AddCommand(NewListCmd(app))
	cmd.AddCommand(NewManCmd(app))
	cmd.AddCommand(NewNewCmd(app))
	cmd.AddCommand(NewRemoveCmd(app))
	cmd.AddCommand(NewSchemaCmd(app))
	cmd.AddCommand(NewUpdateCmd(app))
	cmd.AddCommand(NewVersionCmd(app))

	return cmd
}
