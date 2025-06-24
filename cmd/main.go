// SPDX-License-Identifier: AGPL-3.0-only
package main

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"go.datum.net/datum/cmd/controller"
)

var (
	cfgFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "datum",
	Short: "Datum control plane tooling for AI-native infrastructure orchestration",
	Long: `Datum is building the internet for AI - a neutral, programmable middle layer
where companies can programmatically connect without building the entire stack themselves.

This component provides tooling for extending the Milo control plane with Datum Cloud
specific functionality.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.ExecuteContext(context.Background())
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(controller.NewControllerManagerCommand())
}

func main() {
	Execute()
}
