package utils

import (
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/dtos"
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/models"
)

// convierte la lista de ejercicios rutinadto a la lista para los models
func DTOToEjerciciosRutinaModel(ejerciciosDTO []dtos.EjercicioRutina) []models.EjercicioRutina {
	if ejerciciosDTO == nil {
		return nil
	}

	ejercicios := make([]models.EjercicioRutina, len(ejerciciosDTO))
	for i, ejDTO := range ejerciciosDTO {
		ejercicios[i] = DTOToEjercicioRutinaModel(ejDTO)
	}
	return ejercicios
}

// convierte dtos.EjercicioRutina a models.EjercicioRutina
func DTOToEjercicioRutinaModel(dto dtos.EjercicioRutina) models.EjercicioRutina {
	return models.EjercicioRutina{
		EjercicioID:    dto.EjercicioID,
		Orden:          dto.Orden,
		Series:         SeriesDtoToModel(dto.Series),
		TiempoDescanso: dto.TiempoDescanso,
	}
}

// convierte la lista de models.EjercicioRutina a la lista de dtos.EjercicioRutina
func EjerciciosRutinaModelToDTO(ejerciciosModel []models.EjercicioRutina) []dtos.EjercicioRutina {
	if ejerciciosModel == nil {
		return nil
	}

	ejercicios := make([]dtos.EjercicioRutina, len(ejerciciosModel))
	for i, ejModel := range ejerciciosModel {
		ejercicios[i] = EjercicioRutinaModelToDTO(ejModel)
	}
	return ejercicios
}

// convierte models.EjercicioRutina a dtos.EjercicioRutina
func EjercicioRutinaModelToDTO(model models.EjercicioRutina) dtos.EjercicioRutina {
	return dtos.EjercicioRutina{
		ID:             model.ID,
		EjercicioID:    model.EjercicioID,
		Orden:          model.Orden,
		Series:         SeriesModelToDto(model.Series),
		TiempoDescanso: model.TiempoDescanso,
	}
}
