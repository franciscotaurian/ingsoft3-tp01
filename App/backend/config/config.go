package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	AdminUser      string
	AdminPassword  string
	WhatsAppNumber string
	Port           string
}

func LoadConfig() (*Config, error) {
	// Intentar cargar .env, ignora si no existe
	_ = godotenv.Load()

	cfg := &Config{
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "postgres"),
		DBName:         getEnv("DB_NAME", "catalogo"),
		AdminUser:      getEnv("ADMIN_USER", "admin"),
		AdminPassword:  getEnv("ADMIN_PASSWORD", "admin123"),
		WhatsAppNumber: getEnv("WHATSAPP_NUMBER", "5491112345678"),
		Port:           getEnv("PORT", "8080"),
	}

	return cfg, nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort)
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}
