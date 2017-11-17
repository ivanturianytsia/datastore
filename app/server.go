package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/mux"
)

var index []byte

func init() {
	var err error
	index, err = ioutil.ReadFile(path.Join(getDistDir(), "index.html"))
	if err != nil {
		log.Fatalln(err)
		return
	}
}

type Server struct {
	auth AuthService
	user UserStore
}

func NewServer() (*Server, error) {
	db := NewDatabase("agh-thesis")
	if err := db.Connect(os.Getenv("DB")); err != nil {
		return nil, err
	}
	user, err := NewUserStore(db)
	if err != nil {
		return nil, err
	}
	auth := NewAuthService(user)
	return &Server{
		auth: auth,
		user: user,
	}, nil
}

func (s Server) Route(router *mux.Router) {
	router.Methods("GET").Path("/").HandlerFunc(s.handlePage)
	router.Methods("GET").Path("/data").HandlerFunc(s.handleData)
	router.Methods("GET").Path("/auth/user").HandlerFunc(s.handleUser)
	router.Methods("POST").Path("/auth/login").HandlerFunc(s.handleLogin)
	router.Methods("POST").Path("/auth/register").HandlerFunc(s.handleRegister)

	router.PathPrefix("/static/").Handler(
		http.FileServer(
			http.Dir(path.Join(getDistDir()))))
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type TokenBody struct {
	Token string `json:"token"`
}

func (s Server) handlePage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write(index)
}

func (s Server) handleData(w http.ResponseWriter, r *http.Request) {
	data := []string{"one", "two", "three"}
	Respond(w, r, http.StatusOK, data)
}

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
	token, err := s.auth.LogIn(cred.Email, cred.Password)
	if err != nil {
		RespondErr(w, r, http.StatusForbidden, err)
		return
	}
	Respond(w, r, http.StatusOK, TokenBody{token})
}

func (s Server) handleRegister(w http.ResponseWriter, r *http.Request) {
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
	token, err := s.auth.Register(cred.Email, cred.Password)
	if err != nil {
		RespondErr(w, r, http.StatusForbidden, err)
		return
	}
	Respond(w, r, http.StatusOK, TokenBody{token})
}

func (s Server) handleUser(w http.ResponseWriter, r *http.Request) {
	user, err := s.auth.UserFromRequest(r)
	if err != nil {
		RespondErr(w, r, http.StatusForbidden, err)
		return
	}
	Respond(w, r, http.StatusOK, user)
}
