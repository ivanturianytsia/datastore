package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var index []byte

func init() {
	var err error
	index, err = ioutil.ReadFile("./client/index.html")
	if err != nil {
		log.Fatalln(err)
		return
	}
}

type Server struct {
	auth UserService
	user UserStore
}

func NewServer() (*Server, error) {
	db := NewDatabase("agh-thesis")
	if err := db.Connect(os.Getenv("DB")); err != nil {
		return nil, err
	}
	userstore, err := NewUserStore(db)
	if err != nil {
		return nil, err
	}
	auth := NewUserService(userstore)
	return &Server{
		auth: auth,
		user: userstore,
	}, nil
}

func (s Server) Route(router *mux.Router) {
	router.Methods("GET").Path("/").HandlerFunc(s.handlePage)
	router.Methods("GET").Path("/data").HandlerFunc(s.handleData)
	router.Methods("POST").Path("/auth/login").HandlerFunc(s.handleLogin)
	router.Methods("POST").Path("/auth/register").HandlerFunc(s.handleRegister)

	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/",
			http.FileServer(
				http.Dir(getDistDir()))))
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
	user, err := s.auth.LogIn(cred.Email, cred.Password)
	if err != nil {
		RespondErr(w, r, http.StatusForbidden, err)
		return
	}

	Respond(w, r, http.StatusOK, user)
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
	user, err := s.auth.Register(cred.Email, cred.Password)
	if err != nil {
		RespondErr(w, r, http.StatusForbidden, err)
		return
	}

	Respond(w, r, http.StatusOK, user)
}
