package application

import (
	"context"

	"kwaaka/api/internal/models"
	"kwaaka/api/internal/repository/mongo"
)

type MenuService struct {
	mongo  *mongo.MongoRepository
	logger Logger
}

func NewMenuService(m *mongo.MongoRepository, l Logger) *MenuService {
	return &MenuService{
		mongo:  m,
		logger: l,
	}
}

func (s *MenuService) GetMenu(ctx context.Context, menuID string) (*models.Menu, error) {
	return s.mongo.GetMenuByID(ctx, menuID)
}
