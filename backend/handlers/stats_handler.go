package handlers

import (
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetGlobalStats(c *gin.Context) {
	totalUsers, _ := services.GetTotalUsersCount()
	popularExercises, _ := services.GetMostPopularExercises()

	c.JSON(http.StatusOK, gin.H{
		"total_usuarios":       totalUsers,
		"ejercicios_populares": popularExercises,
	})
}

func GetUserStats(c *gin.Context) {
	userID := c.GetString("user_id")

	stats, err := services.GetUserStats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener las estadísticas personales"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
