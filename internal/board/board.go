package board

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"go.uber.org/zap"
)

// Board control the boards which are available
// in the rooms of office.
type Board struct {
	Connection *autopaho.ConnectionManager
}

// Config contains the configurations are required to
// create connection to the things broker.
type Config struct {
	ServerURL         string
	KeepAlive         uint16
	ConnectRetryDelay time.Duration
}

func New(ctx context.Context, cfg Config, logger *zap.Logger) (Board, error) {
	serverURL, err := url.Parse(cfg.ServerURL)
	if err != nil {
		return Board{}, fmt.Errorf("invalid server url %w", err)
	}

	cliCfg := autopaho.ClientConfig{
		BrokerUrls:        []*url.URL{serverURL},
		KeepAlive:         cfg.KeepAlive,
		ConnectRetryDelay: cfg.ConnectRetryDelay,
		OnConnectionUp: func(cm *autopaho.ConnectionManager, connAck *paho.Connack) {
			logger.Info("mqtt connection up")

			if _, err := cm.Subscribe(context.Background(), &paho.Subscribe{
				Subscriptions: map[string]paho.SubscribeOptions{},
				Properties:    nil,
			}); err != nil {
				logger.Error("failed to subscribe. this is likely to mean no messages will be received.", zap.Error(err))

				return
			}

			logger.Info("mqtt subscription made")
		},
		OnConnectError: func(err error) {
			logger.Error("error whilst attempting connection", zap.Error(err))
		},
		ClientConfig: paho.ClientConfig{
			ClientID: "daftardaar",
			Router: paho.NewSingleHandlerRouter(func(m *paho.Publish) {
			}),
			OnClientError: func(err error) {
				fmt.Printf("server requested disconnect: %s\n", err)
			},
			OnServerDisconnect: func(d *paho.Disconnect) {
				if d.Properties != nil {
					fmt.Printf("server requested disconnect: %s\n", d.Properties.ReasonString)
				} else {
					fmt.Printf("server requested disconnect; reason code: %d\n", d.ReasonCode)
				}
			},
		},
	}

	connectionManager, err := autopaho.NewConnection(ctx, cliCfg)
	if err != nil {
		return Board{}, fmt.Errorf("broker connection failed %w", err)
	}

	return Board{
		Connection: connectionManager,
	}, nil
}
