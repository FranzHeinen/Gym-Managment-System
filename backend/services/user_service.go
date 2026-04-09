package services

import (
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/models"
	"Trabajo-Practico-2025-Programaci-n2-HEINEN-ISAAC-WILLINER/backend/repositories"
)

// GetAllUsers devuelve todos los usarios y se asegura de no exponer las contraseñas
func GetAllUsers() ([]models.User, error) {
	users, err := repositories.GetAllUsers()
	if err != nil {
		return nil, err
	}
	// No devuelve las contraseñas hasheadas
	for i := range users {
		users[i].Password = ""
	}
	return users, nil
}
