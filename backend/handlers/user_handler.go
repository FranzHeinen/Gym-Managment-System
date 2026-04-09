package handlers

import (
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/auth"
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/dtos"
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/repositories"
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetProfile(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, _ := primitive.ObjectIDFromHex(userIDStr)

	user, err := repositories.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	c.JSON(http.StatusOK, dtos.ProfileResponse{
		ID:               user.ID.Hex(),
		Nombre:           user.Nombre,
		Email:            user.Email,
		FechaNacimiento:  user.FechaNacimiento,
		Peso:             user.Peso,
		Altura:           user.Altura,
		NivelExperiencia: user.NivelExperiencia,
		Objetivos:        user.Objetivos,
	})
}

func UpdateProfile(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, _ := primitive.ObjectIDFromHex(userIDStr)

	var req dtos.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := repositories.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	user.Nombre = req.Nombre
	user.FechaNacimiento = req.FechaNacimiento
	user.Peso = req.Peso
	user.Altura = req.Altura
	user.NivelExperiencia = req.NivelExperiencia
	user.Objetivos = req.Objetivos
	user.UpdatedAt = time.Now()

	if err := repositories.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el perfil"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Perfil actualizado exitosamente"})
}

// cambia la contraseña del usuario
func ChangePassword(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, _ := primitive.ObjectIDFromHex(userIDStr)

	var req dtos.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := repositories.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	// Verificar
	if !auth.CheckPasswordHash(req.OldPassword, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "La contraseña actual es incorrecta"})
		return
	}

	// Hashear y guardar nueva contraseña
	hashedPassword, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar la nueva contraseña"})
		return
	}
	user.Password = hashedPassword
	user.UpdatedAt = time.Now()

	if err := repositories.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al cambiar la contraseña"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contraseña actualizada exitosamente"})
}

// GetAllUsers es un handler solo para administradores para obtener todos los usuarios.
func GetAllUsers(c *gin.Context) {
	users, err := services.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los usuarios"})
		return
	}
	c.JSON(http.StatusOK, users)
}
