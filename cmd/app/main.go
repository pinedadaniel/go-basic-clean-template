package main

import (
	"os"

	"github.com/pinedadaniel/go-logger/pkg/log"
	"github.com/pinedadaniel/go-scaffolder-clean-template/internal/app"
	"github.com/pinedadaniel/go-scaffolder-clean-template/internal/config"
	"github.com/pinedadaniel/go-scaffolder-clean-template/pkg/env"
)

const (
	msgConfigErr = "error config load: %s"
	msgAppErr    = "application error: %s"
)

func main() {

	// Configuration
	cfg, err := config.New()

	if err != nil {
		Exit(msgConfigErr, err, false)
	}

	// Global Logger
	log.Config(log.Options{
		Level:  log.ToLevel(cfg.Log.Level),
		Format: log.ToFormat(cfg.Log.Format),
	})

	log.Warn("Starting app",
		log.Field("name", cfg.App.Name),
		log.Field("version", cfg.App.Version),
		log.Field("scope", cfg.App.Scope),
	)

	if err := app.Run(cfg); err != nil {
		Exit(msgAppErr, err, cfg.App.Scope.Is(env.Local))
	}
}

func Exit(msg string, err error, withPanic bool) {
	if withPanic {
		//Stack trace
		log.Panicf(msg, err)
	} else {
		log.Errorf(msg, err)
		os.Exit(1)
	}
}
