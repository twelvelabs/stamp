package cmd

import (
	"github.com/spf13/cobra"

	"github.com/twelvelabs/stamp/internal/stamp"
)

func NewRootCmd(app *stamp.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stamp",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
	examples and usage of using your application. For example:

	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) {
		// 	fmt.Fprintln(cmd.OutOrStdout(), "Hello World!")
		// },
	}

	// cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.stamp.yaml)")
	cmd.AddCommand(NewAddCmd(app))
	cmd.AddCommand(NewListCmd(app))
	cmd.AddCommand(NewNewCmd(app))
	cmd.AddCommand(NewRemoveCmd(app))
	cmd.AddCommand(NewUpdateCmd(app))

	return cmd
}
