package handler

import (
	"github.com/bekzxt/e-commerce/inventory-service/internal/application/usecase"
	"github.com/bekzxt/e-commerce/inventory-service/internal/domain"
	"github.com/bekzxt/e-commerce/inventory-service/pkg/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

type DiscountHandler struct {
	uc *usecase.DiscountUseCase
}

func NewDiscountHandler(uc *usecase.DiscountUseCase) *DiscountHandler {
	return &DiscountHandler{uc}
}

func (h *DiscountHandler) RegisterRoutes(r *gin.Engine) {
	discounts := r.Group("/discounts")
	{
		discounts.POST("", h.Create)
		discounts.DELETE("/:id", h.Delete)
		discounts.GET("", h.List)
	}
}

func (h *DiscountHandler) Create(c *gin.Context) {
	var req dto.Discount
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	discounts := &domain.Discount{
		ID:                 req.ID,
		Name:               req.Name,
		Description:        req.Description,
		DiscountPercentage: req.DiscountPercentage,
		ApplicableProducts: req.ApplicableProducts,
		IsActive:           req.IsActive,
		StartDate:          req.StartDate,
		EndDate:            req.EndDate,
	}

	if _, err := h.uc.Create(discounts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, discounts)
}
func (h *DiscountHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.uc.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *DiscountHandler) List(c *gin.Context) {
	products, err := h.uc.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}
