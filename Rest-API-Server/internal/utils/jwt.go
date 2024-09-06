package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gocql/gocql"
	"net/http"
	"strings"
	"time"
)

var jwtKey = []byte("Baylus")

type Claims struct {
	UserID gocql.UUID `json:"user_id"`
	Email  string     `json:"email"`
	jwt.StandardClaims
}

func GenerateToken(userID gocql.UUID, email string) (string, string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.String(),
		"email":   email,
		"exp":     time.Now().Add(time.Second * 10).Unix(), // 7 days expiration
	})
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.String(),
		"email":   email,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days expiration
	})
	accessTokenString, err := accessToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}
	refreshTokenString, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func ExtractUserIdAndEmailFromContext(tokenString string) (gocql.UUID, string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return gocql.UUID{}, "", err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return gocql.UUID{}, "", errors.New("invalid token")
	}

	userID := claims.UserID
	email := claims.Email
	return userID, email, nil
}

func ExtractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header missing")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("authorization header format must be Bearer {token}")
	}

	return strings.TrimPrefix(authHeader, "Bearer "), nil
}
