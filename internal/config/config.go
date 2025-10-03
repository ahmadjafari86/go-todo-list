package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port          string
	DatabaseURL   string
	JWTSecret     string
	JWTExpMinutes int

	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime int

	ReadTimeout  int
	WriteTimeout int
	IdleTimeout  int
}

func getenvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func New() *Config {
	exp := getenvInt("JWT_EXP_MINUTES", 60)
	return &Config{
		Port:              os.Getenv("PORT"),
		DatabaseURL:       os.Getenv("DATABASE_URL"),
		JWTSecret:         os.Getenv("JWT_SECRET"),
		JWTExpMinutes:     exp,
		DBMaxOpenConns:    getenvInt("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns:    getenvInt("DB_MAX_IDLE_CONNS", 25),
		DBConnMaxLifetime: getenvInt("DB_CONN_MAX_LIFETIME", 300),
		ReadTimeout:       getenvInt("READ_TIMEOUT", 5),
		WriteTimeout:      getenvInt("WRITE_TIMEOUT", 10),
		IdleTimeout:       getenvInt("IDLE_TIMEOUT", 120),
	}
}

func (c *Config) PortOrDefault() string {
	if c.Port == "" {
		return "8282"
	}
	return c.Port
}
