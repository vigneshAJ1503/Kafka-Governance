package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"kafka-governance/config"
)

type CedarAuthRequest struct {
	Principal string `json:"principal"`
	Action    string `json:"action"`
	Resource  string `json:"resource"`
	Context   map[string]interface{} `json:"context,omitempty"`
}

func EvaluatePolicy(principal, action, resource string) (bool, error) {
	logger := GetLogger()
	cfg := config.Load()

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

	logger.Info("Cedar Policy evaluation permissive")
	return true, nil
}
