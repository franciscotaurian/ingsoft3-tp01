package service

// product_service_test.go — Suite de tests para ProductService.
//
// Cubre las siguientes reglas de negocio:
//   R-G: El nombre del producto no puede ser vacío ni solo espacios
//   R-H: El precio debe ser mayor a 0
//   R-I: El stock no puede ser negativo
//   R-J: El producto debe tener una categoría válida y existente
//   R-K: Creación y actualización de producto con persistencia correcta
//
// Técnicas del TP:
//   ✅ Parametrizado : TestCreate_ValidaProducto, TestUpdate_ValidaLosMismosCampos,
//                      TestCreate_ValidaCategoria_RechazaSiEsCeroOInexistente,
//                      TestUpdate_ValidaCategoria_RechazaSiEsCeroOInexistente
//   ✅ Caso de error  : TestUpdate_RetornaErrorSiProductoNoExiste, TestDelete_RetornaErrorSiProductoNoExiste
//   ✅ Mock           : TestCreate_Exitoso_GuardaProducto, TestUpdate_Exitoso

import (
	"errors"
	"testing"

	"realico-comidas-backend/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
// Test Create — Validación de categoría (R-J: cero o inexistente) (AAA)
// ─────────────────────────────────────────────────────────────────────────────

func TestCreate_ValidaCategoria_RechazaSiEsCeroOInexistente(t *testing.T) {
	t.Run("categoria_id_cero", func(t *testing.T) {
		// Arrange
		mockCatRepo := &MockCategoryRepository{}
		mockProdRepo := &MockProductRepository{}
		svc := NewProductService(mockProdRepo, mockCatRepo)
		dto := CreateProductDTO{Name: "Pizza", Price: 1000, Stock: 5, CategoryID: 0}

		// Act
		prod, err := svc.Create(dto)

		// Assert
		assert.Nil(t, prod)
		assert.Error(t, err)
		assert.Equal(t, "debe asignarse una categoría al producto", err.Error())
		mockCatRepo.AssertNotCalled(t, "FindByID")
		mockProdRepo.AssertNotCalled(t, "Create")
	})

	t.Run("categoria_inexistente_en_bd", func(t *testing.T) {
		// Arrange
		mockCatRepo := &MockCategoryRepository{}
		mockProdRepo := &MockProductRepository{}
		mockCatRepo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))
		svc := NewProductService(mockProdRepo, mockCatRepo)
		dto := CreateProductDTO{Name: "Pizza", Price: 1000, Stock: 5, CategoryID: 99}

		// Act
		prod, err := svc.Create(dto)

		// Assert
		assert.Nil(t, prod)
		assert.Error(t, err)
		assert.Equal(t, "la categoría especificada no existe", err.Error())
		mockCatRepo.AssertExpectations(t)
		mockProdRepo.AssertNotCalled(t, "Create")
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// Test Create — Éxito: guarda el producto y lo retorna (R-K) (AAA)
// ─────────────────────────────────────────────────────────────────────────────

func TestCreate_Exitoso_GuardaProducto(t *testing.T) {
	// Arrange
	mockCatRepo := &MockCategoryRepository{}
	mockProdRepo := &MockProductRepository{}

	dto := CreateProductDTO{
		Name:        "  Pizza Napolitana  ",
		Description: "  Con tomate y albahaca  ",
		Price:       1500,
		Stock:       20,
		CategoryID:  1,
		ImageURL:    "  http://img.com/pizza.jpg  ",
	}

	mockCatRepo.On("FindByID", uint(1)).Return(&models.Category{ID: 1, Name: "Pizzas"}, nil)

	// Verifica que se cree con datos limpios (strings con trim)
	mockProdRepo.On("Create", mock.MatchedBy(func(p *models.Product) bool {
		return p.Name == "Pizza Napolitana" &&
			p.Description == "Con tomate y albahaca" &&
			p.Price == 1500 &&
			p.Stock == 20 &&
			p.CategoryID == 1 &&
			p.ImageURL == "http://img.com/pizza.jpg"
	})).Return(nil)

	// FindByID posterior a Create para devolver el producto persistido
	mockProdRepo.On("FindByID", uint(0)).Return(&models.Product{
		ID:          10,
		Name:        "Pizza Napolitana",
		Description: "Con tomate y albahaca",
		Price:       1500,
		Stock:       20,
		CategoryID:  1,
	}, nil)

	svc := NewProductService(mockProdRepo, mockCatRepo)

	// Act
	resultado, err := svc.Create(dto)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resultado)
	assert.Equal(t, uint(10), resultado.ID)
	assert.Equal(t, "Pizza Napolitana", resultado.Name)
	mockCatRepo.AssertExpectations(t)
	mockProdRepo.AssertExpectations(t)
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
// Test Update — Validación de categoría (R-J: cero o inexistente) (AAA)
// ─────────────────────────────────────────────────────────────────────────────

func TestUpdate_ValidaCategoria_RechazaSiEsCeroOInexistente(t *testing.T) {
	t.Run("categoria_id_cero", func(t *testing.T) {
		// Arrange
		mockCatRepo := &MockCategoryRepository{}
		mockProdRepo := &MockProductRepository{}
		svc := NewProductService(mockProdRepo, mockCatRepo)
		dto := UpdateProductDTO{Name: "Pizza", Price: 1000, Stock: 5, CategoryID: 0}

		// Act
		prod, err := svc.Update(1, dto)

		// Assert
		assert.Nil(t, prod)
		assert.Error(t, err)
		assert.Equal(t, "debe asignarse una categoría al producto", err.Error())
		mockProdRepo.AssertNotCalled(t, "FindByID")
		mockCatRepo.AssertNotCalled(t, "FindByID")
	})

	t.Run("categoria_inexistente_en_bd", func(t *testing.T) {
		// Arrange
		mockCatRepo := &MockCategoryRepository{}
		mockProdRepo := &MockProductRepository{}

		// Producto existente con ID 1
		mockProdRepo.On("FindByID", uint(1)).Return(&models.Product{ID: 1, Name: "Pizza", CategoryID: 1}, nil)
		// Categoría 99 no existe
		mockCatRepo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))

		svc := NewProductService(mockProdRepo, mockCatRepo)
		dto := UpdateProductDTO{Name: "Pizza", Price: 1000, Stock: 5, CategoryID: 99}

		// Act
		prod, err := svc.Update(1, dto)

		// Assert
		assert.Nil(t, prod)
		assert.Error(t, err)
		assert.Equal(t, "la categoría especificada no existe", err.Error())
		mockProdRepo.AssertExpectations(t)
		mockCatRepo.AssertExpectations(t)
		mockProdRepo.AssertNotCalled(t, "Update")
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// Test Update — Éxito: actualiza campos y persiste en repositorio (R-K) (AAA)
// ─────────────────────────────────────────────────────────────────────────────

func TestUpdate_Exitoso(t *testing.T) {
	// Arrange
	mockCatRepo := &MockCategoryRepository{}
	mockProdRepo := &MockProductRepository{}

	productoExistente := &models.Product{
		ID:          5,
		Name:        "Pizza Comun",
		Description: "Simple",
		Price:       1000,
		Stock:       10,
		CategoryID:  1,
	}

	dto := UpdateProductDTO{
		Name:        "  Pizza Especial  ",
		Description: "  Con jamón y morrones  ",
		Price:       1800,
		Stock:       15,
		CategoryID:  2,
		ImageURL:    "  http://img.com/especial.jpg  ",
	}

	// 1. FindByID para buscar producto existente
	mockProdRepo.On("FindByID", uint(5)).Return(productoExistente, nil).Once()

	// 2. FindByID para verificar que la nueva categoría existe
	mockCatRepo.On("FindByID", uint(2)).Return(&models.Category{ID: 2, Name: "Pizzas Gourmet"}, nil)

	// 3. Update en BD con datos limpios (strings con trim)
	mockProdRepo.On("Update", mock.MatchedBy(func(p *models.Product) bool {
		return p.ID == 5 &&
			p.Name == "Pizza Especial" &&
			p.Description == "Con jamón y morrones" &&
			p.Price == 1800 &&
			p.Stock == 15 &&
			p.CategoryID == 2 &&
			p.ImageURL == "http://img.com/especial.jpg"
	})).Return(nil)

	// 4. FindByID final para retornar producto actualizado
	mockProdRepo.On("FindByID", uint(5)).Return(&models.Product{
		ID:          5,
		Name:        "Pizza Especial",
		Description: "Con jamón y morrones",
		Price:       1800,
		Stock:       15,
		CategoryID:  2,
		ImageURL:    "http://img.com/especial.jpg",
	}, nil).Once()

	svc := NewProductService(mockProdRepo, mockCatRepo)

	// Act
	resultado, err := svc.Update(5, dto)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resultado)
	assert.Equal(t, "Pizza Especial", resultado.Name)
	assert.Equal(t, float64(1800), resultado.Price)
	assert.Equal(t, uint(2), resultado.CategoryID)
	mockProdRepo.AssertExpectations(t)
	mockCatRepo.AssertExpectations(t)
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
