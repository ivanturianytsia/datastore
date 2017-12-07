package main

import (
	"fmt"
	"net/http"
)

func (s Server) handleUser(w http.ResponseWriter, r *http.Request) {
	user, err := s.auth.UserFromRequest(r)
	if err != nil {
		RespondErr(w, r, http.StatusForbidden, err)
		return
	}
	Respond(w, r, http.StatusOK, user)
}

func (s Server) handlePutUser(w http.ResponseWriter, r *http.Request) {
	user, err := s.auth.UserFromRequest(r)
	if err != nil {
		RespondErr(w, r, http.StatusForbidden, err)
		return
	}
	var body map[string]interface{}
	if err := DecodeBody(r, &body); err != nil {
		RespondErr(w, r, http.StatusBadRequest, err)
		return
	}
	updates := UserUpdates{}
	shouldUpdate := false
	if v, ok := body["twofactor"].(bool); ok {
		updates.TwoFactor(v)
		shouldUpdate = true
	}
	if !shouldUpdate {
		RespondErr(w, r, http.StatusBadRequest, fmt.Errorf("Nothing to update"))
		return
	}
	updated, err := s.user.Update(user.ID.Hex(), updates)
	Respond(w, r, http.StatusOK, updated)
}

func (s Server) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	user, err := s.auth.UserFromRequest(r)
	if err != nil {
		RespondErr(w, r, http.StatusForbidden, err)
		return
	}
	users, err := s.user.SearchByEmail(r.URL.Query().Get("email"))
	if err != nil {
		RespondErr(w, r, http.StatusInternalServerError, err)
		return
	}
	for i, v := range users {
		if v.ID.Hex() == user.ID.Hex() {
			users = append(users[:i], users[i+1:]...)
		}
	}
	Respond(w, r, http.StatusOK, users)
}
