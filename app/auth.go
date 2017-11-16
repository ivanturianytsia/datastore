package main

import (
	"fmt"
)

type UserService interface {
	LogIn(email, password string) (User, error)
	Register(email, password string) (User, error)
}

type userService struct {
	store UserStore
}

func NewUserService(userStore UserStore) UserService {
	return &userService{
		store: userStore,
	}
}

func (us *userService) LogIn(email, password string) (User, error) {
	if email == "" {
		return User{}, fmt.Errorf("Email '%s' is not valid", email)
	}
	if password == "" {
		return User{}, fmt.Errorf("Password is not valid", password)
	}
	user, err := us.store.ReadByEmail(email)
	if err != nil {
		return User{}, err
	}
	if user.ID == "" {
		return User{}, fmt.Errorf("No user found")
	}
	if err := us.store.CheckPassword(user.ID.Hex(), password); err != nil {
		return User{}, err
	}
	return user, nil
}
func (us *userService) Register(email, password string) (User, error) {
	if email == "" {
		return User{}, fmt.Errorf("Email '%s' is not valid", email)
	}
	if password == "" {
		return User{}, fmt.Errorf("Password is not valid", password)
	}
	user, err := us.store.Create(email, password)
	if err != nil {
		return User{}, err
	}
	return user, nil
}
