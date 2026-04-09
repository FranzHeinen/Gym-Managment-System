package services

import (
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/repositories"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Llama al repositorio
func GetTotalUsersCount() (int64, error) {
	return repositories.GetTotalUsersCount()
}
func GetMostPopularExercises() ([]bson.M, error) {
	return repositories.GetMostPopularExercises()
}

func GetMostUsedRoutines() ([]bson.M, error) {
	return repositories.GetMostUsedRoutines()
}

// Reune todas las estadísticas para un usuario específico.
func GetUserStats(userIDStr string) (map[string]interface{}, error) {
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return nil, err
	}

	// Obtener historial de workouts para calcular la frecuencia
	workouts, err := repositories.GetWorkoutsByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Obtener las rutinas más usadas por el usuario
	mostUsedRoutines, err := repositories.GetMostUsedRoutinesByUser(userID)
	if err != nil {
		return nil, err
	}

	// Obtener el progreso de entrenamientos por semana
	progressOverTime, err := repositories.GetWorkoutsPerWeekByUser(userID)
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"frecuencia_entrenamiento": len(workouts),
		"rutinas_mas_utilizadas":   mostUsedRoutines,
		"progreso_en_el_tiempo":    progressOverTime,
	}

	return stats, nil
}
