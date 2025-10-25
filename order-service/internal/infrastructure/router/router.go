package router

import (
	"github.com/bekzxt/e-commerce/order-service/internal/interfaces/http"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, orderHandler http.OrderHandler) {
	orderRoutes := r.Group("/orders")
	{
		orderRoutes.POST("/", orderHandler.CreateOrder)
		orderRoutes.GET("/:id", orderHandler.GetOrder)
		orderRoutes.PATCH("/:id", orderHandler.UpdateOrderStatus)
		orderRoutes.GET("/", orderHandler.ListOrders)
	}
}
