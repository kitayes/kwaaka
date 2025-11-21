package rabbit

import (
	"context"
	"encoding/json"
	"fmt"
	"kwaaka/api/internal/models"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (r *RabbitRepository) PublishMenuParsing(ctx context.Context, body []byte) error {
	if r.client == nil || r.client.Channel == nil {
		return fmt.Errorf("rabbit channel is nil")
	}

	return r.client.Channel.PublishWithContext(
		ctx,
		"",
		MenuParsingQueue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (r *RabbitRepository) PublishProductStatus(ctx context.Context, evt *models.ProductStatusEvent) error {
	if r.client == nil || r.client.Channel == nil {
		return fmt.Errorf("rabbit channel is nil")
	}

	body, err := json.Marshal(evt)
	if err != nil {
		return fmt.Errorf("marshal product status event: %w", err)
	}

	return r.client.Channel.PublishWithContext(
		ctx,
		"",
		ProductStatusQueue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
