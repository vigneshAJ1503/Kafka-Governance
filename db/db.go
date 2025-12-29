package db

import (
	"context"
	"kafka-governance/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var topicCollection *mongo.Collection

func Connect(uri string) (*mongo.Client, *mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, nil, err
	}

	// Optional but good practice
	if err := client.Ping(ctx, nil); err != nil {
		return nil, nil, err
	}

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
	collection := PolicyCollection("kafka_governance")
	_, err := collection.InsertOne(ctx, policy)
	return err
}

func InitTopicRepo(db *mongo.Database) {
	topicCollection = db.Collection("topics")
}

func CreateTopic(ctx context.Context, topic *models.Topic) error {
	topic.CreatedAt = time.Now()
	_, err := topicCollection.InsertOne(ctx, topic)
	return err
}

func ListTopics(ctx context.Context) ([]models.Topic, error) {
	cursor, err := topicCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var topics []models.Topic
	err = cursor.All(ctx, &topics)
	return topics, err
}

func GetTopicByName(ctx context.Context, name string) (*models.Topic, error) {
	var topic models.Topic
	err := topicCollection.FindOne(ctx, bson.M{"name": name}).Decode(&topic)
	if err != nil {
		return nil, err
	}
	return &topic, nil
}

func ApproveTopic(ctx context.Context, name, approvedBy string) error {
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
	return err
}
