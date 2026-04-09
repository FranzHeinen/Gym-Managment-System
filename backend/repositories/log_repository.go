package repositories

import (
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/database"
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateLog(entry models.LogEntry) error {
	collection := database.Database.Collection("logs")
	_, err := collection.InsertOne(context.Background(), entry)
	return err
}

// GetLogs devuelve las últimas 'limit' entradas de log, ordenadas por fecha
func GetLogs(limit int64) ([]bson.M, error) {
	collection := database.Database.Collection("logs")
	var logs []bson.M

	findOptions := options.Find()
	findOptions.SetLimit(limit)
	findOptions.SetSort(bson.D{{Key: "timestamp", Value: -1}}) // Más nuevos primero

	pipeline := mongo.Pipeline{
		// 1. Ordenar por fecha (más nuevos primero)
		{{Key: "$sort", Value: bson.D{{Key: "timestamp", Value: -1}}}},
		// 2. Limitar resultados
		{{Key: "$limit", Value: limit}},
		// 3. Hacer $lookup para buscar el nombre en la colección "users"
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"}, // El nombre de tu colección de usuarios
			{Key: "localField", Value: "user_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "userInfo"},
		}}},
		// 4. "Desarmar" el array de userInfo
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$userInfo"},
			{Key: "preserveNullAndEmptyArrays", Value: true}, // Mantener si es "Invitado"
		}}},
		// 5. Proyectar el formato final
		{{Key: "$project", Value: bson.D{
			{Key: "timestamp", Value: 1},
			{Key: "status_code", Value: 1},
			{Key: "method", Value: 1},
			{Key: "path", Value: 1},
			{Key: "client_ip", Value: 1},
			// Usamos el nombre de userInfo si existe, si no, usamos el user_string
			{Key: "username", Value: bson.D{
				{Key: "$ifNull", Value: bson.A{"$userInfo.nombre", "$user_string"}},
			}},
		}}},
	}

	cursor, err := collection.Aggregate(context.Background(), pipeline, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &logs); err != nil {
		return nil, err
	}
	return logs, nil
}
