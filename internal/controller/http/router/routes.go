package router

import (
	"github.com/gin-gonic/gin"
	"github.com/pinedadaniel/go-basic-clean-template/internal/controller/http/handler"
)

const (
	APIPrefix  = "/api"
	APIVersion = "/v1"

	HealthPath = "/health"
)

func RegisterRoutes(router gin.IRouter, handlers handler.Handlers) {
	api := router.Group(APIPrefix + APIVersion)

	api.GET(HealthPath, handlers.Health.Get)
}
