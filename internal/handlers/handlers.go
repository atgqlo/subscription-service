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
	ServiceName string  `json:"service_name" binding:"required"`
	Price       int     `json:"price" binding:"required,min=0"`
	UserID      string  `json:"user_id" binding:"required"`
	StartDate   string  `json:"start_date" binding:"required,len=7"`
	EndDate     *string `json:"end_date"`
}

// CreateHandler godoc
// @Summary      Создать подписку
// @Description  Новая подписка пользователя
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        request body     CreateRequest true "Данные подписки"
// @Success      201  {object} models.Subscription
// @Failure      400  {object} map[string]string
// @Router       /subscriptions [post]
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
		ctx.JSON(http.StatusCreated, gin.H{})
	}
}

// GetByIDHandler godoc
// @Summary      Получить подписку по ID
// @Description  Детали подписки
// @Tags         subscriptions
// @Produce      json
// @Param        id   path     string  true  "ID подписки"
// @Success      200  {object} models.Subscription
// @Failure      404  {object} map[string]string
// @Router       /subscriptions/{id} [get]
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

// ListHandler godoc
// @Summary      Все подписки
// @Description  Список подписок пользователя
// @Tags         subscriptions
// @Produce      json
// @Success      200  {array}  models.Subscription
// @Router       /subscriptions [get]
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

// UpdateHandler godoc
// @Summary      Обновить подписку
// @Description  Частичное обновление подписки по ID
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id      path     string     true        "ID подписки"
// @Param        request body     UpdateRequest true   "Данные подписки"
// @Success      200     {object} models.Subscription
// @Failure      400     {object} map[string]string
// @Failure      500     {object} map[string]string
// @Router       /subscriptions/{id} [put]
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
			ID:          id,
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

// DeleteHandler godoc
// @Summary      Удалить подписку
// @Description  Полное удаление подписки по ID
// @Tags         subscriptions
// @Param        id   path     string     true        "ID подписки"
// @Success      204  {string}  string     "No Content"
// @Failure      404  {object}  map[string]string  "Подписка не найдена"
// @Failure      500  {object}  map[string]string
// @Router       /subscriptions/{id} [delete]
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

type TotalCostResponse struct {
	TotalCost int `json:"total_cost" example:"498"`
}

// totalCostHandler godoc
// @Summary      Общая стоимость подписок за период
// @Description  Сумма подписок с фильтрацией по user_id, service_name, датам
// @Tags         subscriptions
// @Produce      json
// @Param        user_id      query    string  true   "User ID (UUID)"
// @Param        service_name query    string  false  "Название сервиса (Yandex Plus)"
// @Param        start_date   query    string  true   "Начало периода (MM-YYYY)"
// @Param        end_date     query    string  true   "Конец периода (MM-YYYY)"
// @Success      200          {object} TotalCostResponse
// @Failure      400          {object} map[string]string
// @Router       /subscriptions/total [get]
func (h *Handlers) TotalCostHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userIDStr := ctx.Query("user_id")
		if userIDStr == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
			return
		}

		startDate := ctx.Query("start_date")
		if startDate == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "start_date required (MM-YYYY)"})
			return
		}

		endDate := ctx.Query("end_date")
		if endDate == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "end_date required (MM-YYYY)"})
			return
		}

		userUUID, err := uuid.Parse(userIDStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id UUID"})
			return
		}

		serviceName := ctx.Query("service_name")
		var serviceNamePtr *string
		if serviceName != "" {
			serviceNamePtr = &serviceName
		}

		result, err := h.service.TotalCost(ctx.Request.Context(), userUUID, serviceNamePtr, startDate, endDate)
		if err != nil {
			h.log.Printf("total cost error: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, result)
	}
}
