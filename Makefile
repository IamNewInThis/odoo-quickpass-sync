.PHONY: help build run test clean docker-build docker-run

# Variables
APP_NAME=odoo-quickpass-sync
BUILD_DIR=bin
MAIN_PATH=cmd/server/main.go

help: ## Muestra esta ayuda
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Compila la aplicación
	@echo "🔨 Compilando..."
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "✅ Compilación completada: $(BUILD_DIR)/$(APP_NAME)"

run: ## Ejecuta la aplicación
	@echo "🚀 Iniciando servidor..."
	@go run $(MAIN_PATH)

dev: ## Ejecuta en modo desarrollo con hot reload
	@echo "🔄 Modo desarrollo..."
	@air

test: ## Ejecuta los tests
	@echo "🧪 Ejecutando tests..."
	@go test -v ./...

test-coverage: ## Ejecuta tests con cobertura
	@echo "📊 Generando cobertura..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Reporte de cobertura generado: coverage.html"

clean: ## Limpia archivos generados
	@echo "🧹 Limpiando..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "✅ Limpieza completada"

deps: ## Instala dependencias
	@echo "📦 Instalando dependencias..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependencias instaladas"

lint: ## Ejecuta el linter
	@echo "🔍 Ejecutando linter..."
	@golangci-lint run ./...

docker-build: ## Construye la imagen Docker
	@echo "🐳 Construyendo imagen Docker..."
	@docker build -t $(APP_NAME):latest .
	@echo "✅ Imagen Docker creada"

docker-run: ## Ejecuta el contenedor Docker
	@echo "🐳 Ejecutando contenedor..."
	@docker run -p 8080:8080 --env-file .env $(APP_NAME):latest

migrate-up: ## Ejecuta migraciones
	@echo "⬆️  Ejecutando migraciones..."
	@go run cmd/migrate/main.go up

migrate-down: ## Revierte migraciones
	@echo "⬇️  Revirtiendo migraciones..."
	@go run cmd/migrate/main.go down

.DEFAULT_GOAL := help
