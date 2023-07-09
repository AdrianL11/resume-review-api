package util

import (
	"crypto/sha256"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	pwPreHash := sha256.Sum256([]byte(password))
	bytes, err := bcrypt.GenerateFromPassword(pwPreHash[:], 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	pwPreHash := sha256.Sum256([]byte(password))
	err := bcrypt.CompareHashAndPassword([]byte(hash), pwPreHash[:])
	return err == nil
}
