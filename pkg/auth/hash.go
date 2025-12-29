package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashingPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("gagal untuk hashing password")
	}
	
	return hashedPassword, nil
}