package config

import (
	"github.com/joho/godotenv"
	"github.com/ziliscite/messaging-app/pkg"
	"os"
)

type DatabaseConfig struct {
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string
}

func newDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		DBHost: pkg.MustEnv(os.Getenv("DB_HOST")),
		DBPort: pkg.MustEnv(os.Getenv("DB_PORT")),
		DBUser: pkg.MustEnv(os.Getenv("DB_USER")),
		DBPass: pkg.MustEnv(os.Getenv("DB_PASSWORD")),
		DBName: pkg.MustEnv(os.Getenv("DB_NAME")),
	}
}

type Config struct {
	// Database connection string
	Database *DatabaseConfig

	// Environment is development or production
	Environment string

	// Port server is running on
	Port string

	// Secret JWT key
	Secret string
}

func New() *Config {
	pkg.MustServe(godotenv.Load())

	database := newDatabaseConfig()

	return &Config{
		Database:    database,
		Port:        pkg.MustEnv(os.Getenv("PORT")),
		Environment: pkg.MustEnv(os.Getenv("ENVIRONMENT")),
		Secret:      pkg.MustEnv(os.Getenv("JWT_SECRET")),
	}
}

func (c *Config) Address() string {
	if c.Environment == "production" {
		return "0.0.0.0" + c.Port
	}

	return "localhost" + c.Port
}
