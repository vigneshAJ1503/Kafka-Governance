package routes

import (
	"net/http"

	"kafka-governance/api"
)

func Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /topics", api.CreateTopic)
	mux.HandleFunc("GET /topics", api.ListTopics)
	mux.HandleFunc("GET /topics/{name}", api.GetTopic)
	mux.HandleFunc("POST /topics/{name}/approve", api.ApproveTopic)

	mux.HandleFunc("POST /policies", api.CreatePolicy)
}
