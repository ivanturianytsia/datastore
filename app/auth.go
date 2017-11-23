package main

import (
	"fmt"
	"net/http"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

type AuthService interface {
	LogIn(email, password string) (User, error)
	Register(email, password string) (User, error)
	UserIdToToken(userId bson.ObjectId) (string, error)

	TokenFromRequest(r *http.Request) (string, error)
	UserForToken(token string) (User, error)
	UserFromRequest(r *http.Request) (User, error)
}

type authService struct {
	store UserStore
	token TokenService
}

func NewAuthService(userStore UserStore) AuthService {
	token := NewTokenService(getTokenOptions())
	return &authService{
		store: userStore,
		token: token,
	}
}

func (s *authService) LogIn(email, password string) (User, error) {
	if email == "" {
		return User{}, fmt.Errorf("Email '%s' is not valid", email)
	}
	if password == "" {
		return User{}, fmt.Errorf("Password is not valid", password)
	}
	user, err := s.store.ReadByEmail(email)
	if err != nil {
		return User{}, err
	}
	if user.ID == "" {
		return User{}, fmt.Errorf("No user found")
	}
	if err := s.store.CheckPassword(user.ID.Hex(), password); err != nil {
		return User{}, err
	}
	return user, nil
}
func (s *authService) Register(email, password string) (User, error) {
	if email == "" {
		return User{}, fmt.Errorf("Email '%s' is not valid", email)
	}
	if password == "" {
		return User{}, fmt.Errorf("Password is not valid", password)
	}
	user, err := s.store.Create(email, password)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (s *authService) UserIdToToken(userId bson.ObjectId) (string, error) {
	return s.token.Issue(map[string]string{
		"id": userId.Hex(),
	})
}

func (s *authService) TokenFromRequest(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("No access token in Header.")
	}
	return RemoveWhiteSpace(strings.TrimLeft(authHeader, "Bearer")), nil
}

func (s *authService) UserForToken(token string) (User, error) {
	claims, err := s.token.Parse(token)
	if err != nil {
		return User{}, err
	}
	id, ok := claims["id"]
	if !ok {
		return User{}, fmt.Errorf("No user id in token data")
	}
	return s.store.ReadById(id)
}

func (s *authService) UserFromRequest(r *http.Request) (User, error) {
	token, err := s.TokenFromRequest(r)
	if err != nil {
		return User{}, err
	}
	if token == "" {
		return User{}, fmt.Errorf("No token in the request")
	}
	user, err := s.UserForToken(token)
	if err != nil {
		return User{}, err
	}
	return User{
		ID:    user.ID,
		Email: user.Email,
	}, err
}
