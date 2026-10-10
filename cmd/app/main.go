package main

import (
	"os"

	"github.com/pinedadaniel/go-basic-clean-template/internal/app"
	"github.com/pinedadaniel/go-basic-clean-template/internal/config"
	"github.com/pinedadaniel/go-logger/pkg/log"
)

func main() {

	// Configuration
	cfg, err := config.New()

	if err != nil {
		log.Panicf("config error : %s", err)
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
		log.Errorf("application error: %s", err)
		os.Exit(1)
	}
}
