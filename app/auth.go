package main

import (
	"fmt"
	"net/http"
	"strings"
)

type AuthService interface {
	LogIn(email, password string) (string, error)
	Register(email, password string) (string, error)
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

func (s *authService) LogIn(email, password string) (string, error) {
	if email == "" {
		return "", fmt.Errorf("Email '%s' is not valid", email)
	}
	if password == "" {
		return "", fmt.Errorf("Password is not valid", password)
	}
	user, err := s.store.ReadByEmail(email)
	if err != nil {
		return "", err
	}
	if user.ID == "" {
		return "", fmt.Errorf("No user found")
	}
	if err := s.store.CheckPassword(user.ID.Hex(), password); err != nil {
		return "", err
	}
	token, err := s.token.Issue(map[string]string{
		"id": user.ID.Hex(),
	})
	if err != nil {
		return "", err
	}
	return token, nil
}
func (s *authService) Register(email, password string) (string, error) {
	if email == "" {
		return "", fmt.Errorf("Email '%s' is not valid", email)
	}
	if password == "" {
		return "", fmt.Errorf("Password is not valid", password)
	}
	user, err := s.store.Create(email, password)
	if err != nil {
		return "", err
	}
	token, err := s.token.Issue(map[string]string{
		"id": user.ID.Hex(),
	})
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *authService) TokenFromRequest(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authentication")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("No access token in Cookie or Header.")
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
