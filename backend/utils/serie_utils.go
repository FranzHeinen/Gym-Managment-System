package utils

import (
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/dtos"
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/models"
)

// convierte seriedto a seriemodel
func SerieDtoToModel(serie dtos.Serie) *models.Serie {
	return &models.Serie{
		NumeroSerie:  serie.NumeroSerie,
		Peso:         serie.Peso,
		Repeticiones: serie.Repeticiones,
		Completada:   serie.Completada,
	}
}

// convierte la lista de seriesdto a la lista de seriesmodel
func SeriesDtoToModel(seriesDto []dtos.Serie) []models.Serie {
	if seriesDto == nil {
		return nil
	}

	series := make([]models.Serie, len(seriesDto))
	for i, serieDTO := range seriesDto {
		series[i] = *SerieDtoToModel(serieDTO)
	}
	return series

}

// convierte seriemodel a seriedto
func SerieModelToDto(serie models.Serie) dtos.Serie {
	return dtos.Serie{
		NumeroSerie:  serie.NumeroSerie,
		Peso:         serie.Peso,
		Repeticiones: serie.Repeticiones,
		Completada:   serie.Completada,
	}
}

// convierte la lista de seriesmodel a la lista de seriesdto
func SeriesModelToDto(seriesDto []models.Serie) []dtos.Serie {
	if seriesDto == nil {
		return nil
	}

	series := make([]dtos.Serie, len(seriesDto))
	for i, serieDTO := range seriesDto {
		series[i] = SerieModelToDto(serieDTO)
	}
	return series

}
