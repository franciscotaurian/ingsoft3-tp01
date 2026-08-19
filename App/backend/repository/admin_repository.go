package repository

import (
	"realico-comidas-backend/models"
	"realico-comidas-backend/utils"

	"gorm.io/gorm"
)

// AdminRepository define el contrato para acceso a datos de administradores.
type AdminRepository interface {
	// FindByCredentials busca un admin por usuario Y contraseña.
	// La contraseña recibida se hashea con SHA-256 antes de comparar con la base de datos.
	// Retorna error si no existe ninguna coincidencia.
	FindByCredentials(username, password string) (*models.Admin, error)

	// ExistsByUsername verifica si ya existe un admin con ese nombre de usuario.
	ExistsByUsername(username string) (bool, error)

	// Create inserta un nuevo administrador en la base de datos.
	Create(admin *models.Admin) error
}

type adminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepository{db: db}
}

// FindByCredentials hashea la contraseña recibida con SHA-256 y la compara
// contra el hash almacenado en la base de datos.
// Nunca compara ni almacena contraseñas en texto plano.
func (r *adminRepository) FindByCredentials(username, password string) (*models.Admin, error) {
	hashedPassword := utils.HashSHA256(password)

	var admin models.Admin
	err := r.db.Where("username = ? AND password = ?", username, hashedPassword).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *adminRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Admin{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

func (r *adminRepository) Create(admin *models.Admin) error {
	return r.db.Create(admin).Error
}
