package repositories

import (
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/database"
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateRoutine(rutina models.Rutina) (string, error) {

	collection := database.Database.Collection("rutinas")
	result, err := collection.InsertOne(context.Background(), rutina)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func UpdateRoutine(rutina models.Rutina) (*models.Rutina, error) {
	collection := database.Database.Collection("rutinas")

	update := bson.M{
		"$set": bson.M{
			"fecha_ultima_actualizacion": time.Now(),
			"nombre":                     rutina.Nombre,
			"descripcion":                rutina.Descripcion,
		},
	}

	_, err := collection.UpdateOne(context.Background(), bson.M{
		"_id":     rutina.ID,
		"user_id": rutina.UserID,
	}, update)

	if err != nil {
		return nil, err
	}

	return &rutina, nil
}

func GetRoutineByID(rutinaID, userID primitive.ObjectID) (models.Rutina, error) {
	collection := database.Database.Collection("rutinas")

	var rutina models.Rutina
	err := collection.FindOne(context.Background(), bson.M{
		"_id":     rutinaID,
		"user_id": userID,
	}).Decode(&rutina)

	if err != nil {
		return models.Rutina{}, mongo.ErrNoDocuments
	}

	return rutina, nil
}

func GetRoutinesByUserID(userID primitive.ObjectID) ([]models.Rutina, error) {
	collection := database.Database.Collection("rutinas")
	var rutinas []models.Rutina
	cursor, err := collection.Find(context.Background(), bson.M{
		"user_id": userID,
	})
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &rutinas)
	return rutinas, err
}

func DeleteRoutine(rutinaID, userID primitive.ObjectID) error {
	collection := database.Database.Collection("rutinas")

	result, err := collection.DeleteOne(context.Background(), bson.M{
		"_id":     rutinaID,
		"user_id": userID,
	})

	if err != nil {
		return err
	}

	// Si DeletedCount es 0, significa que no se encontró ningún documento que coincidiera.
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
