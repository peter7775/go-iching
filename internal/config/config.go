package config

import "os"

type Config struct {
	Addr        string
	Storage     string
	SQLitePath  string
	PostgresDSN string
}

func Load() Config {
	return Config{
		Addr:        getenv("APP_ADDR", ":8080"),
		Storage:     getenv("APP_STORAGE", "sqlite"),
		SQLitePath:  getenv("APP_SQLITE_PATH", "./iching.db"),
		PostgresDSN: getenv("APP_PG_DSN", "postgres://postgres:postgres@localhost:5432/iching?sslmode=disable"),
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
