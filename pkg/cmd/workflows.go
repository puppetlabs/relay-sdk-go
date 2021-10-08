package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/puppetlabs/relay-sdk-go/pkg/workflows"
	"github.com/spf13/cobra"
)

func NewWorkflowsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "workflow",
		Short:                 "Manage workflows and their runs",
		DisableFlagsInUseLine: true,
	}

	cmd.AddCommand(NewRunWorkflowCommand())

	return cmd
}

func NewRunWorkflowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "run a workflow by its name",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := workflows.NewDefaultWorkflowsClientFromEnv()
			if err != nil {
				return err
			}

			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return err
			}

			pf, err := cmd.Flags().GetStringSlice("parameter")
			if err != nil {
				return err
			}

			var params map[string]string
			for _, p := range pf {
				parts := strings.Split(p, "=")

				if len(parts) < 2 || len(parts) > 2 {
					return fmt.Errorf("invalid parameter: %s", p)
				}

				params[parts[0]] = parts[1]
			}

			resp, err := client.Run(context.Background(), name, params)
			if err != nil {
				return err
			}

			of, err := cmd.Flags().GetString("output")
			if err != nil {
				return err
			}

			if of == "" {
				fmt.Fprintf(cmd.OutOrStdout(), "Successfully ran workflow %s (%d).\n", name, resp.RunNumber)

				return nil
			}

			if of != "json" {
				return fmt.Errorf("invalid output format: %s.\n", of)
			}

			b, err := json.Marshal(resp)
			if err != nil {
				return err
			}

			cmd.OutOrStdout().Write(b)

			return nil
		},
	}

	cmd.Flags().StringP("name", "n", "", "the workflow name")
	cmd.Flags().StringSliceP("parameter", "p", []string{}, "one or more workflow parameters")
	// Only json right now
	cmd.Flags().StringP("output", "o", "", "the output format to use")

	return cmd
}
