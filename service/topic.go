package service

import (
	"context"
	"errors"

	"kafka-governance/db"
	"kafka-governance/models"
	"kafka-governance/utils"
)

func CreateTopic(ctx context.Context, topic *models.Topic) error {
	logger := utils.GetLogger()
	logger.Info("Processing topic creation request")

	if topic.Name == "" {
		logger.Error("Topic name validation failed")
		return errors.New("topic name required")
	}
	logger.Debug("Topic validation passed")

	topic.Status = models.TopicPending
	err := db.CreateTopic(ctx, topic)
	if err != nil {
		logger.Error("Topic creation failed at database layer")
		return err
	}
	logger.Info("Topic creation completed successfully")
	return nil
}

func ListTopics(ctx context.Context) ([]models.Topic, error) {
	logger := utils.GetLogger()
	logger.Info("Retrieving topics list")

	topics, err := db.ListTopics(ctx)
	if err != nil {
		logger.Error("Failed to retrieve topics list")
		return nil, err
	}
	logger.Infof("Topics list retrieved successfully, count: %d", len(topics))
	return topics, nil
}

func GetTopic(ctx context.Context, name string) (*models.Topic, error) {
	logger := utils.GetLogger()
	logger.Info("Retrieving topic by name")

	topic, err := db.GetTopicByName(ctx, name)
	if err != nil {
		logger.Error("Failed to retrieve topic")
		return nil, err
	}
	logger.Info("Topic retrieved successfully")
	return topic, nil
}

func ApproveTopic(ctx context.Context, name, admin string) error {
	logger := utils.GetLogger()
	logger.Info("Processing topic approval request")

	err := db.ApproveTopic(ctx, name, admin)
	if err != nil {
		logger.Error("Topic approval failed")
		return err
	}
	logger.Info("Topic approved successfully")
	return nil
}
