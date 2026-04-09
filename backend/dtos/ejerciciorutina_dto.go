package dtos

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EjercicioRutina struct {
	ID             primitive.ObjectID `json:"id"`
	EjercicioID    primitive.ObjectID `json:"ejercicio_id" binding:"required"`
	Orden          int                `json:"orden"`
	Series         []Serie            `json:"series"`
	TiempoDescanso int                `json:"tiempo_descanso"`
}
