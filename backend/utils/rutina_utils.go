package utils

import (
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/dtos"
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// convierte CreateRutinaRequest a models.Rutina
func CreaterutinaRequestToRutinaModel(req dtos.CreateRutinaRequest, userID primitive.ObjectID) models.Rutina {
	now := time.Now()

	return models.Rutina{
		ID:                 primitive.NewObjectID(),
		Nombre:             req.Nombre,
		Descripcion:        req.Descripcion,
		Ejercicios:         DTOToEjerciciosRutinaModel(req.Ejercicios),
		UserID:             userID,
		FechaCreacion:      now,
		FechaActualizacion: now,
	}
}

// convierte models.Rutina a dtos.Rutina
func RutinaModelToDTO(model models.Rutina) dtos.Rutina {
	return dtos.Rutina{
		ID:          model.ID,
		Nombre:      model.Nombre,
		Descripcion: model.Descripcion,
		Ejercicios:  EjerciciosRutinaModelToDTO(model.Ejercicios),
		UserID:      model.UserID,
	}
}
