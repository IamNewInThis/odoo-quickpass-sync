package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/IamNewInThis/odoo-quickpass-sync/internal/odoo"
)

type Server struct {
	port       string
	odooClient *odoo.Client
	// quickpassClient *quickpass.Client
	httpServer *http.Server
}

func NewServer(port string, odooClient *odoo.Client) *Server {
	return &Server{
		port:       port,
		odooClient: odooClient,
		httpServer: &http.Server{
			Addr: fmt.Sprintf(":%s", port),
		},
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	// Rutas del servidor
	mux.HandleFunc("/", s.handleHome)
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/odoo/status", s.handleOdooStatus)

	// Rutas de empleados (API v1)
	mux.HandleFunc("/api/v1/employees", s.handleGetEmployees)
	mux.HandleFunc("/api/v1/employees/", s.handleGetEmployeeByID) // Con trailing slash para capturar /employees/{id}

	s.httpServer = &http.Server{
		Addr:         ":" + s.port,
		Handler:      s.loggingMiddleware(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("🚀 Servidor iniciado en http://localhost:%s\n", s.port)
	return s.httpServer.ListenAndServe()
}

// Middleware de logging
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("📥 %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("📤 %s %s - %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// handleHome maneja la ruta principal
func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"service": "Odoo Quickpass Sync Middleware",
		"version": "1.0.0",
		"status":  "running",
		"message": "Welcome to the Odoo Quickpass Sync Middleware Service",
	}
	s.sendJSON(w, http.StatusOK, response)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	}
	s.sendJSON(w, http.StatusOK, response)
}

// handleOdooStatus verifica la conexión con Odoo
func (s *Server) handleOdooStatus(w http.ResponseWriter, r *http.Request) {
	if s.odooClient == nil {
		s.sendJSON(w, http.StatusServiceUnavailable, map[string]interface{}{
			"status":  "error",
			"message": "Cliente Odoo no configurado",
		})
		return
	}

	// Verificar si ya está autenticado
	if s.odooClient.UID == 0 {
		err := s.odooClient.Authenticate()
		if err != nil {
			s.sendJSON(w, http.StatusServiceUnavailable, map[string]interface{}{
				"status":  "error",
				"message": fmt.Sprintf("Error autenticando con Odoo: %v", err),
			})
			return
		}
	}

	response := map[string]interface{}{
		"status":      "connected",
		"client_name": s.odooClient.ClientName,
		"uid":         s.odooClient.UID,
		"database":    s.odooClient.Database,
	}
	s.sendJSON(w, http.StatusOK, response)
}

// sendJSON envía una respuesta JSON
func (s *Server) sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("❌ Error encoding JSON: %v", err)
	}
}

// handleGetEmployees obtiene todos los empleados de Odoo
// GET /api/v1/employees
func (s *Server) handleGetEmployees(w http.ResponseWriter, r *http.Request) {
	// Solo permitir método GET
	if r.Method != http.MethodGet {
		s.sendJSON(w, http.StatusMethodNotAllowed, map[string]interface{}{
			"error": "Método no permitido. Use GET",
		})
		return
	}

	// Verificar cliente de Odoo
	if s.odooClient == nil {
		s.sendJSON(w, http.StatusServiceUnavailable, map[string]interface{}{
			"error": "Cliente Odoo no configurado",
		})
		return
	}

	// Autenticar si es necesario
	if s.odooClient.UID == 0 {
		if err := s.odooClient.Authenticate(); err != nil {
			s.sendJSON(w, http.StatusServiceUnavailable, map[string]interface{}{
				"error": fmt.Sprintf("Error autenticando con Odoo: %v", err),
			})
			return
		}
	}

	// Crear servicio de empleados
	employeeService := odoo.NewEmployeeService(s.odooClient)

	// Obtener todos los empleados
	employees, err := employeeService.GetAllEmployees()
	if err != nil {
		s.sendJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"error": fmt.Sprintf("Error obteniendo empleados: %v", err),
		})
		return
	}

	// Responder con los empleados
	s.sendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"count":   len(employees),
		"data":    employees,
	})
}

// handleGetEmployeeByID obtiene un empleado específico por ID
// GET /api/v1/employees/{id}
func (s *Server) handleGetEmployeeByID(w http.ResponseWriter, r *http.Request) {
	// Solo permitir método GET
	if r.Method != http.MethodGet {
		s.sendJSON(w, http.StatusMethodNotAllowed, map[string]interface{}{
			"error": "Método no permitido. Use GET",
		})
		return
	}

	// Extraer ID de la URL
	// /api/v1/employees/123 -> 123
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/employees/")
	if path == "" || path == r.URL.Path {
		// Si no hay ID, redirigir al listado completo
		s.handleGetEmployees(w, r)
		return
	}

	// Convertir ID a int
	employeeID, err := strconv.Atoi(path)
	if err != nil {
		s.sendJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error": "ID de empleado inválido",
		})
		return
	}

	// Verificar cliente de Odoo
	if s.odooClient == nil {
		s.sendJSON(w, http.StatusServiceUnavailable, map[string]interface{}{
			"error": "Cliente Odoo no configurado",
		})
		return
	}

	// Autenticar si es necesario
	if s.odooClient.UID == 0 {
		if err := s.odooClient.Authenticate(); err != nil {
			s.sendJSON(w, http.StatusServiceUnavailable, map[string]interface{}{
				"error": fmt.Sprintf("Error autenticando con Odoo: %v", err),
			})
			return
		}
	}

	// Crear servicio de empleados
	employeeService := odoo.NewEmployeeService(s.odooClient)

	// Obtener empleado por ID
	employee, err := employeeService.GetEmployeeByID(employeeID)
	if err != nil {
		s.sendJSON(w, http.StatusNotFound, map[string]interface{}{
			"error": fmt.Sprintf("Empleado no encontrado: %v", err),
		})
		return
	}

	// Responder con el empleado
	s.sendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    employee,
	})
}
