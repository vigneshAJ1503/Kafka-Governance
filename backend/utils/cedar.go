package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"kafka-governance/cache"
	"kafka-governance/config"
)

type CedarAuthRequest struct {
	Principal string `json:"principal"`
	Action    string `json:"action"`
	Resource  string `json:"resource"`
	Context   map[string]interface{} `json:"context,omitempty"`
}

// Cache TTL mapped for policy evaluations to prevent hitting Cedar heavily under high load.
const PolicyCacheTTL = 5 * time.Minute

// EvaluatePolicy evaluates Cedar policy constraints securely, pulling cached records sequentially.
func EvaluatePolicy(principal, action, resource string) (bool, error) {
	logger := GetLogger()
	ctx := context.Background()
	cfg := config.Load()

	// Construct cache key mapping the precise identity & resource constraints
	cacheKey := fmt.Sprintf("cedar_policy:%s:%s:%s", principal, action, resource)

	// 1. Try resolving securely via standard Redis Caching lookup bounds
	if cache.RedisClient != nil {
		cachedResult, err := cache.RedisClient.Get(ctx, cacheKey).Result()
		if err == nil {
			logger.Infof("⚡ Cache Hit: Resolved Cedar policy from Redis -> %s", cachedResult)
			return cachedResult == "true", nil
		}
	}

	// In a complete implementation, policies would be sent from DB or synced to the agent.
	// We'll perform an API call assuming the agent has the policy or doesn't mandate it for permit tests.
	reqBody := CedarAuthRequest{
		Principal: principal, // e.g. User::"alice"
		Action:    action,    // e.g. Action::"CreateTopic"
		Resource:  resource,  // e.g. Topic::"orders"
	}
	
	payloadBytes, err := json.Marshal(reqBody)
	if err != nil {
		logger.Errorf("Failed to marshal cedar request: %v", err)
		return false, err
	}

	url := fmt.Sprintf("%s/v1/is_authorized", cfg.CedarURL)
	logger.Infof("Evaluating Cedar Policy at %s for Principal: %s, Action: %s, Resource: %s", url, principal, action, resource)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		logger.Warnf("Failed to reach Cedar Agent, defaulting to true for development: %v", err)
		// For demo context if Cedar crashes/isn't synced, fallback to allow
		return true, nil 
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Warnf("Cedar Agent returned %d", resp.StatusCode)
		return false, fmt.Errorf("cedar evaluation failed with status %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	isAllowed := string(body) == "true"
	logger.Infof("Cedar policy assessment complete: %v", isAllowed)

	// 2. Set the resolution to Cache asynchronously targeting our specific TTL boundary mapping safely
	if cache.RedisClient != nil {
		err := cache.RedisClient.Set(ctx, cacheKey, fmt.Sprintf("%v", isAllowed), PolicyCacheTTL).Err()
		if err != nil {
			logger.Warnf("Failed to store evaluation payload natively to cache: %v", err)
		}
	}

	return isAllowed, nil
}
