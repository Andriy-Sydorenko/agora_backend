package user

import (
	"fmt"
	"github.com/matthewhartstonge/argon2"
)

func hashPassword(password string) (string, error) {
	argon := argon2.DefaultConfig()
	hashedPassword, err := argon.HashEncoded([]byte(password))
	if err != nil {
		return "", fmt.Errorf("argon2 hashing failed: %w", err)
	}
	return string(hashedPassword), nil
}

//nolint:unused
func verifyPassword(rawPassword, passwordHash string) bool {
	ok, _ := argon2.VerifyEncoded([]byte(rawPassword), []byte(passwordHash))
	return ok
}
