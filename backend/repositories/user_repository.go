package repositories

import (
	"context"

	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/database"
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NOTA: Se eliminó la variable "userCollection" de aquí para evitar el pánico al iniciar.

func GetUserByID(userID primitive.ObjectID) (models.User, error) {
	// La colección se define aquí adentro, en lugar de afuera.
	userCollection := database.Database.Collection("users")
	var user models.User
	err := userCollection.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user)
	return user, err
}

func UpdateUser(user models.User) error {
	// La colección se define aquí adentro, en lugar de afuera.
	userCollection := database.Database.Collection("users")
	_, err := userCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": user.ID},
		bson.M{"$set": user},
	)
	return err
}

func GetAllUsers() ([]models.User, error) {
	userCollection := database.Database.Collection("users")
	var users []models.User

	cursor, err := userCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &users); err != nil {
		return nil, err
	}

	return users, nil
}
