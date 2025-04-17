package main

import (
	"log"
	"net/http"
	"strconv"

	inventorypb "github.com/bekzxt/e-commerce/inventory-service/proto"
	orderpb "github.com/bekzxt/e-commerce/order-service/proto"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	orderConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to Order service: %v", err)
	}
	defer orderConn.Close()

	invConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to Inventory service: %v", err)
	}
	defer invConn.Close()

	orderClient := orderpb.NewOrderServiceClient(orderConn)
	invClient := inventorypb.NewInventoryServiceClient(invConn)

	r := gin.Default()

	r.POST("/orders", func(c *gin.Context) {
		req := &orderpb.CreateOrderRequest{}
		if err := c.ShouldBindJSON(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := orderClient.CreateOrder(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	})

	r.GET("/orders/:id", func(c *gin.Context) {
		orderID := c.Param("id")
		req := &orderpb.GetOrderRequest{
			OrderId: orderID,
		}

		resp, err := orderClient.GetOrderByID(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	})
	r.GET("/orders", func(c *gin.Context) {
		req := &orderpb.ListOrdersRequest{}
		if err := c.ShouldBindJSON(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		resp, err := orderClient.ListOrders(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	})
	r.PATCH("/order/:id", func(c *gin.Context) {
		orderID := c.Param("id")
		req := &orderpb.UpdateOrderRequest{OrderId: orderID}
		if err := c.ShouldBindJSON(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		resp, err := orderClient.UpdateOrder(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	})
	r.POST("/products", func(c *gin.Context) {

		req := &inventorypb.CreateProductRequest{}
		if err := c.ShouldBindJSON(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := invClient.CreateProduct(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	})
	r.GET("/products/:id", func(c *gin.Context) {
		productID := c.Param("id")
		id, err := strconv.ParseInt(productID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}
		req := &inventorypb.GetProductRequest{
			Id: id,
		}

		resp, err := invClient.GetProductByID(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	})
	r.GET("/products", func(c *gin.Context) {
		req := &inventorypb.ListProductsRequest{}
		if err := c.ShouldBindJSON(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		resp, err := invClient.ListProducts(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	})

	r.PATCH("/products/:id", func(c *gin.Context) {
		productID := c.Param("id")
		id, err := strconv.ParseInt(productID, 10, 64)
		req := &inventorypb.UpdateProductRequest{Id: id}
		if err := c.ShouldBindJSON(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		resp, err := invClient.UpdateProduct(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	})
	r.DELETE("/products/:id", func(c *gin.Context) {
		productID := c.Param("id")
		id, err := strconv.ParseInt(productID, 10, 64)
		req := &inventorypb.DeleteProductRequest{Id: id}
		if err := c.ShouldBindJSON(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		resp, err := invClient.DeleteProduct(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	})

	log.Println("API Gateway running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start API Gateway: %v", err)
	}
}
