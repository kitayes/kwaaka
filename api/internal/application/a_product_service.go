package application

import (
	"context"
	"time"

	"kwaaka/api/internal/delivery/broker/rabbit"

	"kwaaka/api/internal/models"
)

type ProductService struct {
	rabbit *rabbit.RabbitRepository
	logger Logger
}

func NewProductService(rabbit *rabbit.RabbitRepository, logger Logger) *ProductService {
	return &ProductService{
		rabbit: rabbit,
		logger: logger,
	}
}

func (s *ProductService) ChangeStatus(ctx context.Context, productID, newStatus, reason, userID string) error {
	evt := &models.ProductStatusEvent{
		EventType: "product.status_changed",
		ProductID: productID,
		OldStatus: "", // можно заполнять из БД, если нужно
		NewStatus: newStatus,
		Reason:    reason,
		Timestamp: time.Now().UTC(),
		UserID:    userID,
	}

	return s.rabbit.PublishProductStatus(ctx, evt)
}
