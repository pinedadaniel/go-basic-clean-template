package handler

import "github.com/gin-gonic/gin"

type Health interface {
	Get(c *gin.Context)
}

type Handlers struct {
	Health Health
}
