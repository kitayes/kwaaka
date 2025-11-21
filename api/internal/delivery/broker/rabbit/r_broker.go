package rabbit

import (
	"context"
	"time"
)

type Logger interface {
	Error(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Info(format string, v ...interface{})
	Debug(format string, v ...interface{})
}

type RabbitRepository struct {
	cfg    *Config
	logger Logger

	client *Client
}

func NewRabbitRepository(cfg *Config, logger Logger) *RabbitRepository {
	return &RabbitRepository{
		cfg:    cfg,
		logger: logger,
	}
}

func (r *RabbitRepository) Init() error {
	client, err := NewRabbit(r.cfg)
	if err != nil {
		return err
	}

	r.client = client

	r.logger.Info("init rabbit repository")

	return nil
}

func (r *RabbitRepository) Run(_ context.Context) {
	// заглушка для сервис-менеджера
}

func (r *RabbitRepository) Stop() {
	if r.client == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if r.client.Channel != nil {
		if err := r.client.Channel.Close(); err != nil {
			r.logger.Error("failed to close channel: %v", err)
		} else {
			r.logger.Info("channel closed")
		}
	}

	if r.client.Conn != nil {
		if err := r.client.Conn.Close(); err != nil {
			r.logger.Error("failed to close connection: %v", err)
		} else {
			r.logger.Info("connection closed")
		}
	}

	_ = ctx // тоже заглушка
}
