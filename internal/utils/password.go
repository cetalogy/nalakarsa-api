package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword hashes a plain-text password using Bcrypt with cost 10
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

// ComparePassword compares a hashed password with a plain-text candidate
func ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
