package service

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"realico-comidas-backend/models"
	"realico-comidas-backend/repository"
)

// phoneRegex acepta solo dígitos numéricos (sin espacios, guiones ni caracteres especiales).
var phoneRegex = regexp.MustCompile(`^\d+$`)

var (
	ErrCustomerDataRequired = errors.New("los datos del cliente (nombre, teléfono, dirección) son obligatorios")
	ErrOrderNoItems         = errors.New("el pedido debe contener al menos un producto")
	ErrInvalidStatusChange  = errors.New("transición de estado no permitida")
	ErrOrderNotFound        = errors.New("pedido no encontrado")
	ErrInvalidPhone         = errors.New("el teléfono debe contener solo números (sin espacios ni caracteres especiales)")
)

type CreateOrderItemDTO struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

type CreateOrderDTO struct {
	CustomerName    string               `json:"customer_name"`
	CustomerPhone   string               `json:"customer_phone"`
	CustomerAddress string               `json:"customer_address"`
	Items           []CreateOrderItemDTO `json:"items"`
}

type OrderService interface {
	GetAll() ([]models.Order, error)
	GetByID(id uint) (*models.Order, error)
	Create(dto CreateOrderDTO) (*models.Order, error)
	UpdateStatus(id uint, newStatus string) (*models.Order, error)
	GetMetrics() (*repository.MetricsData, error)
}

type orderService struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
}

func NewOrderService(orderRepo repository.OrderRepository, productRepo repository.ProductRepository) OrderService {
	return &orderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (s *orderService) GetAll() ([]models.Order, error) {
	return s.orderRepo.FindAll()
}

func (s *orderService) GetByID(id uint) (*models.Order, error) {
	return s.orderRepo.FindByID(id)
}

func (s *orderService) Create(dto CreateOrderDTO) (*models.Order, error) {
	// Validar que todos los datos del cliente estén presentes
	cleanName := strings.TrimSpace(dto.CustomerName)
	cleanPhone := strings.TrimSpace(dto.CustomerPhone)
	cleanAddress := strings.TrimSpace(dto.CustomerAddress)

	if cleanName == "" || cleanPhone == "" || cleanAddress == "" {
		return nil, ErrCustomerDataRequired
	}

	// Validar que el teléfono contenga solo dígitos numéricos
	if !phoneRegex.MatchString(cleanPhone) {
		return nil, ErrInvalidPhone
	}

	if len(dto.Items) == 0 {
		return nil, ErrOrderNoItems
	}

	var total float64
	orderItems := make([]models.OrderItem, 0, len(dto.Items))
	stockUpdates := make(map[uint]int)

	for _, itemDTO := range dto.Items {
		if itemDTO.Quantity <= 0 {
			return nil, fmt.Errorf("la cantidad del producto debe ser mayor a cero")
		}

		product, err := s.productRepo.FindByID(itemDTO.ProductID)
		if err != nil {
			return nil, fmt.Errorf("el producto con ID %d no existe", itemDTO.ProductID)
		}

		// Regla R4: Validar que la cantidad no supere el stock disponible
		if itemDTO.Quantity > product.Stock {
			return nil, fmt.Errorf("no hay suficiente stock para '%s' (disponible: %d, solicitado: %d)",
				product.Name, product.Stock, itemDTO.Quantity)
		}

		// Regla R8: Snapshot del precio unitario al momento del pedido
		unitPrice := product.Price
		itemTotal := unitPrice * float64(itemDTO.Quantity)
		total += itemTotal // Regla R5: Cálculo del total

		orderItems = append(orderItems, models.OrderItem{
			ProductID: product.ID,
			Quantity:  itemDTO.Quantity,
			UnitPrice: unitPrice,
		})

		// Regla R3: Calcular el nuevo stock a descontar
		stockUpdates[product.ID] = product.Stock - itemDTO.Quantity
	}

	order := &models.Order{
		CustomerName:    cleanName,
		CustomerPhone:   cleanPhone,
		CustomerAddress: cleanAddress,
		Total:           total,
		Status:          models.OrderStatusPendiente,
		Items:           orderItems,
	}

	err := s.orderRepo.CreateWithTx(order, stockUpdates)
	if err != nil {
		return nil, err
	}

	return s.orderRepo.FindByID(order.ID)
}

func (s *orderService) UpdateStatus(id uint, newStatus string) (*models.Order, error) {
	order, err := s.orderRepo.FindByID(id)
	if err != nil {
		return nil, ErrOrderNotFound
	}

	// Regla R6: Transiciones de estado permitidas:
	// pendiente -> confirmado
	// confirmado -> entregado
	currentStatus := order.Status
	isValid := false

	if currentStatus == models.OrderStatusPendiente && newStatus == models.OrderStatusConfirmado {
		isValid = true
	} else if currentStatus == models.OrderStatusConfirmado && newStatus == models.OrderStatusEntregado {
		isValid = true
	}

	if !isValid {
		return nil, fmt.Errorf("%w: no se puede cambiar de '%s' a '%s'", ErrInvalidStatusChange, currentStatus, newStatus)
	}

	err = s.orderRepo.UpdateStatus(id, newStatus)
	if err != nil {
		return nil, err
	}

	return s.orderRepo.FindByID(id)
}

func (s *orderService) GetMetrics() (*repository.MetricsData, error) {
	return s.orderRepo.GetMetrics()
}
