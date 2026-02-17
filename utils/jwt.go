package utils

import (
	"be-jalan/models"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken creates a new JWT token for the authenticated user
func GenerateToken(user models.User) (string, error) {
	// Get secret key from .env
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return "", fmt.Errorf("JWT_SECRET is not set in environment")
	}

	// Get token expiration from .env
	expirationHoursStr := os.Getenv("JWT_EXPIRATION_HOURS")
	if expirationHoursStr == "" {
		expirationHoursStr = "24" // default to 24 hours
	}

	expirationHours, _ := strconv.Atoi(expirationHoursStr)
	expirationTime := time.Now().Add(time.Duration(expirationHours) * time.Hour)

	// Create claims
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   expirationTime.Unix(),
		"iat":   time.Now().Unix(),
	}

	// Sign the token with the secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken parses and validates a JWT token
func ValidateToken(tokenString string) (*jwt.Token, error) {
	secretKey := os.Getenv("JWT_SECRET")

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify signature method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
}
