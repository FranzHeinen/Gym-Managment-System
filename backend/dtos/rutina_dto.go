package dtos

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Rutina struct {
	ID          primitive.ObjectID `json:"id"`
	Nombre      string             `json:"nombre" binding:"required"`
	Descripcion string             `json:"descripcion" binding:"required"`
	Ejercicios  []EjercicioRutina  `json:"ejercicios"`
	UserID      primitive.ObjectID `json:"user_id" binding:"required"`
}

type CreateRutinaRequest struct {
	Nombre      string            `json:"nombre" binding:"required"`
	Descripcion string            `json:"descripcion" binding:"required"`
	Ejercicios  []EjercicioRutina `json:"ejercicios"`
}

type UpdateRutinaRequest struct {
	Nombre      string `json:"nombre,omitempty"`
	Descripcion string `json:"descripcion,omitempty"`
}

type RutinaResponse struct {
	ID      string `json:"id"`
	Mensaje string `json:"mensaje"`
}
