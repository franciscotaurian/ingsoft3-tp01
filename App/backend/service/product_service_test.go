package service

// product_service_test.go — Suite de tests para ProductService.
//
// Cubre las siguientes reglas de negocio:
//   R-G: El nombre del producto no puede ser vacío ni solo espacios
//   R-H: El precio debe ser mayor a 0
//   R-I: El stock no puede ser negativo
//
// Técnicas del TP:
//   ✅ Parametrizado : TestCreate_ValidaProducto, TestUpdate_ValidaLosMismosCampos
//   ✅ Caso de error  : TestUpdate_RetornaErrorSiProductoNoExiste, TestDelete_RetornaErrorSiProductoNoExiste

import (
	"errors"
	"testing"

	"realico-comidas-backend/models"

	"github.com/stretchr/testify/assert"
)

// ─────────────────────────────────────────────────────────────────────────────
// Test 8 — Parametrizado: validación de campos del producto en Create (R-G, R-H, R-I)
// ─────────────────────────────────────────────────────────────────────────────

func TestCreate_ValidaProducto(t *testing.T) {
	casos := []struct {
		nombre        string
		dto           CreateProductDTO
		errorEsperado error
	}{
		{
			nombre:        "nombre_vacio",
			dto:           CreateProductDTO{Name: "", Price: 100, Stock: 5, CategoryID: 1},
			errorEsperado: ErrProductNameRequired,
		},
		{
			nombre:        "nombre_solo_espacios",
			dto:           CreateProductDTO{Name: "   ", Price: 100, Stock: 5, CategoryID: 1},
			errorEsperado: ErrProductNameRequired,
		},
		{
			nombre:        "precio_cero",
			dto:           CreateProductDTO{Name: "Empanada", Price: 0, Stock: 5, CategoryID: 1},
			errorEsperado: ErrProductInvalidPrice,
		},
		{
			nombre:        "precio_negativo",
			dto:           CreateProductDTO{Name: "Empanada", Price: -1.5, Stock: 5, CategoryID: 1},
			errorEsperado: ErrProductInvalidPrice,
		},
		{
			nombre:        "stock_negativo",
			dto:           CreateProductDTO{Name: "Empanada", Price: 100, Stock: -1, CategoryID: 1},
			errorEsperado: ErrProductInvalidStock,
		},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			svc := NewProductService(&MockProductRepository{}, &MockCategoryRepository{})
			_, err := svc.Create(c.dto)
			assert.ErrorIs(t, err, c.errorEsperado)
		})
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// ProductService.Update — las mismas validaciones que Create (R-G, R-H, R-I)
// ─────────────────────────────────────────────────────────────────────────────

func TestUpdate_ValidaLosMismosCampos(t *testing.T) {
	casos := []struct {
		nombre        string
		dto           UpdateProductDTO
		errorEsperado error
	}{
		{
			nombre:        "nombre_vacio",
			dto:           UpdateProductDTO{Name: "", Price: 100, Stock: 5, CategoryID: 1},
			errorEsperado: ErrProductNameRequired,
		},
		{
			nombre:        "nombre_solo_espacios",
			dto:           UpdateProductDTO{Name: "   ", Price: 100, Stock: 5, CategoryID: 1},
			errorEsperado: ErrProductNameRequired,
		},
		{
			nombre:        "precio_cero",
			dto:           UpdateProductDTO{Name: "Empanada", Price: 0, Stock: 5, CategoryID: 1},
			errorEsperado: ErrProductInvalidPrice,
		},
		{
			nombre:        "precio_negativo",
			dto:           UpdateProductDTO{Name: "Empanada", Price: -10, Stock: 5, CategoryID: 1},
			errorEsperado: ErrProductInvalidPrice,
		},
		{
			nombre:        "stock_negativo",
			dto:           UpdateProductDTO{Name: "Empanada", Price: 100, Stock: -1, CategoryID: 1},
			errorEsperado: ErrProductInvalidStock,
		},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			// Arrange — la validación falla antes de llamar a FindByID
			svc := NewProductService(&MockProductRepository{}, &MockCategoryRepository{})

			// Act
			_, err := svc.Update(1, c.dto)

			// Assert
			assert.ErrorIs(t, err, c.errorEsperado)
		})
	}
}

func TestUpdate_RetornaErrorSiProductoNoExiste(t *testing.T) {
	// Arrange
	mockProductRepo := &MockProductRepository{}
	mockProductRepo.On("FindByID", uint(99)).Return(nil, errors.New("record not found"))

	svc := NewProductService(mockProductRepo, &MockCategoryRepository{})

	// Act
	_, err := svc.Update(99, UpdateProductDTO{Name: "Empanada", Price: 100, Stock: 5, CategoryID: 1})

	// Assert
	assert.ErrorIs(t, err, ErrProductNotFound)
	mockProductRepo.AssertExpectations(t)
}

// ─────────────────────────────────────────────────────────────────────────────
// ProductService.Delete
// ─────────────────────────────────────────────────────────────────────────────

func TestDelete_RetornaErrorSiProductoNoExiste(t *testing.T) {
	// Arrange
	mockProductRepo := &MockProductRepository{}
	mockProductRepo.On("FindByID", uint(99)).Return(nil, errors.New("record not found"))

	svc := NewProductService(mockProductRepo, &MockCategoryRepository{})

	// Act
	err := svc.Delete(99)

	// Assert
	assert.ErrorIs(t, err, ErrProductNotFound)
	mockProductRepo.AssertExpectations(t)
}

func TestDelete_EliminaProductoExistente(t *testing.T) {
	// Arrange
	mockProductRepo := &MockProductRepository{}
	mockProductRepo.On("FindByID", uint(1)).Return(&models.Product{ID: 1, Name: "Empanada"}, nil)
	mockProductRepo.On("Delete", uint(1)).Return(nil)

	svc := NewProductService(mockProductRepo, &MockCategoryRepository{})

	// Act
	err := svc.Delete(1)

	// Assert
	assert.NoError(t, err)
	mockProductRepo.AssertExpectations(t)
}
