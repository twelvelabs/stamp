package main

import (
	"fmt"
	"os"

	"github.com/twelvelabs/stamp/internal/cmd"
	"github.com/twelvelabs/stamp/internal/stamp"
)

var (
	version = "dev"
	commit  = ""
	date    = ""
)

func main() {
	meta := stamp.NewAppMeta(version, commit, date)
	app, err := stamp.NewApp(meta)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	ctx := app.Context()

	command := cmd.NewRootCmd(app)
	if err := command.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
