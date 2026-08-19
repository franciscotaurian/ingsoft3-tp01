package database

import (
	"log"

	"realico-comidas-backend/config"
	"realico-comidas-backend/models"
	"realico-comidas-backend/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.GetDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("Ejecutando AutoMigrate de tablas...")
	err = db.AutoMigrate(
		&models.Admin{},
		&models.Category{},
		&models.Product{},
		&models.Order{},
		&models.OrderItem{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// SeedAdmin garantiza que el usuario administrador exista en la tabla 'admins'
// con la contraseña almacenada como hash SHA-256.
// Si el admin ya existe, actualiza su contraseña al hash actual del .env
// (esto permite migrar de texto plano a hash sin intervención manual).
func SeedAdmin(db *gorm.DB, cfg *config.Config) {
	hashedPassword := utils.HashSHA256(cfg.AdminPassword)

	var admin models.Admin
	result := db.Where("username = ?", cfg.AdminUser).First(&admin)

	if result.Error != nil {
		// El admin no existe, crearlo con la contraseña hasheada
		newAdmin := &models.Admin{
			Username: cfg.AdminUser,
			Password: hashedPassword,
		}
		if err := db.Create(newAdmin).Error; err != nil {
			log.Printf("Advertencia: no se pudo crear el admin '%s': %v", cfg.AdminUser, err)
		} else {
			log.Printf("Usuario administrador '%s' creado con contraseña SHA-256.", cfg.AdminUser)
		}
	} else {
		// El admin ya existe: actualizar la contraseña al hash para asegurar
		// que nunca quede almacenada en texto plano.
		if err := db.Model(&admin).Update("password", hashedPassword).Error; err != nil {
			log.Printf("Advertencia: no se pudo actualizar la contraseña del admin '%s': %v", cfg.AdminUser, err)
		} else {
			log.Printf("Contraseña del administrador '%s' sincronizada como SHA-256.", cfg.AdminUser)
		}
	}
}
