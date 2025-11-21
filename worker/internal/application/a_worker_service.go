package application

import (
	"context"
	"time"

	"kwaaka/worker/internal/models"
	"kwaaka/worker/internal/repository/mongo"
)

type Logger interface {
	Error(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Info(format string, v ...interface{})
	Debug(format string, v ...interface{})
}

type WorkerService struct {
	mongo  *mongo.MongoRepository
	logger Logger
}

func NewWorkerService(m *mongo.MongoRepository, l Logger) *WorkerService {
	return &WorkerService{
		mongo:  m,
		logger: l,
	}
}

func (w *WorkerService) HandleMenuParsing(ctx context.Context, msg *models.MenuParsingMessage) error {
	task, err := w.mongo.GetParsingTask(ctx, msg.TaskID)
	if err != nil {
		w.logger.Error("GetParsingTask err: %v", err)
		return err
	}

	task.Status = models.ParsingStatusProcessing
	task.UpdatedAt = time.Now().UTC()
	if err := w.mongo.UpdateParsingTask(ctx, task); err != nil {
		w.logger.Error("UpdateParsingTask processing err: %v", err)
		return err
	}

	menu := &models.Menu{
		ID:              task.ID,
		Name:            msg.RestaurantName,
		RestaurantID:    msg.RestaurantName,
		Products:        []interface{}{},
		Attributes:      []interface{}{},
		AttributesGroup: []interface{}{},
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}

	menuID, err := w.mongo.CreateMenu(ctx, menu)
	if err != nil {
		w.logger.Error("CreateMenu err: %v", err)

		task.Status = models.ParsingStatusFailed
		task.ErrorMessage = err.Error()
		task.RetryCount++
		task.UpdatedAt = time.Now().UTC()
		_ = w.mongo.UpdateParsingTask(ctx, task)

		return err
	}

	task.Status = models.ParsingStatusCompleted
	task.MenuID = menuID
	task.ErrorMessage = ""
	task.UpdatedAt = time.Now().UTC()
	if err := w.mongo.UpdateParsingTask(ctx, task); err != nil {
		w.logger.Error("UpdateParsingTask completed err: %v", err)
		return err
	}

	return nil
}

func (w *WorkerService) HandleProductStatus(ctx context.Context, evt *models.ProductStatusEvent) error {
	if err := w.mongo.InsertProductStatusAudit(ctx, evt); err != nil {
		w.logger.Error("InsertProductStatusAudit err: %v", err)
		return err
	}
	return nil
}
