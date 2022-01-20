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
		return Board{}, err
	}

	cliCfg := autopaho.ClientConfig{
		BrokerUrls:        []*url.URL{serverURL},
		KeepAlive:         cfg.KeepAlive,
		ConnectRetryDelay: cfg.ConnectRetryDelay,
		OnConnectionUp: func(cm *autopaho.ConnectionManager, connAck *paho.Connack) {
			fmt.Println("mqtt connection up")
			if _, err := cm.Subscribe(context.Background(), &paho.Subscribe{
				Subscriptions: map[string]paho.SubscribeOptions{},
			}); err != nil {
				fmt.Printf("failed to subscribe (%s). This is likely to mean no messages will be received.", err)
				return
			}
			fmt.Println("mqtt subscription made")
		},
		OnConnectError: func(err error) { fmt.Printf("error whilst attempting connection: %s\n", err) },
		ClientConfig: paho.ClientConfig{
			ClientID: "daftardaar",
			Router: paho.NewSingleHandlerRouter(func(m *paho.Publish) {
			}),
			OnClientError: func(err error) { fmt.Printf("server requested disconnect: %s\n", err) },
			OnServerDisconnect: func(d *paho.Disconnect) {
				if d.Properties != nil {
					fmt.Printf("server requested disconnect: %s\n", d.Properties.ReasonString)
				} else {
					fmt.Printf("server requested disconnect; reason code: %d\n", d.ReasonCode)
				}
			},
		},
	}

	//
	// Connect to the broker
	//
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cm, err := autopaho.NewConnection(ctx, cliCfg)
	if err != nil {
		return Board{}, err
	}

	return Board{
		Connection: cm,
	}, nil
}
