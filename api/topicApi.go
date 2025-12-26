package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"kafka-governance/config"
	"kafka-governance/models"
	"kafka-governance/service"
	"kafka-governance/utils"
)

func CreateTopic(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTopicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := service.CreateTopic(ctx, cfg.DBName, models.Topic{
		Name:       req.Name,
		Partitions: req.Partitions,
		Owner:      req.Owner,
	})
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.JSON(w, http.StatusCreated, map[string]string{
		"status": "topic created",
	})
}
