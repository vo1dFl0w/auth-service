package logger

type Config struct {
	env          string
	loggerTimeFormat string
}

func NewLoggerConfig(env string, loggerTimeFormat string) *Config {
	return &Config{env: env, loggerTimeFormat: loggerTimeFormat}
}