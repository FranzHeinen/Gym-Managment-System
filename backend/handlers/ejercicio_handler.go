package handlers

import (
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/models"
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Devuelve todos los ejercicios. Accesible para cualquier usuario logueado.
func GetExercises(c *gin.Context) {
	ejercicios, err := services.GetAllExercises()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los ejercicios"})
		return
	}
	c.JSON(http.StatusOK, ejercicios)
}

// solo para ADMIN
func CreateExercise(c *gin.Context) {
	var req models.CreateEjercicioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetString("user_id")

	ejercicio, err := services.CreateExercise(req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el ejercicio"})
		return
	}
	c.JSON(http.StatusCreated, ejercicio)
}

func UpdateExercise(c *gin.Context) {
	id := c.Param("id")
	var req models.UpdateEjercicioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ejercicio, err := services.UpdateExercise(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el ejercicio"})
		return
	}
	c.JSON(http.StatusOK, ejercicio)
}

func DeleteExercise(c *gin.Context) {
	id := c.Param("id")
	if err := services.DeleteExercise(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el ejercicio"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ejercicio eliminado exitosamente"})
}
