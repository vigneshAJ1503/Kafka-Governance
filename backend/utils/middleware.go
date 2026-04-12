package utils

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JwtValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := GetLogger()
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn("Missing Authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Warn("Invalid Authorization header format")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format. Expected 'Bearer <token>'"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := ValidateJWT(tokenString)
		if err != nil {
			logger.Errorf("JWT validation failed: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set context variables
		c.Set("userId", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}
