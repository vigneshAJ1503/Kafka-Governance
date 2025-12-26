package service

import (
	"context"
	"time"

	"kafka-governance/db"
	"kafka-governance/models"
	"kafka-governance/utils"

	"go.mongodb.org/mongo-driver/bson"
)

func CreateTopic(ctx context.Context, dbName string, topic models.Topic) error {
	topic.Status = "PENDING"
	topic.CreatedAt = time.Now()

	utils.InfoLogger.Printf(
		"Service: creating topic name=%s owner=%s",
		topic.Name,
		topic.Owner,
	)

	collection := db.TopicCollection(dbName)

	_, err := collection.InsertOne(ctx, bson.M{
		"name":        topic.Name,
		"partitions":  topic.Partitions,
		"owner":       topic.Owner,
		"status":      topic.Status,
		"created_at":  topic.CreatedAt,
	})
	if err != nil {
		utils.ErrorLogger.Printf("Mongo insert failed: %v", err)
		return err
	}

	return nil
}
