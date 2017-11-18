package main

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

var defaultTestDB = "mongodb://ivan:ivan@cluster0-shard-00-00-p1nri.mongodb.net:27017,cluster0-shard-00-01-p1nri.mongodb.net:27017,cluster0-shard-00-02-p1nri.mongodb.net:27017/agh-datastore-test?ssl=true&replicaSet=Cluster0-shard-0&authSource=admin"

var tUserData = struct {
	store     UserStore
	email     string
	secret    string
	newSecret string
	id        bson.ObjectId
}{
	email:     "test1@mail.com",
	secret:    "secret",
	newSecret: "secret1",
}

var tAuthData = struct {
	service   AuthService
	email     string
	secret    string
	badSecret string
	token     string
}{
	email:     "test2@mail.com",
	secret:    "secret",
	badSecret: "secret1",
}

func init() {
	testDB := os.Getenv("TEST_DB")
	if testDB == "" {
		testDB = defaultTestDB
	}

	db := NewDatabase("agh-datastore-test")
	if err := db.Connect(testDB); err != nil {
		panic(err)
		return
	}

	if err := db.Drop(); err != nil {
		panic(err)
		return
	}
	var err error
	if tUserData.store, err = NewUserStore(db); err != nil {
		panic(err)
		return
	}

	tAuthData.service = NewAuthService(tUserData.store)
}

func Test_CreateUser(t *testing.T) {
	user, err := tUserData.store.Create(tUserData.email, tUserData.secret)
	if err != nil {
		t.Error(err)
		return
	}
	if user.ID.Hex() == "" {
		t.Error("No user id")
		return
	}
	if user.Email != tUserData.email {
		t.Error("Wrong email, got '%s', expected '%s'", user.Email, tUserData.email)
		return
	}

	tUserData.id = user.ID
}

func Test_ReadUserByEmailAndCheckPassword(t *testing.T) {
	user, err := tUserData.store.ReadByEmail(tUserData.email)
	if err != nil {
		t.Error(err)
		return
	}
	if user.ID.Hex() == "" {
		t.Error("No user id")
		return
	}
	if err := tUserData.store.CheckPassword(user.ID.Hex(), tUserData.secret); err != nil {
		t.Error(err)
		return
	}
}
func Test_UpdateUser(t *testing.T) {
	user, err := tUserData.store.Update(tUserData.id.Hex(), NewUpdates().Password(tUserData.newSecret))
	if err != nil {
		t.Error(err)
		return
	}
	if err := tUserData.store.CheckPassword(user.ID.Hex(), tUserData.secret); err == nil {
		t.Error("Old password should not work")
		return
	}
}

func Test_DeleteUser(t *testing.T) {
	if err := tUserData.store.Delete(tUserData.id.Hex()); err != nil {
		t.Error(err)
		return
	}
	if _, err := tUserData.store.ReadById(tUserData.id.Hex()); err == nil {
		t.Error("Deleted user should not be found")
		return
	}
}

// Service

func Test_AuthRegister(t *testing.T) {
	token, err := tAuthData.service.Register(tAuthData.email, tAuthData.secret)
	if err != nil {
		t.Error(err)
		return
	}
	if token == "" {
		t.Errorf("No token")
		return
	}
	tAuthData.token = token
}
func Test_AuthDoubleRegister(t *testing.T) {
	if _, err := tAuthData.service.Register(tAuthData.email, tAuthData.secret); err == nil {
		t.Error("Should not register same email")
		return
	}
}
func Test_AuthLogIn(t *testing.T) {
	token, err := tAuthData.service.LogIn(tAuthData.email, tAuthData.secret)
	if err != nil {
		t.Error(err)
		return
	}
	if token == "" {
		t.Errorf("No token")
		return
	}
}
func Test_AuthBadLogIn(t *testing.T) {
	if _, err := tAuthData.service.LogIn(tAuthData.email, tAuthData.badSecret); err == nil {
		t.Errorf("Should not log in with invalid password")
		return
	}
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
