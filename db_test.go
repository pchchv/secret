package main

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestInserter(t *testing.T) {
	// set up test environment
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database("test_db")
	keys_collection := db.Collection("keys")
	secrets_collection := db.Collection("secrets")

	// create a Secret instance for testing
	key, err := genKey()
	if err != nil {
		t.Fatal(err)
	}
	encryptedText, _, err := encrypt("test message")
	if err != nil {
		t.Fatal(err)
	}
	hash, err := hashPassword("password123")
	if err != nil {
		t.Fatal(err)
	}
	secret := Secret{
		encryptedtext: EncryptedText{
			text: encryptedText,
			hash: hash,
		},
		key: Key{
			key:  key,
			hash: hash,
		},
		password: "testpassword123",
	}

	// test inserting the Secret into the database
	password, err := inserter(secret)
	if err != nil {
		t.Fatal(err)
	}

	// check that the password string has the correct format
	expectedPrefix := fmt.Sprint(1) + "{" + "testpassword123" + "}"
	if !strings.HasPrefix(password, expectedPrefix) {
		t.Errorf("unexpected password prefix, expected %v, got %v", expectedPrefix, password)
	}

	// check that the keys and secrets were inserted into their respective collections
	var retrievedKey Key
	err = keys_collection.FindOne(context.Background(), bson.M{}).Decode(&retrievedKey)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(key, retrievedKey.key) {
		t.Errorf("unexpected key, expected %v, got %v", key, retrievedKey.key)
	}
	if !checkPasswordHash(hash, retrievedKey.hash) {
		t.Errorf("unexpected key hash, expected %v, got %v", hash, retrievedKey.hash)
	}

	var retrievedText EncryptedText
	err = secrets_collection.FindOne(context.Background(), bson.M{}).Decode(&retrievedText)
	if err != nil {
		t.Fatal(err)
	}
	if encryptedText != retrievedText.text {
		t.Errorf("unexpected text, expected %v, got %v", encryptedText, retrievedText.text)
	}
	if !checkPasswordHash(hash, retrievedText.hash) {
		t.Errorf("unexpected text hash, expected %v, got %v", hash, retrievedText.hash)
	}
}
