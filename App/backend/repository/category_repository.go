package repository

import (
	"realico-comidas-backend/models"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	FindAll() ([]models.Category, error)
	FindByID(id uint) (*models.Category, error)
	Create(category *models.Category) error
	Update(category *models.Category) error
	Delete(id uint) error
	HasProducts(id uint) (bool, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) FindAll() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) FindByID(id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.Category{}, id).Error
}

func (r *categoryRepository) HasProducts(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Product{}).Where("category_id = ?", id).Count(&count).Error
	return count > 0, err
}
