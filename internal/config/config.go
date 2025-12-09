package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/go-playground/validator/v10"
)

type Config struct {
	AppEnv               string `env:"APP_ENV" envDefault:"development" validate:"oneof=development qa production"`
	Port                 int    `env:"PORT" envDefault:"8080"`
	DatabaseUrl          string `env:"DATABASE_URL" required:"true"`
	JwtSecret            string `env:"JWT_SECRET" required:"true"`
	JwtExpiredMin        int    `env:"JWT_EXPIRE_MIN" envDefault:"15"`
	JwtRefreshExpireDays int    `env:"JWT_REFRESH_EXPIRE_DAYS" envDefault:"15"`
	LogLevel             string `env:"LOG_LEVEL" envDefault:"info"`
	AllowOrigins         string `env:"ALLOW_ORIGINS" envDefault:"*"`
}

func CargarVariables() (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	v := validator.New()

	if err := v.Struct(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
