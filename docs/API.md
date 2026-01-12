# 📚 API Endpoints - Odoo Quickpass Sync

## 🚀 Servidor

Inicia el servidor:
```bash
go run cmd/api/main.go
```

El servidor estará disponible en: `http://localhost:8080`

---

## 📡 Endpoints Disponibles

### 1. Health Check
Verifica que el servidor esté funcionando.

**Request:**
```bash
GET http://localhost:8080/health
```

**Response:**
```json
{
  "status": "healthy",
  "time": "2026-01-12T15:30:00Z"
}
```

---

### 2. Estado de Odoo
Verifica la conexión con Odoo.

**Request:**
```bash
GET http://localhost:8080/odoo/status
```

**Response:**
```json
{
  "status": "connected",
  "client_name": "Default Client",
  "uid": 2,
  "database": "bokatocl-bokato-staging-27079827"
}
```

---

### 3. Obtener Todos los Empleados
Obtiene la lista completa de empleados desde Odoo.

**Request:**
```bash
GET http://localhost:8080/api/v1/employees
```

**Response:**
```json
{
  "success": true,
  "count": 25,
  "data": [
    {
      "id": 1,
      "identification_id": "12345678-9",
      "name": "Juan Pablo Pérez González",
      "first_name": "Juan",
      "surname": "Pablo",
      "second_surname": "Pérez",
      "nationality": {
        "id": 46,
        "code": "",
        "name": "Chile"
      },
      "work_email": "juan.perez@bokato.cl",
      "private_email": "juan@gmail.com",
      "work_phone": "+56912345678",
      "private_phone": "+56987654321",
      "private_address": {
        "street": "Av. Libertador Bernardo O'Higgins 123",
        "city": "Santiago",
        "state": "Región Metropolitana"
      },
      "commune": {
        "id": 1,
        "name": "Santiago"
      },
      "photo_url": "/web/image?model=hr.employee&id=1&field=image_1920",
      "birthday_parsed": "1990-05-15T00:00:00Z",
      "gender": "male"
    }
  ]
}
```

**Ejemplo con curl:**
```bash
curl -X GET http://localhost:8080/api/v1/employees
```

**Ejemplo con HTTPie:**
```bash
http GET http://localhost:8080/api/v1/employees
```

---

### 4. Obtener Empleado por ID
Obtiene información detallada de un empleado específico.

**Request:**
```bash
GET http://localhost:8080/api/v1/employees/{id}
```

**Parámetros:**
- `id` (path) - ID del empleado en Odoo

**Ejemplo:**
```bash
GET http://localhost:8080/api/v1/employees/1
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "identification_id": "12345678-9",
    "name": "Juan Pablo Pérez González",
    "first_name": "Juan",
    "surname": "Pablo",
    "second_surname": "Pérez",
    "nationality": {
      "id": 46,
      "code": "",
      "name": "Chile"
    },
    "work_email": "juan.perez@bokato.cl",
    "private_email": "juan@gmail.com",
    "work_phone": "+56912345678",
    "private_phone": "+56987654321",
    "private_address": {
      "street": "Av. Libertador Bernardo O'Higgins 123",
      "city": "Santiago",
      "state": "Región Metropolitana"
    },
    "commune": {
      "id": 1,
      "name": "Santiago"
    },
    "photo_url": "/web/image?model=hr.employee&id=1&field=image_1920",
    "birthday_parsed": "1990-05-15T00:00:00Z",
    "gender": "male"
  }
}
```

**Ejemplo con curl:**
```bash
curl -X GET http://localhost:8080/api/v1/employees/1
```

**Response de error (404):**
```json
{
  "error": "Empleado no encontrado: empleado no encontrado con ID: 999"
}
```

---

## 🧪 Probar con Postman

1. **Importar colección:**
   - Crea una nueva colección llamada "Odoo Quickpass Sync"
   
2. **Agregar requests:**
   - GET Health: `http://localhost:8080/health`
   - GET Odoo Status: `http://localhost:8080/odoo/status`
   - GET All Employees: `http://localhost:8080/api/v1/employees`
   - GET Employee by ID: `http://localhost:8080/api/v1/employees/1`

3. **Headers:**
   - No se requieren headers especiales (por ahora)

---

## 🔐 Autenticación (Próximamente)

En futuras versiones, los endpoints requerirán un API Key:

```bash
curl -X GET http://localhost:8080/api/v1/employees \
  -H "X-API-Key: your-api-key-here"
```

---

## 📊 Códigos de Estado HTTP

- `200 OK` - Solicitud exitosa
- `400 Bad Request` - Parámetros inválidos
- `404 Not Found` - Recurso no encontrado
- `405 Method Not Allowed` - Método HTTP no permitido
- `500 Internal Server Error` - Error del servidor
- `503 Service Unavailable` - Servicio no disponible (ej: Odoo desconectado)

---

## 🐛 Debugging

Para ver los logs del servidor:
```bash
# Los logs aparecerán en la consola
📥 GET /api/v1/employees
👥 Obteniendo todos los empleados de Odoo...
✅ Se obtuvieron 25 empleados
📤 GET /api/v1/employees - 1.234s
```

---

## 🚀 Próximos Endpoints

- `POST /api/v1/employees` - Crear empleado
- `PUT /api/v1/employees/{id}` - Actualizar empleado
- `GET /api/v1/payrolls/{employee_id}` - Obtener liquidaciones
- `POST /api/v1/attendances` - Registrar asistencia
- `GET /api/v1/time-off/{employee_id}` - Obtener solicitudes de tiempo libre
