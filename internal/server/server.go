package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"subscriptons-service/internal/handlers"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	logger *log.Logger
	srv    *http.Server
}

func NewServer(h *handlers.Handlers, logger *log.Logger) *Server {
	r := gin.Default()
	r.Use(
		handlers.LoggerMiddleware(logger),
		handlers.RecoveryMiddleware(),
	)
	api := r.Group("/subscriptions")
	{
		api.POST("", h.CreateHandler())
		api.GET("", h.ListHandler())
		api.GET("/total", h.TotalCostHandler())
		api.GET("/:id", h.GetByIDHandler())
		api.PUT("/:id", h.UpdateHandler())
		api.DELETE("/:id", h.DeleteHandler())
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	return &Server{
		router: r,
		logger: logger,
		srv:    srv,
	}
}

func (s *Server) Run(addr string) error {
	s.logger.Printf("server starting on: %s ", addr)
	s.srv.Addr = addr

	go func() {
		if err := s.srv.ListenAndServe(); err != nil {
			s.logger.Printf("listed : %v", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	s.logger.Println("shtting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		s.logger.Fatal("server forced to shutdown: %w", err)
		return err
	}
	s.logger.Println("server stopped")
	return nil
}
