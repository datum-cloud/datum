package cmd

import "github.com/spf13/cobra"

func Webhook() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "authorization-webhook",
		Short: "An authorization webhook backed by the Datum IAM service",
	}

	cmd.AddCommand(
		serveCommand(),
	)

	return cmd
}
