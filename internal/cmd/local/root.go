package local

import (
	"context"
	"os"

	"github.com/r523/dafdardaar/internal/board"
	"github.com/r523/dafdardaar/internal/config"
	"github.com/r523/dafdardaar/internal/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// ExitFailure status code.
const ExitFailure = 1

func main(ctx context.Context, cfg config.Config, logger *zap.Logger) {
	board, err := board.New(ctx, cfg.BoardConnection, logger)
	if err != nil {
		logger.Fatal("board initiation failed", zap.Error(err))
	}

	<-board.Connection.Done()
	logger.Info("up and running connection to mqtt broker")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cfg := config.New()

	logger := logger.New(cfg.Logger)

	// nolint: exhaustivestruct
	root := &cobra.Command{
		Use:   "dafdardaar-local",
		Short: "Office Local Controller",
		Run: func(cmd *cobra.Command, args []string) {
			main(cmd.Context(), cfg, logger)
		},
	}

	if err := root.Execute(); err != nil {
		logger.Error("failed to execute root command", zap.Error(err))

		os.Exit(ExitFailure)
	}
}
