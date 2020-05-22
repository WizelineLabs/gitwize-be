package cypher

import (
	"testing"
)

const passPhase = "abcdef 123456"
const testString = "Hello World"

func Test_CreateHash(t *testing.T) {
	s := createHash(passPhase)
	if len(s) != 32 {
		t.Errorf("expected %d, got %d", 32, len(s))
	}
}

func Test_Encrypt(t *testing.T) {
	encrypted := encrypt([]byte(testString), passPhase)
	if len(encrypted) != 39 {
		t.Errorf("expected encrypted length than %d, got %d", 39, len(encrypted))
	}
}

func Test_Decrypt(t *testing.T) {
	encrypted := encrypt([]byte(testString), passPhase)
	decrypted := decrypt(encrypted, passPhase)
	if string(decrypted) != testString {
		t.Errorf("expected %s, got %s", testString, decrypted)
	}
}

func Test_EncryptAndDecryptString(t *testing.T) {
	encrypted := EncryptString(testString, passPhase)
	decrypted := DecryptString(encrypted, passPhase)
	if len(encrypted) != 52 {
		t.Errorf("expected encrypted length than %d, got %d", 52, len(encrypted))
	}
	if decrypted != testString {
		t.Errorf("expected %s, got %s", testString, decrypted)
	}
}
