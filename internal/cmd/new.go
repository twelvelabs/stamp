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
		Config:   app.Config,
		IO:       app.IO,
		Prompter: app.Prompter,
		Store:    app.Store,
	}

	cmd := &cobra.Command{
		Use:   "new <name>",
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
	Config   *core.Config
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
	// Load the generator
	generator, err := a.Store.Load(a.Name)
	if err != nil {
		return err
	}

	for _, val := range generator.Values.All() {
		// viper forces all config keys to lowercase,
		// so users have to store defaults by flag name :shrug:
		// See: https://github.com/spf13/viper/issues/1014
		if def, ok := a.Config.Defaults[val.FlagName()]; ok {
			val.Default = def
		}
	}

	// Re-configure the command (now that we have the generator)
	a.cmd.Use = strings.ReplaceAll(a.cmd.Use, "<name>", generator.Name())
	for _, v := range generator.Values.Args() {
		a.cmd.Use += fmt.Sprintf(" [<%s>]", v.FlagName())
	}
	a.cmd.DisableFlagParsing = false

	// Add and parse the generator's flags
	for _, val := range generator.Values.Flags() {
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
	remaining, err := generator.Values.SetArgs(nonFlagArgs)
	if err != nil {
		return err
	}
	a.args = remaining
	// TODO: should we error or warn if there are extra pos args left over?

	if err := generator.Values.Prompt(a.Prompter); err != nil {
		return err
	}
	if err := generator.Values.Validate(); err != nil {
		return err
	}

	fmt.Fprintln(a.IO.Err, "")
	fmt.Fprintln(a.IO.Err, "Running:", generator.Name())
	fmt.Fprintln(a.IO.Err, "")

	ctx := gen.NewTaskContext(a.IO, a.Prompter, a.Store)
	ctx.DryRun = dryRun
	values := generator.Values.GetAll()

	if err := generator.Tasks.Execute(ctx, values); err != nil {
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
