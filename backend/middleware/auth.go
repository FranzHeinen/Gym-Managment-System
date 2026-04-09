package middleware

import (
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/auth"
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/repositories"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token de autorización requerido"})
			c.Abort()
			return
		}

		// Verificar que el header tenga el formato "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de token inválido"})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		sessionID, err := primitive.ObjectIDFromHex(claims.SessionID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token con sesión corrupta"})
			c.Abort()
			return
		}

		// Buscamos la sesión en la base de datos
		_, err = repositories.GetSessionByID(sessionID)
		if err != nil {
			if err == mongo.ErrNoDocuments {

				c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesión no válida (cerrada o expirada)"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al verificar la sesión"})
			}
			c.Abort()
			return
		}

		// Inyectar el user_id en el contexto
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_rol", claims.Rol)
		c.Set("session_id", claims.SessionID)
		c.Next()
	}
}

// CheckAdmin - Middleware para verificar el rol de Administrador (Ruta: /admin)
func CheckAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtenemos el rol del contexto (es instantáneo, ya se leyó del JWT).
		rol, exists := c.Get("user_rol")
		if !exists || rol != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Acceso denegado. Se requiere rol de administrador."})
			c.Abort()
			return
		}

		c.Next()
	}
}

// CheckUser - Middleware para verificar el rol de Usuario (Ruta: /user)
func CheckUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtenemos el rol del contexto.
		rol, exists := c.Get("user_rol")
		// Un usuario normal puede acceder a rutas de usuario, y un admin también.
		if !exists || (rol != "USER" && rol != "ADMIN") {
			c.JSON(http.StatusForbidden, gin.H{"error": "Acceso denegado. Se requiere un rol de usuario válido."})
			c.Abort()
			return
		}

		c.Next()
	}
}
