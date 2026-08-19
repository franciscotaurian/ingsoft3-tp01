package service

import (
	"errors"
	"strings"

	"realico-comidas-backend/models"
	"realico-comidas-backend/repository"
)

var (
	ErrProductNameRequired = errors.New("el nombre del producto no puede estar vacío")
	ErrProductInvalidPrice = errors.New("el precio debe ser mayor a 0")
	ErrProductInvalidStock = errors.New("el stock no puede ser negativo")
	ErrProductNotFound     = errors.New("producto no encontrado")
)

type CreateProductDTO struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	CategoryID  uint    `json:"category_id"`
	ImageURL    string  `json:"image_url"`
}

type UpdateProductDTO struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	CategoryID  uint    `json:"category_id"`
	ImageURL    string  `json:"image_url"`
}

type ProductService interface {
	GetAll() ([]models.Product, error)
	GetAvailable() ([]models.Product, error)
	GetByID(id uint) (*models.Product, error)
	Create(dto CreateProductDTO) (*models.Product, error)
	Update(id uint, dto UpdateProductDTO) (*models.Product, error)
	Delete(id uint) error
}

type productService struct {
	productRepo  repository.ProductRepository
	categoryRepo repository.CategoryRepository
}

func NewProductService(productRepo repository.ProductRepository, categoryRepo repository.CategoryRepository) ProductService {
	return &productService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *productService) GetAll() ([]models.Product, error) {
	return s.productRepo.FindAll()
}

func (s *productService) GetAvailable() ([]models.Product, error) {
	// Regla R1: Solo productos con stock > 0
	return s.productRepo.FindAvailable()
}

func (s *productService) GetByID(id uint) (*models.Product, error) {
	return s.productRepo.FindByID(id)
}

func (s *productService) Create(dto CreateProductDTO) (*models.Product, error) {
	cleanName := strings.TrimSpace(dto.Name)
	if cleanName == "" {
		return nil, ErrProductNameRequired
	}
	if dto.Price <= 0 {
		return nil, ErrProductInvalidPrice
	}
	if dto.Stock < 0 {
		return nil, ErrProductInvalidStock
	}

	// Validación: un producto debe tener una categoría asignada
	if dto.CategoryID == 0 {
		return nil, errors.New("debe asignarse una categoría al producto")
	}

	// Verificar que la categoría exista en la base de datos
	_, err := s.categoryRepo.FindByID(dto.CategoryID)
	if err != nil {
		return nil, errors.New("la categoría especificada no existe")
	}

	product := &models.Product{
		Name:        cleanName,
		Description: strings.TrimSpace(dto.Description),
		Price:       dto.Price,
		Stock:       dto.Stock,
		CategoryID:  dto.CategoryID,
		ImageURL:    strings.TrimSpace(dto.ImageURL),
	}

	err = s.productRepo.Create(product)
	if err != nil {
		return nil, err
	}

	return s.productRepo.FindByID(product.ID)
}

func (s *productService) Update(id uint, dto UpdateProductDTO) (*models.Product, error) {
	cleanName := strings.TrimSpace(dto.Name)
	if cleanName == "" {
		return nil, ErrProductNameRequired
	}
	if dto.Price <= 0 {
		return nil, ErrProductInvalidPrice
	}
	if dto.Stock < 0 {
		return nil, ErrProductInvalidStock
	}

	// Validación: un producto debe tener una categoría asignada
	if dto.CategoryID == 0 {
		return nil, errors.New("debe asignarse una categoría al producto")
	}

	existing, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, ErrProductNotFound
	}

	_, err = s.categoryRepo.FindByID(dto.CategoryID)
	if err != nil {
		return nil, errors.New("la categoría especificada no existe")
	}

	existing.Name = cleanName
	existing.Description = strings.TrimSpace(dto.Description)
	existing.Price = dto.Price
	existing.Stock = dto.Stock
	existing.CategoryID = dto.CategoryID
	existing.ImageURL = strings.TrimSpace(dto.ImageURL)

	err = s.productRepo.Update(existing)
	if err != nil {
		return nil, err
	}

	return s.productRepo.FindByID(id)
}

func (s *productService) Delete(id uint) error {
	_, err := s.productRepo.FindByID(id)
	if err != nil {
		return ErrProductNotFound
	}
	return s.productRepo.Delete(id)
}
