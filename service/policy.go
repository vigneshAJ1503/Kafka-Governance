package service

import (
	"context"
	"time"

	"kafka-governance/db"
	"kafka-governance/models"
)

func CreatePolicy(
	ctx context.Context,
	policy models.Policy,
) error {

	policy.CreatedAt = time.Now()
	return db.InsertPolicy(ctx, policy)
}
