package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/puppetlabs/leg/encoding/transfer"
	"github.com/puppetlabs/relay-core/pkg/model"
	outputsclient "github.com/puppetlabs/relay-sdk-go/pkg/outputs"
	"github.com/spf13/cobra"
	"github.com/puppetlabs/leg/timeutil/pkg/retry"
)

func NewOutputCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "output",
		Short:                 "Manage data that needs to be accessible to other tasks",
		DisableFlagsInUseLine: true,
	}

	cmd.AddCommand(NewSetOutputCommand())

	return cmd
}

func NewSetOutputCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "set",
		Short:                 "Set a value to a key that can be fetched by another task",
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := outputsclient.NewDefaultOutputsClientFromNebulaEnv()
			if err != nil {
				return err
			}

			key, err := cmd.Flags().GetString("key")
			if err != nil {
				return err
			}

			if isSensitive, err := cmd.Flags().GetBool("sensitive"); err != nil {
				return err
			} else if isSensitive {
				metadata := &model.StepOutputMetadata{
					Sensitive: true,
				}
				if err := client.SetOutputMetadata(context.Background(), key, metadata); err != nil {
					return err
				}
			}

			var value interface{}

			valueString, err := cmd.Flags().GetString("value")
			if err != nil {
				return err
			}

			if asJSON, err := cmd.Flags().GetBool("json"); err != nil {
				return err
			} else if asJSON {
				var encoded transfer.JSONInterface
				if err := json.Unmarshal([]byte(valueString), &encoded); err != nil {
					return fmt.Errorf("JSON decoding error: %+v", err)
				}

				value = encoded.Data
			} else {
				value = valueString
			}

			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, DefaultMetadataTimeout)
			defer cancel()
			waitOptions := []retry.WaitOption{}
			retryError := retry.Wait(ctx, func(ctx context.Context) (bool, error) {
				outputError := client.SetOutput(ctx, key, value)

				if outputError != nil {
					return retry.Repeat(outputError)
				}

				return retry.Done(nil)
			}, waitOptions...)

			if retryError != nil {
				return retryError
			}
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Successfully set output for %q.\n", key); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringP("key", "k", "", "the output key")
	cmd.Flags().StringP("value", "v", "", "the output value")
	cmd.Flags().Bool("sensitive", false, "flag the output value as sensitive")
	cmd.Flags().Bool("json", false, "whether the value should be interpreted as a JSON string")

	return cmd
}
