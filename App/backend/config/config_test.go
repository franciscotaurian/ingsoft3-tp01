package config

import (
	"os"
	"testing"
)

func TestGetDSN(t *testing.T) {
	t.Run("Prioritizes DATABASE_URL when set", func(t *testing.T) {
		cfg := &Config{
			DatabaseURL: "postgres://user:pass@neon.tech/dbname?sslmode=require",
			DBHost:      "localhost",
			DBPort:      "5432",
			DBUser:      "postgres",
			DBPassword:  "postgres",
			DBName:      "catalogo",
		}
		expected := "postgres://user:pass@neon.tech/dbname?sslmode=require"
		if got := cfg.GetDSN(); got != expected {
			t.Errorf("GetDSN() = %v, expected %v", got, expected)
		}
	})

	t.Run("Falls back to individual fields when DATABASE_URL is empty", func(t *testing.T) {
		cfg := &Config{
			DatabaseURL: "",
			DBHost:      "localhost",
			DBPort:      "5432",
			DBUser:      "postgres",
			DBPassword:  "postgres",
			DBName:      "catalogo",
		}
		expected := "host=localhost user=postgres password=postgres dbname=catalogo port=5432 sslmode=disable"
		if got := cfg.GetDSN(); got != expected {
			t.Errorf("GetDSN() = %v, expected %v", got, expected)
		}
	})
}

func TestLoadConfig_DatabaseURL(t *testing.T) {
	testURL := "postgres://testuser:testpass@localhost:5432/testdb"
	os.Setenv("DATABASE_URL", testURL)
	defer os.Unsetenv("DATABASE_URL")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() returned unexpected error: %v", err)
	}

	if cfg.DatabaseURL != testURL {
		t.Errorf("cfg.DatabaseURL = %v, expected %v", cfg.DatabaseURL, testURL)
	}

	if cfg.GetDSN() != testURL {
		t.Errorf("cfg.GetDSN() = %v, expected %v", cfg.GetDSN(), testURL)
	}
}
