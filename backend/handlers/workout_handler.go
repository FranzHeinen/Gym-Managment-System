package handlers

import (
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/models"
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterWorkout(c *gin.Context) {
	userID := c.GetString("user_id")
	var req models.RegisterWorkoutRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.RegisterWorkout(req, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registrar el entrenamiento"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Entrenamiento registrado exitosamente"})
}

func GetWorkoutHistory(c *gin.Context) {
	userID := c.GetString("user_id")

	workouts, err := services.GetWorkoutHistory(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener el historial de entrenamientos"})
		return
	}

	c.JSON(http.StatusOK, workouts)
}
