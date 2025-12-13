package database

import (
	"fmt"
	"log"
	"net/url"

	"github.com/Andriy-Sydorenko/agora_backend/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func buildDBDSN(cfg *config.DatabaseConfig) string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.DBUser, cfg.DBPassword),
		Host:   fmt.Sprintf("%s:%d", cfg.DBHost, cfg.DBPort),
		Path:   cfg.DBName,
	}
	q := u.Query()
	q.Add("sslmode", "disable")
	u.RawQuery = q.Encode()
	return u.String()
}

func Connect(cfg *config.DatabaseConfig) *gorm.DB {
	dsn := buildDBDSN(cfg)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		log.Fatalln("failed to connect to database:", err)
		return nil
	}
	return db
}
