package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pinedadaniel/go-logger/pkg/log"
	"github.com/pinedadaniel/go-scaffolder-clean-template/internal/config"
	"github.com/pinedadaniel/go-scaffolder-clean-template/internal/config/profile"
	"github.com/pinedadaniel/go-scaffolder-clean-template/internal/controller/http"
	"github.com/pinedadaniel/go-scaffolder-clean-template/internal/controller/http/handler"
	"github.com/pinedadaniel/go-scaffolder-clean-template/internal/controller/http/handler/health"
	"github.com/pinedadaniel/go-scaffolder-clean-template/internal/controller/http/router"
	usecaseHealth "github.com/pinedadaniel/go-scaffolder-clean-template/internal/core/usecase/health"
	healthRepository "github.com/pinedadaniel/go-scaffolder-clean-template/internal/repository/health"
)

const (
	defaultReadTimeout     = 5 * time.Second
	defaultWriteTimeout    = 10 * time.Second
	defaultIdleTimeout     = 60 * time.Second
	defaultShutdownTimeout = 10 * time.Second
)

type repositories struct {
	health usecaseHealth.Repository
}

type useCases struct {
	health *usecaseHealth.UseCase
}

type dependencies struct {
	repositories repositories
	useCases     useCases
}

func Run(cfg *config.Config) error {
	scopeProfile, err := profile.New(cfg.ProfileDir, cfg.App.Scope)
	if err != nil {
		return err
	}

	handlers, err := buildHandlers(cfg, scopeProfile)

	if err != nil {
		return err
	}

	server := http.New(
		router.New(handlers, cfg.App.Scope),
		cfg.HTTP.Port,
		defaultReadTimeout,
		defaultWriteTimeout,
		defaultIdleTimeout,
		defaultShutdownTimeout,
	)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	serverErr := make(chan error, 1)

	go func() {
		log.Info("Starting HTTP server",
			log.Field("port", cfg.HTTP.Port),
		)

		serverErr <- server.Start()
	}()

	select {
	case err := <-serverErr:
		if err != nil {
			return fmt.Errorf("start HTTP server: %w", err)
		}
		return nil

	case <-ctx.Done():
		log.Info("Shutdown signal received")
	}

	if err := server.Shutdown(context.Background()); err != nil {
		return fmt.Errorf("shutdown HTTP server: %w", err)
	}

	return nil
}

func buildHandlers(cfg *config.Config, scopeProfile profile.Profile) (handler.Handlers, error) {

	deps, err := buildDependencies(scopeProfile, cfg)
	if err != nil {
		return handler.Handlers{}, err
	}

	healthHandler := health.New(deps.useCases.health)

	return handler.Handlers{
		Health: healthHandler,
	}, nil
}

func initializeRepositories(cfg profile.Profile) (repositories, error) {
	healthRepo, err := healthRepository.New()

	if err != nil {
		return repositories{}, fmt.Errorf("could not initialize health repository: %w", err)
	}

	return repositories{
		health: healthRepo,
	}, nil
}

func initializeUseCases(repos repositories, env *config.Config) useCases {
	return useCases{
		health: usecaseHealth.New(repos.health, env.App.Version),
	}
}

func buildDependencies(cfg profile.Profile, env *config.Config) (dependencies, error) {
	repos, err := initializeRepositories(cfg)
	if err != nil {
		return dependencies{}, fmt.Errorf("could not initialize repositories: %w", err)
	}

	uscs := initializeUseCases(repos, env)

	return dependencies{
		repositories: repos,
		useCases:     uscs,
	}, nil
}
