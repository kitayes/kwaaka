package rabbit

import (
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	Host     string `env:"RABBIT_HOST,required"`
	Port     string `env:"RABBIT_PORT"`
	Username string `env:"RABBIT_USERNAME,required"`
	Password string `env:"RABBIT_PASSWORD,required"`
	VHost    string `env:"RABBIT_VHOST"`
}

const (
	MenuParsingQueue   = "menu-parsing"
	ProductStatusQueue = "product-status"
)

type Client struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewRabbit(cfg *Config) (*Client, error) {
	dsn := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.VHost,
	)

	conn, err := amqp.DialConfig(dsn, amqp.Config{
		Dial:      amqp.DefaultDial(5 * time.Second),
		Heartbeat: 10 * time.Second,
		Locale:    "en_US",
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to RabbitMQ: %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("rabbit channel open error: %w", err)
	}

	client := &Client{
		Conn:    conn,
		Channel: ch,
	}

	if err := declareQueues(client); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("rabbit declareQueues error: %w", err)
	}

	return client, nil
}

func declareQueues(c *Client) error {
	if _, err := c.Channel.QueueDeclare(
		MenuParsingQueue,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	); err != nil {
		return fmt.Errorf("declare queue %s: %w", MenuParsingQueue, err)
	}

	if _, err := c.Channel.QueueDeclare(
		ProductStatusQueue,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("declare queue %s: %w", ProductStatusQueue, err)
	}

	return nil
}
