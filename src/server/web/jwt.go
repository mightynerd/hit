package web

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (web *Web) createSignedUserJWT(userId string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userId,
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(web.jwtSecret)
	return tokenString, err
}

func (web *Web) getUserIdFromJWT(tokenString string) (string, error) {
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("invalid jwt signing alg")
		}

		return web.jwtSecret, nil
	})

	if err != nil {
		return "", err
	}

	if !parsedToken.Valid {
		return "", fmt.Errorf("invalid token")
	}

	userId, err := parsedToken.Claims.GetSubject()

	if err != nil {
		return "", fmt.Errorf("missing user id")
	}

	return userId, nil
}

type contextKey string

const userContextKey = contextKey("user")

func (web *Web) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing auth header"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid auth header format"})
			return
		}

		token := parts[1]

		userId, err := web.getUserIdFromJWT(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		user, err := web.db.GetUserById(userId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing user"})
			return
		}

		c.Set("user", &user)
		c.Next()
	}
}
