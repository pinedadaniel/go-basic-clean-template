package health

import (
	"context"

	"github.com/pinedadaniel/go-scaffolder-clean-template/internal/core/domain"
)

type Repository interface {
	GetHealth(ctx context.Context, version string) (domain.Health, error)
}
