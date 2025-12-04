package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string `yaml:"env"`
	HttpAddr string `yaml:"http_addr"`
	DB       struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		DBname   string `yaml:"dbname"`
		Sslmode  string `yaml:"sslmode"`
	} `yaml:"db"`
	JWTsecret    string `yaml:"jwt_secret"`
	CookieSecure bool   `yaml:"cookie_secure"`
}

func LoadConfig() (*Config, error) {
	var cfg Config

	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		return nil, fmt.Errorf("config path not set")
	}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	return &cfg, nil
}
