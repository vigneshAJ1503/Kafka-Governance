package api

import (
	"encoding/json"
	"net/http"

	"kafka-governance/models"
	"kafka-governance/service"
	"kafka-governance/utils"
)

func CreateTopic(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLogger()

	var topic models.Topic
	if err := json.NewDecoder(r.Body).Decode(&topic); err != nil {
		logger.Error("Failed to decode topic request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Debug("Topic request body decoded successfully")

	topic.RequestedBy = r.Header.Get("X-User-Id")

	if err := service.CreateTopic(r.Context(), &topic); err != nil {
		logger.Error("Failed to create topic")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Info("Topic created successfully")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(topic)
}

func ListTopics(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLogger()

	topics, err := service.ListTopics(r.Context())
	if err != nil {
		logger.Error("Failed to list topics")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.Infof("Successfully retrieved topics list, count: %d", len(topics))
	json.NewEncoder(w).Encode(topics)
}

func GetTopic(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLogger()
	name := r.PathValue("name")

	topic, err := service.GetTopic(r.Context(), name)
	if err != nil {
		logger.Error("Topic not found")
		http.Error(w, "topic not found", http.StatusNotFound)
		return
	}
	logger.Info("Topic retrieved successfully")
	json.NewEncoder(w).Encode(topic)
}

func ApproveTopic(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLogger()
	name := r.PathValue("name")
	admin := r.Header.Get("X-User-Id")

	if err := service.ApproveTopic(r.Context(), name, admin); err != nil {
		logger.Error("Failed to approve topic")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.Info("Topic approved successfully")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"approved"}`))
}
