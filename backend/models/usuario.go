package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Nombre           string             `bson:"nombre" json:"nombre"`
	Email            string             `bson:"email" json:"email" binding:"required,email"`
	Password         string             `bson:"password" json:"-"`
	Rol              string             `bson:"rol" json:"rol"` // "ADMIN" o "USER"
	FechaNacimiento  time.Time          `bson:"fecha_nacimiento,omitempty" json:"fecha_nacimiento,omitempty"`
	Peso             float64            `bson:"peso,omitempty" json:"peso,omitempty"`
	Altura           float64            `bson:"altura,omitempty" json:"altura,omitempty"`
	NivelExperiencia string             `bson:"nivel_experiencia,omitempty" json:"nivel_experiencia,omitempty"` // "Principiante", "Intermedio", "Avanzado"
	Objetivos        []string           `bson:"objetivos,omitempty" json:"objetivos,omitempty"`                 // "Perder peso", "Ganar músculo", "Mantenerse"
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Nombre   string `json:"nombre" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
