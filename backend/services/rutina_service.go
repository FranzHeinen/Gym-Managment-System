package services

import (
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/dtos"
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/repositories"
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/utils"

	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateRoutine(req dtos.CreateRutinaRequest, userID string) (*dtos.RutinaResponse, error) {

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	rutinaModel := utils.CreaterutinaRequestToRutinaModel(req, userObjectID)

	rutinaID, err := repositories.CreateRoutine(rutinaModel)
	if err != nil {
		return nil, err
	}

	return &dtos.RutinaResponse{
		ID:      rutinaID,
		Mensaje: "Rutina creada exitosamente",
	}, nil
}

func UpdateRoutine(rutinaID, userIDStr string, req dtos.UpdateRutinaRequest) (*dtos.RutinaResponse, error) {

	// Convertir IDs
	rutinaObjectID, err := primitive.ObjectIDFromHex(rutinaID)
	if err != nil {
		return nil, err
	}

	userObjectID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return nil, err
	}

	rutina, err := repositories.GetRoutineByID(rutinaObjectID, userObjectID)
	if err != nil || rutina.ID.IsZero() {
		return nil, err
	}

	if req.Nombre != "" {
		rutina.Nombre = req.Nombre
	}
	if req.Descripcion != "" {
		rutina.Descripcion = req.Descripcion
	}

	// Actualizar en repository
	updatedRutina, err := repositories.UpdateRoutine(rutina)
	if err != nil {
		return nil, err
	}

	// Convertir a DTO de response
	return &dtos.RutinaResponse{
		ID:      updatedRutina.ID.String(),
		Mensaje: "Rutina actualizada exitosamente",
	}, nil
}

func GetRoutines(userIDstr string) ([]dtos.Rutina, error) {

	userObjectID, err := primitive.ObjectIDFromHex(userIDstr)
	if err != nil {
		return nil, err
	}

	rutinasModel, err := repositories.GetRoutinesByUserID(userObjectID)
	if err != nil {
		return nil, err
	}

	var rutinas []dtos.Rutina
	for i := 0; i < len(rutinasModel); i++ {

		rutinas = append(rutinas, utils.RutinaModelToDTO(rutinasModel[i]))
	}

	return rutinas, nil

}

func GetRoutineByID(rutinaIDStr, userIDStr string) (*dtos.Rutina, error) {
	rutinaID, err := primitive.ObjectIDFromHex(rutinaIDStr)
	if err != nil {
		return nil, err
	}

	rutinas, err := GetRoutines(userIDStr)
	if err != nil {

		return nil, err
	}
	for _, rutina := range rutinas {
		if rutina.ID == rutinaID {
			return &rutina, nil
		}
	}

	return nil, mongo.ErrNoDocuments
}

func DeleteRoutine(rutinaIDStr, userIDStr string) error {
	rutinaID, err := primitive.ObjectIDFromHex(rutinaIDStr)
	if err != nil {
		return err
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return err
	}

	return repositories.DeleteRoutine(rutinaID, userID)
}

func DuplicateRoutine(rutinaIDStr, userIDStr string) (*dtos.RutinaResponse, error) {
	rutinaID, err := primitive.ObjectIDFromHex(rutinaIDStr)
	if err != nil {
		return nil, err
	}
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return nil, err
	}

	rutinaOriginal, err := repositories.GetRoutineByID(rutinaID, userID)
	if err != nil {
		return nil, err
	}

	nuevaRutina := rutinaOriginal
	nuevaRutina.ID = primitive.NewObjectID()
	nuevaRutina.Nombre = rutinaOriginal.Nombre + " (Copia)"
	nuevaRutina.FechaCreacion = time.Now()
	nuevaRutina.FechaActualizacion = time.Now()

	nuevoID, err := repositories.CreateRoutine(nuevaRutina)
	if err != nil {
		return nil, err
	}

	return &dtos.RutinaResponse{
		ID:      nuevoID,
		Mensaje: "Rutina duplicada exitosamente",
	}, nil
}
