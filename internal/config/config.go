package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"time"
)

type Config struct {
	App      AppConfig     `yaml:"app"`
	Server   ServerConfig  `yaml:"server"`
	Logging  LoggingConfig `yaml:"logging"`
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Project  ProjectConfig
	Google   GoogleConfig
}
type AppConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type ServerConfig struct {
	ReadTimeout  time.Duration `yaml:"read_timeout"` // TODO: consider mapstructure instead of yml
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
	Cors         CorsConfig
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

func Load(path string) *Config {
	cfg := new(Config)

	corsCfg, dbCfg, redisCfg, jwtCfg, projectCfg, googleCfg := loadEnv()

	yamlFile, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln("Error reading config.yaml:", err.Error())
		return nil
	}
	err = yaml.Unmarshal(yamlFile, cfg)
	if err != nil {
		log.Fatalln("Error parsing config.yaml:", err.Error())
		return nil
	}

	cfg.Server.Cors = corsCfg
	cfg.Database = dbCfg
	cfg.Redis = redisCfg
	cfg.JWT = jwtCfg
	cfg.Project = projectCfg
	cfg.Google = googleCfg

	return cfg
}
