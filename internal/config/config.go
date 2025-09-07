package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type ServerConfig struct {
	Port string
	Mode string
}

func Load() *Config {
	godotenv.Load()

	host := getEnv("DB_HOST", "db")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "todouser")
	password := getEnv("DB_PASSWORD", "todopass123")
	dbName := getEnv("DB_NAME", "todoapi")
	sslMode := getEnv("DB_SSLMODE", "disable")

	serverPort := getEnv("SERVER_PORT", "8080")
	serverMode := getEnv("SERVER_MODE", "debug")

	config := &Config{
		Database: DatabaseConfig{
			Host:     host,
			Port:     port,
			User:     user,
			Password: password,
			Name:     dbName,
			SSLMode:  sslMode,
		},
		Server: ServerConfig{
			Port: serverPort,
			Mode: serverMode,
		},
	}

	return config
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode)
}
