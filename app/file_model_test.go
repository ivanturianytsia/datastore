package main

import (
	"fmt"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

var tFile = struct {
	store    FileStore
	fileId   bson.ObjectId
	ownerId  bson.ObjectId
	userId   bson.ObjectId
	filename string
	path     string
}{
	ownerId:  bson.NewObjectId(),
	userId:   bson.NewObjectId(),
	filename: "test.jpeg",
}

func Test_Create(t *testing.T) {
	file, err := tFile.store.Create(tFile.filename, tFile.ownerId.Hex())
	if err != nil {
		t.Error(err)
		return
	}
	if file.Id.Hex() == "" {
		t.Errorf("File has no id.")
		return
	}
	if file.OwnerId.Hex() != tFile.ownerId.Hex() {
		t.Errorf("File owner id invalid, expected '%s', got '%s'.", tFile.ownerId, file.OwnerId)
		return
	}
	if file.Filename != tFile.filename {
		t.Errorf("File name invalid, expected '%s', got '%s'.", tFile.filename, file.Filename)
		return
	}
	tFile.path = fmt.Sprintf("%s/%s", tFile.ownerId.Hex(), tFile.filename)
	if file.Path != tFile.path {
		t.Errorf("File path invalid, expected '%s', got '%s'.", tFile.path, file.Path)
		return
	}
	tFile.fileId = file.Id
}
func Test_AddAlowedIdById(t *testing.T) {
	file, err := tFile.store.AddAlowedIdById(tFile.fileId.Hex(), tFile.userId.Hex())
	if err != nil {
		t.Error(err)
		return
	}
	if _, ok := file.AllowedIds[tFile.userId.Hex()]; !ok {
		t.Errorf("File allowed ids doesn't contain '%s'.", tFile.userId)
		return
	}
}

func Test_GetByAllowedIdGood(t *testing.T) {
	files, err := tFile.store.GetByAllowedId(tFile.userId.Hex())
	if err != nil {
		t.Error(err)
		return
	}
	if len(files) != 1 {
		t.Errorf("File slice should contain 1 item, shared with user '%s', but contains %s.", tFile.userId, len(files))
		return
	}
}
func Test_RemoveAlowedIdById(t *testing.T) {
	file, err := tFile.store.RemoveAlowedIdById(tFile.fileId.Hex(), tFile.userId.Hex())
	if err != nil {
		t.Error(err)
		return
	}
	if _, ok := file.AllowedIds[tFile.userId.Hex()]; ok {
		t.Errorf("File allowed ids should not contain '%s'.", tFile.userId)
		return
	}
}
func Test_GetByAllowedIdBad(t *testing.T) {
	files, err := tFile.store.GetByAllowedId(tFile.userId.Hex())
	if err != nil {
		t.Error(err)
		return
	}
	if len(files) != 0 {
		t.Errorf("File slice should not contain items, but contains %s.", len(files))
		return
	}
}
func Test_GetById(t *testing.T) {
	file, err := tFile.store.GetById(tFile.fileId.Hex())
	if err != nil {
		t.Error(err)
		return
	}
	if file.Id.Hex() == "" {
		t.Errorf("File has no id.")
		return
	}
	if file.OwnerId.Hex() != tFile.ownerId.Hex() {
		t.Errorf("File owner id invalid, expected '%s', got '%s'.", tFile.ownerId, file.OwnerId)
		return
	}
	if file.Filename != tFile.filename {
		t.Errorf("File name invalid, expected '%s', got '%s'.", tFile.filename, file.Filename)
		return
	}
	if file.Path != tFile.path {
		t.Errorf("File path invalid, expected '%s', got '%s'.", tFile.path, file.Path)
		return
	}
}
func Test_GetByPath(t *testing.T) {
	if _, err := tFile.store.GetByPath(tFile.path); err != nil {
		t.Error(err)
		return
	}
}
func Test_GetByOwnerId(t *testing.T) {
	files, err := tFile.store.GetByOwnerId(tFile.ownerId.Hex())
	if err != nil {
		t.Error(err)
		return
	}
	if len(files) != 1 {
		t.Errorf("File slice should contain 1 item, owned by '%s', but contains %s.", tFile.ownerId, len(files))
		return
	}
}
