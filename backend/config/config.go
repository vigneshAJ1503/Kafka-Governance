package config

import (
	"log"
	"os"
)

type Config struct {
	AppPort          string
	MongoURI         string
	DBName           string
	CedarURL         string
	JWTSecret        string
	UserCollection   string
	TopicCollection  string
	PolicyCollection string
}

func Load() *Config {
	cfg := &Config{
		AppPort:          getEnv("APP_PORT", "8080"),
		MongoURI:         getEnv("MONGO_URI", "mongodb://localhost:27017"),
		DBName:           getEnv("MONGO_DB", "kafkaGovernance"),
		UserCollection:   getEnv("USER_COLLECTION", "users"),
		TopicCollection:  getEnv("TOPIC_COLLECTION", "topics"),
		PolicyCollection: getEnv("POLICY_COLLECTION", "policies"),
		CedarURL:         getEnv("CEDAR_URL", "http://localhost:8180"),
		JWTSecret:        getEnv("JWT_SECRET", "dev-secret"),
	}

	log.Println("Config loaded")
	return cfg
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
