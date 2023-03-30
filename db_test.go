package main

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func TestFinder(t *testing.T) {
	// Set up a mock key and encrypted text in the collections
	keyId := primitive.NewObjectID()
	textId := primitive.NewObjectID()
	key := Key{
		key:  []byte("secretkey"),
		hash: "hashofsecretkey",
	}
	encryptedText := EncryptedText{
		text: "encryptedtext",
		hash: "hashofencryptedtext",
	}
	k, _ := bson.Marshal(key)
	tex, _ := bson.Marshal(encryptedText)
	keys_collection.InsertOne(context.Background(), k)
	secrets_collection.InsertOne(context.Background(), tex)

	// Test case: valid password
	pass := fmt.Sprint(keyId) + "{password}" + fmt.Sprint(textId)
	s, err := finder(pass)
	if err != nil {
		t.Errorf("finder() error = %v, wantErr nil", err)
	}
	if !reflect.DeepEqual(s.key, key) {
		t.Errorf("finder() s.key = %v, want %v", s.key, key)
	}
	if !reflect.DeepEqual(s.encryptedtext, encryptedText) {
		t.Errorf("finder() s.encryptedtext = %v, want %v", s.encryptedtext, encryptedText)
	}

	// Test case: invalid password
	pass = "invalidpassword"
	_, err = finder(pass)
	if err == nil {
		t.Errorf("finder() error = nil, wantErr")
	}
}
