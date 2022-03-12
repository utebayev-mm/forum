package model

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(password string) string{
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Println(err)
	}
	return string(bytes)
}

func CheckPassword(hash, password string) bool{
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}