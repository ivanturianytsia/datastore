package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/mux"
)

type Server struct {
	auth  AuthService
	user  UserStore
	index []byte
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
	index, err := ioutil.ReadFile(path.Join(getDistDir(), "index.html"))
	if err != nil {
		return nil, err
	}
	return &Server{
		auth:  auth,
		user:  user,
		index: index,
	}, nil
}

func (s Server) Route(router *mux.Router) {
	router.Methods("GET").Path("/").HandlerFunc(s.handlePage)
	router.Methods("GET").Path("/data").HandlerFunc(s.handleData)
	router.Methods("GET").Path("/auth/user").HandlerFunc(s.handleUser)
	router.Methods("POST").Path("/auth/login").HandlerFunc(s.handleLogin)
	router.Methods("POST").Path("/auth/register").HandlerFunc(s.handleRegister)

	router.Methods("POST").Path("/upload").HandlerFunc(s.handleUpload)
	router.Methods("GET").Path("/files").HandlerFunc(s.handleGetFiles)
	router.Methods("GET").Path("/files/{filename}").HandlerFunc(s.handleGetFile)
	router.Methods("DELETE").Path("/files/{filename}").HandlerFunc(s.handleDeleteFile)

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
	w.Write(s.index)
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

func (s Server) handleUpload(w http.ResponseWriter, r *http.Request) {
	user, err := s.auth.UserFromRequest(r)
	if err != nil {
		RespondErr(w, r, http.StatusForbidden, err)
		return
	}
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		RespondErr(w, r, http.StatusBadRequest, err)
		return
	}
	defer file.Close()

	filename := path.Join(getDataDir(), user.ID.Hex(), handler.Filename)
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		RespondErr(w, r, http.StatusInternalServerError, err)
		return
	}
	defer f.Close()
	if _, err := io.Copy(f, file); err != nil {
		RespondErr(w, r, http.StatusInternalServerError, err)
		return
	}
	Respond(w, r, http.StatusOK, []string{handler.Filename})
}
func (s Server) handleGetFiles(w http.ResponseWriter, r *http.Request) {
	user, err := s.auth.UserFromRequest(r)
	if err != nil {
		RespondErr(w, r, http.StatusForbidden, err)
		return
	}
	dir := path.Join(getDataDir(), user.ID.Hex())
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		RespondErr(w, r, http.StatusBadRequest, err)
		return
	}
	filelist := []map[string]interface{}{}
	for _, v := range files {
		if !v.IsDir() {
			filelist = append(filelist, map[string]interface{}{
				"name":     v.Name(),
				"modified": v.ModTime(),
				"size":     v.Size(),
			})
		}
	}
	Respond(w, r, http.StatusOK, filelist)
}
func (s Server) handleGetFile(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	user, err := s.auth.UserForToken(token)
	if err != nil {
		RespondErr(w, r, http.StatusForbidden, err)
		return
	}
	filename := path.Join(getDataDir(), user.ID.Hex(), mux.Vars(r)["filename"])
	f, err := os.Open(filename)
	if err != nil {
		RespondErr(w, r, http.StatusInternalServerError, err)
		return
	}
	defer f.Close()
	if _, err := io.Copy(w, f); err != nil {
		RespondErr(w, r, http.StatusInternalServerError, err)
		return
	}
}
func (s Server) handleDeleteFile(w http.ResponseWriter, r *http.Request) {
	user, err := s.auth.UserFromRequest(r)
	if err != nil {
		RespondErr(w, r, http.StatusForbidden, err)
		return
	}
	filename := path.Join(getDataDir(), user.ID.Hex(), mux.Vars(r)["filename"])
	if err := os.Remove(filename); err != nil {
		RespondErr(w, r, http.StatusInternalServerError, err)
		return
	}
	Respond(w, r, http.StatusNoContent, nil)
}
