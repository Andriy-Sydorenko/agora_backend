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

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
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

type GoogleConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	SMTPUseTLS   bool

	ClientID          string
	ClientSecret      string
	ClientRedirectURL string
}

type CorsConfig struct {
	AllowedOrigins []string
}

func loadEnv() (CorsConfig, DatabaseConfig, RedisConfig, JWTConfig, ProjectConfig, GoogleConfig) {
	if _, ok := os.LookupEnv("IS_DOCKER"); !ok {
		if err := godotenv.Load(); err != nil {
			log.Fatalln("⚠️ No .env file found, falling back to OS envs. Details:", err.Error())
		}
	}

	corsCfg := CorsConfig{
		AllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", []string{"*"}, parseStringSlice),
	}
	dbCfg := DatabaseConfig{
		DBHost:     getEnv("POSTGRES_HOST", "localhost", parseString),
		DBPort:     getEnv("POSTGRES_PORT", 5432, parseInt),
		DBUser:     getEnv("POSTGRES_USER", "postgres", parseString),
		DBName:     getEnv("POSTGRES_DB", "agora_db", parseString),
		DBPassword: getEnv("POSTGRES_PASSWORD", "password", parseString),
	}
	redisCfg := RedisConfig{
		Host:     getEnv("REDIS_HOST", "localhost", parseString),
		Port:     getEnv("REDIS_PORT", 6379, parseInt),
		Password: getEnv("REDIS_PASSWORD", "", parseString),
		DB:       getEnv("REDIS_DB_NUM", 0, parseInt),
	}
	jwtCfg := JWTConfig{
		Secret:            getEnv("JWT_SECRET_KEY", "supadupasecret", parseString),
		AccessLifetime:    getEnv("JWT_ACCESS_TOKEN_LIFETIME_SECONDS", 15*time.Minute, parseDuration),
		RefreshLifetime:   getEnv("JWT_REFRESH_TOKEN_LIFETIME_SECONDS", 24*time.Hour, parseDuration),
		JwtTokenCookieKey: getEnv("JWT_TOKEN_COOKIE_KEY", "token", parseString),
	}
	projectCfg := ProjectConfig{
		IsProduction: getEnv("IS_PRODUCTION", false, parseBool),
		AppPort:      getEnv("APP_PORT", 8080, parseInt),
	}
	googleCfg := GoogleConfig{
		SMTPHost:          getEnv("GOOGLE_SMTP_HOST", "smtp.gmail.com", parseString),
		SMTPPort:          getEnv("GOOGLE_SMTP_PORT", 587, parseInt),
		SMTPUsername:      getEnv("GOOGLE_SMTP_USERNAME", "email@gmail.com", parseString),
		SMTPPassword:      getEnv("GOOGLE_SMTP_PASSWORD", "somepassword", parseString),
		SMTPUseTLS:        getEnv("GOOGLE_SMTP_USE_TLS", true, parseBool),
		ClientID:          getEnv("GOOGLE_CLIENT_ID", "google_client_id", parseString),
		ClientSecret:      getEnv("GOOGLE_CLIENT_SECRET", "supadupasecret", parseString),
		ClientRedirectURL: getEnv("GOOGLE_REDIRECT_URL", "someurl.com", parseString),
	}

	return corsCfg, dbCfg, redisCfg, jwtCfg, projectCfg, googleCfg
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
