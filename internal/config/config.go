package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
)

type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBname   string `yaml:"dbname"`
	Sslmode  string `yaml:"sslmode"`
}

type CorsConfig struct {
	AllowedOrigins   []string `yaml:"allowed_origins"`
	AllowCredentials bool     `yaml:"allowed_credentials"`
	AllowedMethods   []string `yaml:"allowed_methods"`
	AllowedHeaders   []string `yaml:"allowed_headers"`
	ExposedHeaders   []string `yaml:"exposed_headers"`
	MaxAge           int      `yaml:"max_age"`
}

type CookieConfig struct {
	CookieSecure bool `yaml:"cookie_secure"`
}

type Config struct {
	Env       string         `yaml:"env"`
	Server    ServerConfig   `yaml:"server"`
	Postgres  PostgresConfig `yaml:"postgres"`
	Cors      CorsConfig     `yaml:"cors"`
	Cookie    CookieConfig   `yaml:"cookie"`
	JWTsecret string         `yaml:"jwt_secret"`
}

func LoadConfig() (*Config, error) {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		return nil, fmt.Errorf("config path not set")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	if v := os.Getenv("JWT_SECRET"); v != "" {
		cfg.JWTsecret = v
	}

	if v := os.Getenv("POSTGRES_HOST"); v != "" {
		cfg.Postgres.Host = v
	}
	if v := os.Getenv("POSTGRES_PORT"); v != "" {
		cfg.Postgres.Port = v
	}
	if v := os.Getenv("POSTGRES_USER"); v != "" {
		cfg.Postgres.Username = v
	}
	if v := os.Getenv("POSTGRES_PASSWORD"); v != "" {
		cfg.Postgres.Password = v
	}
	if v := os.Getenv("POSTGRES_DB"); v != "" {
		cfg.Postgres.DBname = v
	}
	if v := os.Getenv("POSTGRES_SSLMODE"); v != "" {
		cfg.Postgres.Sslmode = v
	}

	if v := os.Getenv("SERVER_HOST"); v != "" {
		cfg.Server.Host = v
	}
	if v := os.Getenv("SERVER_PORT"); v != "" {
		cfg.Server.Port = v
	}

	if v := os.Getenv("COOKIE_SECURE"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			cfg.Cookie.CookieSecure = b
		}
	}

	if v := os.Getenv("CORS_ALLOWED_ORIGINS"); v != "" {
		parts := strings.Split(v, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		cfg.Cors.AllowedOrigins = parts
	}

	if v := os.Getenv("CORS_EXPOSE_HEADERS"); v != "" {
		parts := strings.Split(v, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		cfg.Cors.ExposedHeaders = parts
	}

	if v := os.Getenv("CORS_ALLOW_CREDENTIALS"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			cfg.Cors.AllowCredentials = b
		}
	}

	if cfg.JWTsecret == "" {
		return nil, fmt.Errorf("JWT_SECRET not set")
	}
	if cfg.Postgres.Password == "" {
		return nil, fmt.Errorf("POSTGRES_PASSWORD not set")
	}

	return &cfg, nil
}
