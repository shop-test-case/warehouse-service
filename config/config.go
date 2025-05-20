package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	DBUser    string
	DBPass    string
	DBHost    string
	DBPort    string
	DBName    string
	JWTSecret string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		Port:      os.Getenv("PORT"),
		DBUser:    os.Getenv("DB_USER"),
		DBPass:    os.Getenv("DB_PASSWORD"),
		DBHost:    os.Getenv("DB_HOST"),
		DBPort:    os.Getenv("DB_PORT"),
		DBName:    os.Getenv("DB_NAME"),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}
}
