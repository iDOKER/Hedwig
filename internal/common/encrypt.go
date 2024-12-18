package common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
)

func GenerateKey(token string) []byte {
	hash := sha256.Sum256([]byte(token))
	key := hash[:16] // 128 位（16 字节）
	return key
}

// 使用 AES-GCM 对数据进行加密
func UuidEncrypt(plainText []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, 12) // GCM 推荐 nonce 大小为 12 字节
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	cipherText := aesGCM.Seal(nil, nonce, []byte(plainText), nil)
	result := append(nonce, cipherText...)
	return base64.StdEncoding.EncodeToString(result), nil
}

// 使用 AES-GCM 对数据进行解密
func UuidDecrypt(encryptedText string, key []byte) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	if len(data) < 12 {
		return "", fmt.Errorf("invalid data size")
	}

	nonce := data[:12]
	cipherText := data[12:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plainText, err := aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}
