package main

import (
	"context"
	"log"
	"net/http"

	"kafka-governance/config"
	"kafka-governance/db"
	"kafka-governance/routes"
	"kafka-governance/utils"
)

func main() {
	// Initialize logger with config from environment
	utils.InitLoggerFromConfig()
	logger := utils.GetLogger()

	logger.Info("Starting Kafka Governance application")

	cfg := config.Load()
	logger.Info("Configuration loaded successfully")

	client, database, err := db.Connect(cfg.MongoURI)
	if err != nil {
		logger.Error("Failed to connect to database")
		log.Fatal(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			logger.Error("Failed to disconnect from database")
			log.Fatal(err)
		}
	}()

	logger.Info("Database connection established")

	db.InitTopicRepo(database)
	logger.Info("Topic repository initialized")

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	routes.Register(mux)
	logger.Info("Routes registered successfully")

	logger.Infof("Server running on port: %s", cfg.AppPort)
	if err := http.ListenAndServe(":"+cfg.AppPort, mux); err != nil {
		logger.Error("Server failed to start")
		log.Fatal(err)
	}
}

