package handlers

import (
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/repositories"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// LogView es una struct limpia para pasar datos a la plantilla
type LogView struct {
	Timestamp string
	Status    int
	StatusCSS string
	Method    string
	Path      string
	Username  string
	ClientIP  string
}

// GetLogsPage renderiza la página que muestra los logs
func GetLogsPage(c *gin.Context) {
	// 1. Obtener los logs "sucios" como BSON
	logsBSON, err := repositories.GetLogs(100)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "admin_logs.html", gin.H{
			"TemplateName": "admin_logs",
			"Error":        "Error al cargar los logs",
		})
		return
	}

	// 2. Convertir los datos BSON "sucios" a una struct "limpia"
	var logsView []LogView = make([]LogView, 0)

	for _, log := range logsBSON {
		view := LogView{
			Method:   getString(log, "method"),
			Path:     getString(log, "path"),
			Username: getString(log, "username"),
			ClientIP: getString(log, "client_ip"),
		}

		// Convertir el timestamp de MongoDB (primitive.DateTime)
		if ts, ok := log["timestamp"].(primitive.DateTime); ok {
			view.Timestamp = ts.Time().In(time.Local).Format("2006/01/02 - 15:04:05")
		}

		// Convertir y clasificar el status code (puede ser int32, int64, float64)
		status := getInt(log, "status_code")
		view.Status = status

		if status > 499 {
			view.StatusCSS = "bg-danger"
		} else if status > 399 {
			view.StatusCSS = "bg-warning text-dark"
		} else if status > 299 {
			view.StatusCSS = "bg-secondary"
		} else if status > 199 {
			view.StatusCSS = "bg-success"
		} else {
			view.StatusCSS = "bg-info"
		}

		logsView = append(logsView, view)
	}

	// 3. Pasar los datos "limpios" a la plantilla
	c.HTML(http.StatusOK, "admin_logs.html", gin.H{
		"TemplateName": "admin_logs",
		"Logs":         logsView,
	})
}

// getString convierte de forma segura un campo de bson.M a string
func getString(log bson.M, key string) string {
	if val, ok := log[key].(string); ok {
		return val
	}
	return "" // Devuelve vacío si no es un string
}

// getInt convierte de forma segura varios tipos numéricos de bson.M a int
func getInt(log bson.M, key string) int {
	if val, ok := log[key].(int32); ok {
		return int(val)
	}
	if val, ok := log[key].(int64); ok {
		return int(val)
	}
	if val, ok := log[key].(float64); ok {
		return int(val)
	}
	if val, ok := log[key].(int); ok {
		return val
	}
	return 0 // Devuelve 0 si no es un número
}
