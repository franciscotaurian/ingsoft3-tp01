package repository

import (
	"realico-comidas-backend/models"

	"gorm.io/gorm"
)

type ProductRepository interface {
	FindAll() ([]models.Product, error)
	FindAvailable() ([]models.Product, error)
	FindByID(id uint) (*models.Product, error)
	Create(product *models.Product) error
	Update(product *models.Product) error
	Delete(id uint) error
	UpdateStock(tx *gorm.DB, id uint, newStock int) error
	GetDB() *gorm.DB
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *productRepository) FindAll() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Category").Find(&products).Error
	return products, err
}

func (r *productRepository) FindAvailable() ([]models.Product, error) {
	var products []models.Product
	// R1: Stock > 0
	err := r.db.Preload("Category").Where("stock > ?", 0).Find(&products).Error
	return products, err
}

func (r *productRepository) FindByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Category").First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}

func (r *productRepository) UpdateStock(tx *gorm.DB, id uint, newStock int) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.Model(&models.Product{}).Where("id = ?", id).Update("stock", newStock).Error
}
