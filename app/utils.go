package main

import (
	"os"

	"golang.org/x/crypto/bcrypt"
)

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return ":3000"
	}
	return ":" + port
}
func getDistDir() string {
	filename := os.Getenv("DIST_DIR")
	if filename == "" {
		return "./client/dist"
	}
	return filename
}

type Hasher interface {
	Hash(password string) (string, error)
	Check(password string, hash string) error
}

type BCryptHasher struct{}

func NewBCryptHasher() Hasher {
	return BCryptHasher{}
}

func (h BCryptHasher) Hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (h BCryptHasher) Check(password string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
