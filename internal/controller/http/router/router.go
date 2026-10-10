package router

import (
	"github.com/gin-gonic/gin"
	"github.com/pinedadaniel/go-scaffolder-clean-template/internal/controller/http/handler"
	"github.com/pinedadaniel/go-scaffolder-clean-template/pkg/env"
)

func New(handlers handler.Handlers, scope env.Scope) *gin.Engine {
	ginMode := gin.ReleaseMode
	if scope.Is(env.Local) {
		ginMode = gin.DebugMode
	}
	gin.SetMode(ginMode)

	router := gin.New()

	router.Use(
		gin.Logger(),
		gin.Recovery(),
	)

	RegisterRoutes(router, handlers)

	return router
}
