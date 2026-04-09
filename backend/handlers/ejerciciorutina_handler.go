package handlers

import (
	"net/http"

	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/dtos"
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddExerciseToRoutine(c *gin.Context) {

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}
	userIDStr := userID.(string)

	rutinaID := c.Param("id")

	var req dtos.EjercicioRutina
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.AddExerciseToRoutine(rutinaID, userIDStr, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func UpdateExerciseInRoutine(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}
	userIDStr := userID.(string)

	rutinaID := c.Param("id")
	ejercicioID := c.Param("ejercicioId")

	if rutinaID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se ingresó la rutina"})
		return
	}
	if ejercicioID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se ingresó el ejercicio"})
		return
	}
	var req dtos.EjercicioRutina
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.UpdateExerciseInRoutine(rutinaID, ejercicioID, userIDStr, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func GetExercisesFromRoutine(c *gin.Context) {

	rutinaID := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}
	userIDStr := userID.(string)

	ejercicios, err := services.GetExercisesFromRoutine(rutinaID, userIDStr)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Rutina no encontrada"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los ejercicios de la rutina"})
		}
		return
	}

	c.JSON(http.StatusOK, ejercicios)
}

func GetExerciseFromRoutine(c *gin.Context) {

	rutinaID := c.Param("id")
	ejercicioID := c.Param("ejercicioId")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}
	userIDStr := userID.(string)

	ejercicio, err := services.GetExerciseFromRoutine(rutinaID, ejercicioID, userIDStr)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Ejercicio no encontrado en esta rutina"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener el ejercicio"})
		}
		return
	}

	c.JSON(http.StatusOK, ejercicio)
}

func RemoveExerciseFromRoutine(c *gin.Context) {
	rutinaID := c.Param("id")
	ejercicioID := c.Param("ejercicioId")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}
	userIDStr := userID.(string)

	err := services.RemoveExerciseFromRoutine(rutinaID, ejercicioID, userIDStr)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Rutina o ejercicio no encontrado"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el ejercicio de la rutina"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ejercicio eliminado de la rutina exitosamente"})
}
