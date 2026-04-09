package repositories

import (
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/database"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Cuenta todos los documentos en la colección de usuarios
func GetTotalUsersCount() (int64, error) {
	collection := database.Database.Collection("users")
	count, err := collection.CountDocuments(context.Background(), bson.M{})
	return count, err
}

// Cuenta las rutinas más completadas para un usuario específico.
func GetMostUsedRoutinesByUser(userID primitive.ObjectID) ([]bson.M, error) {
	collection := database.Database.Collection("workouts")
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "user_id", Value: userID}}}}, // Filtra los workouts solo para el usuario actual
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$rutina_id"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}}},
		{{Key: "$limit", Value: 5}},
		// Hacemos un $lookup para obtener la información de la rutina
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "rutinas"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "rutinaInfo"},
		}}},
		// $lookup devuelve un array, lo "desarmamos"
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$rutinaInfo"},
			{Key: "preserveNullAndEmptyArrays", Value: false},
		}}},
		// Proyectamos solo los campos que nos interesan
		{{Key: "$project", Value: bson.D{
			{Key: "count", Value: 1},
			{Key: "nombre", Value: "$rutinaInfo.nombre"},
			{Key: "_id", Value: 0},
		}}},
	}
	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	var results []bson.M
	if err = cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}
	return results, nil
}

// Obtener los ejercicios más utilizados
// Key especifica explícitamente el nombre de cada campo junto con su valor
func GetMostPopularExercises() ([]bson.M, error) {
	collection := database.Database.Collection("rutinas")
	pipeline := mongo.Pipeline{
		{{Key: "$unwind", Value: "$ejercicios"}}, // "Desarma" el array de ejercicios en documentos individuales
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$ejercicios.ejercicio_id"},        // Agrupa por el ID del ejericicio
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}, // Cuenta cuántas veces aparece
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}}}, // Ordena de más popular a menos popular
		{{Key: "$limit", Value: 5}},                                // Devuelve solo el top 5
		// Hacemos un $lookup para obtener la información del ejercicio (incluido el nombre)
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "exercises"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "ejercicioInfo"},
		}}},
		// $lookup devuelve un array, lo "desarmamos"
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$ejercicioInfo"},
			{Key: "preserveNullAndEmptyArrays", Value: false}, // Omitir si no hay info
		}}},
		// Proyectamos solo los campos que nos interesan
		{{Key: "$project", Value: bson.D{
			{Key: "count", Value: 1},
			{Key: "nombre", Value: "$ejercicioInfo.nombre"},
			{Key: "_id", Value: 0}, // Ocultamos el _id
		}}},
	}
	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	var results []bson.M
	if err = cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}
	return results, nil
}

// Obtener las rutinas más usadas
// (Esta función no se modifica, ya que no se usará en el dashboard global)
func GetMostUsedRoutines() ([]bson.M, error) {
	collection := database.Database.Collection("workouts")
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$rutina_id"},                      // Agrupa por el ID de la rutina
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}, // Cuenta las apariciones
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}}},
		{{Key: "$limit", Value: 5}},
	}
	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	var results []bson.M
	if err = cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}
	return results, nil
}

// Agrupa los entrenamientos por semana para un usuario específico.
func GetWorkoutsPerWeekByUser(userID primitive.ObjectID) ([]bson.M, error) {
	collection := database.Database.Collection("workouts")
	pipeline := mongo.Pipeline{
		// Filtra los documentos para obtener solo los del usuario actual
		{{Key: "$match", Value: bson.D{{Key: "user_id", Value: userID}}}},
		// Agrupa los documentos por la semana del año en que se registraron
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "year", Value: bson.D{{Key: "$year", Value: "$fecha"}}},    // Extrae el año de la fecha
				{Key: "week", Value: bson.D{{Key: "$isoWeek", Value: "$fecha"}}}, // Extrae el número de la semana
			}},
			{Key: "total_workouts", Value: bson.D{{Key: "$sum", Value: 1}}}, // Cuenta cuántos workouts hay en ese grupo
		}}},
		// Ordena los resultados por año y semana para tener una línea de tiempo
		{{Key: "$sort", Value: bson.D{{Key: "_id.year", Value: 1}, {Key: "_id.week", Value: 1}}}},
	}
	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	var results []bson.M
	if err = cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}
	return results, nil
}
