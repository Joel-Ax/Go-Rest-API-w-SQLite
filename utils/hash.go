package utils

import "golang.org/x/crypto/bcrypt"

func Hash(password string) (string, error) {
	bcrypt, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bcrypt), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
