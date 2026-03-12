package server

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"subscriptons-service/internal/handlers"
)

type Server struct {
	router *gin.Engine
	logger *log.Logger
}

func NewServer(h *handlers.Handlers, logger *log.Logger) *Server {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/subscriptions")
	{
		api.POST("", h.CreateHandler())
		api.GET("", h.ListHandler())
		api.GET("/total", h.TotalCostHandler())
		api.GET("/:id", h.GetByIDHandler())
		api.PUT("/:id", h.UpdateHandler())
		api.DELETE("/:id", h.DeleteHandler())
	}

	return &Server{
		router: r,
		logger: logger,
	}
}

func (s *Server) Run(port string) error {
	s.logger.Printf("server starting on %s", port)
	return s.router.Run(port)
}
