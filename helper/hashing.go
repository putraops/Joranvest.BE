package helper

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

type Hashing interface {
	HashAndSalt(pwd []byte) string
}

func HashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		panic("Failed to hash a password")
	}
	return string(hash)
}
