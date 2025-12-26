package main

import (
	"log"
	"net/http"

	"kafka-governance/config"
	"kafka-governance/db"
	"kafka-governance/routes"
	"kafka-governance/utils"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	// load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on env vars")
	}

	// init logger
	utils.InitLogger()
	utils.InfoLogger.Println("Starting Kafka Governance Control Plane")

	// load config
	cfg := config.Load()

	// connect mongo
	if err := db.Connect(cfg.MongoURI); err != nil {
		utils.ErrorLogger.Fatalf("MongoDB connection failed: %v", err)
	}
	utils.InfoLogger.Println("Connected to MongoDB")

	// setup router
	r := chi.NewRouter()
	routes.Register(r)

	// start server
	utils.InfoLogger.Printf("Server running on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
