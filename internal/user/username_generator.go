package user

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"strings"
)

var (
	prefixes = []string{
		"cosmic", "cyber", "neon", "swift", "shadow", "pixel", "quantum",
		"ultra", "hyper", "nova", "astro", "turbo", "apex", "prime",
		"phantom", "stellar", "zero", "omega", "delta", "fusion",
	}

	suffixes = []string{
		"x", "io", "ix", "os", "ax", "ex", "is", "us", "on", "an",
	}

	emailCleanupRegex = regexp.MustCompile(`[^a-zA-Z0-9]`)
)

// GenerateUsernameFromEmail creates a unique, cool username from an email
// Format: {prefix}_{base}_{4-char-hex} or {base}_{suffix}_{4-char-hex}
func GenerateUsernameFromEmail(email string) string {
	base := extractBase(email)
	salt := generateShortSalt()

	if randomBool() {
		prefix := prefixes[randomIndex(len(prefixes))]
		return truncateUsername(prefix + "_" + base + "_" + salt)
	}

	suffix := suffixes[randomIndex(len(suffixes))]
	return truncateUsername(base + suffix + "_" + salt)
}

func extractBase(email string) string {
	parts := strings.Split(email, "@")
	local := parts[0]

	local = strings.Split(local, "+")[0] // Remove +alias
	local = emailCleanupRegex.ReplaceAllString(local, "")
	local = strings.ToLower(local)

	if len(local) > 12 {
		local = local[:12]
	}

	return local
}

func generateShortSalt() string {
	bytes := make([]byte, 2)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func randomIndex(max int) int {
	bytes := make([]byte, 1)
	rand.Read(bytes)
	return int(bytes[0]) % max
}

func randomBool() bool {
	bytes := make([]byte, 1)
	rand.Read(bytes)
	return bytes[0]%2 == 0
}

func truncateUsername(username string) string {
	if len(username) > 50 { // UsernameMaxLen from validator
		return username[:50]
	}
	return username
}
