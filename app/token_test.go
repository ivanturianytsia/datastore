package main

import (
	"testing"
	"time"
)

var testTokenClaims = map[string]string{
	"foo": "bar",
	"key": "val",
}
var testSecret = []byte("secret1")
var testToken string
var testTokenService TokenService

func init() {
	testTokenService = NewTokenService(testSecret, time.Minute*5)
}

func Test_IssueToken(t *testing.T) {
	var err error
	testToken, err = testTokenService.Issue(testTokenClaims)
	if err != nil {
		t.Error(err)
		return
	}
	if testToken == "" {
		t.Error("No Token")
		return
	}
}

func Test_ParseToken(t *testing.T) {
	claims, err := testTokenService.Parse(testToken)
	if err != nil {
		t.Error(err)
		return
	}
	for key, val := range testTokenClaims {
		if claims[key] != val {
			t.Errorf("Key '%s' expected to be '%s' but got '%s'", key, val, claims[key])
			return
		}
	}
}
