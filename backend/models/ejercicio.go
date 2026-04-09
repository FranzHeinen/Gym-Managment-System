package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Ejercicio struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Nombre             string             `bson:"nombre" json:"nombre" binding:"required"`
	Descripcion        string             `bson:"descripcion" json:"descripcion" binding:"required"`
	Categoria          string             `bson:"categoria" json:"categoria" binding:"required"`
	GrupoMuscular      string             `bson:"grupo_muscular" json:"grupo_muscular" binding:"required"`
	Dificultad         string             `bson:"dificultad" json:"dificultad" binding:"required"`
	Demostracion       string             `bson:"demostracion" json:"demostracion" binding:"required"`
	Instruccion        string             `bson:"instruccion" json:"instruccion" binding:"required"`
	UserID             primitive.ObjectID `bson:"user_id" json:"user_id" binding:"required"`
	FechaCreacion      time.Time          `bson:"fecha_creacion" json:"fecha_creacion"`
	FechaActualizacion time.Time          `bson:"fecha_ultima_actualizacion" json:"fecha_ultima_actualizacion"`
}

type CreateEjercicioRequest struct {
	Nombre        string `json:"nombre" binding:"required"`
	Descripcion   string `json:"descripcion" binding:"required"`
	Categoria     string `json:"categoria" binding:"required"`
	GrupoMuscular string `json:"grupo_muscular" binding:"required"`
	Dificultad    string `json:"dificultad" binding:"required"`
	Demostracion  string `json:"demostracion" binding:"required"`
	Instruccion   string `json:"instruccion" binding:"required"`
}

type UpdateEjercicioRequest struct {
	Nombre        string `json:"nombre"`
	Descripcion   string `json:"descripcion"`
	Categoria     string `json:"categoria"`
	GrupoMuscular string `json:"grupo_muscular"`
	Dificultad    string `json:"dificultad"`
	Demostracion  string `json:"demostracion"`
	Instruccion   string `json:"instruccion"`
}
