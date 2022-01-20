package local

import (
	"os"

	"github.com/r523/dafdardaar/internal/config"
	"github.com/spf13/cobra"
)

// ExitFailure status code.
const ExitFailure = 1

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	_ = config.New()

	// logger := logger.New(cfg.Logger)

	// nolint: exhaustivestruct
	root := &cobra.Command{
		Use:   "saf",
		Short: "Queue with NATS Jetstream to remove all the erlangs from cloud",
	}

	if err := root.Execute(); err != nil {
		// logger.Error("failed to execute root command", zap.Error(err))

		os.Exit(ExitFailure)
	}
}
