package main

import (
	"fmt"
	"os"

	"github.com/twelvelabs/stamp/internal/cmd"
	"github.com/twelvelabs/stamp/internal/core"
)

func main() {
	app, err := core.NewApp()
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
