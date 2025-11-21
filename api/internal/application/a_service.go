package application

import (
	"context"
	"kwaaka/api/internal/delivery/broker/rabbit"
	"kwaaka/api/internal/models"
	"kwaaka/api/internal/repository/mongo"
)

type Logger interface {
	Error(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Info(format string, v ...interface{})
	Debug(format string, v ...interface{})
}

type Parse interface {
	CreateTask(ctx context.Context, spreadsheetID, restaurantName string) (string, error)
	GetTask(ctx context.Context, taskID string) (*models.ParsingTask, error)
}

type Menu interface {
	GetMenu(ctx context.Context, menuID string) (*models.Menu, error)
}

type Product interface {
	ChangeStatus(ctx context.Context, productID, newStatus, reason, userID string) error
}

type Service struct {
	Parse   Parse
	Menu    Menu
	Product Product
	logger  Logger
}

func NewService(mongoRepo *mongo.MongoRepository, rabbitRepo *rabbit.RabbitRepository, logger Logger) *Service {
	parseSvc := NewParseService(mongoRepo, rabbitRepo, logger)
	menuSvc := NewMenuService(mongoRepo, logger)
	productSvc := NewProductService(rabbitRepo, logger)

	return &Service{
		Parse:   parseSvc,
		Menu:    menuSvc,
		Product: productSvc,
		logger:  logger,
	}
}

func (a *Service) Init() error {
	return nil
}

func (a *Service) Run(_ context.Context) {}

func (a *Service) Stop() {}
