package service

// category_service_test.go — Suite de tests para CategoryService.
//
// Cubre:
//   - Create: nombre vacío → error; nombre con espacios trimmeado → repo.Create llamado
//   - Update: nombre vacío → error
//   - Delete: R7 — no borrar con productos; happy path sin productos

import (
	"errors"
	"testing"

	"realico-comidas-backend/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ─────────────────────────────────────────────────────────────────────────────
// CategoryService.Create
// ─────────────────────────────────────────────────────────────────────────────

func TestCategoryCreate_ValidaNombreVacio(t *testing.T) {
	casos := []struct {
		nombre string
		input  string
	}{
		{"nombre_vacio", ""},
		{"nombre_solo_espacios", "   "},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			// Arrange — repo no se llama si la validación falla antes
			svc := NewCategoryService(&MockCategoryRepository{})

			// Act
			_, err := svc.Create(c.input)

			// Assert
			assert.ErrorIs(t, err, ErrCategoryNameRequired)
		})
	}
}

func TestCategoryCreate_GuardaConNombreTrimmeado(t *testing.T) {
	// Arrange
	mockRepo := &MockCategoryRepository{}
	// Verifica que Create fue llamado con el nombre sin espacios extra
	mockRepo.On("Create", mock.MatchedBy(func(cat *models.Category) bool {
		return cat.Name == "Pizzas"
	})).Return(nil)

	svc := NewCategoryService(mockRepo)

	// Act
	resultado, err := svc.Create("  Pizzas  ")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "Pizzas", resultado.Name)
	mockRepo.AssertExpectations(t)
}

// ─────────────────────────────────────────────────────────────────────────────
// CategoryService.Update
// ─────────────────────────────────────────────────────────────────────────────

func TestCategoryUpdate_ValidaNombreVacio(t *testing.T) {
	// Arrange — repo no se llama: la validación falla antes de FindByID
	svc := NewCategoryService(&MockCategoryRepository{})

	// Act
	_, err := svc.Update(1, "")

	// Assert
	assert.ErrorIs(t, err, ErrCategoryNameRequired)
}

func TestCategoryUpdate_RetornaErrorSiCategoriaNoExiste(t *testing.T) {
	// Arrange
	mockRepo := &MockCategoryRepository{}
	mockRepo.On("FindByID", uint(99)).Return(nil, errors.New("record not found"))

	svc := NewCategoryService(mockRepo)

	// Act
	_, err := svc.Update(99, "Sushi")

	// Assert
	assert.ErrorIs(t, err, ErrCategoryNotFound)
	mockRepo.AssertExpectations(t)
}

// ─────────────────────────────────────────────────────────────────────────────
// CategoryService.Delete — Regla R7
// ─────────────────────────────────────────────────────────────────────────────

func TestCategoryDelete_R7_NoPuedeEliminarConProductosAsociados(t *testing.T) {
	// Arrange
	mockRepo := &MockCategoryRepository{}
	mockRepo.On("FindByID", uint(1)).Return(&models.Category{ID: 1, Name: "Pizzas"}, nil)
	// La categoría tiene productos — R7 debe bloquear la eliminación
	mockRepo.On("HasProducts", uint(1)).Return(true, nil)

	svc := NewCategoryService(mockRepo)

	// Act
	err := svc.Delete(1)

	// Assert
	assert.ErrorIs(t, err, ErrCategoryHasProducts)
	mockRepo.AssertExpectations(t)
}

func TestCategoryDelete_PermiteEliminarSinProductos(t *testing.T) {
	// Arrange
	mockRepo := &MockCategoryRepository{}
	mockRepo.On("FindByID", uint(1)).Return(&models.Category{ID: 1, Name: "Pizzas"}, nil)
	mockRepo.On("HasProducts", uint(1)).Return(false, nil)
	mockRepo.On("Delete", uint(1)).Return(nil)

	svc := NewCategoryService(mockRepo)

	// Act
	err := svc.Delete(1)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCategoryDelete_RetornaErrorSiCategoriaNoExiste(t *testing.T) {
	// Arrange
	mockRepo := &MockCategoryRepository{}
	mockRepo.On("FindByID", uint(99)).Return(nil, errors.New("record not found"))

	svc := NewCategoryService(mockRepo)

	// Act
	err := svc.Delete(99)

	// Assert
	assert.ErrorIs(t, err, ErrCategoryNotFound)
	mockRepo.AssertExpectations(t)
}
