package config

import (
	"fmt"
	"path/filepath"

	"github.com/joho/godotenv"
)

// LoadEnvFile loads environment variables from a specified file
func LoadEnvFile(relativePath string) error {
	absPath, err := filepath.Abs(relativePath)
	if err != nil {
		return fmt.Errorf("error obtaining absolute path: %v", err)
	}

	if err := godotenv.Load(absPath); err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}

	return nil
}
