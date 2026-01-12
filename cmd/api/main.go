package main

import (
	"log"
	"os"

	"github.com/IamNewInThis/odoo-quickpass-sync/internal/config"
	"github.com/IamNewInThis/odoo-quickpass-sync/internal/odoo"
	"github.com/IamNewInThis/odoo-quickpass-sync/internal/server"
)

func main() {
	// Cargar variables de entorno
	if err := config.LoadEnv(); err != nil {
		log.Fatalf("❌ Error al cargar las variables de entorno: %v", err)
	}

	// Configurar cliente de Odoo
	odooConfig, err := odoo.NewConfigFromEnv()
	if err != nil {
		log.Printf("⚠️ Error configurando Odoo: %v", err)
		log.Println("ℹ️ El servidor iniciará sin conexión a Odoo")
	}

	var odooClient *odoo.Client
	if odooConfig != nil {
		odooClient = odoo.NewClient(odooConfig)

		// Intentar autenticar al inicio
		if err := odooClient.Authenticate(); err != nil {
			log.Printf("⚠️ Error autenticando con Odoo: %v", err)
			log.Println("ℹ️ El servidor iniciará, pero la conexión a Odoo no está disponible")
		}
	}

	// Obtener puerto del entorno o usar 8081 por defecto
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// Crear e iniciar servidor
	srv := server.NewServer(port, odooClient)

	log.Printf("🎯 Odoo Quickpass Service - Middleware Odoo/Quickpass")
	log.Printf("🌐 Escuchando en puerto %s", port)

	if err := srv.Start(); err != nil {
		log.Fatalf("❌ Error iniciando servidor: %v", err)
	}
}
