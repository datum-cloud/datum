package app

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"go.datumapis.com/datum/pkg/cmd"
)

func NewWebhook() *cobra.Command {
	webhook := &cobra.Command{
		Use:   "datum-authorization-webhook",
		Short: "An authorization webhook backed by the Datum IAM service",
		PersistentPreRunE: func(webhook *cobra.Command, args []string) error {
			logger, err := cmd.SetupLogging(webhook)
			if err != nil {
				return fmt.Errorf("failed to configure logging: %w", err)
			}

			slog.SetDefault(logger)
			return nil
		},
	}

	cmd.AddLoggingFlags(webhook.PersistentFlags())

	webhook.AddCommand(
		serveCommand(),
	)

	return webhook
}
