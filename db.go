package main

import (
	"context"
	"time"

	"github.com/pchchv/golog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func saver(s Secret) string {
	// TODO: Implement data transfer to the database
	return s.password
}

func getter(pass string) (s Secret, err error) {
	// TODO: Implement data retrieval from the database
	return
}

func database() {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(getEnvValue("MONGO")).
		SetServerAPIOptions(serverAPIOptions)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		golog.Fatal(err.Error())
	}

	golog.Info("Connected to MongoDB!")
	keys_collection = client.Database(getEnvValue("DATABASE")).Collection("keys")
	secrets_collection = client.Database(getEnvValue("DATABASE")).Collection("secrets")
}
