package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	mrand "math/rand"
	"os"

	"github.com/joho/godotenv"
	"github.com/pchchv/golog"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

const charSet = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "!@#$%&*" + "0123456789"

var (
	keys_collection    *mongo.Collection
	secrets_collection *mongo.Collection
)

type Secret struct {
	encryptedtext EncryptedText
	key           Key
	password      string
}

type Key struct {
	key  []byte
	hash string
}

type EncryptedText struct {
	text string
	hash string
}

func init() {
	// Load values from .env into the system
	if err := godotenv.Load(); err != nil {
		golog.Panic("No .env file found")
	}
}

func getEnvValue(v string) string {
	value, exist := os.LookupEnv(v)
	if !exist {
		golog.Panic("Value %v does not exist", v)
	}
	return value
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func genPassword() (password string) {
	length := mrand.Intn(17-7) + 7
	for i := 0; i < length; i++ {
		password += string([]rune(charSet)[mrand.Intn(len(charSet))])
	}
	return
}

func genKey() ([]byte, error) {
	key := make([]byte, 32)

	if _, err := rand.Read(key); err != nil {
		return nil, err
	}

	return key, nil
}

func encrypt(text string) (encoded string, key []byte, err error) {
	// create byte array from the input string
	plainText := []byte(text)

	// getting the 32-bit passphrase
	key, err = genKey()
	if err != nil {
		return
	}

	// create a new AES cipher using the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	// make the cipher text a byte array of size BlockSize + the length of the message
	cipherText := make([]byte, aes.BlockSize+len(plainText))

	// iv is the ciphertext up to the blocksize
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	// encrypt the data:
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	// return string encoded in base64
	return base64.RawStdEncoding.EncodeToString(cipherText), key, err
}

func decrypt(key []byte, secure string) (decoded string, err error) {
	// remove base64 encoding:
	cipherText, err := base64.RawStdEncoding.DecodeString(secure)
	if err != nil {
		return
	}

	// create a new AES cipher with the key and encrypted message
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	if len(cipherText) < aes.BlockSize {
		err = errors.New("Ciphertext block size is too short!")
		return
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	// decrypt the message
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), err
}

func reverse(s string) (result string) {
	for _, v := range s {
		result = string(v) + result
	}
	return
}

func encryptor(text string) (pass string, err error) {
	s := Secret{}
	s.encryptedtext.text, s.key.key, err = encrypt(text)
	if err != nil {
		return
	}

	s.password = genPassword()
	s.encryptedtext.hash, err = hashPassword(s.password)
	s.key.hash, err = hashPassword(reverse(s.password))
	if err != nil {
		return
	}

	return inserter(s)
}

func decryptor(password string) (text string, err error) {
	secret, err := finder(password)
	if err != nil {
		return
	}

	if !checkPasswordHash(reverse(password), secret.key.hash) || !checkPasswordHash(password, secret.encryptedtext.hash) {
		return "", errors.New("invalid password")
	}

	text, err = decrypt(secret.key.key, secret.encryptedtext.text)
	if err != nil {
		return
	}

	return
}

func main() {
	database()
	tgbot()
}
