package models

import "time"

const (
	OrderStatusPendiente  = "pendiente"
	OrderStatusConfirmado = "confirmado"
	OrderStatusEntregado  = "entregado"
)

type Order struct {
	ID              uint        `gorm:"primaryKey" json:"id"`
	CustomerName    string      `gorm:"not null" json:"customer_name"`
	CustomerPhone   string      `gorm:"not null" json:"customer_phone"`
	CustomerAddress string      `gorm:"not null" json:"customer_address"`
	Total           float64     `gorm:"not null" json:"total"`
	Status          string      `gorm:"not null;default:'pendiente'" json:"status"`
	Items           []OrderItem `gorm:"foreignKey:OrderID" json:"items"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}
