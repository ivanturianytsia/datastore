package main

import (
	"net/http"
)

func (s Server) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	if _, err := s.auth.UserFromRequest(r); err != nil {
		RespondErr(w, r, http.StatusForbidden, err)
		return
	}
	users, err := s.user.SearchByEmail(r.URL.Query().Get("email"))
	if err != nil {
		RespondErr(w, r, http.StatusInternalServerError, err)
		return
	}
	Respond(w, r, http.StatusOK, users)
}
