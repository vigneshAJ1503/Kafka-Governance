package routes

import (
	"net/http"

	"kafka-governance/api"
	"kafka-governance/utils"

	"github.com/go-chi/chi/v5"
)

func Register(r chi.Router) {
	// request logging middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			utils.InfoLogger.Printf(
				"%s %s from %s",
				req.Method,
				req.URL.Path,
				req.RemoteAddr,
			)
			next.ServeHTTP(w, req)
		})
	})


	r.Post("/topics", api.CreateTopic)
	r.Post("/policies", api.CreatePolicy)
}
