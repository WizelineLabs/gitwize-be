package cypher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
)

func createSHA256Hash(key string) []byte {
	hash := sha256.Sum256([]byte(key))
	return hash[:]
}

func aes256Encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher(createSHA256Hash(passphrase))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func aes256Decrypt(data []byte, passphrase string) []byte {
	key := createSHA256Hash(passphrase)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}

//EncryptString export
func EncryptString(s, passphase string) string {
	data := []byte(s)
	return EncodeBase64(aes256Encrypt(data, passphase))
}

//DecryptString export
func DecryptString(s, passphase string) string {
	data := DecodeBase64(s)
	return string(aes256Decrypt(data, passphase))
}

//EncodeBase64 export
func EncodeBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

//DecodeBase64 export
func DecodeBase64(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}
