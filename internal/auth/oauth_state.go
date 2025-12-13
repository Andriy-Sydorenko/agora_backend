package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
)

func GenerateState(jwtSecret string) (string, error) {
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	randomB64 := base64.URLEncoding.EncodeToString(randomBytes)

	mac := hmac.New(sha256.New, []byte(jwtSecret))
	mac.Write([]byte(randomB64))
	signature := base64.URLEncoding.EncodeToString(mac.Sum(nil))

	state := fmt.Sprintf("%s.%s", randomB64, signature)
	return state, nil
}

func ValidateState(state, jwtSecret string) (bool, error) {
	parts := strings.SplitN(state, ".", 2)
	if len(parts) != 2 {
		return false, fmt.Errorf("invalid state format")
	}

	randomB64, signatureB64 := parts[0], parts[1]

	mac := hmac.New(sha256.New, []byte(jwtSecret))
	mac.Write([]byte(randomB64))
	expectedSignature := base64.URLEncoding.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(signatureB64), []byte(expectedSignature)) {
		return false, fmt.Errorf("invalid state signature")
	}

	return true, nil
}
