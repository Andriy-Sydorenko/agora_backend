package utils

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var (
	ErrInvalidToken  = errors.New("invalid JWT token")
	ErrExpiredToken  = errors.New("JWT token has expired")
	ErrInvalidClaims = errors.New("invalid JWT claims")
)

func GenerateJWT(jwtSecret string, accessTokenLifetime time.Duration, userID string) (string, error) {
	// TODO: Using default algorithm, can be changed later
	jwtSecretSlice := []byte(jwtSecret)
	tokenExpiry := time.Now().Add(accessTokenLifetime)
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": tokenExpiry.Unix(),
	})
	token, err := tokenObj.SignedString(jwtSecretSlice)
	if err != nil {
		return "", errors.New(fmt.Sprintln("failed to generate JWT token:", err))
	}
	return token, nil
}

func DecryptJWT(tokenString string, jwtSecret string) (string, error) {
	if tokenString == "" {
		return "", ErrInvalidToken
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", ErrExpiredToken
		}
		return "", ErrInvalidToken
	}

	if !token.Valid {
		return "", ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrInvalidClaims
	}
	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		return "", ErrInvalidClaims
	}
	return userID, nil
}
