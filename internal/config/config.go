package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type ServerConfig struct {
	Env              string        `yaml:"env" env:"SERVER_ENV" env-required:"true"`
	Host             string        `yaml:"host" env:"SERVER_HTTP_HOST" env-required:"true"`
	Port             string        `yaml:"port" env:"SERVER_HTTP_PORT" env-required:"true"`
	LoggerTimeFormat string        `yaml:"logger_time_fromat"`
	RequestTimeout   time.Duration `yaml:"request_duration" env-default:"5s"`
}

type PostgresConfig struct {
	Host     string `yaml:"host" env:"POSTGRES_HOST" env-required:"true"`
	Port     string `yaml:"port" env:"POSTGRES_PORT" env-required:"true"`
	Username string `yaml:"username" env:"POSTGRES_USER" env-required:"true"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD" env-required:"true"`
	DBname   string `yaml:"dbname" env:"POSTGRES_DB" env-required:"true"`
	Sslmode  string `yaml:"sslmode" env:"POSTGRES_SSLMODE" env-required:"true"`
}

type CorsConfig struct {
	AllowedOrigins   []string `yaml:"allowed_origins" env:"CORS_ALLOWED_ORIGINS" env-required:"true"`
	AllowCredentials bool     `yaml:"allowed_credentials" env:"CORS_ALLOW_CREDENTIALS" env-default:"true"`
	AllowedMethods   []string `yaml:"allowed_methods" env:"CORS_ALLOWED_METHODS" env-required:"true"`
	AllowedHeaders   []string `yaml:"allowed_headers" env:"CORS_ALLOWED_HEADERS" env-required:"true"`
	ExposedHeaders   []string `yaml:"exposed_headers" env:"CORS_EXPOSE_HEADERS" env-required:"true"`
	MaxAge           int      `yaml:"max_age" env:"CORS_MAX_AGE" env-required:"true"`
}

type CookieConfig struct {
	MaxAge       int  `yaml:"cookie_max_age" env-default:"300"`
	CookieSecure bool `yaml:"cookie_secure" env:"COOKIE_SECURE" env-required:"true"`
}

type GoogleOAuthConfig struct {
	ClientID        string `yaml:"client_id" env:"GOOGLE_OAUTH_CLIENT_ID" env-required:"true"`
	ClientSecret    string `yaml:"client_secret" env:"GOOGLE_OAUTH_CLIENT_SECRET" env-required:"true"`
	RedirectURL     string `yaml:"redirect_url" env:"GOOGLE_OAUTH_REDIRECT_URL" env-required:"true"`
	BaseUserInfoURL string `yaml:"base_user_info_url" env:"GOOGLE_OAUTH_BASE_USER_INFO_URL" env-required:"true"`
}

type Config struct {
	Server      ServerConfig      `yaml:"server"`
	Postgres    PostgresConfig    `yaml:"postgres"`
	Cors        CorsConfig        `yaml:"cors"`
	Cookie      CookieConfig      `yaml:"cookie"`
	GoogleOAuth GoogleOAuthConfig `yaml:"google_oauth"`
	JWTsecret   string            `yaml:"jwt_secret" env:"JWT_SECRET" env-required:"true"`
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

	/*
		if v := os.Getenv("ENV"); v != "" {
			cfg.Env = v
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
	*/

	if v := os.Getenv("COOKIE_SECURE"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			cfg.Cookie.CookieSecure = b
		}
	}

	if v := os.Getenv("CORS_ALLOWED_HEADERS"); v != "" {
		parts := strings.Split(v, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		cfg.Cors.AllowedHeaders = parts
	}

	if v := os.Getenv("CORS_ALLOWED_METHODS"); v != "" {
		parts := strings.Split(v, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		cfg.Cors.AllowedMethods = parts
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

	/*
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
	*/

	return &cfg, nil
}
