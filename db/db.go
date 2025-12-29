package db

import (
	"context"
	"errors"
	"kafka-governance/models"
	"kafka-governance/utils"
	"time"

	"github.com/google/uuid"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var topicCollection *mongo.Collection

func Connect(uri string) (*mongo.Client, *mongo.Database, error) {
	logger := utils.GetLogger()
	logger.Info("Attempting to connect to MongoDB")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		logger.Error("MongoDB connection failed")
		return nil, nil, err
	}
	logger.Debug("MongoDB client created")

	// Optional but good practice
	if err := client.Ping(ctx, nil); err != nil {
		logger.Error("MongoDB ping failed")
		return nil, nil, err
	}
	logger.Info("MongoDB connection successful")

	db := client.Database("kafka_governance")

	return client, db, nil
}

func PolicyCollection(dbName string) *mongo.Collection {
	return Client.Database(dbName).Collection("policies")
}

func InsertPolicy(
	ctx context.Context,
	policy models.Policy,
) error {
	logger := utils.GetLogger()
	logger.Debug("Inserting policy into database")

	collection := PolicyCollection("kafka_governance")
	_, err := collection.InsertOne(ctx, policy)
	if err != nil {
		logger.Error("Failed to insert policy into database")
		return err
	}
	logger.Info("Policy inserted successfully")
	return nil
}

func InitTopicRepo(db *mongo.Database) {
	logger := utils.GetLogger()
	logger.Debug("Initializing topic repository")
	topicCollection = db.Collection("topics")
	logger.Info("Topic repository initialized")
}

func CreateTopic(ctx context.Context, topic *models.Topic) (*models.Topic, error) {
	logger := utils.GetLogger()
	logger.Debug("Creating topic in database")

	var existingTopic models.Topic
	err := topicCollection.FindOne(ctx, bson.M{"name": topic.Name}).Decode(&existingTopic)
	if err == nil {
		logger.Error("Topic with same name already exists")
		return nil, errors.New("topic with same name already exists")
	}

	topic.ID = uuid.New().String()

	topic.CreatedAt = time.Now()
	_, err = topicCollection.InsertOne(ctx, topic)
	if err != nil {
		logger.Error("Failed to create topic in database")
		return nil, err
	}
	logger.Info("Topic created in database successfully")
	return topic, nil
}

func ListTopics(ctx context.Context) ([]models.Topic, error) {
	logger := utils.GetLogger()
	logger.Debug("Fetching topics from database")

	cursor, err := topicCollection.Find(ctx, bson.M{})
	if err != nil {
		logger.Error("Failed to query topics from database")
		return nil, err
	}
	defer cursor.Close(ctx)

	var topics []models.Topic
	err = cursor.All(ctx, &topics)
	if err != nil {
		logger.Error("Failed to decode topics from cursor")
		return nil, err
	}
	logger.Infof("Successfully fetched topics from database, count: %d", len(topics))
	return topics, nil
}

func GetTopicByName(ctx context.Context, name string) (*models.Topic, error) {
	logger := utils.GetLogger()
	logger.Debug("Fetching topic by name from database")

	var topic models.Topic
	err := topicCollection.FindOne(ctx, bson.M{"name": name}).Decode(&topic)
	if err != nil {
		logger.Error("Topic not found in database")
		return nil, err
	}
	logger.Info("Topic found in database")
	return &topic, nil
}

func ApproveTopic(ctx context.Context, name, approvedBy string) error {
	logger := utils.GetLogger()
	logger.Debug("Updating topic approval status in database")

	now := time.Now()
	_, err := topicCollection.UpdateOne(
		ctx,
		bson.M{"name": name},
		bson.M{
			"$set": bson.M{
				"status":     models.TopicApproved,
				"approvedBy": approvedBy,
				"approvedAt": now,
			},
		},
	)
	if err != nil {
		logger.Error("Failed to update topic approval status")
		return err
	}
	logger.Info("Topic approval status updated successfully")
	return nil
}
