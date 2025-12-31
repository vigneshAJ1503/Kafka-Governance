package api

import (
	"net/http"

	"kafka-governance/models"
	"kafka-governance/service"
	"kafka-governance/utils"

	"github.com/gin-gonic/gin"
)

var policies []models.Policy

func CreatePolicy(c *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("Received a request to create a policy")

	var p models.Policy
	if err := c.ShouldBindJSON(&p); err != nil {
		logger.Error("Failed to decode policy request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	logger.Debug("Policy request body decoded successfully")

	// Validate required fields
	if p.Principal == "" {
		logger.Error("Principal is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Principal is required"})
		return
	}

	if p.Action == "" {
		logger.Error("Action is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Action is required"})
		return
	}

	if p.Resource == "" {
		logger.Error("Resource is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Resource is required"})
		return
	}

	if p.Effect != "permit" && p.Effect != "forbid" {
		logger.Error("Effect must be either 'permit' or 'forbid'")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Effect must be either 'permit' or 'forbid'"})
		return
	}
	logger.Debug("Policy validation passed")

	if err := service.CreatePolicy(c.Request.Context(), p); err != nil {
		logger.Error("Failed to create policy")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create policy"})
		return
	}
	logger.Info("Policy created successfully")

	c.JSON(http.StatusCreated, p)
}
