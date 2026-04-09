package main

import (
	"log"

	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/database"
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/middleware"
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	if err := database.Connect(); err != nil {
		log.Fatal("Error conectando a MongoDB:", err)
	}
	defer database.Disconnect()

	// Configurar Gin
	r := gin.Default()
	r.Use(middleware.LoggerMiddleware())

	// Cargar templates HTML (la forma simple y original)
	r.LoadHTMLGlob("templates/*")

	// Servir archivos estáticos
	r.Static("/static", "./static")

	// Configurar rutas
	routes.SetupRoutes(r)

	// Iniciar servidor
	log.Println("Servidor iniciado en http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Error iniciando servidor:", err)
	}
}
