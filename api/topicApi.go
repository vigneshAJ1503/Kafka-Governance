package api

import (
	"encoding/json"
	"net/http"

	"kafka-governance/models"
	"kafka-governance/service"
)

func CreateTopic(w http.ResponseWriter, r *http.Request) {
	var topic models.Topic
	if err := json.NewDecoder(r.Body).Decode(&topic); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	topic.RequestedBy = r.Header.Get("X-User-Id")

	if err := service.CreateTopic(r.Context(), &topic); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(topic)
}

func ListTopics(w http.ResponseWriter, r *http.Request) {
	topics, err := service.ListTopics(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(topics)
}

func GetTopic(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	topic, err := service.GetTopic(r.Context(), name)
	if err != nil {
		http.Error(w, "topic not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(topic)
}

func ApproveTopic(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	admin := r.Header.Get("X-User-Id")

	if err := service.ApproveTopic(r.Context(), name, admin); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"approved"}`))
}
