package utils

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var (
	ErrInvalidAccessToken  = errors.New("invalid access token")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrInvalidToken        = errors.New("invalid JWT token")
	ErrExpiredToken        = errors.New("JWT token has expired")
	ErrInvalidClaims       = errors.New("invalid JWT claims")
	ErrInvalidTokenType    = errors.New("invalid token type")
)

const (
	TokenTypeAccess             = "access"
	TokenTypeRefresh            = "refresh"
	RefreshTokenBlacklistPrefix = "refresh_token_blacklist:"
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

func GenerateJWT(jwtSecret string, tokenType string, tokenLifetime time.Duration, userID string) (string, error) {
	if tokenType != TokenTypeAccess && tokenType != TokenTypeRefresh {
		return "", ErrInvalidTokenType
	}

	tokenExpiry := time.Now().Add(tokenLifetime)

	// TODO: Using default algorithm, can be changed later
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  userID,
		"exp":  tokenExpiry.Unix(),
		"type": tokenType,
	})
	token, err := tokenObj.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", errors.New(fmt.Sprintln("failed to generate JWT token:", err))
	}
	return token, nil
}

func DecryptJWT(tokenString string, jwtSecret string, expectedTokenType string) (string, jwt.MapClaims, error) {
	if tokenString == "" {
		return "", nil, ErrInvalidToken
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", nil, ErrExpiredToken
		}
		return "", nil, ErrInvalidToken
	}

	if !token.Valid {
		return "", nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", nil, ErrInvalidClaims
	}

	if tokenType, ok := claims["type"]; !ok || tokenType != expectedTokenType {
		return "", nil, ErrInvalidTokenType
	}

	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		return "", nil, ErrInvalidClaims
	}
	return userID, claims, nil
}
