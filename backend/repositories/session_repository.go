package repositories

import (
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/database"
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateSession(session models.Session) (primitive.ObjectID, error) {
	collection := database.Database.Collection("sessions")
	result, err := collection.InsertOne(context.Background(), session)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func GetSessionByID(sessionID primitive.ObjectID) (models.Session, error) {
	collection := database.Database.Collection("sessions")
	var session models.Session
	err := collection.FindOne(context.Background(), bson.M{"_id": sessionID}).Decode(&session)
	return session, err
}

func DeleteSession(sessionID primitive.ObjectID) error {
	collection := database.Database.Collection("sessions")
	result, err := collection.DeleteOne(context.Background(), bson.M{"_id": sessionID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
