# Odoo-Quickpass Sync Middleware

Middleware de integración bilateral entre Odoo y Quickpass para la gestión de recursos humanos y asistencia.

## 📋 Descripción

Sistema de integración que permite la sincronización bidireccional de información entre:
- **Odoo**: Sistema ERP para gestión de empleados, contratos y liquidaciones
- **Quickpass**: Sistema de asistencia biométrica y portal de empleados

## 🎯 Funcionalidades

### Odoo → Quickpass
- Liquidaciones de empleados
- Creación de empleados + contratos
- Tiempos personales (vacaciones, licencias, justificaciones)

### Quickpass → Odoo
- Registro de asistencias con horas trabajadas
- Tiempos personales (vacaciones, licencias, justificaciones)

### Bilateral
- Consulta y gestión de tiempos personales
- Sincronización de datos de empleados

## 🏗️ Arquitectura

```
.
├── cmd/
│   └── server/              # Punto de entrada de la aplicación
├── internal/
│   ├── config/              # Configuración de la aplicación
│   ├── domain/              # Modelos y entidades del dominio
│   ├── handlers/            # Handlers HTTP
│   ├── middleware/          # Middleware HTTP (auth, logging, etc.)
│   ├── odoo/                # Cliente y servicios de Odoo
│   ├── quickpass/           # Cliente y servicios de Quickpass
│   ├── repository/          # Capa de persistencia
│   ├── server/              # Configuración del servidor HTTP
│   ├── services/            # Lógica de negocio
│   └── utils/               # Utilidades compartidas
├── pkg/
│   ├── logger/              # Logger personalizado
│   └── validator/           # Validadores
├── migrations/              # Migraciones de base de datos
├── docs/                    # Documentación
└── tests/                   # Tests de integración
```

## 🚀 Instalación

### Requisitos
- Go 1.21+
- PostgreSQL 13+ (opcional, para caché/logs)
- Acceso a APIs de Odoo y Quickpass

### Configuración

1. Clonar el repositorio:
```bash
git clone https://github.com/IamNewInThis/odoo-quickpass-sync.git
cd odoo-quickpass-sync
```

2. Copiar y configurar variables de entorno:
```bash
cp .env.example .env
# Editar .env con tus credenciales
```

3. Instalar dependencias:
```bash
go mod download
```

4. Ejecutar el servidor:
```bash
go run cmd/api/main.go
```

## 📝 Variables de Entorno

Ver archivo `.env.example` para la configuración completa.

## 🔌 Endpoints API

### Empleados
- `POST /api/v1/employees` - Crear empleado
- `GET /api/v1/employees/:id` - Obtener empleado
- `PUT /api/v1/employees/:id` - Actualizar empleado

### Liquidaciones
- `GET /api/v1/payrolls/:employee_id` - Obtener liquidaciones
- `POST /api/v1/payrolls/sync` - Sincronizar liquidaciones

### Asistencias
- `POST /api/v1/attendances` - Registrar asistencia
- `GET /api/v1/attendances/:employee_id` - Obtener asistencias

### Tiempos Personales
- `POST /api/v1/time-off` - Solicitar tiempo libre
- `GET /api/v1/time-off/:employee_id` - Consultar solicitudes
- `PUT /api/v1/time-off/:id` - Actualizar solicitud

### Webhooks
- `POST /webhooks/odoo` - Webhook de Odoo
- `POST /webhooks/quickpass` - Webhook de Quickpass

## 🧪 Testing

```bash
# Ejecutar todos los tests
go test ./...

# Tests con cobertura
go test -cover ./...

# Tests de integración
go test -tags=integration ./tests/...
```

## 📚 Documentación

Ver carpeta `docs/` para documentación detallada de:
- Arquitectura del sistema
- Flujos de integración
- Especificaciones de API
- Guías de desarrollo

## 🤝 Contribución

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📄 Licencia

Este proyecto es privado y propietario.

## 👥 Autores

- Bokato Team

## 🐛 Reporte de Issues

Para reportar bugs o solicitar features, por favor crea un issue en el repositorio.
