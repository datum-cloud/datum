package cmd

import (
	"github.com/spf13/cobra"
	"go.datumapis.com/datum/internal/authorization-webhook/cmd"
)

var webhook = &cobra.Command{
	Use:   "datum",
	Short: "Datum Cloud",
}

func init() {
	webhook.AddCommand(cmd.Webhook())
}

func Execute() error {
	return webhook.Execute()
}
