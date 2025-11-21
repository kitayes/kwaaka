package mongo

import (
	"context"
	"time"

	"kwaaka/worker/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	collectionMenus        = "menus"
	collectionParsingTasks = "parsing_tasks"
	collectionAudit        = "product_status_audit"
)

func (r *MongoRepository) GetParsingTask(ctx context.Context, id string) (*models.ParsingTask, error) {
	coll := r.db.Collection(collectionParsingTasks)

	var task models.ParsingTask
	if err := coll.FindOne(ctx, bson.M{"_id": id}).Decode(&task); err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *MongoRepository) UpdateParsingTask(ctx context.Context, task *models.ParsingTask) error {
	if r.db == nil {
		return mongo.ErrClientDisconnected
	}

	task.UpdatedAt = time.Now().UTC()

	coll := r.db.Collection(collectionParsingTasks)
	_, err := coll.UpdateByID(ctx, task.ID, bson.M{
		"$set": task,
	})
	return err
}

func (r *MongoRepository) CreateMenu(ctx context.Context, menu *models.Menu) (string, error) {
	if r.db == nil {
		return "", mongo.ErrClientDisconnected
	}

	now := time.Now().UTC()
	if menu.CreatedAt.IsZero() {
		menu.CreatedAt = now
	}
	menu.UpdatedAt = now

	coll := r.db.Collection(collectionMenus)
	_, err := coll.InsertOne(ctx, menu)
	if err != nil {
		return "", err
	}

	return menu.ID, nil
}

func (r *MongoRepository) InsertProductStatusAudit(ctx context.Context, evt *models.ProductStatusEvent) error {
	if r.db == nil {
		return mongo.ErrClientDisconnected
	}

	doc := bson.M{
		"product_id": evt.ProductID,
		"event_type": evt.EventType,
		"old_status": evt.OldStatus,
		"new_status": evt.NewStatus,
		"reason":     evt.Reason,
		"user_id":    evt.UserID,
		"timestamp":  evt.Timestamp,
		"created_at": time.Now().UTC(),
	}

	coll := r.db.Collection(collectionAudit)
	_, err := coll.InsertOne(ctx, doc)
	return err
}
