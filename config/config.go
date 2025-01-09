package config

import (
	"github.com/joho/godotenv"
	"github.com/ziliscite/messaging-app/pkg/must"
	"os"
	"strconv"
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

type TokenConfig struct {
	Secret                     string
	AccessTokenExpirationTime  int64
	RefreshTokenExpirationTime int64
}

func newTokenConfig() *TokenConfig {
	accessExp := must.Must(strconv.Atoi(must.MustEnv(os.Getenv("ACCESS_TOKEN_EXPIRATION_MINUTES"))))
	refreshExp := must.Must(strconv.Atoi(must.MustEnv(os.Getenv("REFRESH_TOKEN_EXPIRATION_MINUTES"))))

	return &TokenConfig{
		Secret:                     must.MustEnv(os.Getenv("JWT_SECRET")),
		AccessTokenExpirationTime:  int64(accessExp),
		RefreshTokenExpirationTime: int64(refreshExp),
	}
}

type Config struct {
	// Database connection string
	Database *DatabaseConfig

	// Mongo connection string
	Mongo string

	// Elastic elk connection string
	Elastic string

	// Environment is development or production
	Environment string

	// Port server is running on
	Port string

	// WebsocketPort server is running on
	WebsocketPort string

	// Token config for JWT
	Token *TokenConfig
}

func New() *Config {
	must.MustServe(godotenv.Load())

	database := newDatabaseConfig()
	token := newTokenConfig()

	return &Config{
		Database:      database,
		Mongo:         must.MustEnv(os.Getenv("MONGO_URI")),
		Elastic:       must.MustEnv(os.Getenv("ELASTIC_APM_SERVER_URL")),
		Port:          must.MustEnv(os.Getenv("PORT")),
		WebsocketPort: must.MustEnv(os.Getenv("WEB_SOCKET_PORT")),
		Environment:   must.MustEnv(os.Getenv("ENVIRONMENT")),
		Token:         token,
	}
}

func (c *Config) Address() string {
	if c.Environment == "production" {
		return "0.0.0.0" + c.Port
	}

	return "localhost" + c.Port
}

func (c *Config) WebsocketAddress() string {
	if c.Environment == "production" {
		return "0.0.0.0" + c.WebsocketPort
	}

	return "localhost" + c.WebsocketPort
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
