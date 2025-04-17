package http

import (
	"github.com/bekzxt/e-commerce/order-service/internal/application/usecase"
	"github.com/bekzxt/e-commerce/order-service/internal/domain"
	"github.com/bekzxt/e-commerce/order-service/internal/interfaces/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OrderHandler struct {
	usecase *usecase.OrderUseCase
}

func NewOrderHandler(uc *usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{usecase: uc}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}
	order, err := h.usecase.CreateOrder(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Could not create order: " + err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}
	c.JSON(http.StatusCreated, toOrderResponse(order))
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	orderID := c.Param("id")
	order, err := h.usecase.GetOrderByID(orderID)
	if err != nil {
		if err == domain.ErrOrderNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Message: err.Error(),
				Code:    http.StatusNotFound,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Error fetching order: " + err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, toOrderResponse(order))
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}
	newStatus := domain.OrderStatus(req.Status)
	if err := h.usecase.UpdateOrderStatus(orderID, newStatus); err != nil {
		if err == domain.ErrOrderNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Message: err.Error(),
				Code:    http.StatusNotFound,
			})
			return
		} else if err == domain.ErrInvalidStatus {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Could not update order: " + err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
}

func (h *OrderHandler) ListOrders(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "user_id query parameter is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	orders, err := h.usecase.ListOrdersByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Could not list orders: " + err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	var responses []dto.OrderResponse
	for _, order := range orders {
		responses = append(responses, toOrderResponse(order))
	}

	c.JSON(http.StatusOK, responses)
}

func toOrderResponse(o *domain.Order) dto.OrderResponse {
	var items []dto.OrderItemResponse
	for _, item := range o.Items {
		items = append(items, dto.OrderItemResponse{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
	}

	return dto.OrderResponse{
		ID:         o.ID,
		UserID:     o.UserID,
		Status:     string(o.Status),
		TotalPrice: o.TotalPrice,
		Items:      items,
	}
}
