package api

import (
	"encoding/json"
	"net/http"

	"kafka-governance/models"
	"kafka-governance/service"
	"kafka-governance/utils"

	"github.com/gin-gonic/gin"
)

// sendErrorResponse sends a JSON error response with proper status code
func sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

// sendSuccessResponse sends a JSON success response
func sendSuccessResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func CreateTopic(c *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("Processing topic creation")

	// Validate required header
	requestedBy := c.GetHeader("X-User-Id")
	if requestedBy == "" {
		logger.Error("X-User-Id header missing")
		c.JSON(http.StatusForbidden, gin.H{"error": "X-User-Id header is required"})
		return
	}
	logger.Debug("X-User-Id header validated")

	var topic models.Topic
	if err := c.ShouldBindJSON(&topic); err != nil {
		logger.Error("Failed to decode request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if topic.Name == "" {
		logger.Error("Topic name validation failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Topic name is required"})
		return
	}

	if topic.Cluster == "" {
		logger.Error("Cluster name validation failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cluster name is required"})
		return
	}

	if topic.Partitions <= 0 {
		logger.Error("Partitions validation failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Partitions must be greater than 0"})
		return
	}

	if topic.Replicas <= 0 {
		logger.Error("Replicas validation failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Replicas must be greater than 0"})
		return
	}

	topic.RequestedBy = requestedBy
	createdTopic, err := service.CreateTopic(c.Request.Context(), &topic)
	if err != nil {
		logger.Error("Service layer returned error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create topic"})
		return

	}

	logger.Info("Topic created successfully")
	c.JSON(http.StatusCreated, gin.H{"message": "Topic created successfully", "topic": createdTopic})
}

func ListTopics(c *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("Received a request to list topics")

	topics, err := service.ListTopics(c.Request.Context())
	if err != nil {
		logger.Error("Failed to list topics")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve topics"})
		return
	}

	if topics == nil {
		topics = []models.Topic{}
	}

	logger.Infof("Successfully retrieved topics list, count: %d", len(topics))
	c.JSON(http.StatusOK, topics)
}

func GetTopic(c *gin.Context) {
	logger := utils.GetLogger()
	name := c.Param("name")
	logger.Info("Received a request to get topic")

	if name == "" {
		logger.Error("Topic name is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Topic name is required"})
		return
	}

	topic, err := service.GetTopic(c.Request.Context(), name)
	if err != nil {
		logger.Error("Topic not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Topic not found"})
		return
	}
	logger.Info("Topic retrieved successfully")
	c.JSON(http.StatusOK, topic)
}

func ApproveTopic(c *gin.Context) {
	logger := utils.GetLogger()
	name := c.Param("name")
	logger.Info("Received a request to approve topic")

	if name == "" {
		logger.Error("Topic name is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Topic name is required"})
		return
	}

	// Validate admin header
	admin := c.GetHeader("X-User-Id")
	if admin == "" {
		logger.Error("X-User-Id header is required for approval")
		c.JSON(http.StatusForbidden, gin.H{"error": "X-User-Id header is required"})
		return
	}
	logger.Debug("X-User-Id header validated")

	if err := service.ApproveTopic(c.Request.Context(), name, admin); err != nil {
		logger.Error("Failed to approve topic")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve topic"})
		return
	}
	logger.Info("Topic approved successfully")

	c.JSON(http.StatusOK, gin.H{"status": "approved"})
}
