package service

import (
	"context"
	"errors"

	"kafka-governance/db"
	"kafka-governance/models"
)

func CreateTopic(ctx context.Context, topic *models.Topic) error {
	if topic.Name == "" {
		return errors.New("topic name required")
	}
	topic.Status = models.TopicPending
	return db.CreateTopic(ctx, topic)
}

func ListTopics(ctx context.Context) ([]models.Topic, error) {
	return db.ListTopics(ctx)
}

func GetTopic(ctx context.Context, name string) (*models.Topic, error) {
	return db.GetTopicByName(ctx, name)
}

func ApproveTopic(ctx context.Context, name, admin string) error {
	return db.ApproveTopic(ctx, name, admin)
}
