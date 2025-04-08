package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port        int
	Environment string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret     string
	Expiration int // in minutes
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigName("config")        // name of config file (without extension)
	viper.SetConfigType("yaml")          // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(path)            // path to look for the config file in
	viper.AddConfigPath(".")             // look in current directory
	viper.AddConfigPath("/")             // look in root directory
	viper.AddConfigPath("/home/nonroot") // look in nonroot home directory
	viper.AutomaticEnv()

	// Environment variables override
	// Map environment variables to config keys
	viper.SetEnvPrefix("APP")
	viper.BindEnv("server.port", "SERVER_PORT")
	viper.BindEnv("server.environment", "SERVER_ENVIRONMENT")
	viper.BindEnv("database.host", "DATABASE_HOST")
	viper.BindEnv("database.port", "DATABASE_PORT")
	viper.BindEnv("database.user", "DATABASE_USER")
	viper.BindEnv("database.password", "DATABASE_PASSWORD")
	viper.BindEnv("database.dbname", "DATABASE_DBNAME")
	viper.BindEnv("database.sslmode", "DATABASE_SSLMODE")
	viper.BindEnv("jwt.secret", "JWT_SECRET")
	viper.BindEnv("jwt.expiration", "JWT_EXPIRATION")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}
