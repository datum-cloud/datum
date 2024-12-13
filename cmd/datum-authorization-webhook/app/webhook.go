package app

import "github.com/spf13/cobra"

func NewWebhook() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "datum-authorization-webhook",
		Short: "An authorization webhook backed by the Datum IAM service",
	}

	cmd.AddCommand(
		serveCommand(),
	)

	return cmd
}
