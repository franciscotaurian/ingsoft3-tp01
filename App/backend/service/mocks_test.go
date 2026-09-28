package service

// mocks_test.go — Implementaciones de test doubles para los repositorios.
// Todos los mocks usan testify/mock y satisfacen las interfaces definidas
// en el paquete repository.

import (
	"realico-comidas-backend/models"
	"realico-comidas-backend/repository"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// ─── MockOrderRepository ─────────────────────────────────────────────────────

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) FindAll() ([]models.Order, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Order), args.Error(1)
}

func (m *MockOrderRepository) FindByID(id uint) (*models.Order, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderRepository) CreateWithTx(order *models.Order, stockUpdates map[uint]int) error {
	args := m.Called(order, stockUpdates)
	return args.Error(0)
}

func (m *MockOrderRepository) UpdateStatus(id uint, status string) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *MockOrderRepository) UpdateTotal(id uint, newTotal float64) error {
	args := m.Called(id, newTotal)
	return args.Error(0)
}

func (m *MockOrderRepository) GetMetrics() (*repository.MetricsData, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.MetricsData), args.Error(1)
}

// ─── MockProductRepository ───────────────────────────────────────────────────

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) FindAll() ([]models.Product, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *MockProductRepository) FindAvailable() ([]models.Product, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *MockProductRepository) FindByID(id uint) (*models.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductRepository) Create(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) Update(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockProductRepository) UpdateStock(tx *gorm.DB, id uint, newStock int) error {
	args := m.Called(tx, id, newStock)
	return args.Error(0)
}

func (m *MockProductRepository) GetDB() *gorm.DB {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*gorm.DB)
}

// ─── MockCategoryRepository ──────────────────────────────────────────────────

type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) FindAll() ([]models.Category, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Category), args.Error(1)
}

func (m *MockCategoryRepository) FindByID(id uint) (*models.Category, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *MockCategoryRepository) Create(category *models.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockCategoryRepository) Update(category *models.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockCategoryRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockCategoryRepository) HasProducts(id uint) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}
