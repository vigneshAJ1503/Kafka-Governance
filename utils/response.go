package utils

import (
	"encoding/json"
	"errors"
	"net/http"
)

var ErrUnauthorized = errors.New("unauthorized")

func JSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}
