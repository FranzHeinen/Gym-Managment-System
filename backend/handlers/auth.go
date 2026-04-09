package handlers

import (
	"context"
	"net/http"
	"time"

	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/auth"
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/database"
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/models"
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/repositories"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar si el usuario ya existe
	collection := database.Database.Collection("users")
	var existingUser models.User
	err := collection.FindOne(context.Background(), bson.M{"email": req.Email}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "El usuario ya existe"})
		return
	}

	// Hash de la contraseña
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar la contraseña"})
		return
	}

	// Crear nuevo usuario
	user := models.User{
		Nombre:    req.Nombre,
		Email:     req.Email,
		Password:  hashedPassword,
		Rol:       "USER",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el usuario"})
		return
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	user.Password = "" // No devolver la contraseña

	expirationTime := time.Now().Add(24 * time.Hour)
	session := models.Session{
		UserID:    user.ID,
		UserAgent: c.Request.UserAgent(),
		ClientIP:  c.ClientIP(),
		ExpiresAt: expirationTime,
		CreatedAt: time.Now(),
	}
	sessionID, err := repositories.CreateSession(session)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear la sesión"})
		return
	}

	// Generar token
	token, err := auth.GenerateToken(user.ID, user.Email, user.Rol, sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar el token"})
		return
	}

	c.JSON(http.StatusCreated, models.AuthResponse{
		Token: token,
		User:  user,
	})
}

func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Buscar usuario
	collection := database.Database.Collection("users")
	var user models.User
	err := collection.FindOne(context.Background(), bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar el usuario"})
		}
		return
	}

	// Verificar contraseña
	if !auth.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas"})
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	session := models.Session{
		UserID:    user.ID,
		UserAgent: c.Request.UserAgent(),
		ClientIP:  c.ClientIP(),
		ExpiresAt: expirationTime,
		CreatedAt: time.Now(),
	}
	sessionID, err := repositories.CreateSession(session)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear la sesión"})
		return
	}

	// Generar token
	token, err := auth.GenerateToken(user.ID, user.Email, user.Rol, sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar el token"})
		return
	}

	user.Password = "" // No devolver la contraseña

	c.JSON(http.StatusOK, models.AuthResponse{
		Token: token,
		User:  user,
	})
}

func Logout(c *gin.Context) {

	sessionIDHex, exists := c.Get("session_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudo encontrar la sesión del token"})
		return
	}

	sessionID, err := primitive.ObjectIDFromHex(sessionIDHex.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de sesión inválido"})
		return
	}

	// Eliminar la sesión de la base de datos
	err = repositories.DeleteSession(sessionID)
	if err != nil {
		// Si no la encuentra, es porque ya se borró. No es un error grave.
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusOK, gin.H{"message": "Sesión ya estaba cerrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al cerrar la sesión"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sesión cerrada exitosamente"})
}
