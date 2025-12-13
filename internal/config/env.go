package config

import (
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBName     string
	DBPassword string
}

type ProjectConfig struct {
	IsProduction bool
	AppPort      int
}

type JWTConfig struct {
	Secret            string
	AccessLifetime    time.Duration
	RefreshLifetime   time.Duration
	JwtTokenCookieKey string
}

type CorsConfig struct {
	AllowedOrigins []string
}

func loadEnv() (*CorsConfig, *DatabaseConfig, *JWTConfig, *ProjectConfig) {
	if _, ok := os.LookupEnv("IS_DOCKER"); !ok {
		if err := godotenv.Load(); err != nil {
			log.Fatalln("⚠️ No .env file found, falling back to OS envs. Details:", err.Error())
		}
	}

	corsCfg := &CorsConfig{
		AllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", []string{"*"}, parseStringSlice),
	}
	dbCfg := &DatabaseConfig{
		DBHost:     getEnv("POSTGRES_HOST", "localhost", parseString),
		DBPort:     getEnv("POSTGRES_PORT", 5432, parseInt),
		DBUser:     getEnv("POSTGRES_USER", "postgres", parseString),
		DBName:     getEnv("POSTGRES_DB", "agora_db", parseString),
		DBPassword: getEnv("POSTGRES_PASSWORD", "password", parseString),
	}
	jwtCfg := &JWTConfig{
		Secret:            getEnv("JWT_SECRET_KEY", "supadupasecret", parseString),
		AccessLifetime:    getEnv("JWT_ACCESS_TOKEN_LIFETIME_SECONDS", 15*time.Minute, parseDuration),
		RefreshLifetime:   getEnv("JWT_REFRESH_TOKEN_LIFETIME_SECONDS", 24*time.Hour, parseDuration),
		JwtTokenCookieKey: getEnv("JWT_TOKEN_COOKIE_KEY", "token", parseString),
	}
	projectCfg := &ProjectConfig{
		IsProduction: getEnv("IS_PRODUCTION", false, parseBool),
		AppPort:      getEnv("APP_PORT", 8080, parseInt),
	}

	return corsCfg, dbCfg, jwtCfg, projectCfg
}

type parseFunc[T any] func(string) (T, error)

func getEnv[T any](key string, fallback T, parser parseFunc[T]) T {
	valueStr, ok := os.LookupEnv(key)

	if !ok || valueStr == "" {
		return fallback
	}

	val, err := parser(valueStr)
	if err != nil {
		log.Printf("⚠️ Invalid value for %s: %v. Using fallback: %v", key, err, fallback)
		return fallback
	}
	return val
}

func parseString(val string) (string, error) {
	return val, nil
}

func parseInt(val string) (int, error) {
	return strconv.Atoi(val)
}

//nolint:unused
func parseBool(val string) (bool, error) {
	return strconv.ParseBool(val)
}

func parseDuration(val string) (time.Duration, error) {
	return time.ParseDuration(val)
}

func parseStringSlice(val string) ([]string, error) {
	if val == "" {
		return nil, errors.New("value is empty")
	}
	parts := strings.Split(val, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts, nil
}
