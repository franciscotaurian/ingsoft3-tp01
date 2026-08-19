package models

// Admin representa un usuario administrador del sistema.
// Las credenciales se almacenan en esta tabla y se validan
// en cada request autenticado contra el backend.
type Admin struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"unique;not null" json:"username"`
	Password string `gorm:"not null" json:"-"` // no se serializa en respuestas JSON
}
