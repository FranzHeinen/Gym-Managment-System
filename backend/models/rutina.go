package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Rutina struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Nombre             string             `bson:"nombre" json:"nombre" binding:"required"`
	Descripcion        string             `bson:"descripcion" json:"descripcion" binding:"required"`
	Ejercicios         []EjercicioRutina  `bson:"ejercicios" json:"ejercicios"`
	UserID             primitive.ObjectID `bson:"user_id" json:"user_id" binding:"required"`
	FechaCreacion      time.Time          `bson:"fecha_creacion" json:"fecha_creacion"`
	FechaActualizacion time.Time          `bson:"fecha_ultima_actualizacion" json:"fecha_ultima_actualizacion"`
}
