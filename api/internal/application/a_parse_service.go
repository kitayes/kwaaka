package application

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"kwaaka/api/internal/delivery/broker/rabbit"
	"kwaaka/api/internal/models"
	"kwaaka/api/internal/repository/mongo"
)

type ParseService struct {
	mongo  *mongo.MongoRepository
	rabbit *rabbit.RabbitRepository
	logger Logger
}

func NewParseService(m *mongo.MongoRepository, r *rabbit.RabbitRepository, l Logger) *ParseService {
	return &ParseService{
		mongo:  m,
		rabbit: r,
		logger: l,
	}
}

func (s *ParseService) CreateTask(ctx context.Context, spreadsheetID, restaurantName string) (string, error) {
	taskID := uuid.NewString()
	now := time.Now().UTC()

	task := &models.ParsingTask{
		ID:             taskID,
		Status:         models.ParsingStatusQueued,
		SpreadsheetID:  spreadsheetID,
		RestaurantName: restaurantName,
		RetryCount:     0,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.mongo.CreateParsingTask(ctx, task); err != nil {
		return "", err
	}

	msg := struct {
		TaskID         string `json:"task_id"`
		SpreadsheetID  string `json:"spreadsheet_id"`
		RestaurantName string `json:"restaurant_name"`
	}{
		TaskID:         taskID,
		SpreadsheetID:  spreadsheetID,
		RestaurantName: restaurantName,
	}

	body, _ := json.Marshal(msg)

	if err := s.rabbit.PublishMenuParsing(ctx, body); err != nil {
		return "", err
	}

	return taskID, nil
}

func (s *ParseService) GetTask(ctx context.Context, taskID string) (*models.ParsingTask, error) {
	return s.mongo.GetParsingTask(ctx, taskID)
}
