package mongo

import (
	"context"
	"kwaaka/api/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// можно устранить лишнюю зависимость но попозже

const (
	collectionMenus        = "menus"
	collectionParsingTasks = "parsing_tasks"
	collectionAudit        = "product_status_audit"
)

func (r *MongoRepository) EnsureIndexes(ctx context.Context) error {
	if r.db == nil {
		return mongo.ErrClientDisconnected
	}

	if err := r.ensureIndexes(ctx); err != nil {
		return err
	}
	r.logger.Info("ensure indexes")
	return nil
}

func (r *MongoRepository) ensureIndexes(ctx context.Context) error {
	menus := r.db.Collection(collectionMenus)
	if _, err := menus.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "restaurant_id", Value: 1}}},
		{Keys: bson.D{{Key: "products.ext_id", Value: 1}}},
	}); err != nil {
		return err
	}

	tasks := r.db.Collection(collectionParsingTasks)
	if _, err := tasks.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "_id", Value: 1}}},
		{Keys: bson.D{{Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "created_at", Value: 1}}},
	}); err != nil {
		return err
	}
	audit := r.db.Collection(collectionAudit)
	if _, err := audit.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "product_id", Value: 1}}},
		{Keys: bson.D{{Key: "timestamp", Value: 1}}},
	}); err != nil {
		return err
	}

	return nil
}

func (r *MongoRepository) CreateParsingTask(ctx context.Context, task *models.ParsingTask) error {
	coll := r.db.Collection(collectionParsingTasks)
	_, err := coll.InsertOne(ctx, task)
	return err
}

func (r *MongoRepository) GetParsingTask(ctx context.Context, id string) (*models.ParsingTask, error) {
	coll := r.db.Collection(collectionParsingTasks)

	var task models.ParsingTask
	if err := coll.FindOne(ctx, bson.M{"_id": id}).Decode(&task); err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *MongoRepository) GetMenuByID(ctx context.Context, id string) (*models.Menu, error) {
	coll := r.db.Collection(collectionMenus)

	var menu models.Menu
	if err := coll.FindOne(ctx, bson.M{"_id": id}).Decode(&menu); err != nil {
		return nil, err
	}
	return &menu, nil
}
