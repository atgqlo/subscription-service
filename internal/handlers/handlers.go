package handlers

import (
	"log"
	"net/http"
	"subscriptons-service/internal/models"
	"subscriptons-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handlers struct {
	service *service.SubscriptionService
	log     *log.Logger
}

func NewSubscriptionHandler(service *service.SubscriptionService, log *log.Logger) *Handlers {
	return &Handlers{
		service: service,
		log:     log,
	}
}

type CreateRequest struct {
	ServiceName string    `json:"service_name" binding:"required"`
	Price       int       `json:"price" binding:"required,min=0"`
	UserID      uuid.UUID `json:"user_id" binding:"required,uuid"`
	StartDate   string    `json:"start_date" binding:"required,len=7"`
	EndDate     *string   `json:"end_date"`
}

func (h *Handlers) CreateHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req CreateRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			h.log.Printf("bad request:%v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		created := &models.Subscription{
			ServiceName: req.ServiceName,
			Price:       req.Price,
			UserID:      req.UserID,
			StartDate:   req.StartDate,
			EndDate:     req.EndDate,
		}
		if err := h.service.Create(ctx.Request.Context(), created); err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		h.log.Printf("created subscription id=%s service=%s", created.ID, created.ServiceName)
		ctx.JSON(http.StatusCreated, created)

	}
}

func (h *Handlers) GetByIDHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")

		id, err := uuid.Parse(idStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
			h.log.Printf("invalid uuid:%v", err)
			return
		}
		sub, err := h.service.GetByID(ctx.Request.Context(), id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "sub not found"})
			return
		}
		h.log.Printf("getted subscription id=%s service=%s", sub.ID, sub.ServiceName)
		ctx.JSON(http.StatusOK, sub)
	}
}

func (h *Handlers) ListHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		list, err := h.service.List(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		h.log.Printf("listed %d subscriptions", len(list))
		ctx.JSON(http.StatusOK, list)
	}
}

type UpdateRequest struct {
	ServiceName string  `json:"service_name" binding:"required"`
	Price       int     `json:"price" binding:"required,min=0"`
	StartDate   string  `json:"start_date" binding:"required,len=7"`
	EndDate     *string `json:"end_date"`
}

func (h *Handlers) UpdateHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			h.log.Printf("invalid UUID:%v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var req UpdateRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			h.log.Printf("bad request:%v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		updated := &models.Subscription{
			ServiceName: req.ServiceName,
			Price:       req.Price,
			StartDate:   req.StartDate,
			EndDate:     req.EndDate,
		}
		if err := h.service.Update(ctx.Request.Context(), id, updated); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		h.log.Printf("updated subscription id=%s", id.String())
		ctx.JSON(http.StatusOK, updated)
	}
}
func (h *Handlers) DeleteHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")

		id, err := uuid.Parse(idStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
			h.log.Printf("invalid uuid: %v", err)
			return
		}

		if err := h.service.Delete(ctx.Request.Context(), id); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		h.log.Printf("delete subsription id: %s", id.String())
		ctx.JSON(http.StatusOK, nil)
	}
}

type TotalCostRequest struct {
	UserID      uuid.UUID `form:"user_id" binding:"required,uuid"`
	ServiceName *string   `form:"service_name"`
	StartDate   string    `form:"start_date" binding:"required,len=7"`
	EndDate     *string   `form:"end_date"`
}

func (h *Handlers) TotalCostHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req TotalCostRequest

		if err := ctx.ShouldBindQuery(&req); err != nil {
			h.log.Printf("bad request: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := h.service.TotalCost(ctx.Request.Context(), req.UserID, req.ServiceName, req.StartDate, *req.EndDate)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, result)
	}
}
