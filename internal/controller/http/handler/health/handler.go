package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pinedadaniel/go-basic-clean-template/internal/controller/http/handler"
	"github.com/pinedadaniel/go-logger/pkg/log"
)

type Handler struct {
	healthUseCase UseCase
}

func New(useCase UseCase) *Handler {
	return &Handler{
		healthUseCase: useCase,
	}
}

func (h *Handler) Get(ctx *gin.Context) {

	var headers RequestHeaders

	if err := ctx.ShouldBindHeader(&headers); err != nil {
		log.Error("Invalid request headers", log.Err(err))
		if headers.XCallerId == "" {
			handleError(ctx, ErrCallerIDRequired)
			return
		}
		handleError(ctx, ErrInvalidRequestHeaders)
		return
	}

	health, err := h.healthUseCase.Execute(ctx.Request.Context())
	if err != nil {
		log.Error("Error HealthUseCase", log.Err(err))
		handleError(ctx, err)
		return
	}

	res := Response{
		Status:    string(health.Status),
		Timestamp: health.Timestamp,
		Version:   health.Version,
	}

	ctx.JSON(http.StatusOK, handler.Success(res))

}
