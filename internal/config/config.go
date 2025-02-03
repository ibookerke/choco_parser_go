package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

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

	return config, nil
}
