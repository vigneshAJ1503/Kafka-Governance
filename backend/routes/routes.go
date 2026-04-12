package routes

import (
	"kafka-governance/api"
	"kafka-governance/utils"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	r.Use(utils.GinLoggingMiddleware())

	auth := r.Group("/api/v1/auth")
	{
		auth.POST("/signup", api.Signup)
		auth.POST("/signin", api.Signin)
	}

	v1 := r.Group("/api/v1")
	v1.Use(utils.JwtValidationMiddleware())
	{
		v1.POST("/topics", api.CreateTopic)
		v1.GET("/topics", api.ListTopics)
		v1.GET("/topics/:name", api.GetTopic)
		v1.POST("/topics/:name/approve", api.ApproveTopic)
		v1.POST("/policies", api.CreatePolicy)
	}
}
