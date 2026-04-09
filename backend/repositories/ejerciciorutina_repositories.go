package repositories

import (
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/database"
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/models"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddExerciseToRoutine(rutinaID, userID primitive.ObjectID, ejercicio *models.EjercicioRutina) error {
	collection := database.Database.Collection("rutinas")

	filter := bson.M{
		"_id":     rutinaID,
		"user_id": userID,
	}

	update := bson.M{
		"$push": bson.M{"ejercicios": ejercicio},
		"$set":  bson.M{"fecha_ultima_actualizacion": time.Now()},
	}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("rutina no encontrada")
	}
	return nil
}

func GetExerciseFromRoutine(rutinaID, userID, ejercicioRutinaID primitive.ObjectID) (models.EjercicioRutina, error) {
	rutina, err := GetRoutineByID(rutinaID, userID)
	if err != nil {
		return models.EjercicioRutina{}, err
	}

	for _, ejercicio := range rutina.Ejercicios {
		if ejercicio.ID == ejercicioRutinaID {
			return ejercicio, nil
		}
	}

	return models.EjercicioRutina{}, errors.New("ejercicio no encontrado")
}

func UpdateExerciseInRoutine(rutinaID, userID, ejercicioRutinaID primitive.ObjectID, ejercicio *models.EjercicioRutina) error {
	collection := database.Database.Collection("rutinas")

	filter := bson.M{
		"_id":            rutinaID,
		"user_id":        userID,
		"ejercicios._id": ejercicioRutinaID,
	}

	update := bson.M{
		"$set": bson.M{
			"ejercicios.$.orden":           ejercicio.Orden,
			"ejercicios.$.series":          ejercicio.Series,
			"ejercicios.$.tiempo_descanso": ejercicio.TiempoDescanso,
			"fecha_ultima_actualizacion":   time.Now(),
		},
	}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("no se encontró la rutina o el ejercicio para actualizar")
	}

	return nil
}

func RemoveExerciseFromRoutine(rutinaID, userID, ejercicioID primitive.ObjectID) error {
	collection := database.Database.Collection("rutinas")

	update := bson.M{
		"$pull": bson.M{
			"ejercicios": bson.M{"_id": ejercicioID},
		},
	}

	filter := bson.M{
		"_id":     rutinaID,
		"user_id": userID,
	}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func RemoveExerciseFromAllRoutines(ejercicioID primitive.ObjectID) error {
	collection := database.Database.Collection("rutinas")

	// Filtra todas las rutinas que contengan el ejercicio_id en su array "ejercicios"
	filter := bson.M{"ejercicios.ejercicio_id": ejercicioID}

	// Define la operación de $pull para eliminar el subdocumento que coincida
	update := bson.M{
		"$pull": bson.M{
			"ejercicios": bson.M{"ejercicio_id": ejercicioID},
		},
		// Opcional: actualiza la fecha de modificación de las rutinas afectadas
		"$set": bson.M{"fecha_ultima_actualizacion": time.Now()},
	}

	// UpdateMany se asegura de que se aplique a TODAS las rutinas, no solo a la primera
	_, err := collection.UpdateMany(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
