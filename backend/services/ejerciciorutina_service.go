package services

import (
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/dtos"
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/repositories"
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/utils"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddExerciseToRoutine(rutinaID, userIDStr string, req dtos.EjercicioRutina) (*dtos.RutinaResponse, error) {
	rutinaObjectID, err := primitive.ObjectIDFromHex(rutinaID)
	if err != nil {
		return nil, errors.New("ID de rutina inválido")
	}

	userObjectID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return nil, errors.New("ID de usuario inválido")
	}

	var nuevoEjercicio = utils.DTOToEjercicioRutinaModel(req)

	nuevoEjercicio.ID = primitive.NewObjectID()

	err = repositories.AddExerciseToRoutine(rutinaObjectID, userObjectID, &nuevoEjercicio)
	if err != nil {
		return nil, err
	}

	return &dtos.RutinaResponse{
		ID:      rutinaID,
		Mensaje: "Ejercicio agregado exitosamente",
	}, nil
}

func UpdateExerciseInRoutine(rutinaIDStr, ejercicioRutinaIDStr, userIDStr string, req dtos.EjercicioRutina) (*dtos.RutinaResponse, error) {
	if rutinaIDStr == "" || ejercicioRutinaIDStr == "" || userIDStr == "" {
		return nil, errors.New("IDs de rutina, ejercicio y usuario requeridos")
	}

	rutinaObjectID, err := primitive.ObjectIDFromHex(rutinaIDStr)
	if err != nil {
		return nil, errors.New("ID de rutina inválido")
	}

	ejercicioRutinaObjectID, err := primitive.ObjectIDFromHex(ejercicioRutinaIDStr)
	if err != nil {
		return nil, errors.New("ID de ejercicio inválido")
	}

	userObjectID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return nil, errors.New("ID de usuario inválido")
	}

	ejercicioActual, err := repositories.GetExerciseFromRoutine(rutinaObjectID, userObjectID, ejercicioRutinaObjectID)
	if err != nil {
		return nil, errors.New("ejercicio no encontrado en la rutina")
	}

	if req.Orden != 0 {
		ejercicioActual.Orden = req.Orden
	}
	if req.TiempoDescanso != 0 {
		ejercicioActual.TiempoDescanso = req.TiempoDescanso
	}

	if req.Series != nil {
		ejercicioActual.Series = utils.SeriesDtoToModel(req.Series)
	}

	err = repositories.UpdateExerciseInRoutine(rutinaObjectID, userObjectID, ejercicioRutinaObjectID, &ejercicioActual)
	if err != nil {
		return nil, err
	}

	return &dtos.RutinaResponse{
		ID:      rutinaIDStr,
		Mensaje: "Ejercicio actualizado exitosamente",
	}, nil

}

func GetExerciseFromRoutine(rutinaIDStr, ejercicioIDStr, userIDStr string) (*dtos.EjercicioRutina, error) {
	ejercicios, err := GetExercisesFromRoutine(rutinaIDStr, userIDStr)
	if err != nil {
		return nil, err
	}

	ejercicioID, err := primitive.ObjectIDFromHex(ejercicioIDStr)
	if err != nil {
		return nil, err
	}

	for _, ejercicio := range ejercicios {
		if ejercicio.EjercicioID == ejercicioID {
			return &ejercicio, nil // Lo encontramos y lo devolvemos.
		}
	}

	// Si el bucle termina, significa que no lo encontramos.
	return nil, mongo.ErrNoDocuments
}
func GetExercisesFromRoutine(rutinaIDStr, userIDStr string) ([]dtos.EjercicioRutina, error) {

	rutina, err := GetRoutineByID(rutinaIDStr, userIDStr)
	if err != nil {
		return nil, err
	}

	return rutina.Ejercicios, nil
}

func RemoveExerciseFromRoutine(rutinaIDStr, ejercicioIDStr, userIDStr string) error {

	rutinaID, err := primitive.ObjectIDFromHex(rutinaIDStr)
	if err != nil {
		return errors.New("formato de ID de rutina inválido")
	}
	ejercicioID, err := primitive.ObjectIDFromHex(ejercicioIDStr)
	if err != nil {
		return errors.New("formato de ID de ejercicio inválido")
	}
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return errors.New("formato de ID de usuario inválido")
	}

	return repositories.RemoveExerciseFromRoutine(rutinaID, userID, ejercicioID)
}
