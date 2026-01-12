package odoo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Servicio para integración con Odoo ERP usando JSON-RPC
type Client struct {
	// Configuración de conexión
	URL      string
	Database string
	Username string
	Password string
	APIKey   string // API Key para autenticación

	// Información del cliente (multi-tenant)
	ClientID   string
	ClientName string

	// Estado de autenticación
	UID        int
	httpClient *http.Client
}

// Inicializa el servicio con configuración de un cliente específico
func NewClient(config *Config) *Client {
	return &Client{
		URL:        config.URL,
		Database:   config.Database,
		Username:   config.Username,
		Password:   config.Password,
		APIKey:     config.APIKey,
		ClientID:   config.ClientID,
		ClientName: config.ClientName,
		httpClient: &http.Client{},
	}
}

// GetAuthPassword devuelve la contraseña o API Key para autenticación
// Con API Key, Odoo requiere usar la API Key como "password" en las llamadas
func (c *Client) GetAuthPassword() string {
	if c.APIKey != "" {
		return c.APIKey
	}
	return c.Password
}

// jsonRPCRequest representa una petición JSON-RPC a Odoo
type jsonRPCRequest struct {
	JSONRPC string                 `json:"jsonrpc"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
	ID      int                    `json:"id"`
}

// jsonRPCResponse representa una respuesta JSON-RPC de Odoo
type jsonRPCResponse struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Result  interface{}   `json:"result,omitempty"`
	Error   *jsonRPCError `json:"error,omitempty"`
}

// jsonRPCError representa un error en la respuesta JSON-RPC
type jsonRPCError struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// Authenticate autentica con Odoo y obtiene el UID
func (c *Client) Authenticate() error {
	fmt.Printf("🔑 Intentando autenticar con Odoo (Cliente: %s)...\n", c.ClientName)

	// Si tenemos API Key, usarla directamente (método preferido)
	if c.APIKey != "" {
		fmt.Println("🔐 Usando autenticación con API Key")
		// Con API Key, debemos hacer authenticate usando el username y API Key como password
		payload := jsonRPCRequest{
			JSONRPC: "2.0",
			Method:  "call",
			Params: map[string]interface{}{
				"service": "common",
				"method":  "authenticate",
				"args":    []interface{}{c.Database, c.Username, c.APIKey, map[string]interface{}{}},
			},
			ID: 1,
		}

		response, err := c.doRequest(payload)
		if err != nil {
			return fmt.Errorf("error en la petición de autenticación con API Key: %w", err)
		}

		if response.Error != nil {
			return fmt.Errorf("error de autenticación con API Key: %s - Verifica que la API Key sea válida", response.Error.Message)
		}

		// El resultado debe ser un número (UID)
		uid, ok := response.Result.(float64)
		if !ok || uid == 0 {
			return fmt.Errorf("API Key inválida o respuesta inesperada")
		}

		c.UID = int(uid)
		fmt.Printf("✅ Autenticado con API Key. UID: %d (Cliente: %s)\n", c.UID, c.ClientName)
		return nil
	}

	// Fallback a autenticación tradicional con usuario/contraseña
	fmt.Println("🔐 Usando autenticación con usuario/contraseña")
	payload := jsonRPCRequest{
		JSONRPC: "2.0",
		Method:  "call",
		Params: map[string]interface{}{
			"service": "common",
			"method":  "authenticate",
			"args":    []interface{}{c.Database, c.Username, c.Password, map[string]interface{}{}},
		},
		ID: 1,
	}

	response, err := c.doRequest(payload)
	if err != nil {
		return fmt.Errorf("error en la petición de autenticación: %w", err)
	}

	if response.Error != nil {
		return fmt.Errorf("error de autenticación: %s", response.Error.Message)
	}

	// El resultado debe ser un número (UID)
	uid, ok := response.Result.(float64)
	if !ok || uid == 0 {
		return fmt.Errorf("credenciales inválidas o respuesta inesperada")
	}

	c.UID = int(uid)
	fmt.Printf("✅ Autenticado correctamente. UID: %d (Cliente: %s)\n", c.UID, c.ClientName)

	return nil
}

// doRequest realiza una petición JSON-RPC a Odoo
func (c *Client) doRequest(payload jsonRPCRequest) (*jsonRPCResponse, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error al serializar la petición: %w", err)
	}

	req, err := http.NewRequest("POST", c.URL+"/jsonrpc", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error al crear la petición HTTP: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al realizar la petición HTTP: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer la respuesta: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error HTTP %d: %s", resp.StatusCode, string(body))
	}

	var response jsonRPCResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error al deserializar la respuesta: %w", err)
	}

	return &response, nil
}
