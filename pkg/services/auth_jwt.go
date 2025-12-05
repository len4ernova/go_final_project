package services

import (
	"crypto/aes"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	expTime = 8
	Key     = "1234567891234567"
)

var JwtKey = []byte("my_secret_key")

type Claims struct {
	CheckSum string `json:"checksum"`
	jwt.RegisteredClaims
}

func GenerateJWT(checksum string) (string, error) {
	experationTime := time.Now().Add(expTime * time.Hour)
	claims := &Claims{
		CheckSum: checksum,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(experationTime),
			Issuer:    "Planner",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtKey)
}

// GetPass - зашифровать пароль с помощью ключа.
// Длина пароля должна быть 16 знаков.
// Если короче, тогда дополнить нулями.
// При последующем сравнении нужно учитывать и дополнять короткие пароли.
// Returns: закодированный ключ в base64.
func EncryptPass(key, text string) (string, error) {
	bc, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", fmt.Errorf("Error: create a new [cipher.Block]: %v", err)
	}
	var dst = make([]byte, 16)
	//

	//var src []byte
	if len(text) < 16 {
		text = fmt.Sprintf("%016s", text)
	} else if len(text) > 16 {
		return "", fmt.Errorf("Error: expected the lengths password (<= 16 characters): %v", err)
	}
	src := []byte(text)
	bc.Encrypt(dst, src)

	return base64.StdEncoding.EncodeToString(dst), nil
}

// tODO del
// сравнение введенного пароля и сохраненного в БД
func ComparePass(key, passInp, passDB string) (bool, error) {
	passG, err := EncryptPass(key, passInp)
	if err != nil {
		return false, err
	}
	fmt.Println(passG, passDB, passG == passDB)
	return passG == passDB, nil

}

// DecryptPass - декодирование пароля
func DecryptPass(key, passDb string) (string, error) {
	text, err := base64.StdEncoding.DecodeString(passDb)
	if err != nil {
		return "", fmt.Errorf("could not base64 decode: %v", err)
	}

	bc, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", fmt.Errorf("could not create new cipher: %v", err)
	}
	if len(text) < aes.BlockSize {
		return "", fmt.Errorf("invalid ciphertext block size")
	}

	var res = make([]byte, 16)
	bc.Decrypt(res, text)
	fmt.Println("ffff", string(res))
	return string(res), nil
}
