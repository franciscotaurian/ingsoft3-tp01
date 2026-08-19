package repository

import (
	"realico-comidas-backend/models"

	"gorm.io/gorm"
)

type MetricsData struct {
	TotalBilled       float64           `json:"total_billed"`
	OrdersByStatus    map[string]int64  `json:"orders_by_status"`
	MostSoldProduct   string            `json:"most_sold_product"`
	TotalProducts     int64             `json:"total_products"`
	InventoryValue    float64           `json:"inventory_value"`
}

type OrderRepository interface {
	FindAll() ([]models.Order, error)
	FindByID(id uint) (*models.Order, error)
	CreateWithTx(order *models.Order, stockUpdates map[uint]int) error
	UpdateStatus(id uint, status string) error
	GetMetrics() (*MetricsData, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) FindAll() ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Preload("Items.Product").Order("created_at desc").Find(&orders).Error
	return orders, err
}

func (r *orderRepository) FindByID(id uint) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("Items.Product").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) CreateWithTx(order *models.Order, stockUpdates map[uint]int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Guardar la orden y sus items
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// 2. Actualizar stocks descontados (R3)
		for productID, newStock := range stockUpdates {
			if err := tx.Model(&models.Product{}).Where("id = ?", productID).Update("stock", newStock).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *orderRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&models.Order{}).Where("id = ?", id).Update("status", status).Error
}

func (r *orderRepository) GetMetrics() (*MetricsData, error) {
	metrics := &MetricsData{
		OrdersByStatus: make(map[string]int64),
	}

	// 1. Total facturado
	var totalBilled float64
	r.db.Model(&models.Order{}).Select("COALESCE(SUM(total), 0)").Scan(&totalBilled)
	metrics.TotalBilled = totalBilled

	// 2. Pedidos por estado
	rows, err := r.db.Model(&models.Order{}).Select("status, COUNT(*) as count").Group("status").Rows()
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var status string
			var count int64
			if err := rows.Scan(&status, &count); err == nil {
				metrics.OrdersByStatus[status] = count
			}
		}
	}

	// 3. Producto más vendido
	type MostSold struct {
		Name string
	}
	var mostSold MostSold
	r.db.Table("order_items").
		Select("products.name as name").
		Joins("JOIN products ON products.id = order_items.product_id").
		Group("products.id, products.name").
		Order("SUM(order_items.quantity) DESC").
		Limit(1).
		Scan(&mostSold)
	metrics.MostSoldProduct = mostSold.Name
	if metrics.MostSoldProduct == "" {
		metrics.MostSoldProduct = "Ninguno"
	}

	// 4. Total productos
	var totalProducts int64
	r.db.Model(&models.Product{}).Count(&totalProducts)
	metrics.TotalProducts = totalProducts

	// 5. Valor del inventario (sum(price * stock))
	var inventoryValue float64
	r.db.Model(&models.Product{}).Select("COALESCE(SUM(price * stock), 0)").Scan(&inventoryValue)
	metrics.InventoryValue = inventoryValue

	return metrics, nil
}
