package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/puppetlabs/relay-sdk-go/pkg/decorators"
	"github.com/spf13/cobra"
)

func NewDecoratorCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "decorator",
		Short:                 "Manage step decorators",
		DisableFlagsInUseLine: true,
	}

	cmd.AddCommand(NewSetDecoratorCommand())

	return cmd
}

func NewSetDecoratorCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "set a step decorator",
		Args:  cobra.ExactArgs(1),
		RunE:  doSetDecorator,
	}

	cmd.Flags().StringP("name", "n", "", "the decorator name")

	cmd.Flags().StringSliceP("value", "v", []string{}, "one or more decorator values")

	// lint: this is instructing the linter to ignore the error check here
	// since it's not strictly necessary (it's a potential flag does not exist
	// error from pflag's SetAnnotation).
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("value")

	return cmd
}

func doSetDecorator(cmd *cobra.Command, args []string) error {
	client, err := decorators.NewDefaultClientFromEnv()
	if err != nil {
		return err
	}

	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return err
	}

	rawValues, err := cmd.Flags().GetStringSlice("value")
	if err != nil {
		return err
	}

	values := map[string]string{}
	for _, v := range rawValues {
		parts := strings.SplitN(v, "=", 2)

		if len(parts) == 2 {
			key, value := parts[0], parts[1]
			if key != "" {
				values[key] = value

				continue
			}
		}

		return fmt.Errorf("invalid value: %s", v)
	}

	switch args[0] {
	case "link":
		if _, ok := values["description"]; !ok {
			return fmt.Errorf("link decorator: description field is required")
		}

		if _, ok := values["uri"]; !ok {
			return fmt.Errorf("link decorator: uri field is required")
		}

		values["type"] = "link"
	default:
		return fmt.Errorf("invalid decorator type: %s", args[0])
	}

	return client.Set(context.Background(), name, values)
}
