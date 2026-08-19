package handler

import (
	"net/http"
	"strconv"

	"realico-comidas-backend/service"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) GetAvailable(c *gin.Context) {
	// Catálogo público: productos con stock > 0
	products, err := h.productService.GetAvailable()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) GetAll(c *gin.Context) {
	// Panel de admin: todos los productos
	products, err := h.productService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de producto inválido"})
		return
	}

	product, err := h.productService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) Create(c *gin.Context) {
	var dto service.CreateProductDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos del producto inválidos"})
		return
	}

	product, err := h.productService.Create(dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de producto inválido"})
		return
	}

	var dto service.UpdateProductDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos del producto inválidos"})
		return
	}

	product, err := h.productService.Update(uint(id), dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de producto inválido"})
		return
	}

	err = h.productService.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Producto eliminado correctamente"})
}
