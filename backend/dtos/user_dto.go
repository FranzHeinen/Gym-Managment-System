package dtos

import "time"

type ProfileResponse struct {
	ID               string    `json:"id"`
	Nombre           string    `json:"nombre"`
	Email            string    `json:"email"`
	FechaNacimiento  time.Time `json:"fecha_nacimiento,omitempty"`
	Peso             float64   `json:"peso,omitempty"`
	Altura           float64   `json:"altura,omitempty"`
	NivelExperiencia string    `json:"nivel_experiencia,omitempty"`
	Objetivos        []string  `json:"objetivos,omitempty"`
}

type UpdateProfileRequest struct {
	Nombre           string    `json:"nombre"`
	FechaNacimiento  time.Time `json:"fecha_nacimiento"`
	Peso             float64   `json:"peso"`
	Altura           float64   `json:"altura"`
	NivelExperiencia string    `json:"nivel_experiencia"`
	Objetivos        []string  `json:"objetivos"`
}

// DTO para la solicitud de cambio de contraseña
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}
