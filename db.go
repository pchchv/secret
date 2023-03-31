package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pchchv/golog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func inserter(s Secret) (password string, err error) {
	keyBson, err := bson.Marshal(s.key)
	if err != nil {
		return "", err
	}

	res, err := keys_collection.InsertOne(context.Background(), keyBson)
	if err != nil {
		return "", err
	}

	password = fmt.Sprintf("%v{%v}", res.InsertedID, s.password)

	textBson, err := bson.Marshal(s.encryptedtext)
	if err != nil {
		return "", err
	}

	res, err = secrets_collection.InsertOne(context.Background(), textBson)
	if err != nil {
		return "", err
	}

	password += fmt.Sprintf("%v", res.InsertedID)

	return password, nil
}

func finder(pass string) (s Secret, err error) {
	splitPass := strings.Split(pass, "{")
	keyID := splitPass[0]
	password := strings.TrimSuffix(splitPass[1], "}")

	objectID, err := primitive.ObjectIDFromHex(keyID)
	if err != nil {
		return s, err
	}

	res := keys_collection.FindOneAndDelete(context.Background(), bson.M{"_id": objectID})
	if err = res.Decode(&s.key); err != nil {
		return s, err
	}

	objectID, err = primitive.ObjectIDFromHex(password)
	if err != nil {
		return s, err
	}

	res = secrets_collection.FindOneAndDelete(context.Background(), bson.M{"_id": objectID})
	if err = res.Decode(&s.encryptedtext); err != nil {
		return s, err
	}

	s.password = password
	return s, nil
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

	err = client.Ping(ctx, nil)
	if err != nil {
		golog.Fatal(err.Error())
	}

	golog.Info("Connected to MongoDB!")
	keys_collection = client.Database(getEnvValue("DATABASE")).Collection("keys")
	secrets_collection = client.Database(getEnvValue("DATABASE")).Collection("secrets")
}
