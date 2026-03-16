package handlers

import (
	"log"
	"net/http"
	"strconv"
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
	ServiceName string  `json:"service_name" binding:"required"`
	Price       int     `json:"price" binding:"required,min=0"`
	UserID      string  `json:"user_id" binding:"required"`
	StartDate   string  `json:"start_date" binding:"required,len=7"`
	EndDate     *string `json:"end_date"`
}

func (h *Handlers) CreateHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req CreateRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			h.log.Printf("bind error: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userUUID, err := uuid.Parse(req.UserID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id UUID"})
			return
		}
		sub := &models.Subscription{
			ServiceName: req.ServiceName,
			Price:       req.Price,
			UserID:      userUUID,
			StartDate:   req.StartDate,
			EndDate:     req.EndDate,
		}

		if err := h.service.Create(ctx.Request.Context(), sub); err != nil {
			h.log.Printf("create service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusCreated, sub)
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

type ListResponse struct {
	Subscriptions []models.Subscription `json:"subscriptions"`
	Pagination    struct {
		Total  int `json:"total"`
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
		Pages  int `json:"pages"`
	} `json:"pagination"`
}

func (h *Handlers) ListHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))
		offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

		list, total, err := h.service.List(ctx.Request.Context(), limit, offset)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := ListResponse{
			Subscriptions: make([]models.Subscription, 0, len(list)),
		}
		for _, sub := range list {
			resp.Subscriptions = append(resp.Subscriptions, *sub)
		}
		resp.Pagination.Total = total
		resp.Pagination.Limit = limit
		resp.Pagination.Offset = offset
		if limit > 0 {
			resp.Pagination.Pages = (total + limit - 1) / limit
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

type UpdateRequest struct {
	ServiceName *string `json:"service_name" binding:"required"`
	Price       *int    `json:"price" binding:"required,min=0"`
	StartDate   *string `json:"start_date" binding:"required,len=7"`
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
			ID: id,
		}
		if req.ServiceName != nil {
			updated.ServiceName = *req.ServiceName
		}
		if req.Price != nil {
			updated.Price = *req.Price
		}
		if req.StartDate != nil {
			updated.StartDate = *req.StartDate
		}
		if req.EndDate != nil {
			updated.EndDate = req.EndDate
		}
		if err := h.service.Update(ctx.Request.Context(), updated); err != nil {
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
		ctx.JSON(http.StatusNoContent, nil)
	}
}

type TotalCostRequest struct {
	UserID      string  `form:"user_id" binding:"required"`
	ServiceName *string `form:"service_name"`
	StartDate   string  `form:"start_date" binding:"required,len=7"`
	EndDate     *string `form:"end_date"`
}

func (h *Handlers) TotalCostHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req TotalCostRequest

		if err := ctx.ShouldBindQuery(&req); err != nil {
			h.log.Printf("bad request: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		h.log.Printf("DEBUG UserID=%q len=%d", req.UserID, len(req.UserID))

		userUUID, err := uuid.Parse(req.UserID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id UUID"})
			return
		}

		endDateValue := ""
		if req.EndDate != nil {
			endDateValue = *req.EndDate
		}

		result, err := h.service.TotalCost(ctx.Request.Context(), userUUID, req.ServiceName, req.StartDate, endDateValue)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"total_cost": result})
	}
}
