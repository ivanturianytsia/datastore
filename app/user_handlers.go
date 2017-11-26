package main

import (
	"net/http"
)

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
