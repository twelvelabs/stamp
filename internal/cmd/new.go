package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/twelvelabs/stamp/internal/stamp"
)

func NewNewCmd(app *stamp.App) *cobra.Command {
	action := NewNewAction(app)

	cmd := &cobra.Command{
		Use:   "new [name]",
		Short: "Run the named generator",
		Args:  cobra.ArbitraryArgs,
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

	cmd.Flags().BoolVar(&app.Config.DryRun, "dry-run", app.Config.DryRun, "Show generator tasks without taking action.")
	cmd.Flags().Lookup("dry-run").NoOptDefVal = "true"

	cmd.Flags().SortFlags = false
	cmd.DisableFlagParsing = true
	cmd.SilenceUsage = true

	return cmd
}

func NewNewAction(app *stamp.App) *NewAction {
	return &NewAction{
		App: app,
	}
}

type NewAction struct {
	*stamp.App

	Name string

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
	if a.Name == "" {
		if a.showHelp() {
			// User ran `stamp new --help`.
			return pflag.ErrHelp
		}
	}
	return nil
}

func (a *NewAction) Run() error {
	a.Name = strings.Trim(a.Name, " ")
	if a.Name == "" {
		all, err := a.Store.LoadAll()
		if err != nil {
			return err
		}

		if len(all) == 0 {
			return pflag.ErrHelp
		}

		names := []string{}
		for _, g := range all {
			names = append(names, g.Name())
		}

		a.Name, err = a.UI.Select("Generator", names, names[0])
		if err != nil {
			return err
		}
	}

	generator, err := a.Store.Load(a.Name)
	if err != nil {
		return err
	}

	// Update usage text w/ info from generator.
	a.setUsage(generator)
	// Set any user supplied default values from the config file.
	// Needs to be done prior to flag registration so that
	// the correct defaults are shown in usage.
	a.setDefaults(generator)

	// Add and parse the generator's flags.
	a.registerFlags(generator)
	if err := a.parseFlags(); err != nil {
		return err
	}

	// Show generator specific help (now that flags are parsed).
	if a.showHelp() {
		return pflag.ErrHelp
	}

	// Set the positional args
	if err := a.setArgs(generator); err != nil {
		return err
	}

	if err := generator.Values.Prompt(a.UI); err != nil {
		return err
	}
	if err := generator.Values.Validate(); err != nil {
		return err
	}

	a.UI.Out("\n")
	a.UI.Out("Running: %s\n", generator.Name())
	a.UI.Out("\n")

	// Prepare everything needed to execute.
	ctx := stamp.NewTaskContext(a.App)
	values := generator.Values.GetAll()

	// And finally... Release the hounds™
	if err := generator.Tasks.Execute(ctx, values); err != nil {
		return err
	}

	a.UI.Out("\n")

	return nil
}

func (a *NewAction) setUsage(generator *stamp.Generator) {
	a.cmd.Use = strings.ReplaceAll(a.cmd.Use, "[name]", generator.Name())
	for _, v := range generator.Values.Args() {
		a.cmd.Use += fmt.Sprintf(" [<%s>]", v.FlagName())
	}
	if desc := generator.Description(); desc != "" {
		a.cmd.Long = desc
	}
}

func (a *NewAction) setDefaults(generator *stamp.Generator) {
	for _, val := range generator.Values.All() {
		// viper forces all config keys to lowercase,
		// so users have to store defaults by flag name :shrug:
		// See: https://github.com/spf13/viper/issues/1014
		if def, ok := a.Config.Defaults[val.FlagName()]; ok {
			val.Default = def
		}
		_ = val.Get() // hack to force value set cache update
	}
}

func (a *NewAction) registerFlags(generator *stamp.Generator) {
	for _, val := range generator.Values.Flags() {
		a.cmd.Flags().Var(val, val.FlagName(), val.Help)
		if val.IsBoolFlag() {
			a.cmd.Flags().Lookup(val.FlagName()).NoOptDefVal = "true"
		}
	}
}

func (a *NewAction) parseFlags() error {
	a.cmd.DisableFlagParsing = false
	return a.cmd.ParseFlags(a.args)
}

func (a *NewAction) setArgs(generator *stamp.Generator) error {
	nonFlagArgs := a.cmd.Flags().Args()
	remaining, err := generator.Values.SetArgs(nonFlagArgs)
	if err != nil {
		return err
	}
	a.args = remaining
	// TODO: should we error or warn if there are extra pos args left over?
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
