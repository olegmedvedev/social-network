package config

import (
	"log"
	"os"
)

type DbConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type Config struct {
	DB        DbConfig
	JWTSecret string
}

func Load() *Config {
	dbCfg := DbConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}
	if dbCfg.Host == "" || dbCfg.User == "" || dbCfg.Password == "" || dbCfg.Name == "" {
		log.Fatal("Database environment variables are not set")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}
	return &Config{
		DB:        dbCfg,
		JWTSecret: jwtSecret,
	}
}
