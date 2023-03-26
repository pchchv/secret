package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

const plaintext = `Lorem ipsum dolor sit amet, consectetur adipiscing elit.
	Phasellus non purus nec sem fringilla vestibulum. Phasellus ornare justo quis enim euismod,
	vitae posuere augue malesuada. Duis nec lacus risus. Integer vulputate tincidunt diam,
	eget tristique nulla iaculis nec. Praesent auctor eget orci in tempus. Cras lorem elit,
	posuere sed velit non, pulvinar aliquet mauris. Cras porttitor leo vitae venenatis sagittis.
	Integer in pharetra quam. Donec faucibus, dolor non condimentum fringilla, ante ex eleifend purus,
	ac sollicitudin diam augue congue erat. Vestibulum malesuada dolor vitae purus faucibus tincidunt.
	Mauris hendrerit lacinia sapien eget pellentesque. Praesent eget volutpat neque, elementum fermentum libero.
	Suspendisse consequat maximus laoreet. Phasellus nec volutpat tellus, vel egestas nisl.
	Suspendisse diam augue, egestas sit amet laoreet vel, aliquam nec quam. Nunc lobortis tellus urna,
	ut rutrum dolor commodo nec. Morbi dignissim nulla eget tortor tempus, eu iaculis felis dapibus.
	Maecenas sapien turpis, iaculis eu faucibus eu, blandit a ex. Proin egestas, turpis ac cursus pulvinar,
	elit orci fermentum mi, nec ornare eros leo a neque. Proin vitae gravida mi. Ut at nunc nisl.
	Nulla consequat aliquet magna. Aenean venenatis a ante non varius. Cras sed nisl tortor.
	Duis sagittis ac quam et hendrerit. Sed congue faucibus tortor suscipit convallis. Suspendisse potenti.
	Nulla consectetur vestibulum nunc, non volutpat odio lobortis mollis.
	Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus.
	Maecenas in ex sodales, euismod risus eu, dapibus elit. Fusce finibus placerat vulputate.
	Integer molestie erat id convallis fringilla. Vivamus eu diam sed ex blandit volutpat. Etiam at convallis nibh.
	In a magna nisi. Sed non nisi justo. Vestibulum quis orci orci. Mauris et nibh nisi. Cras et rhoncus nunc.
	Aenean maximus nibh eu nisi tempor gravida. Maecenas interdum venenatis elit ac convallis.
	Nam nulla risus, consectetur at finibus ut, vestibulum congue nunc. Curabitur dignissim et velit et faucibus.
	Nullam vel mauris ac neque fermentum laoreet. Maecenas ipsum justo, pulvinar at posuere pellentesque, semper sed erat.
	Pellentesque eu tortor et eros pretium volutpat eu sed quam. Maecenas volutpat sem lacus, vel elementum sapien vehicula sed.
	Mauris a lobortis sem, sit amet viverra massa. Proin porta nibh sed metus mattis, at iaculis sapien malesuada.
	Nulla iaculis tortor a augue imperdiet, id tincidunt neque feugiat. Duis egestas eu lacus a feugiat.
	Nam commodo nunc quis arcu hendrerit, non blandit mi volutpat.
	Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae;
	Nunc a nisi sed mauris vestibulum condimentum. Nam mauris justo, lobortis quis nunc ut, ultrices dignissim leo.
	Quisque quis quam nec libero pharetra scelerisque in et dui. Praesent turpis nunc, varius sit amet iaculis nec, iaculis in diam.
	Proin eu elementum diam, id condimentum nunc. Sed pulvinar ipsum augue, in cursus leo consectetur sed.
	Aliquam posuere, lectus vel sagittis molestie, lectus odio consequat ipsum, in facilisis ipsum metus vel augue.
	Ut ut ligula ut diam suscipit ultrices. Vivamus mattis porttitor sagittis. Mauris vel est dui.
	Curabitur posuere libero eu pharetra porta. Sed porttitor rhoncus sagittis. Vivamus efficitur dictum velit at congue.
	Ut vitae leo sit amet purus consequat eleifend. Aliquam a hendrerit sem. Sed non odio at odio cursus rutrum non auctor turpis.
	Sed congue vitae erat quis commodo. Sed eleifend pretium augue, in posuere sapien feugiat id. Donec commodo, ex eu ultricies imperdiet,
	dolor libero bibendum ante, vitae sodales tortor diam a orci. Vivamus consectetur arcu ut congue.`

func TestPasswordHash(t *testing.T) {
	pass := "qewhy#fcu3!rt"

	_, err := hashPassword(pass)

	if err != nil {
		t.Fatal()
	}
}

func TestHashPassword(t *testing.T) {
	password := "test123"

	hashedPassword, err := hashPassword(password)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err != nil {
		t.Errorf("Password does not match hash: %v", err)
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

func TestCheckHashPassword(t *testing.T) {
	password := "test123"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	match := checkPasswordHash(password, string(hashedPassword))

	if !match {
		t.Error("Password does not match hash")
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

func TestGenKey2(t *testing.T) {
	key, err := genKey()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(key) != 32 {
		t.Errorf("Unexpected key length: expected 32, got %v", len(key))
	}
}

func TestEncrypt(t *testing.T) {
	ct, key, err := encrypt(plaintext)
	if err != nil {
		t.Fatal(err)
	}

	if len(ct) < 2 {
		t.Fatal()
	}

	if len([]byte(key)) < aes.BlockSize {
		t.Fatal()
	}
}

func TestEncrypt2(t *testing.T) {
	encoded, key, err := encrypt(plaintext)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(key) != 32 {
		t.Errorf("Unexpected key length: expected 32, got %v", len(key))
	}

	decoded, err := base64.RawStdEncoding.DecodeString(encoded)
	if err != nil {
		t.Errorf("Unexpected error decoding base64: %v", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		t.Errorf("Unexpected error creating cipher block: %v", err)
	}

	if len(decoded) < aes.BlockSize {
		t.Errorf("Unexpected decoded length: expected at least %v, got %v", aes.BlockSize, len(decoded))
	}

	iv := decoded[:aes.BlockSize]
	decoded = decoded[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(decoded, decoded)

	if string(decoded) != plaintext {
		t.Errorf("Unexpected decrypted message: expected %q, got %q", plaintext, string(decoded))
	}
}

func TestDecrypt(t *testing.T) {
	ct, key, err := encrypt(plaintext)
	if err != nil {
		t.Fatal(err)
	}

	text, err := decrypt(key, ct)
	if err != nil {
		t.Fatal(err)
	}

	if text != plaintext {
		t.Fatal()
	}
}
