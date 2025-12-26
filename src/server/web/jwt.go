package web

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/golang-jwt/jwt/v5"
)

func (web *Web) parseAndValidateJWT(tokenString string) (*jwt.Token, bool) {
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("invalid jwt signing alg")
		}

		return web.JWTSecret, nil
	})

	if err != nil {
		return nil, false
	}

	if !parsedToken.Valid {
		return nil, false
	}

	return parsedToken, true
}

func (web *Web) signJWT(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(web.JWTSecret)
	return tokenString, err
}

func (web *Web) createSignedStateJWT(redirectTo string) (string, error) {
	tokenString, err := web.signJWT(jwt.MapClaims{
		"redirect_to": redirectTo,
		"exp":         time.Now().Add(time.Hour).Unix(),
		"iat":         time.Now(),
	})
	return tokenString, err
}

func (web *Web) getRedirectToFromJWT(tokenString string) (string, error) {
	parsedToken, valid := web.parseAndValidateJWT(tokenString)

	if !valid {
		return "", fmt.Errorf("invalid token")
	}

	claims := parsedToken.Claims.(jwt.MapClaims)
	redirectTo := claims["redirect_to"].(string)

	return redirectTo, nil
}

func (web *Web) createSignedUserJWT(userId string) (string, error) {
	tokenString, err := web.signJWT(jwt.MapClaims{
		"sub": userId,
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now(),
	})

	return tokenString, err
}

func (web *Web) getUserIdFromJWT(tokenString string) (string, error) {
	parsedToken, valid := web.parseAndValidateJWT(tokenString)

	if !valid {
		return "", fmt.Errorf("invalid token")
	}

	userId, err := parsedToken.Claims.GetSubject()

	if err != nil {
		return "", fmt.Errorf("missing user id")
	}

	return userId, nil
}

func (web *Web) AuthMiddleware(ctx huma.Context, next func(huma.Context)) {
	authHeader := ctx.Header("Authorization")

	if authHeader == "" {
		huma.WriteErr(web.API, ctx, http.StatusUnauthorized, "missing auth header")
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		huma.WriteErr(web.API, ctx, http.StatusUnauthorized, "invalid auth header format")
		return
	}

	token := parts[1]

	userId, err := web.getUserIdFromJWT(token)
	if err != nil {
		huma.WriteErr(web.API, ctx, http.StatusUnauthorized, "invalid token")
		return
	}

	user, err := web.DB.GetUserById(userId)
	if err != nil {
		huma.WriteErr(web.API, ctx, http.StatusUnauthorized, "missing user")
		return
	}

	newCtx := context.WithValue(ctx.Context(), "user", &user)
	ctx = huma.WithContext(ctx, newCtx)
	next(ctx)
}
