package bootstrap

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Server   ServerConfig
	DB       DBConfig
	Redis    RedisConfig
	JWT      JWTConfig
	CORS     CORSConfig
}

type ServerConfig struct {
	Port string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func (c DBConfig) DSN() string {
	return "postgres://" + c.User + ":" + c.Password +
		"@" + c.Host + ":" + c.Port + "/" + c.Name +
		"?sslmode=" + c.SSLMode
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func (c RedisConfig) Addr() string {
	return c.Host + ":" + c.Port
}

type JWTConfig struct {
	Secret     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type CORSConfig struct {
	Origins []string
}

func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: envStr("SERVER_PORT", "8080"),
		},
		DB: DBConfig{
			Host:     envStr("DB_HOST", "localhost"),
			Port:     envStr("DB_PORT", "5432"),
			User:     envStr("DB_USER", "miku"),
			Password: envStr("DB_PASSWORD", "miku_secret"),
			Name:     envStr("DB_NAME", "miku_blog"),
			SSLMode:  envStr("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     envStr("REDIS_HOST", "localhost"),
			Port:     envStr("REDIS_PORT", "6379"),
			Password: envStr("REDIS_PASSWORD", "miku_redis"),
			DB:       envInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:     envStr("JWT_SECRET", "change-me-in-production"),
			AccessTTL:  envDuration("JWT_ACCESS_TTL", 15*time.Minute),
			RefreshTTL: envDuration("JWT_REFRESH_TTL", 7*24*time.Hour),
		},
		CORS: CORSConfig{
			Origins: strings.Split(envStr("CORS_ORIGINS", "http://localhost:4321"), ","),
		},
	}
}

func envStr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func envDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
