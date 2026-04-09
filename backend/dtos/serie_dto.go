package dtos

type Serie struct {
	NumeroSerie  int     `json:"numero_serie"`
	Repeticiones int     `json:"repeticiones"`
	Peso         float64 `json:"peso"`
	Completada   bool    `json:"completada"`
}
