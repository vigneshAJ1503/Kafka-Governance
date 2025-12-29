package routes

import (
	"net/http"

	"kafka-governance/api"
	"kafka-governance/utils"
)

// LoggingMiddleware wraps handlers to log entry and exit points

func Register(mux *http.ServeMux) {
	apiV1 := http.NewServeMux()

	apiV1.HandleFunc("POST /topics", utils.LoggingMiddleware("POST /api/v1/topics", api.CreateTopic))
	apiV1.HandleFunc("GET /topics", utils.LoggingMiddleware("GET /api/v1/topics", api.ListTopics))
	apiV1.HandleFunc("GET /topics/{name}", utils.LoggingMiddleware("GET /api/v1/topics/{name}", api.GetTopic))
	apiV1.HandleFunc("POST /topics/{name}/approve", utils.LoggingMiddleware("POST /api/v1/topics/{name}/approve", api.ApproveTopic))
	apiV1.HandleFunc("POST /policies", utils.LoggingMiddleware("POST /api/v1/policies", api.CreatePolicy))

	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", apiV1))
}
