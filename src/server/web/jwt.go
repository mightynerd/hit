package web

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

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

func (web *Web) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing auth header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid auth header format", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		userId, err := web.getUserIdFromJWT(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		user, err := web.db.GetUserById(userId)
		if err != nil {
			http.Error(w, "Missing user", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, &user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
