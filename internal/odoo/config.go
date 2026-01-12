package odoo

import (
	"fmt"
	"os"
)

// Config contiene la configuración para conectarse a Odoo
type Config struct {
	// Configuración de conexión Odoo
	URL      string
	Database string
	Username string
	Password string
	APIKey   string // API Key de Odoo (recomendado sobre usuario/contraseña)

	// Información del cliente (multi-tenant)
	ClientID   string
	ClientName string
}

// NewConfigFromEnv crea una configuración desde variables de entorno
// Soporta autenticación con API Key (recomendado) o usuario/contraseña (legacy)
func NewConfigFromEnv() (*Config, error) {
	url := os.Getenv("ODOO_URL")
	if url == "" {
		return nil, fmt.Errorf("ODOO_URL no está configurado")
	}

	database := os.Getenv("ODOO_DATABASE")
	if database == "" {
		return nil, fmt.Errorf("ODOO_DATABASE no está configurado")
	}

	apiKey := os.Getenv("ODOO_API_KEY")
	username := os.Getenv("ODOO_USERNAME")
	password := os.Getenv("ODOO_PASSWORD")

	// Validar que tengamos al menos un método de autenticación
	if apiKey == "" && (username == "" || password == "") {
		return nil, fmt.Errorf("debe configurar ODOO_API_KEY o ODOO_USERNAME+ODOO_PASSWORD")
	}

	return &Config{
		URL:        url,
		Database:   database,
		Username:   username,
		Password:   password,
		APIKey:     apiKey,
		ClientID:   "default",
		ClientName: "Default Client",
	}, nil
}
