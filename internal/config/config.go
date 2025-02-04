package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

// Rest is a configuration for REST API.
type Rest struct {
	Port               int      `env:"PORT" env-default:"8004"`
	Host               string   `env:"HOST" env-default:"0.0.0.0"`
	AllowedCORSOrigins []string `env:"ALLOWED_CORS_ORIGINS" envSeparator:"," env-default:"*"`
}

// Project - project configuration.
type Project struct {
	Debug       bool   `env:"DEBUG" env-default:"false"`
	Name        string `env:"PROJECT_NAME" env-default:"choco_parser_go"`
	Environment string `env:"ENVIRONMENT" env-default:"development"`
	ServiceName string `env:"SERVICE_NAME" env-default:"monolith"`
}

// Database is a configuration for database.
type Database struct {
	DSN string `env:"DATABASE_DSN" env-default:"postgres://postgres:postgres@0.0.0.0:5499/choco?sslmode=disable"`
}

type Choco struct {
	ClientId    int64  `env:"CHOCO_CLIENT_ID" env-default:"" env-description:"Choco API client id"`
	FingerPrint string `env:"CHOCO_X_FINGERPRINT" env-default:"" env-description:"Choco API fingerprint"`
	ChocoToken  string `env:"CHOCO_AUTH_TOKEN"`
}

// Logger is a configuration for logger.
type Logger struct {
	LogLevel string `env:"LOG_LEVEL" env-default:"info"`
	DevMode  bool   `env:"DEV_MODE" env-default:"false"`
	Encoder  string `env:"ENCODER" env-default:"console"`
}

// Config - contains all configuration parameters in config package.
type Config struct {
	Project  Project
	Database Database
	Logger   Logger
	Choco    Choco
	Rest     Rest
}

func Get() (Config, error) {
	config := Config{}

	if err := cleanenv.ReadEnv(&config.Project); err != nil {
		return config, fmt.Errorf("error reading project config: %w", err)
	}

	if err := cleanenv.ReadEnv(&config.Database); err != nil {
		return config, fmt.Errorf("error reading database config: %w", err)
	}

	if err := cleanenv.ReadEnv(&config.Logger); err != nil {
		return config, fmt.Errorf("error reading logger config: %w", err)
	}

	if err := cleanenv.ReadEnv(&config.Choco); err != nil {
		return config, fmt.Errorf("error reading choco config: %w", err)
	}

	if err := cleanenv.ReadEnv(&config.Rest); err != nil {
		return config, fmt.Errorf("error reading rest config: %w", err)
	}

	return config, nil
}
