package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/twelvelabs/stamp/internal/core"
	"github.com/twelvelabs/stamp/internal/gen"
	"github.com/twelvelabs/stamp/internal/iostreams"
	"github.com/twelvelabs/stamp/internal/value"
)

func NewNewCmd(app *core.App) *cobra.Command {
	action := &NewAction{
		IO:       app.IO,
		Prompter: app.Prompter,
		Store:    app.Store,
	}

	cmd := &cobra.Command{
		Use:   "new [name]",
		Short: "Run the named generator",
		Long:  "TODO",
		Args:  cobra.ArbitraryArgs,
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

	cmd.Flags().BoolVar(&action.DryRun, "dry-run", true, "Show generator tasks without taking action.")
	cmd.Flags().Lookup("dry-run").NoOptDefVal = "true"

	cmd.Flags().SortFlags = false
	cmd.DisableFlagParsing = true
	cmd.SilenceUsage = true

	return cmd
}

type NewAction struct {
	IO       *iostreams.IOStreams
	Prompter value.Prompter
	Store    *gen.Store

	Name   string
	DryRun bool

	cmd  *cobra.Command
	args []string
}

func (a *NewAction) Setup(cmd *cobra.Command, args []string) error {
	a.cmd = cmd
	a.args = args

	// Since we're manually parsing flags they have yet to be removed from `args`.
	if len(args) >= 1 && !strings.HasPrefix(args[0], "-") {
		// strip name out of the args
		a.Name, a.args = args[0], args[1:]
	}

	return nil
}

func (a *NewAction) Validate() error {
	a.Name = strings.Trim(a.Name, " ")
	if a.Name == "" {
		if a.showHelp() {
			// User ran `stamp new --help`.
			return pflag.ErrHelp
		} else {
			return errors.New("name must not be blank")
		}
	}
	return nil
}

func (a *NewAction) Run() error {
	// Load the gen
	gen, err := a.Store.Load(a.Name)
	if err != nil {
		return err
	}

	// Re-configure the command (now that we have the generator)
	a.cmd.Use = strings.ReplaceAll(a.cmd.Use, "[name]", gen.Name())
	a.cmd.DisableFlagParsing = false

	// Add and parse the generator's flags
	for _, val := range gen.Values.Flags() {
		a.cmd.Flags().Var(val, val.FlagName(), val.Help)
		if val.IsBoolFlag() {
			a.cmd.Flags().Lookup(val.FlagName()).NoOptDefVal = "true"
		}
	}
	if err := a.cmd.ParseFlags(a.args); err != nil {
		return err
	}

	// Show generator specific help (now that flags are parsed)
	if a.showHelp() {
		return pflag.ErrHelp
	}

	dryRun, err := a.cmd.Flags().GetBool("dry-run")
	if err != nil {
		return err
	}

	// Set the positional args
	nonFlagArgs := a.cmd.Flags().Args()
	remaining, err := gen.Values.SetArgs(nonFlagArgs)
	if err != nil {
		return err
	}
	a.args = remaining
	// TODO: should we error or warn if there are extra pos args left over?

	if err := gen.Values.Prompt(a.Prompter); err != nil {
		return err
	}
	if err := gen.Values.Validate(); err != nil {
		return err
	}

	values := gen.Values.GetAll()
	// for k, v := range values {
	// 	fmt.Fprintln(a.IO.Out, " -", k, ":", v)
	// }

	fmt.Fprintln(a.IO.Err, "")
	fmt.Fprintln(a.IO.Err, "Running:", gen.Name())
	fmt.Fprintln(a.IO.Err, "")

	if err := gen.Tasks.Execute(values, a.IO, a.Prompter, dryRun); err != nil {
		return err
	}

	fmt.Fprintln(a.IO.Err, "")

	return nil
}

func (a *NewAction) showHelp() bool {
	// may not have parsed args yet, so manually check first
	for _, arg := range a.args {
		if arg == "-h" || arg == "--help" {
			return true
		}
	}
	// otherwise, query the command
	return a.cmd.Flags().Changed("help")
}
