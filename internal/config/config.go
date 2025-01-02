package config

import (
	"github.com/joho/godotenv"
	"github.com/ziliscite/messaging-app/internal/util"
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
		DBHost: util.MustEnv(os.Getenv("DB_HOST")),
		DBPort: util.MustEnv(os.Getenv("DB_PORT")),
		DBUser: util.MustEnv(os.Getenv("DB_USER")),
		DBPass: util.MustEnv(os.Getenv("DB_PASSWORD")),
		DBName: util.MustEnv(os.Getenv("DB_NAME")),
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
	util.MustServe(godotenv.Load())

	database := newDatabaseConfig()

	return &Config{
		Database:    database,
		Port:        util.MustEnv(os.Getenv("PORT")),
		Environment: util.MustEnv(os.Getenv("ENVIRONMENT")),
		Secret:      util.MustEnv(os.Getenv("JWT_SECRET")),
	}
}

func (c *Config) Address() string {
	if c.Environment == "production" {
		return "0.0.0.0" + c.Port
	}

	return "localhost" + c.Port
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
