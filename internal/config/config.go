package config
package config




























}	return defaultValue	}		return value	if value := os.Getenv(key); value != "" {func getEnv(key, defaultValue string) string {}	}, nil		RateLimitRPM: 1000,		JWTSecret:    getEnv("JWT_SECRET", "your-secret-key"),		RedisURL:     getEnv("REDIS_URL", "redis://localhost:6379"),		DatabaseURL:  getEnv("DATABASE_URL", "postgres://user:pass@localhost:5432/events?sslmode=disable"),		Port:         getEnv("PORT", "8080"),	return &Config{func Load() (*Config, error) {}	RateLimitRPM int	JWTSecret    string	RedisURL     string	DatabaseURL  string	Port         stringtype Config struct {)	"os"import (