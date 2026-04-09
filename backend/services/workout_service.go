package services

import (
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/models"
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/repositories"

	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RegisterWorkout(req models.RegisterWorkoutRequest, userIDStr string) error {
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return err
	}
	rutinaID, err := primitive.ObjectIDFromHex(req.RutinaID)
	if err != nil {
		return err
	}

	workout := models.Workout{
		UserID:   userID,
		RutinaID: rutinaID,
		Fecha:    time.Now(),
	}
	return repositories.RegisterWorkout(workout)
}

func GetWorkoutHistory(userIDStr string) ([]bson.M, error) {
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return nil, err
	}
	return repositories.GetWorkoutsByUserID(userID)
}
