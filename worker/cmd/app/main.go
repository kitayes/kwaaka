package main

import (
	"context"
	"log/slog"

	"kwaaka/worker/internal/application"
	"kwaaka/worker/internal/delivery/broker/rabbit"
	workerhttp "kwaaka/worker/internal/delivery/http/v1"
	"kwaaka/worker/internal/repository/mongo"
	"kwaaka/worker/pkg/config"
	"kwaaka/worker/pkg/logger"
	"kwaaka/worker/pkg/service"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

type Config struct {
	Mongo  mongo.Config
	Rabbit rabbit.Config

	HTTP   workerhttp.Config `envPrefix:"WORKER_HTTP_"`
	Logger logger.Config     `envPrefix:"LOGGER_"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Error("error loading env variables: %s", err.Error())
	}

	cfg := Config{}
	if err := config.ReadEnvConfig(&cfg); err != nil {
		slog.Error("error initializing configs: %s", err.Error())
		return
	}

	log := logger.NewLogger(&cfg.Logger)

	mongoRepo := mongo.NewMongoRepository(&cfg.Mongo, log)
	rabbitRepo := rabbit.NewRabbitRepository(&cfg.Rabbit, log)

	workerSvc := application.NewWorkerService(mongoRepo, log)
	consumer := rabbit.NewWorkerConsumer(rabbitRepo, workerSvc, log)

	httpHandler := workerhttp.NewHandler(&cfg.HTTP, log)

	srv := service.NewManager(log)
	srv.AddService(
		mongoRepo,
		rabbitRepo,
		consumer,
		httpHandler,
	)

	ctx := context.Background()
	if err := srv.Run(ctx); err != nil {
		err = errors.Wrap(err, "worker srv.Run(...) err:")
		log.Error(err.Error())
		return
	}

	log.Info("Kwaaka WORKER started")
}
