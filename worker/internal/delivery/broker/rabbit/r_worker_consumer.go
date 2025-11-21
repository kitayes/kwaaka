package rabbit

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"kwaaka/worker/internal/models"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	maxRetryCount = 3
	retryHeader   = "x-retry-count"

	menuParsingDLQ   = "menu-parsing-dlq"
	productStatusDLQ = "product-status-dlq"
)

type WorkerHandler interface {
	HandleMenuParsing(ctx context.Context, msg *models.MenuParsingMessage) error
	HandleProductStatus(ctx context.Context, evt *models.ProductStatusEvent) error
}

type WorkerConsumer struct {
	rabbit  *RabbitRepository
	handler WorkerHandler
	logger  Logger
}

func NewWorkerConsumer(r *RabbitRepository, h WorkerHandler, l Logger) *WorkerConsumer {
	return &WorkerConsumer{
		rabbit:  r,
		handler: h,
		logger:  l,
	}
}

func (w *WorkerConsumer) Init() error {
	if w.rabbit.client == nil || w.rabbit.client.Channel == nil {
		return fmt.Errorf("rabbit client/channel is nil")
	}

	ch := w.rabbit.client.Channel

	if _, err := ch.QueueDeclare(
		menuParsingDLQ,
		true, false, false, false, nil,
	); err != nil {
		return err
	}

	if _, err := ch.QueueDeclare(
		productStatusDLQ,
		true, false, false, false, nil,
	); err != nil {
		return err
	}

	w.logger.Info("worker consumer initialized")
	return nil
}

func (w *WorkerConsumer) Run(ctx context.Context) {
	ch := w.rabbit.client.Channel

	go w.consume(ctx, ch, MenuParsingQueue, menuParsingDLQ, w.handleMenuParsingDelivery)
	go w.consume(ctx, ch, ProductStatusQueue, productStatusDLQ, w.handleProductStatusDelivery)
}

func (w *WorkerConsumer) Stop() {}

type deliveryHandler func(ctx context.Context, ch *amqp.Channel, d amqp.Delivery) error

func (w *WorkerConsumer) consume(
	ctx context.Context,
	ch *amqp.Channel,
	queue string,
	dlq string,
	handler deliveryHandler,
) {
	msgs, err := ch.Consume(
		queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		w.logger.Error("consume %s err: %v", queue, err)
		return
	}

	w.logger.Info("consumer started for queue=%s", queue)

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("consumer for %s stopped by context", queue)
			return
		case d, ok := <-msgs:
			if !ok {
				w.logger.Info("consumer for %s: channel closed", queue)
				return
			}
			_ = w.handleWithRetry(ctx, ch, d, dlq, handler)
		}
	}
}

func (w *WorkerConsumer) handleWithRetry(
	ctx context.Context,
	ch *amqp.Channel,
	d amqp.Delivery,
	dlq string,
	handler deliveryHandler,
) error {
	retryCount := 0
	if v, ok := d.Headers[retryHeader]; ok {
		switch vv := v.(type) {
		case int32:
			retryCount = int(vv)
		case int64:
			retryCount = int(vv)
		case string:
			if n, err := strconv.Atoi(vv); err == nil {
				retryCount = n
			}
		}
	}

	if err := handler(ctx, ch, d); err != nil {
		w.logger.Error("handler err (retries=%d): %v", retryCount, err)

		if retryCount+1 >= maxRetryCount {
			if err := ch.PublishWithContext(
				ctx,
				"",
				dlq,
				false,
				false,
				amqp.Publishing{
					ContentType:   d.ContentType,
					Body:          d.Body,
					Headers:       d.Headers,
					DeliveryMode:  amqp.Persistent,
					CorrelationId: d.CorrelationId,
				},
			); err != nil {
				w.logger.Error("publish to DLQ %s err: %v", dlq, err)
			}
			_ = d.Ack(false)
			return err
		}

		retryCount++
		if d.Headers == nil {
			d.Headers = amqp.Table{}
		}
		d.Headers[retryHeader] = retryCount

		time.Sleep(time.Second * time.Duration(retryCount))

		if err := ch.PublishWithContext(
			ctx,
			"",
			d.RoutingKey,
			false,
			false,
			amqp.Publishing{
				ContentType:   d.ContentType,
				Body:          d.Body,
				Headers:       d.Headers,
				DeliveryMode:  amqp.Persistent,
				CorrelationId: d.CorrelationId,
			},
		); err != nil {
			w.logger.Error("requeue err: %v", err)
		}
		_ = d.Ack(false)
		return err
	}

	return d.Ack(false)
}

func (w *WorkerConsumer) handleMenuParsingDelivery(ctx context.Context, _ *amqp.Channel, d amqp.Delivery) error {
	var msg models.MenuParsingMessage
	if err := json.Unmarshal(d.Body, &msg); err != nil {
		w.logger.Error("menu-parsing unmarshal err: %v", err)
		return err
	}

	return w.handler.HandleMenuParsing(ctx, &msg)
}

func (w *WorkerConsumer) handleProductStatusDelivery(ctx context.Context, _ *amqp.Channel, d amqp.Delivery) error {
	var evt models.ProductStatusEvent
	if err := json.Unmarshal(d.Body, &evt); err != nil {
		w.logger.Error("product-status unmarshal err: %v", err)
		return err
	}

	return w.handler.HandleProductStatus(ctx, &evt)
}
