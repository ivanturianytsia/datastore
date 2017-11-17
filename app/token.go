package main

import (
	jwt "github.com/dgrijalva/jwt-go"

	"fmt"
	"time"
)

type TokenService interface {
	Issue(data map[string]string) (string, error)
	Parse(token string) (map[string]string, error)
}

type tokenService struct {
	secret     []byte
	expiration time.Duration
}

func NewTokenService(secret []byte, expiresIn time.Duration) TokenService {
	return &tokenService{
		secret:     secret,
		expiration: expiresIn,
	}
}

func (ts tokenService) Issue(data map[string]string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	for key, val := range data {
		token.Claims.(jwt.MapClaims)[key] = val
	}
	token.Claims.(jwt.MapClaims)["exp"] = time.Now().Add(ts.expiration).Unix()

	tokenString, err := token.SignedString(ts.secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (ts tokenService) Parse(token string) (map[string]string, error) {

	parsedToken, err := jwt.Parse(token, func(tok *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := tok.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", tok.Header["alg"])
		}
		return ts.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("Error parsing token: Token is invalid")
	}

	result := map[string]string{}
	for key, val := range claims {
		if key == "exp" {
			continue
		}
		result[key] = val.(string)
	}
	return result, nil
}
