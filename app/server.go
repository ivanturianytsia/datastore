package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/mux"
)

type Server struct {
	auth         AuthService
	user         UserStore
	smscode      CodeAuthService
	emailcode    CodeAuthService
	passwordless PasswordlessRequestStore
	files        FileStore
	index        []byte
}

func NewServer() (*Server, error) {
	db := NewDatabase("agh-datastore")
	if err := db.Connect(os.Getenv("DB")); err != nil {
		return nil, err
	}
	user, err := NewUserStore(db)
	if err != nil {
		return nil, err
	}
	auth := NewAuthService(user)
	passwordless, err := NewPasswordlessRequestStore(db)
	if err != nil {
		return nil, err
	}
	files, err := NewFileStore(db)
	if err != nil {
		return nil, err
	}
	emailcode := NewMailgunService()
	smscode := NewSMSApiService()
	index, err := ioutil.ReadFile(path.Join(getDistDir(), "index.html"))
	if err != nil {
		return nil, err
	}
	return &Server{
		auth:         auth,
		user:         user,
		files:        files,
		index:        index,
		emailcode:    emailcode,
		smscode:      smscode,
		passwordless: passwordless,
	}, nil
}

func (s Server) Route(router *mux.Router) {
	router.Methods("GET").Path("/").HandlerFunc(s.handlePage)
	router.Methods("GET").Path("/auth/user").HandlerFunc(s.handleUser)
	router.Methods("PUT").Path("/user").HandlerFunc(s.handlePutUser)
	router.Methods("POST").Path("/auth/login").HandlerFunc(s.handleLogin)
	router.Methods("POST").Path("/auth/register").HandlerFunc(s.handleRegister)
	router.Methods("POST").Path("/auth/code").HandlerFunc(s.handleCode)

	router.Methods("GET").Path("/users").HandlerFunc(s.handleGetUsers)

	router.Methods("POST").Path("/upload").HandlerFunc(s.handleUpload)
	router.Methods("GET").Path("/files").HandlerFunc(s.handleGetFiles)
	router.Methods("PUT").Path("/file/{fileid}").HandlerFunc(s.handleFileUpdate)
	router.Methods("GET").Path("/files/{ownerid}/{filename}").HandlerFunc(s.handleGetFile)
	router.Methods("DELETE").Path("/files/{filename}").HandlerFunc(s.handleDeleteFile)

	router.PathPrefix("/static/").Handler(
		http.FileServer(
			http.Dir(path.Join(getDistDir()))))

}

func (s Server) handlePage(w http.ResponseWriter, r *http.Request) {
	// TODO: remove
	s.index, _ = ioutil.ReadFile(path.Join(getDistDir(), "index.html"))

	w.WriteHeader(http.StatusOK)
	w.Write(s.index)
}
