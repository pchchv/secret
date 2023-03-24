package main

import "testing"

func TestPasswordHash(t *testing.T) {
	pass := "qewhy#fcu3!rt"

	_, err := hashPassword(pass)

	if err != nil {
		t.Fatal()
	}
}

func TestCheckPasswordHash(t *testing.T) {
	pass := "qewhy#fcu3!rt"

	hash, err := hashPassword(pass)

	if err != nil {
		t.Fatal()
	}

	if !checkPasswordHash(pass, hash) {
		t.Fatal()
	}
}

func TestGenKey(t *testing.T) {
	key, err := genKey()
	if err != nil {
		t.Fatal(err)
	}

	if key[0] == 0 && key[1] == 0 && key[2] == 0 {
		t.Fatal()
	}
}
