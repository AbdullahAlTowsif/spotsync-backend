package server

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	"spotsync-backend/internal/config"
)

type Server struct {
	echo *echo.Echo
	cfg  *config.Config

	// API is the /api/v1 group; domains register their routes on it.
	API *echo.Group
}

func New(cfg *config.Config) *Server {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.RequestLogger())

	e.GET("/health", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	return &Server{
		echo: e,
		cfg:  cfg,
		API:  e.Group("/api/v1"),
	}
}

// Run starts the server and shuts it down gracefully when ctx is cancelled.
func (s *Server) Run(ctx context.Context) error {
	sc := echo.StartConfig{
		Address:         ":" + s.cfg.Port,
		GracefulTimeout: 10 * time.Second,
	}
	return sc.Start(ctx, s.echo)
}
