package http

import (
	"context"
	"errors"
	stdhttp "net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	server          *stdhttp.Server
	shutdownTimeout time.Duration
}

func New(
	router *gin.Engine,
	port string,
	readTimeout time.Duration,
	writeTimeout time.Duration,
	shutdownTimeout time.Duration,
) *Server {
	return &Server{
		server: &stdhttp.Server{
			Addr:         ":" + port,
			Handler:      router,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},
		shutdownTimeout: shutdownTimeout,
	}
}

func (s *Server) Start() error {
	err := s.server.ListenAndServe()
	if errors.Is(err, stdhttp.ErrServerClosed) {
		return nil
	}

	return err
}

func (s *Server) Shutdown(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(shutdownCtx)
}
