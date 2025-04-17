package router

import (
	"github.com/bekzxt/e-commerce/order-service/internal/interfaces/http"
	"github.com/gin-gonic/gin"
)

func SetupRouter(orderHandler http.OrderHandler) *gin.Engine {
	r := gin.Default()

	orderRoutes := r.Group("/orders")
	{
		orderRoutes.POST("/", orderHandler.CreateOrder)
		orderRoutes.GET("/:id", orderHandler.GetOrder)
		orderRoutes.PATCH("/:id", orderHandler.UpdateOrderStatus)
		orderRoutes.GET("/", orderHandler.ListOrders)
	}

	return r
}
