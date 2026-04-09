package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Session representa una sesión de usuario activa en la base de datos
type Session struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	UserAgent string             `bson:"user_agent" json:"user_agent"`
	ClientIP  string             `bson:"client_ip" json:"client_ip"`
	ExpiresAt time.Time          `bson:"expires_at" json:"expires_at"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
