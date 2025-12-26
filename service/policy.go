package service

import (
	"context"
	"time"

	"kafka-governance/db"
	"kafka-governance/models"
	"kafka-governance/utils"

	"go.mongodb.org/mongo-driver/bson"
)

func CreatePolicy(ctx context.Context, dbName string, policy models.Policy) error {
	utils.InfoLogger.Printf(
		"Service: creating policy principal=%s resource=%s",
		policy.Principal,
		policy.Resource,
	)

	allowed, err := utils.EvaluatePolicy("input.json")
	if err != nil {
		return err
	}
	if !allowed {
		return utils.ErrUnauthorized
	}

	policy.CreatedAt = time.Now()

	collection := db.PolicyCollection(dbName)

	_, err = collection.InsertOne(ctx, bson.M{
		"principal":  policy.Principal,
		"resource":   policy.Resource,
		"action":     policy.Action,
		"created_at": policy.CreatedAt,
	})
	if err != nil {
		utils.ErrorLogger.Printf("Mongo insert failed: %v", err)
		return err
	}

	return nil
}
