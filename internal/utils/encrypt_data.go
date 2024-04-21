package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"os"
)

func Encrypt(text []byte) (string, error) {
	bytes := []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}
	block, err := aes.NewCipher([]byte(os.Getenv("ENCRYPTION_SECRET")))
	if err != nil {
		return "", err
	}

	cfb := cipher.NewCFBEncrypter(block, bytes)
	cipherText := make([]byte, len(text))
	cfb.XORKeyStream(cipherText, text)
	return base64.URLEncoding.EncodeToString(cipherText), nil
}
