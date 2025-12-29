package api

import (
	"encoding/json"
	"net/http"

	"kafka-governance/models"
	"kafka-governance/service"
)

var policies []models.Policy

func CreatePolicy(w http.ResponseWriter, r *http.Request) {
	var p models.Policy

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := service.CreatePolicy(r.Context(), p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}
