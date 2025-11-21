package main

import (
	"context"
	"kwaaka/api/internal/application"
	"kwaaka/api/internal/delivery/broker/rabbit"
	v1 "kwaaka/api/internal/delivery/http/v1"
	"kwaaka/api/internal/repository/mongo"
	"kwaaka/api/pkg/config"
	"kwaaka/api/pkg/logger"
	"kwaaka/api/pkg/service"
	"log/slog"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

type Config struct {
	Mongo  mongo.Config
	Rabbit rabbit.Config

	HTTP   v1.Config     `envPrefix:"HTTP_"`
	Logger logger.Config `envPrefix:"LOGGER_"`
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

	servicesApp := application.NewService(mongoRepo, rabbitRepo, log)

	handlers := v1.NewHandler(servicesApp, &cfg.HTTP, log)

	srv := service.NewManager(log)
	srv.AddService(
		mongoRepo,
		rabbitRepo,
		handlers,
	)

	ctx := context.Background()
	if err := srv.Run(ctx); err != nil {
		err = errors.Wrap(err, "srv.Run(...) err:")
		log.Error(err.Error())
		return
	}

	log.Info("Kwaaka API started")
}
