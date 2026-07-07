package main

import (
	"context"
	"log"

	"kafka-governance/config"
	"kafka-governance/db"
	"kafka-governance/cache"
	"kafka-governance/routes"
	"kafka-governance/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize logger with config from environment
	utils.InitLoggerFromConfig()
	logger := utils.GetLogger()

	cfg := config.Load()
	logger.Info("Configuration loaded successfully")

	client, database, err := db.Connect(cfg.MongoURI)
	if err != nil {
		logger.Errorf("Failed to connect to MongoDB: %v", err)
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	if err := cache.InitRedis(); err != nil {
		logger.Errorf("Failed to connect to Redis: %v", err)
		log.Fatalf("Failed to connect to Redis: %v", err)
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
    db.InitUserRepo(database)
    logger.Info("User repository initialized")
    db.InitPolicyRepo(database)
    logger.Info("Policy repository initialized")

	r := gin.New()
	r.Use(gin.Recovery())

	
	r.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	routes.Register(r)
	logger.Info("Routes registered successfully")

	addr := ":" + cfg.AppPort
	logger.Infof("Server running on port: %s", cfg.AppPort)
	if err := r.Run(addr); err != nil {
		logger.Error("Server failed to start")
		log.Fatal(err)
	}
}
