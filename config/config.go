package config

import (
	"github.com/joho/godotenv"
	"github.com/ziliscite/messaging-app/pkg/must"
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
		DBHost: must.MustEnv(os.Getenv("DB_HOST")),
		DBPort: must.MustEnv(os.Getenv("DB_PORT")),
		DBUser: must.MustEnv(os.Getenv("DB_USER")),
		DBPass: must.MustEnv(os.Getenv("DB_PASSWORD")),
		DBName: must.MustEnv(os.Getenv("DB_NAME")),
	}
}

func (c *DatabaseConfig) ConnectionString() string {
	return "host=" + c.DBHost + " port=" + c.DBPort + " user=" + c.DBUser + " password=" + c.DBPass + " dbname=" + c.DBName + " sslmode=disable"
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
	must.MustServe(godotenv.Load())

	database := newDatabaseConfig()

	return &Config{
		Database:    database,
		Port:        must.MustEnv(os.Getenv("PORT")),
		Environment: must.MustEnv(os.Getenv("ENVIRONMENT")),
		Secret:      must.MustEnv(os.Getenv("JWT_SECRET")),
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
