package api

import (
	"net/http"
	"strings"

	"kafka-governance/db"
	"kafka-governance/models"
	"kafka-governance/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	logger := utils.GetLogger()
	var req models.SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to decode signup request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Failed to hash password")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	role := "USER"
	if strings.ToUpper(req.Role) == "ADMIN" {
		role = "ADMIN"
	}

	user := &models.User{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		Role:         role,
	}

	createdUser, err := db.CreateUser(c.Request.Context(), user)
	if err != nil {
		logger.Errorf("Failed to create user: %v", err)
		c.JSON(http.StatusConflict, gin.H{"error": "Failed to create user, username might exist"})
		return
	}

	// generate token
	token, err := utils.GenerateJWT(createdUser.ID, createdUser.Username, createdUser.Role)
	if err != nil {
		logger.Error("Failed to generate JWT")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User created but failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"token":   token,
		"user": gin.H{
			"id":       createdUser.ID,
			"username": createdUser.Username,
			"role":     createdUser.Role,
		},
	})
}

func Signin(c *gin.Context) {
	logger := utils.GetLogger()
	var req models.SigninRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to decode signin request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user, err := db.GetUserByUsername(c.Request.Context(), req.Username)
	if err != nil {
		logger.Warn("Signin failed: user not found")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		logger.Warn("Signin failed: incorrect password")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Username, user.Role)
	if err != nil {
		logger.Error("Failed to generate JWT on signin")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}
