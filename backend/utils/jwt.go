package utils

import (
	"errors"
	"time"

	"kafka-governance/config"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaim struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID, username, role string) (string, error) {
	cfg := config.Load()
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &JWTClaim{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateJWT(signedToken string) (*JWTClaim, error) {
	cfg := config.Load()

	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok || !token.Valid {
		return nil, errors.New("invalid jwt token")
	}

	return claims, nil
}
