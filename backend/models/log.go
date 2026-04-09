package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// LogEntry representa una entrada de log en la base de datos
type LogEntry struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Timestamp  time.Time          `bson:"timestamp" json:"timestamp"`
	StatusCode int                `bson:"status_code" json:"status_code"`
	Latency    time.Duration      `bson:"latency" json:"latency"`
	ClientIP   string             `bson:"client_ip" json:"client_ip"`
	Method     string             `bson:"method" json:"method"`
	Path       string             `bson:"path" json:"path"`
	UserID     primitive.ObjectID `bson:"user_id,omitempty" json:"user_id,omitempty"`
	UserString string             `bson:"user_string,omitempty" json:"user_string,omitempty"`
}
