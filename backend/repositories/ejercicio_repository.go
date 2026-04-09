package repositories

import (
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/database"
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetExerciseByID(id primitive.ObjectID) (models.Ejercicio, error) {
	collection := database.Database.Collection("exercises")
	var ejercicio models.Ejercicio
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&ejercicio)
	return ejercicio, err
}

func GetAllExercises() ([]models.Ejercicio, error) {
	collection := database.Database.Collection("exercises")
	var ejercicios []models.Ejercicio
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &ejercicios)
	return ejercicios, err
}

func CreateExercise(ejercicio models.Ejercicio) (primitive.ObjectID, error) {
	collection := database.Database.Collection("exercises")
	result, err := collection.InsertOne(context.Background(), ejercicio)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func UpdateExercise(ejercicio models.Ejercicio) error {
	collection := database.Database.Collection("exercises")
	_, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": ejercicio.ID},
		bson.M{"$set": ejercicio},
	)
	return err
}

func DeleteExercise(id primitive.ObjectID) error {
	collection := database.Database.Collection("exercises")
	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}
