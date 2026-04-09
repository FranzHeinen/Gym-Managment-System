package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Workout representa un entrenamiento completado por un usuario en una fecha específica.
type Workout struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID   primitive.ObjectID `bson:"user_id" json:"user_id"`
	RutinaID primitive.ObjectID `bson:"rutina_id" json:"rutina_id"`
	Fecha    time.Time          `bson:"fecha" json:"fecha"`
}

// RegisterWorkoutRequest es la estructura que esperamos en el body del POST.
type RegisterWorkoutRequest struct {
	RutinaID string `json:"rutina_id" binding:"required"`
}
