package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv carga las variables de entorno desde el archivo .env
func LoadEnv() error {
	// Intentar cargar .env, pero no fallar si no existe
	if err := godotenv.Load(); err != nil {
		// Si el archivo no existe, no es un error crítico
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("error al cargar archivo .env: %w", err)
	}
	return nil
}
