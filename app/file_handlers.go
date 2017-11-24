package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/mux"
)

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
