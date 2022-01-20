package config

import (
	"time"

	"github.com/r523/dafdardaar/internal/board"
)

// Default return default configuration.
func Default() Config {
	return Config{
		BoardConnection: board.Config{
			ServerURL:         "tcp://127.0.0.1:1883",
			KeepAlive:         0,
			ConnectRetryDelay: time.Millisecond * 100,
		},
	}
}
