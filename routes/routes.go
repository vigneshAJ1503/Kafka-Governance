package routes

import (
	"kafka-governance/api"
	"kafka-governance/utils"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	r.Use(utils.GinLoggingMiddleware())

	v1 := r.Group("/api/v1")
	{
		v1.POST("/topics", api.CreateTopic)
		v1.GET("/topics", api.ListTopics)
		v1.GET("/topics/:name", api.GetTopic)
		v1.POST("/topics/:name/approve", api.ApproveTopic)
		v1.POST("/policies", api.CreatePolicy)
	}
}
