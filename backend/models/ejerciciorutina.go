package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type EjercicioRutina struct {
	ID             primitive.ObjectID `json:"id" bson:"_id, omitempty"`
	EjercicioID    primitive.ObjectID `json:"ejercicio_id" binding:"required" bson:"ejercicio_id"`
	Orden          int                `json:"orden"`
	Series         []Serie            `json:"series" binding:"required"`
	TiempoDescanso int                `json:"tiempo_descanso" bson:"tiempo_descanso"`
}
