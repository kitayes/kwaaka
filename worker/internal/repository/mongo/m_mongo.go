package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Host     string `env:"MONGO_HOST,required"`
	Port     string `env:"MONGO_PORT"`
	Username string `env:"MONGO_USERNAME,required"`
	Password string `env:"MONGO_PASSWORD,required"`
	DBName   string `env:"MONGO_DBNAME,required"`
}

type Mongo struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func NewMongoDB(cfg *Config) (*Mongo, error) {
	uri := fmt.Sprintf(
		"mongodb://%s:%s@%s:%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)

	clientOpts := options.Client().
		ApplyURI(uri).
		SetServerSelectionTimeout(5 * time.Second).
		SetConnectTimeout(5 * time.Second).
		SetMaxPoolSize(100).
		SetMinPoolSize(10)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return &Mongo{
		Client: client,
		DB:     client.Database(cfg.DBName),
	}, nil
}
