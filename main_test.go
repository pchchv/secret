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
