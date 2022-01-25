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
	Logger     *zap.Logger
}

// Config contains the configurations are required to
// create connection to the things broker.
type Config struct {
	ServerURL         string
	KeepAlive         uint16
	ConnectRetryDelay time.Duration
}

func New(ctx context.Context, cfg Config, logger *zap.Logger) (*Board, error) {
	b := new(Board)

	serverURL, err := url.Parse(cfg.ServerURL)
	if err != nil {
		return nil, fmt.Errorf("invalid server url %w", err)
	}

	router := paho.NewStandardRouter()

	router.RegisterHandler("/+/detect", b.onDetect)

	// nolint: exhaustivestruct
	cliCfg := autopaho.ClientConfig{
		BrokerUrls:        []*url.URL{serverURL},
		KeepAlive:         cfg.KeepAlive,
		ConnectRetryDelay: cfg.ConnectRetryDelay,
		OnConnectionUp: func(cm *autopaho.ConnectionManager, connAck *paho.Connack) {
			logger.Info("mqtt connection up")

			if _, err := cm.Subscribe(context.Background(), &paho.Subscribe{
				Subscriptions: map[string]paho.SubscribeOptions{
					"/+/detect": {
						QoS:               1,
						RetainHandling:    0,
						NoLocal:           false,
						RetainAsPublished: false,
					},
				},
				Properties: nil,
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
			OnClientError: func(err error) {
				logger.Info("server requested disconnect", zap.Error(err))
			},
			Router: router,
			OnServerDisconnect: func(d *paho.Disconnect) {
				if d.Properties != nil {
					logger.Info("server requested disconnect", zap.String("reason", d.Properties.ReasonString))
				} else {
					logger.Info("server requested disconnect", zap.Int("reason code", int(d.ReasonCode)))
				}
			},
		},
	}

	connectionManager, err := autopaho.NewConnection(ctx, cliCfg)
	if err != nil {
		return nil, fmt.Errorf("broker connection failed %w", err)
	}

	b.Connection = connectionManager
	b.Logger = logger

	return b, nil
}

// onDetect subscribes on detection event from ultrasonice sensor.
// it checks on/off period from database for turning on the door lamp.
func (b *Board) onDetect(m *paho.Publish) {
	b.Logger.Info("detect", zap.String("topic", m.Topic))

	var roomID int

	fmt.Sscanf(m.Topic, "/%d/detect", &roomID)

	b.Logger.Info("someone is behind the door", zap.Int("room", roomID))

	b.Connection.Publish(context.Background(), &paho.Publish{
		QoS:     1,
		Topic:   fmt.Sprintf("/%d/lamp", roomID),
		Payload: []byte("on"),
	})
}
