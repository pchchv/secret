package main

import (
	"context"
	"fmt"
	"time"

	"github.com/pchchv/golog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func inserter(s Secret) (password string) {
	key := s.key
	text := s.encryptedtext

	k, err := bson.Marshal(key)
	if err != nil {
		golog.Panic(err.Error())
	}

	result, err := keys_collection.InsertOne(context.TODO(), k)
	if err != nil {
		golog.Panic(err.Error())
	}
	password = fmt.Sprint(result.InsertedID) + "{" + s.password + "}"

	t, err := bson.Marshal(text)
	if err != nil {
		golog.Panic(err.Error())
	}

	result, err = secrets_collection.InsertOne(context.TODO(), t)
	if err != nil {
		golog.Panic(err.Error())
	}

	password += fmt.Sprint(result.InsertedID)

	return
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
