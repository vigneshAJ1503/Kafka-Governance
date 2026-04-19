package service

import (
	"context"
	"time"

	"kafka-governance/db"
	"kafka-governance/models"
	"kafka-governance/utils"
)

func CreatePolicy(
	ctx context.Context,
	policy models.Policy,
) error {
	logger := utils.GetLogger()
	logger.Info("Creating new policy")

	policy.CreatedAt = time.Now()
	err := db.InsertPolicy(ctx, policy)
	if err != nil {
		logger.Error("Policy creation failed")
		return err
	}
	logger.Info("Policy created successfully")
	return nil
}
