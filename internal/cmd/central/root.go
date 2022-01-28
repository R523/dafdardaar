package central

import (
	"context"
	"os"

	"github.com/r523/dafdardaar/internal/config"
	"github.com/r523/dafdardaar/internal/db"
	"github.com/r523/dafdardaar/internal/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// ExitFailure status code.
const ExitFailure = 1

func main(ctx context.Context, cfg config.Config, logger *zap.Logger) {
	db, err := db.New(cfg.DB)
	if err != nil {
		logger.Fatal("database connection failed", zap.Error(err))
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cfg := config.New()

	logger := logger.New(cfg.Logger)

	// nolint: exhaustivestruct
	root := &cobra.Command{
		Use:   "dafdardaar-central",
		Short: "Central Controller",
		Run: func(cmd *cobra.Command, args []string) {
			main(cmd.Context(), cfg, logger)
		},
	}

	if err := root.Execute(); err != nil {
		logger.Error("failed to execute root command", zap.Error(err))

		os.Exit(ExitFailure)
	}
}
