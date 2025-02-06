package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func AddLoggingFlags(cmd *pflag.FlagSet) {

	cmd.String("log-level", "INFO", "The level of logs that should be emitted from the service. Supports: 'ERROR', 'WARN', 'INFO', 'DEBUG'")
	cmd.String("log-format", "json", "The format of logs that should be emitted from the service. Supports: 'json' and 'text'")
}

func SetupLogging(cmd *cobra.Command) (*slog.Logger, error) {
	level, err := cmd.Flags().GetString("log-level")
	if err != nil {
		return nil, fmt.Errorf("failed to get `--log-level` flag: %w", err)
	}

	format, err := cmd.Flags().GetString("log-format")
	if err != nil {
		return nil, fmt.Errorf("failed to get `--log-format` flag: %w", err)
	}

	var handler slog.Handler

	handlerOpts := &slog.HandlerOptions{
		AddSource: true,
	}

	switch level {
	case "ERROR":
		handlerOpts.Level = slog.LevelError
	case "WARN":
		handlerOpts.Level = slog.LevelWarn
	case "INFO":
		handlerOpts.Level = slog.LevelInfo
	case "DEBUG":
		handlerOpts.Level = slog.LevelDebug
	default:
		return nil, fmt.Errorf("log level '%s' is invalid. Supported options are: 'ERROR', 'WARN', 'INFO', 'DEBUG'", level)
	}

	switch format {
	case "json":
		handler = slog.NewJSONHandler(os.Stderr, handlerOpts)
	case "text":
		handler = slog.NewTextHandler(os.Stderr, handlerOpts)
	default:
		return nil, fmt.Errorf("log format '%s' is invalid. Supported options are 'json', 'text'", format)
	}

	return slog.New(handler), nil
}
