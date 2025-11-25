package utils

import (
	"fmt"

	"github.com/matthewhartstonge/argon2"
)

// TODO: implementing password service interface seems like overhead rn, making functions exportable instead

func HashPassword(password string) (string, error) {
	argon := argon2.DefaultConfig()
	hashedPassword, err := argon.HashEncoded([]byte(password))
	if err != nil {
		return "", fmt.Errorf("argon2 hashing failed: %w", err)
	}
	return string(hashedPassword), nil
}

func VerifyPassword(rawPassword, passwordHash string) bool {
	ok, _ := argon2.VerifyEncoded([]byte(rawPassword), []byte(passwordHash))
	return ok
}
