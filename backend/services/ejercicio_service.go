package services

import (
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/models"
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/repositories"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllExercises() ([]models.Ejercicio, error) {
	return repositories.GetAllExercises()
}

func CreateExercise(req models.CreateEjercicioRequest, userID string) (models.Ejercicio, error) {
	userObjID, _ := primitive.ObjectIDFromHex(userID)
	now := time.Now()

	ejercicio := models.Ejercicio{
		Nombre:             req.Nombre,
		Descripcion:        req.Descripcion,
		Categoria:          req.Categoria,
		GrupoMuscular:      req.GrupoMuscular,
		Dificultad:         req.Dificultad,
		Demostracion:       req.Demostracion,
		Instruccion:        req.Instruccion,
		UserID:             userObjID,
		FechaCreacion:      now,
		FechaActualizacion: now,
	}

	newID, err := repositories.CreateExercise(ejercicio)
	if err != nil {
		return models.Ejercicio{}, err
	}
	ejercicio.ID = newID
	return ejercicio, nil
}

func UpdateExercise(id string, req models.UpdateEjercicioRequest) (models.Ejercicio, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Ejercicio{}, err
	}

	ejercicio, err := repositories.GetExerciseByID(objID)
	if err != nil {
		return models.Ejercicio{}, err
	}

	// Actualizamos los campos del request
	if req.Nombre != "" {
		ejercicio.Nombre = req.Nombre
	}
	if req.Descripcion != "" {
		ejercicio.Descripcion = req.Descripcion
	}
	if req.Categoria != "" {
		ejercicio.Categoria = req.Categoria
	}
	if req.GrupoMuscular != "" {
		ejercicio.GrupoMuscular = req.GrupoMuscular
	}
	if req.Dificultad != "" {
		ejercicio.Dificultad = req.Dificultad
	}
	if req.Demostracion != "" {
		ejercicio.Demostracion = req.Demostracion
	}
	if req.Instruccion != "" {
		ejercicio.Instruccion = req.Instruccion
	}
	ejercicio.FechaActualizacion = time.Now()

	err = repositories.UpdateExercise(ejercicio)
	return ejercicio, err
}

func DeleteExercise(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	err = repositories.DeleteExercise(objID)
	if err != nil {
		return err
	}
	err = repositories.RemoveExerciseFromAllRoutines(objID)
	if err != nil {
		return err
	}

	return nil
}
