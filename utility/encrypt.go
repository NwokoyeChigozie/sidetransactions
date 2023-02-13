package utility

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func Hash(str string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	return string(hashed), err
}

func CompareHash(str string, hashed string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(str)) == nil
}

func ShaHash(str string) (string, error) {
	passSha1 := sha1.New()
	_, err := passSha1.Write([]byte(str))
	if err != nil {
		return str, err
	}

	getSha1 := passSha1.Sum(nil)
	return fmt.Sprintf("%x", getSha1), nil
}

func Md5(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
