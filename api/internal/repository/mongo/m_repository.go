package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type Logger interface {
	Error(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Info(format string, v ...interface{})
	Debug(format string, v ...interface{})
}

type MongoRepository struct {
	cfg    *Config
	logger Logger

	mongo *Mongo
	db    *mongo.Database
}

func NewMongoRepository(cfg *Config, logger Logger) *MongoRepository {
	return &MongoRepository{
		cfg:    cfg,
		logger: logger,
	}
}

func (r *MongoRepository) Init() error {
	m, err := NewMongoDB(r.cfg)
	if err != nil {
		return err
	}

	r.mongo = m
	r.db = m.DB

	r.logger.Info("Mongo connected, db=%s", r.cfg.DBName)

	if err := r.EnsureIndexes(context.Background()); err != nil {
		r.logger.Error("EnsureIndexes err: ", err)
		return err
	}

	return nil
}

func (r *MongoRepository) Run(_ context.Context) {
	// заглушка для сервис-менеджера
}

func (r *MongoRepository) Stop() {
	if r.mongo == nil || r.mongo.Client == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := r.mongo.Client.Disconnect(ctx); err != nil {
		r.logger.Error("mongo disconnect fail: %v", err)
	} else {
		r.logger.Info("mongo disconnect success")
	}
}
