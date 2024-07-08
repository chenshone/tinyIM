package tools

import (
	"golang.org/x/crypto/bcrypt"
)

func HashWithSalt(data string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
	return string(hashedPassword)
}
