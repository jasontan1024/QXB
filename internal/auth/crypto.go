package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/argon2"
)

// argon2 参数
const (
	argonTime    uint32 = 1
	argonMemory  uint32 = 64 * 1024
	argonThreads uint8  = 1
	argonKeyLen  uint32 = 32
)

func deriveKey(password string, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)
}

func randomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	return b, err
}

// hashPassword 返回 hash 和 salt（均 base64）
func hashPassword(password string) (hashB64, saltB64 string, err error) {
	salt, err := randomBytes(16)
	if err != nil {
		return "", "", err
	}
	key := deriveKey(password, salt)
	return base64.StdEncoding.EncodeToString(key), base64.StdEncoding.EncodeToString(salt), nil
}

func verifyPassword(password, hashB64, saltB64 string) bool {
	salt, err := base64.StdEncoding.DecodeString(saltB64)
	if err != nil {
		return false
	}
	expected, err := base64.StdEncoding.DecodeString(hashB64)
	if err != nil {
		return false
	}
	key := deriveKey(password, salt)
	if len(key) != len(expected) {
		return false
	}
	var diff byte
	for i := range key {
		diff |= key[i] ^ expected[i]
	}
	return diff == 0
}

func encryptPrivateKey(password string, plaintext []byte) (cipherB64, saltB64 string, err error) {
	salt, err := randomBytes(16)
	if err != nil {
		return "", "", err
	}
	key := deriveKey(password, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", "", err
	}
	nonce, err := randomBytes(gcm.NonceSize())
	if err != nil {
		return "", "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), base64.StdEncoding.EncodeToString(salt), nil
}

func decryptPrivateKey(password, cipherB64, saltB64 string) ([]byte, error) {
	salt, err := base64.StdEncoding.DecodeString(saltB64)
	if err != nil {
		return nil, err
	}
	ciphertext, err := base64.StdEncoding.DecodeString(cipherB64)
	if err != nil {
		return nil, err
	}
	key := deriveKey(password, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}
	nonce := ciphertext[:gcm.NonceSize()]
	body := ciphertext[gcm.NonceSize():]
	return gcm.Open(nil, nonce, body, nil)
}
