package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pchchv/golog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func inserter(s Secret) (password string, err error) {
	key := s.key
	text := s.encryptedtext

	k, err := bson.Marshal(key)
	if err != nil {
		return
	}

	result, err := keys_collection.InsertOne(context.TODO(), k)
	if err != nil {
		return
	}
	password = fmt.Sprint(result.InsertedID) + "{" + s.password + "}"

	t, err := bson.Marshal(text)
	if err != nil {
		return
	}

	result, err = secrets_collection.InsertOne(context.TODO(), t)
	if err != nil {
		return
	}

	return password + fmt.Sprint(result.InsertedID), nil
}

func finder(pass string) (s Secret, err error) {
	p := strings.Split(pass, "{")
	keyId := p[0]
	p = strings.Split(p[1], "}")
	s.password = p[0]
	textId := p[1]

	res := keys_collection.FindOneAndDelete(context.TODO(), bson.M{"_id": keyId})
	err = res.Decode(s.key)
	if err != nil {
		return
	}

	res = secrets_collection.FindOneAndDelete(context.TODO(), bson.M{"_id": textId})
	err = res.Decode(s.encryptedtext)
	if err != nil {
		return
	}

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
