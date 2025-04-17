package handler

import (
	"github.com/bekzxt/e-commerce/inventory-service/internal/application/usecase"
	"github.com/bekzxt/e-commerce/inventory-service/internal/domain"
	"github.com/bekzxt/e-commerce/inventory-service/pkg/dto"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ProductHandler struct {
	uc *usecase.ProductUseCase
}

func NewProductHandler(uc *usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{uc}
}

func (h *ProductHandler) RegisterRoutes(r *gin.Engine) {
	products := r.Group("/products")
	{
		products.POST("", h.Create)
		products.GET("/:id", h.GetByID)
		products.PATCH("/:id", h.Update)
		products.DELETE("/:id", h.Delete)
		products.GET("", h.List)
	}
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := &domain.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       int32(req.Stock),
		CategoryID:  int32(req.CategoryID),
	}

	if _, err := h.uc.CreateUC(product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	product, err := h.uc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.uc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Stock != nil {
		product.Stock = int32(*req.Stock)
	}
	if req.CategoryID != nil {
		product.CategoryID = int32(*req.CategoryID)
	}

	if _, err := h.uc.Update(product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.uc.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *ProductHandler) List(c *gin.Context) {
	products, err := h.uc.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}
