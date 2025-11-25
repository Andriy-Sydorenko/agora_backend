package utils

import (
	"errors"
	"fmt"
	"github.com/Andriy-Sydorenko/agora_backend/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func GenerateJWT(cfg config.JWTConfig, userID string) (string, error) {
	// TODO: Using default algorithm, can be changed later
	jwtSecret := []byte(cfg.Secret)
	tokenExpiry := time.Now().Add(cfg.AccessLifetime)
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": tokenExpiry.Unix(),
	})
	token, err := tokenObj.SignedString(jwtSecret)
	if err != nil {
		return "", errors.New(fmt.Sprintln("failed to generate JWT token:", err))
	}
	return token, nil
}
