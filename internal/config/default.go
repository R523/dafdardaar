package config

import (
	"time"

	"github.com/r523/dafdardaar/internal/board"
	"github.com/r523/dafdardaar/internal/logger"
)

// Default return default configuration.
func Default() Config {
	return Config{
		BoardConnection: board.Config{
			ServerURL:         "tcp://127.0.0.1:1883",
			KeepAlive:         1,
			ConnectRetryDelay: time.Millisecond * 100,
		},
		Logger: logger.Config{
			Level: "debug",
		},
	}
}
