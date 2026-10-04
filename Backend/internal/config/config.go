package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort        string
	GinMode        string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBSSLMode      string
	DBTimezone     string
	JWTSecret      string
	JWTExpireHours int
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	hours, err := strconv.Atoi(os.Getenv("JWT_EXPIRE_HOURS"))
	if err != nil {
		hours = 8
	}

	cfg := &Config{
		AppPort:        os.Getenv("APP_PORT"),
		GinMode:        os.Getenv("GIN_MODE"),
		DBHost:         os.Getenv("DB_HOST"),
		DBPort:         os.Getenv("DB_PORT"),
		DBUser:         os.Getenv("DB_USER"),
		DBPassword:     os.Getenv("DB_PASSWORD"),
		DBName:         os.Getenv("DB_NAME"),
		DBSSLMode:      os.Getenv("DB_SSL_MODE"),
		DBTimezone:     os.Getenv("DB_TIMEZONE"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		JWTExpireHours: hours,
	}

	if cfg.JWTSecret == "" {
		log.Fatal("JWT secret not set")
	}
	return cfg
}
