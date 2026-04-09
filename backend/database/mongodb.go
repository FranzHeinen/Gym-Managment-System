package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var Database *mongo.Database

func Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return err
	}

	// Verificar la conexión
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	Client = client
	Database = client.Database("prueba_frontend")

	go func() {
		collection := Database.Collection("sessions")
		indexModel := mongo.IndexModel{
			Keys:    bson.M{"expires_at": 1},
			Options: options.Index().SetExpireAfterSeconds(0),
		}

		_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
		if err != nil {
			log.Printf("Error al crear índice TTL para sesiones: %v\n", err)
		} else {
			log.Println("Índice TTL para 'sessions' asegurado.")
		}
	}()

	log.Println("Conectado a MongoDB exitosamente")
	return nil
}

func Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return Client.Disconnect(ctx)
}
