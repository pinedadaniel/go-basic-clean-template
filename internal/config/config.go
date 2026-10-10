package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/pinedadaniel/go-basic-clean-template/pkg/env"
	config "github.com/pinedadaniel/go-load-env/pkg/env"
)

type (
	Config struct {
		App     App
		HTTP    HTTP
		Log     Log
		Metrics Metrics
		Swagger Swagger
		Tracing Tracing
	}

	// App -.
	App struct {
		Name      string    `env:"APP_NAME,required"`
		Version   string    `env:"APP_VERSION,required"`
		Scope     env.Scope `env:"APP_SCOPE,required"`
		Component string    `env:"APP_COMPONENT"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env:"HTTP_PORT,required"`
	}

	// Log -.
	Log struct {
		Level  string `env:"LOG_LEVEL,required"`
		Format string `env:"LOG_FORMAT,required"`
	}

	// Metrics -.
	Metrics struct {
		Enabled bool `env:"METRICS_ENABLED" envDefault:"false"`
	}

	// Swagger -.
	Swagger struct {
		Enabled bool `env:"SWAGGER_ENABLED" envDefault:"true"`
	}

	// Tracing -.
	Tracing struct {
		Enabled bool `env:"TRACING_ENABLED" envDefault:"false"`
	}
)

// New returns config env.
func New() (*Config, error) {
	cfg := &Config{}

	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	if err := config.Parse(cfg); err != nil {
		return nil, fmt.Errorf("error parse config: %w", err)
	}

	if isValid, err := cfg.App.Scope.IsScopeValid(); !isValid {
		return nil, err
	}

	return cfg, nil
}
