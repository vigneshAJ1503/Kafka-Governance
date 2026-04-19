package db

import (
	"context"
	"errors"
	"time"

	"kafka-governance/models"
	"kafka-governance/utils"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection

func InitUserRepo(db *mongo.Database) {
	logger := utils.GetLogger()
	logger.Debug("Initializing user repository")
	userCollection = db.Collection("users")
}

func CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	logger := utils.GetLogger()
	logger.Debug("Creating user in database")

	var existingUser models.User
	err := userCollection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&existingUser)
	if err == nil {
		logger.Error("User with same username already exists")
		return nil, errors.New("username already exists")
	}

	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()

	_, err = userCollection.InsertOne(ctx, user)
	if err != nil {
		logger.Error("Failed to create user in database")
		return nil, err
	}
	logger.Info("User created in database successfully")
	return user, nil
}

func GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	logger := utils.GetLogger()
	logger.Debug("Fetching user by username")

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		logger.Error("User not found in database")
		return nil, err
	}
	return &user, nil
}
