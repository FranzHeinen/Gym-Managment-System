package middleware

import (
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/models"
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/repositories"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Registra información sobre cada solicitud recibida.
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next() // Procesa la petición
		latency := time.Since(startTime)

		userIDStr, exists := c.Get("user_id")
		userIdentifier := "Invitado"
		if exists {
			userIdentifier = userIDStr.(string)
		}

		// Imprimir en la consola (como antes)
		log.Printf("[LOG] | %3d | %13v | %15s | %-7s %s | Usuario: %s",
			c.Writer.Status(),
			latency,
			c.ClientIP(),
			c.Request.Method,
			c.Request.URL.Path,
			userIdentifier,
		)

		// Crear la entrada de Log para la base de datos
		logEntry := models.LogEntry{
			ID:         primitive.NewObjectID(),
			Timestamp:  startTime,
			StatusCode: c.Writer.Status(),
			Latency:    latency,
			ClientIP:   c.ClientIP(),
			Method:     c.Request.Method,
			Path:       c.Request.URL.Path,
		}

		objID, err := primitive.ObjectIDFromHex(userIdentifier)
		if err == nil {
			logEntry.UserID = objID
		} else {
			logEntry.UserString = userIdentifier
		}
		go saveLogToDB(logEntry)
	}
}

// saveLogToDB es una función helper para guardar el log sin bloquear la respuesta.
func saveLogToDB(entry models.LogEntry) {
	err := repositories.CreateLog(entry)
	if err != nil {
		log.Printf("[ERROR] No se pudo guardar el log en la base de datos: %v\n", err)
	}
}
