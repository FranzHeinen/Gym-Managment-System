package handlers

import (
	"net/http"

	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/dtos"
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateRoutine(c *gin.Context) {
	var req dtos.CreateRutinaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Validaciones

	if req.Nombre == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El nombre es requerido para crear una rutina"})
		return
	}

	userID := c.GetString("user_id")
	RutinaResponse, err := services.CreateRoutine(req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, RutinaResponse)
}

func UpdateRoutine(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	userIDStr := userID.(string)
	rutinaID := c.Param("id")

	// Validaciones de negocio
	if rutinaID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se ingresó la rutina"})
		return
	}

	var req dtos.UpdateRutinaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rutinaResponse, err := services.UpdateRoutine(rutinaID, userIDStr, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rutinaResponse)

}

func GetRoutines(c *gin.Context) {

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	userIDStr := userID.(string)

	rutinas, err := services.GetRoutines(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener las rutinas"})
		return
	}

	c.JSON(http.StatusOK, rutinas)
}

func GetRoutineByID(c *gin.Context) {
	rutinaID := c.Param("id")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}
	userIDStr := userID.(string)

	rutina, err := services.GetRoutineByID(rutinaID, userIDStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Error al obtener la rutina"})
		return
	}

	c.JSON(http.StatusOK, rutina)
}

func DeleteRoutine(c *gin.Context) {
	rutinaID := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}
	userIDStr := userID.(string)

	err := services.DeleteRoutine(rutinaID, userIDStr)
	if err != nil {

		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Rutina no encontrada o no pertenece al usuario"})
		} else {
			// Para cualquier otro error, devolvemos 500
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar la rutina"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rutina eliminada exitosamente"})
}

func DuplicateRoutine(c *gin.Context) {
	rutinaID := c.Param("id")
	userID := c.GetString("user_id")

	response, err := services.DuplicateRoutine(rutinaID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}
