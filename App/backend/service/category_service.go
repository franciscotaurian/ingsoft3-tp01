package service

import (
	"errors"
	"strings"

	"realico-comidas-backend/models"
	"realico-comidas-backend/repository"
)

var (
	ErrCategoryNameRequired = errors.New("el nombre de la categoría no puede estar vacío")
	ErrCategoryHasProducts  = errors.New("no se puede eliminar la categoría porque tiene productos asociados")
	ErrCategoryNotFound     = errors.New("categoría no encontrada")
)

type CategoryService interface {
	GetAll() ([]models.Category, error)
	GetByID(id uint) (*models.Category, error)
	Create(name string) (*models.Category, error)
	Update(id uint, name string) (*models.Category, error)
	Delete(id uint) error
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) GetAll() ([]models.Category, error) {
	return s.repo.FindAll()
}

func (s *categoryService) GetByID(id uint) (*models.Category, error) {
	return s.repo.FindByID(id)
}

func (s *categoryService) Create(name string) (*models.Category, error) {
	cleanName := strings.TrimSpace(name)
	if cleanName == "" {
		return nil, ErrCategoryNameRequired
	}

	category := &models.Category{
		Name: cleanName,
	}

	err := s.repo.Create(category)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (s *categoryService) Update(id uint, name string) (*models.Category, error) {
	cleanName := strings.TrimSpace(name)
	if cleanName == "" {
		return nil, ErrCategoryNameRequired
	}

	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrCategoryNotFound
	}

	existing.Name = cleanName
	err = s.repo.Update(existing)
	if err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *categoryService) Delete(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return ErrCategoryNotFound
	}

	// Regla R7: No se puede eliminar una categoría con productos asociados
	hasProducts, err := s.repo.HasProducts(id)
	if err != nil {
		return err
	}
	if hasProducts {
		return ErrCategoryHasProducts
	}

	return s.repo.Delete(id)
}
