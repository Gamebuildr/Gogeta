package gbcrypto

import "testing"

// TODO: TestEncrypt, the IV is created on the fly so non-trivial
// func TestEncrypt(t *testing.T) {
// 	secretKey := "SecretKey"
// 	message := "encrypt me"
// 	expected := "Pz3cJ61A4Dw4vTFcqWN6t39pyWyOGCQe9F0="

// 	encrypted := Encrypt(secretKey, message)

// 	if encrypted != expected {
// 		t.Errorf("Expected %v, got %v", expected, encrypted)
// 	}
// }

func TestDecrypt(t *testing.T) {
	secretKey := "SecretKey"
	message := "Pz3cJ61A4Dw4vTFcqWN6t39pyWyOGCQe9F0="
	expected := "encrypt me"

	cryptography := Cryptography{}
	decrypted := cryptography.Decrypt(secretKey, message)

	if decrypted != expected {
		t.Errorf("Expected %v, got %v", expected, decrypted)
	}
}
