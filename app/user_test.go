package main

import (
	"testing"

	"gopkg.in/mgo.v2/bson"
)

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

func Test_UserCreate(t *testing.T) {
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

func Test_UserReadByEmailAndCheckPassword(t *testing.T) {
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
func Test_UserUpdate(t *testing.T) {
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

func Test_UserDelete(t *testing.T) {
	if err := tUserData.store.Delete(tUserData.id.Hex()); err != nil {
		t.Error(err)
		return
	}
	if _, err := tUserData.store.ReadById(tUserData.id.Hex()); err == nil {
		t.Error("Deleted user should not be found")
		return
	}
}
