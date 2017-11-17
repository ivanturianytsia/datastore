package main

import (
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

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

func getTokenOptions() ([]byte, time.Duration) {
	secret := os.Getenv("TOKEN_SECRET")
	if secret == "" {
		secret = "SUPER_TOEKN_SECPET"
	}
	expiration, err := strconv.Atoi(os.Getenv("TOKEN_EXPIRATION"))
	if err != nil {
		expiration = 24
	}
	return []byte(secret), (time.Hour * time.Duration(expiration))
}

func RemoveWhiteSpace(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
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
