package router

import (
	"github.com/gin-gonic/gin"
	"github.com/pinedadaniel/go-basic-clean-template/internal/controller/http/handler"
	"github.com/pinedadaniel/go-basic-clean-template/pkg/env"
)

func New(handlers handler.Handlers, scope env.Scope) *gin.Engine {
	router := gin.New()

	if isLocal := scope.Is(env.Local); !isLocal {
		gin.SetMode(gin.ReleaseMode)
	}

	router.Use(
		gin.Logger(),
		gin.Recovery(),
	)

	RegisterRoutes(router, handlers)

	return router
}
