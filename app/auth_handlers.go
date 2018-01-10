package main

import (
	"fmt"
	"net/http"
	"regexp"
)

type Credentials struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phonenumber"`
}
type PasswordlessBody struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}
type TokenBody struct {
	Token string `json:"token"`
}

var validPhone = regexp.MustCompile(`^\+\d{11}$`)

func (s Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var cred Credentials
	if err := DecodeBody(r, &cred); err != nil {
		RespondErr(w, r, http.StatusBadRequest, err)
		return
	}
	if cred.Email == "" {
		RespondErr(w, r, http.StatusBadRequest, fmt.Errorf("Invalid email"))
		return
	}
	if cred.Password == "" {
		RespondErr(w, r, http.StatusBadRequest, fmt.Errorf("Invalid password"))
		return
	}
	user, err := s.auth.LogIn(cred.Email, cred.Password)
	if err != nil {
		RespondErr(w, r, http.StatusForbidden, err)
		return
	}

	// Check if should TwoFactor
	if user.TwoFactor {
		// If yes - send email
		request, err := s.passwordless.Add(user.ID.Hex())
		if err != nil {
			RespondErr(w, r, http.StatusForbidden, err)
			return
		}
		if err := s.emailcode.SendCode(user, request.Code); err != nil {
			RespondErr(w, r, http.StatusInternalServerError, err)
			return
		}
		Respond(w, r, http.StatusOK, map[string]string{"email": user.Email})
		return
	}
	// If no - respond with token
	token, err := s.auth.UserIdToToken(user.ID)
	if err != nil {
		RespondErr(w, r, http.StatusForbidden, err)
		return
	}
	Respond(w, r, http.StatusOK, TokenBody{token})
	return
}

func (s Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	var cred Credentials
	if err := DecodeBody(r, &cred); err != nil {
		RespondErr(w, r, http.StatusBadRequest, err)
		return
	}
	// TODO: validate email
	if cred.Email == "" {
		RespondErr(w, r, http.StatusBadRequest, fmt.Errorf("Invalid email '%s'", cred.Email))
		return
	}
	if cred.Password == "" {
		RespondErr(w, r, http.StatusBadRequest, fmt.Errorf("Invalid password"))
		return
	}
	if cred.PhoneNumber == "" || !validPhone.MatchString(cred.PhoneNumber) {
		RespondErr(w, r, http.StatusBadRequest, fmt.Errorf("Invalid phone number '%s'", cred.PhoneNumber))
		return
	}

	user, err := s.auth.Register(cred.Email, cred.Password)
	if err != nil {
		RespondErr(w, r, http.StatusForbidden, err)
		return
	}
	updatedUser, err := s.user.Update(user.ID.Hex(), NewUserUpdates().PhoneNumber(cred.PhoneNumber))
	if err != nil {
		RespondErr(w, r, http.StatusForbidden, err)
		return
	}
	request, err := s.passwordless.Add(updatedUser.ID.Hex())
	if err != nil {
		RespondErr(w, r, http.StatusForbidden, err)
		return
	}
	if err := s.emailcode.SendCode(user, request.Code); err != nil {
		RespondErr(w, r, http.StatusInternalServerError, err)
		return
	}
	Respond(w, r, http.StatusOK, map[string]string{"email": updatedUser.Email})
}

func (s Server) handleCode(w http.ResponseWriter, r *http.Request) {
	var body PasswordlessBody
	if err := DecodeBody(r, &body); err != nil {
		RespondErr(w, r, http.StatusBadRequest, err)
		return
	}
	if body.Email == "" {
		RespondErr(w, r, http.StatusBadRequest, fmt.Errorf("Invalid email"))
		return
	}
	if body.Code == "" {
		RespondErr(w, r, http.StatusBadRequest, fmt.Errorf("Invalid code"))
		return
	}
	user, err := s.user.ReadByEmail(body.Email)
	if err != nil {
		RespondErr(w, r, http.StatusForbidden, err)
		return
	}
	if err := s.passwordless.Verify(body.Code, user.ID.Hex()); err != nil {
		RespondErr(w, r, http.StatusForbidden, err)
		return
	}
	token, err := s.auth.UserIdToToken(user.ID)
	if err != nil {
		RespondErr(w, r, http.StatusForbidden, err)
		return
	}
	Respond(w, r, http.StatusOK, TokenBody{token})
}
