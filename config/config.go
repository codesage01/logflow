package config

import "os"

type Config struct {
    Port        string
    Env         string // e.g., "production" or "development"
    DatabaseURL string // Railway provides this automatically if you add a DB
}

func Load() *Config {
    return &Config{
        Port:        getEnv("PORT", "8080"),
        Env:         getEnv("ENV", "development"),
        DatabaseURL: os.Getenv("DATABASE_URL"), // Added for future-proofing
    }
}

// Helper function to keep code clean
func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}
