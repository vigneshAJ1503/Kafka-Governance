package api

import (
	"encoding/json"
	"net/http"

	"kafka-governance/models"
	"kafka-governance/service"
	"kafka-governance/utils"
)

var policies []models.Policy

func CreatePolicy(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLogger()

	var p models.Policy

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		logger.Error("Failed to decode policy request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Debug("Policy request body decoded successfully")

	if err := service.CreatePolicy(r.Context(), p); err != nil {
		logger.Error("Failed to create policy")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.Info("Policy created successfully")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}
