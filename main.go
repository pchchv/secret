package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	mrand "math/rand"

	"golang.org/x/crypto/bcrypt"
)

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

func main() {}
