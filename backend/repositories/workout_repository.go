package repositories

import (
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/database"
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// RegisterWorkout inserta un nuevo documento de workout en la base de datos.
func RegisterWorkout(workout models.Workout) error {
	collection := database.Database.Collection("workouts")
	_, err := collection.InsertOne(context.Background(), workout)
	return err
}

// GetWorkoutsByUserID devuelve todos los workouts de un usuario específico.
func GetWorkoutsByUserID(userID primitive.ObjectID) ([]bson.M, error) {
	collection := database.Database.Collection("workouts")

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "user_id", Value: userID}}}},
		{{Key: "$sort", Value: bson.D{{Key: "fecha", Value: -1}}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "rutinas"},
			{Key: "localField", Value: "rutina_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "rutinaInfo"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$rutinaInfo"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 1},
			{Key: "fecha", Value: 1},
			{Key: "rutina_id", Value: 1},
			{Key: "rutina_nombre", Value: bson.D{
				{Key: "$ifNull", Value: bson.A{"$rutinaInfo.nombre", "Rutina Eliminada"}},
			}},
		}}},
	}

	var results []bson.M
	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}

	return results, nil
}
