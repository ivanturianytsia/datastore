package main

import (
	"fmt"
	"io"
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

	// Save file
	dir := path.Join(getDataDir(), user.ID.Hex())
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, os.ModePerm)
	}
	filename := path.Join(dir, handler.Filename)
	f, err := os.Create(filename)
	if err != nil {
		RespondErr(w, r, http.StatusInternalServerError, err)
		return
	}
	defer f.Close()
	written, err := io.Copy(f, file)
	if err != nil {
		RespondErr(w, r, http.StatusInternalServerError, err)
		return
	}

	// Create entry in db
	filedb, err := s.files.Create(handler.Filename, user.ID.Hex(), written)
	if err != nil {
		RespondErr(w, r, http.StatusBadRequest, err)
		return
	}
	Respond(w, r, http.StatusOK, filedb)
}
func (s Server) handleGetFiles(w http.ResponseWriter, r *http.Request) {
	user, err := s.auth.UserFromRequest(r)
	if err != nil {
		RespondErr(w, r, http.StatusForbidden, err)
		return
	}
	files, err := s.files.GetByOwnerId(user.ID.Hex())
	if err != nil {
		RespondErr(w, r, http.StatusBadRequest, err)
		return
	}
	shared, err := s.files.GetByAllowedId(user.ID.Hex())
	if err != nil {
		RespondErr(w, r, http.StatusBadRequest, err)
		return
	}
	Respond(w, r, http.StatusOK, map[string][]File{
		"files":  files,
		"shared": shared,
	})
}
func (s Server) handleGetFile(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	user, err := s.auth.UserForToken(token)
	if err != nil {
		RespondErr(w, r, http.StatusForbidden, err)
		return
	}
	filepath := fmt.Sprintf("%s/%s", mux.Vars(r)["ownerid"], mux.Vars(r)["filename"])
	file, err := s.files.GetByPath(filepath)
	if err != nil {
		RespondErr(w, r, http.StatusInternalServerError, err)
		return
	}
	// Check permissions
	if _, ok := file.AllowedIds[user.ID.Hex()]; (file.OwnerId.Hex() != user.ID.Hex()) && !ok {
		RespondErr(w, r, http.StatusForbidden, fmt.Errorf("You are not allowed to view this file"))
		return
	}

	filename := path.Join(getDataDir(), filepath)
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
	filepath := fmt.Sprintf("%s/%s", user.ID.Hex(), mux.Vars(r)["filename"])
	if err := s.files.DeleteByPath(filepath); err != nil {
		RespondErr(w, r, http.StatusInternalServerError, err)
		return
	}
	filename := path.Join(getDataDir(), filepath)
	if err := os.Remove(filename); err != nil {
		RespondErr(w, r, http.StatusInternalServerError, err)
		return
	}
	Respond(w, r, http.StatusNoContent, nil)
}

func (s Server) handleFileUpdate(w http.ResponseWriter, r *http.Request) {
	user, err := s.auth.UserFromRequest(r)
	if err != nil {
		RespondErr(w, r, http.StatusForbidden, err)
		return
	}

	// Check owner
	oldFile, err := s.files.GetById(mux.Vars(r)["fileid"])
	if err != nil {
		RespondErr(w, r, http.StatusInternalServerError, err)
		return
	}
	if oldFile.OwnerId.Hex() != user.ID.Hex() {
		RespondErr(w, r, http.StatusForbidden, fmt.Errorf("You have no permission to edit this file"))
		return
	}

	var body map[string]interface{}
	if err := DecodeBody(r, &body); err != nil {
		RespondErr(w, r, http.StatusBadRequest, err)
		return
	}
	updates := NewFileUpdates()
	shouldUpdate := false
	if v, ok := body["allowedids"].(map[string]interface{}); ok {
		var allowedIds = map[string]struct{}{}
		for i, _ := range v {
			allowedIds[i] = struct{}{}
		}
		updates.AllowedIds(allowedIds)
		shouldUpdate = true
	}
	if v, ok := body["description"].(string); ok {
		updates.Description(v)
		shouldUpdate = true
	}
	if !shouldUpdate {
		RespondErr(w, r, http.StatusBadRequest, fmt.Errorf("Nothing to update"))
		return
	}
	updated, err := s.files.UpdateById(oldFile.Id.Hex(), updates)
	Respond(w, r, http.StatusOK, updated)
}
