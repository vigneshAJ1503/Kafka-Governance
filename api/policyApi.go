package api

import (
	"net/http"

	"kafka-governance/utils"
)

func CreatePolicy(w http.ResponseWriter, r *http.Request) {
	utils.JSON(w, http.StatusNotImplemented, map[string]string{
		"message": "CreatePolicy not implemented yet",
	})
}
