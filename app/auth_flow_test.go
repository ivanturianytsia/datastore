package main

import (
	"fmt"
	"net/http"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

var tAuthData = struct {
	service   AuthService
	userid    bson.ObjectId
	email     string
	secret    string
	badSecret string
	code      string
	token     string
}{
	email:     "test2@mail.com",
	secret:    "secret",
	badSecret: "secret1",
}

var tPasswordlessRequestStore PasswordlessRequestStore

// Service

func Test_AuthRegister(t *testing.T) {
	user, err := tAuthData.service.Register(tAuthData.email, tAuthData.secret)
	if err != nil {
		t.Error(err)
		return
	}
	if user.ID.Hex() == "" {
		t.Errorf("No user id")
		return
	}
	if user.Email != tAuthData.email {
		t.Errorf("Invalid user email, expected '%s', got '%s'.", tAuthData.email, user.Email)
		return
	}
}
func Test_AuthDoubleRegister(t *testing.T) {
	if _, err := tAuthData.service.Register(tAuthData.email, tAuthData.secret); err == nil {
		t.Error("Should not register same email")
		return
	}
}
func Test_AuthLogIn(t *testing.T) {
	user, err := tAuthData.service.LogIn(tAuthData.email, tAuthData.secret)
	if err != nil {
		t.Error(err)
		return
	}
	if user.ID.Hex() == "" {
		t.Errorf("No user id")
		return
	}
	if user.Email != tAuthData.email {
		t.Errorf("Invalid user email, expected '%s', got '%s'.", tAuthData.email, user.Email)
		return
	}
	tAuthData.userid = user.ID
}
func Test_AuthBadLogIn(t *testing.T) {
	if _, err := tAuthData.service.LogIn(tAuthData.email, tAuthData.badSecret); err == nil {
		t.Errorf("Should not log in with invalid password")
		return
	}
}

func Test_PasswordlessRequestAdd(t *testing.T) {
	request, err := tPasswordlessRequestStore.Add(tAuthData.userid.Hex())
	if err != nil {
		t.Error(err)
		return
	}
	if request.Code == "" {
		t.Error("No passcode")
		return
	}
}

func Test_PasswordlessRequestDoubleAdd(t *testing.T) {
	request, err := tPasswordlessRequestStore.Add(tAuthData.userid.Hex())
	if err != nil {
		t.Error(err)
		return
	}
	if request.Code == "" {
		t.Error("No passcode")
		return
	}

	tAuthData.code = request.Code
}

func Test_PasswordlessRequestVerify(t *testing.T) {
	if err := tPasswordlessRequestStore.Verify(tAuthData.code, tAuthData.userid.Hex()); err != nil {
		t.Error(err)
		return
	}
}

func Test_AuthGenerateToken(t *testing.T) {
	token, err := tAuthData.service.UserIdToToken(tAuthData.userid)
	if err != nil {
		t.Error(err)
		return
	}
	if token == "" {
		t.Errorf("Token missing")
		return
	}
	tAuthData.token = token
}

func Test_AuthUserForToken(t *testing.T) {
	user, err := tAuthData.service.UserForToken(tAuthData.token)
	if err != nil {
		t.Error(err)
		return
	}
	if user.ID.Hex() == "" {
		t.Errorf("No user id")
		return
	}
	if user.Email != tAuthData.email {
		t.Errorf("Invalid email, expected '%s', got '%s'", tAuthData.email, user.Email)
		return
	}
}

func Test_AuthTokenFromRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tAuthData.token))

	token, err := tAuthData.service.TokenFromRequest(req)
	if err != nil {
		t.Error(err)
		return
	}
	if token != tAuthData.token {
		t.Errorf("Invalid token, expected '%s', got '%s'", tAuthData.token, token)
		return
	}
}
