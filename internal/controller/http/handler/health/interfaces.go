package health

import (
	"context"

	"github.com/pinedadaniel/go-scaffolder-clean-template/internal/core/usecase/health"
)

type UseCase interface {
	Execute(ctx context.Context) (health.Output, error)
}
