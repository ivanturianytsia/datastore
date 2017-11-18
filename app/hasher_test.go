package main

import (
	"testing"
)

var testPassword = "secret"
var testBadPassword = "secret-is-wrong"
var testHash string
var testHasher Hasher

func init() {
	testHasher = NewBCryptHasher()
}

func Test_HashPassword(t *testing.T) {
	var err error
	testHash, err = testHasher.Hash(testPassword)
	if err != nil {
		t.Error(err)
		return
	}
	if testHash == "" {
		t.Error("No hash returned by .Hash()")
		return
	}
}

func Test_CheckPassword(t *testing.T) {
	err := testHasher.Check(testPassword, testHash)
	if err != nil {
		t.Error(err)
		return
	}
}
func Test_CheckBadPassword(t *testing.T) {
	err := testHasher.Check(testBadPassword, testHash)
	if err == nil {
		t.Error("Bad password expected to returns no error")
		return
	}
}
