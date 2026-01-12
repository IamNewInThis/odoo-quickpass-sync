package odoo

import (
	"fmt"
	"strings"
	"time"
)

// Country representa un país en Odoo
type Country struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// Commune representa una comuna en Odoo
type Commune struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Address representa una dirección
type Address struct {
	Street string `json:"street"`
	City   string `json:"city"`
	State  string `json:"state"`
}

// HrEmployee representa un empleado en Odoo (modelo hr.employee)
type HrEmployee struct {
	ID               int           `json:"id"`
	IdentificationID string        `json:"identification_id"`
	Name             string        `json:"name"`           // Nombre completo
	FirstName        string        `json:"first_name"`     // Parseado de name
	Surname          string        `json:"surname"`        // Parseado de name
	SecondSurname    string        `json:"second_surname"` // Parseado de name
	CountryID        []interface{} `json:"country_id"`     // [id, name]
	Nationality      *Country      `json:"nationality"`    // Parseado
	WorkEmail        string        `json:"work_email"`
	PrivateEmail     string        `json:"private_email"`
	WorkPhone        string        `json:"work_phone"`
	PrivatePhone     string        `json:"private_phone"`
	PrivateStreet    string        `json:"private_street"`
	PrivateCity      string        `json:"private_city"`
	PrivateStateID   []interface{} `json:"private_state_id"` // [id, name]
	PrivateAddress   *Address      `json:"private_address"`  // Parseado
	HrCommuneID      []interface{} `json:"hr_commune"`       // [id, name]
	HrCommune        *Commune      `json:"commune"`          // Parseado
	Image1920        interface{}   `json:"image_1920"`       // Base64 o false
	PhotoURL         string        `json:"photo_url"`        // Parseado
	Birthday         interface{}   `json:"birthday"`         // Fecha como string o false
	BirthdayParsed   *time.Time    `json:"birthday_parsed"`  // Parseado
	Gender           string        `json:"gender"`
}

// EmployeeService proporciona operaciones para empleados
type EmployeeService struct {
	client *Client
}

// NewEmployeeService crea un nuevo servicio de empleados
func NewEmployeeService(client *Client) *EmployeeService {
	return &EmployeeService{
		client: client,
	}
}

// GetAllEmployees obtiene todos los empleados de Odoo
func (s *EmployeeService) GetAllEmployees() ([]*HrEmployee, error) {
	if s.client.UID == 0 {
		return nil, fmt.Errorf("cliente no autenticado")
	}

	fmt.Println("👥 Obteniendo todos los empleados de Odoo...")

	// Ejecutar método search_read en hr.employee
	payload := jsonRPCRequest{
		JSONRPC: "2.0",
		Method:  "call",
		Params: map[string]interface{}{
			"service": "object",
			"method":  "execute_kw",
			"args": []interface{}{
				s.client.Database,
				s.client.UID,
				s.client.GetAuthPassword(), // Usa API Key si está disponible
				"hr.employee",
				"search_read",
				[]interface{}{
					[]interface{}{}, // Sin filtros, obtener todos
				},
				map[string]interface{}{
					"fields": []string{
						"id",
						"identification_id",
						"name",
						"country_id",
						"work_email",
						"private_email",
						"work_phone",
						"private_phone",
						"private_street",
						"private_city",
						"private_state_id",
						"hr_commune",
						"image_1920",
						"birthday",
						"gender",
					},
				},
			},
		},
		ID: 1,
	}

	response, err := s.client.doRequest(payload)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo empleados: %w", err)
	}

	if response.Error != nil {
		return nil, fmt.Errorf("error de Odoo: %s", response.Error.Message)
	}

	// Procesar resultado
	resultSlice, ok := response.Result.([]interface{})
	if !ok {
		return nil, fmt.Errorf("formato de respuesta inválido")
	}

	employees := make([]*HrEmployee, 0, len(resultSlice))

	for _, item := range resultSlice {
		empData, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		employee := s.parseEmployeeData(empData)
		employees = append(employees, employee)
	}

	fmt.Printf("✅ Se obtuvieron %d empleados\n", len(employees))
	return employees, nil
}

// GetEmployeeByID obtiene un empleado específico por su ID
func (s *EmployeeService) GetEmployeeByID(employeeID int) (*HrEmployee, error) {
	if s.client.UID == 0 {
		return nil, fmt.Errorf("cliente no autenticado")
	}

	fmt.Printf("🔍 Buscando empleado ID: %d\n", employeeID)

	// Ejecutar método read en hr.employee
	payload := jsonRPCRequest{
		JSONRPC: "2.0",
		Method:  "call",
		Params: map[string]interface{}{
			"service": "object",
			"method":  "execute_kw",
			"args": []interface{}{
				s.client.Database,
				s.client.UID,
				s.client.GetAuthPassword(), // Usa API Key si está disponible
				"hr.employee",
				"read",
				[]interface{}{[]int{employeeID}},
				map[string]interface{}{
					"fields": []string{
						"id",
						"identification_id",
						"name",
						"country_id",
						"work_email",
						"private_email",
						"work_phone",
						"private_phone",
						"private_street",
						"private_city",
						"private_state_id",
						"hr_commune",
						"image_1920",
						"birthday",
						"gender",
					},
				},
			},
		},
		ID: 1,
	}

	response, err := s.client.doRequest(payload)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo empleado: %w", err)
	}

	if response.Error != nil {
		return nil, fmt.Errorf("error de Odoo: %s", response.Error.Message)
	}

	// Procesar resultado
	resultSlice, ok := response.Result.([]interface{})
	if !ok || len(resultSlice) == 0 {
		return nil, fmt.Errorf("empleado no encontrado con ID: %d", employeeID)
	}

	empData, ok := resultSlice[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("formato de respuesta inválido")
	}

	employee := s.parseEmployeeData(empData)
	fmt.Printf("✅ Empleado encontrado: %s\n", employee.Name)

	return employee, nil
}

// parseEmployeeData convierte los datos crudos de Odoo a estructura HrEmployee
func (s *EmployeeService) parseEmployeeData(data map[string]interface{}) *HrEmployee {
	employee := &HrEmployee{}

	// ID
	if id, ok := data["id"].(float64); ok {
		employee.ID = int(id)
	}

	// Identification ID
	if identID, ok := data["identification_id"].(string); ok {
		employee.IdentificationID = identID
	}

	// Name y división en partes
	if name, ok := data["name"].(string); ok {
		employee.Name = name
		parts := strings.Split(name, " ")
		if len(parts) > 0 {
			employee.FirstName = parts[0]
		}
		if len(parts) > 1 {
			employee.Surname = parts[1]
		}
		if len(parts) > 2 {
			employee.SecondSurname = parts[2]
		}
	}

	// Country/Nationality
	if countryID, ok := data["country_id"].([]interface{}); ok && len(countryID) == 2 {
		employee.CountryID = countryID
		if id, ok := countryID[0].(float64); ok {
			if name, ok := countryID[1].(string); ok {
				employee.Nationality = &Country{
					ID:   int(id),
					Name: name,
					Code: "", // Odoo no devuelve el código en read, habría que hacer otra llamada
				}
			}
		}
	}

	// Emails
	if email, ok := data["work_email"].(string); ok {
		employee.WorkEmail = email
	}
	if email, ok := data["private_email"].(string); ok {
		employee.PrivateEmail = email
	}

	// Phones
	if phone, ok := data["work_phone"].(string); ok {
		employee.WorkPhone = phone
	}
	if phone, ok := data["private_phone"].(string); ok {
		employee.PrivatePhone = phone
	}

	// Address
	if street, ok := data["private_street"].(string); ok {
		employee.PrivateStreet = street
	}
	if city, ok := data["private_city"].(string); ok {
		employee.PrivateCity = city
	}

	// State
	var stateName string
	if stateID, ok := data["private_state_id"].([]interface{}); ok && len(stateID) == 2 {
		employee.PrivateStateID = stateID
		if name, ok := stateID[1].(string); ok {
			stateName = name
		}
	}

	// Construir dirección completa
	if employee.PrivateStreet != "" {
		employee.PrivateAddress = &Address{
			Street: employee.PrivateStreet,
			City:   employee.PrivateCity,
			State:  stateName,
		}
	}

	// Commune
	if communeID, ok := data["hr_commune"].([]interface{}); ok && len(communeID) == 2 {
		employee.HrCommuneID = communeID
		if id, ok := communeID[0].(float64); ok {
			if name, ok := communeID[1].(string); ok {
				employee.HrCommune = &Commune{
					ID:   int(id),
					Name: name,
				}
			}
		}
	}

	// Photo URL
	employee.Image1920 = data["image_1920"]
	if data["image_1920"] != false && data["image_1920"] != nil {
		employee.PhotoURL = fmt.Sprintf("/web/image?model=hr.employee&id=%d&field=image_1920", employee.ID)
	}

	// Birthday
	employee.Birthday = data["birthday"]
	if birthdayStr, ok := data["birthday"].(string); ok && birthdayStr != "" {
		// Odoo devuelve fechas en formato YYYY-MM-DD
		if t, err := time.Parse("2006-01-02", birthdayStr); err == nil {
			employee.BirthdayParsed = &t
		}
	}

	// Gender
	if gender, ok := data["gender"].(string); ok {
		employee.Gender = gender
	}

	return employee
}
