package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/twelvelabs/stamp/internal/stamp"
)

func NewSchemaCmd(app *stamp.App) *cobra.Command {
	action := NewSchemaAction(app)

	cmd := &cobra.Command{
		Use:    "schema",
		Short:  "Create the generator.yaml JSON schema file",
		Hidden: true,
		Args:   cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := action.Validate(args); err != nil {
				return err
			}
			return action.Run(cmd.Context())
		},
	}

	return cmd
}

func NewSchemaAction(app *stamp.App) *SchemaAction {
	return &SchemaAction{
		App: app,
	}
}

type SchemaAction struct {
	*stamp.App
}

func (a *SchemaAction) Validate(_ []string) error {
	return nil
}

func (a *SchemaAction) Run(_ context.Context) error {
	metadata := &stamp.GeneratorMetadata{}
	schema, err := metadata.ReflectSchema()
	if err != nil {
		return fmt.Errorf("schema reflect: %w", err)
	}

	// Convert the schema to JSON and output.
	buf, err := json.MarshalIndent(schema, "", "    ")
	if err != nil {
		return fmt.Errorf("schema marshal: %w", err)
	}
	a.UI.Out(string(buf) + "\n")

	return nil
}
