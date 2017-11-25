package main

import (
	"fmt"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// As owner I want to get whitelisted users for my file
// As owner I want to add users to whitelist
// As owner I want to remove users from whitelist
// As owner I want to download the file if I am owner
// As user I want to download the file if I am whitelisted
// As user I want to view list of files shared with me

type File struct {
	Id         bson.ObjectId       `bson:"_id,omitempty" json:"id,omitempty"`
	Filename   string              `bson:"filename,omitempty" json:"filename,omitempty"`
	Size       int64               `bson:"size,omitempty" json:"size,omitempty"`
	OwnerId    bson.ObjectId       `bson:"ownerid,omitempty" json:"ownerid,omitempty"`
	Path       string              `bson:"path,omitempty" json:"path,omitempty"`
	AllowedIds map[string]struct{} `bson:"allowedids" json:"allowedids"`
	CreatedOn  time.Time           `bson:"created" json:"created"`
	UpdatedOn  time.Time           `bson:"updated" json:"updated"`
}

type FileStore interface {
	Create(filename, ownerId string, size int64) (File, error)
	DeleteByPath(path string) error
	UpdateById(id string, allowedIds map[string]struct{}) (File, error)
	GetById(id string) (File, error)
	GetByPath(path string) (File, error)
	GetByOwnerId(ownerId string) ([]File, error)
	GetByAllowedId(userId string) ([]File, error)
}

type fileStore struct {
	db         *Database
	collection string
}

func NewFileStore(db *Database) (FileStore, error) {
	var collection = "files"
	if err := db.EnsureUnique(collection, []string{"path"}); err != nil {
		return nil, err
	}
	return &fileStore{
		db:         db,
		collection: collection,
	}, nil
}

func (store *fileStore) Create(filename, ownerId string, size int64) (File, error) {
	if !bson.IsObjectIdHex(ownerId) {
		return File{}, fmt.Errorf("Owner id '%s' is invalid", ownerId)
	}
	now := time.Now()
	file := File{
		Id:         bson.NewObjectId(),
		Filename:   filename,
		Size:       size,
		OwnerId:    bson.ObjectIdHex(ownerId),
		Path:       fmt.Sprintf("%s/%s", ownerId, filename),
		AllowedIds: map[string]struct{}{},
		CreatedOn:  now,
		UpdatedOn:  now,
	}
	if err := store.db.WithCollection(store.collection, func(c *mgo.Collection) error {
		return c.Insert(file)
	}); err != nil {
		return File{}, err
	}
	return file, nil
}

func (store *fileStore) DeleteByPath(path string) error {
	return store.db.WithCollection(store.collection, func(c *mgo.Collection) error {
		return c.Remove(bson.M{"path": path})
	})
}
func (store *fileStore) UpdateById(id string, allowedIds map[string]struct{}) (File, error) {
	if !bson.IsObjectIdHex(id) {
		return File{}, fmt.Errorf("File id '%s' is invalid", id)
	}
	var file File
	if err := store.db.WithCollection(store.collection, func(c *mgo.Collection) error {
		return c.FindId(bson.ObjectIdHex(id)).One(&file)
	}); err != nil {
		return File{}, err
	}
	file.AllowedIds = allowedIds
	file.UpdatedOn = time.Now()
	if err := store.db.WithCollection(store.collection, func(c *mgo.Collection) error {
		return c.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": file})
	}); err != nil {
		return File{}, err
	}
	return file, nil
}
func (store *fileStore) GetById(id string) (File, error) {
	if !bson.IsObjectIdHex(id) {
		return File{}, fmt.Errorf("File id '%s' is invalid", id)
	}
	var file File
	if err := store.db.WithCollection(store.collection, func(c *mgo.Collection) error {
		return c.FindId(bson.ObjectIdHex(id)).One(&file)
	}); err != nil {
		return File{}, err
	}
	return file, nil
}
func (store *fileStore) GetByPath(path string) (File, error) {
	var file File
	if err := store.db.WithCollection(store.collection, func(c *mgo.Collection) error {
		return c.Find(bson.M{"path": path}).One(&file)
	}); err != nil {
		return File{}, err
	}
	return file, nil
}
func (store *fileStore) GetByOwnerId(ownerId string) ([]File, error) {
	if !bson.IsObjectIdHex(ownerId) {
		return []File{}, fmt.Errorf("Owner id '%s' is invalid", ownerId)
	}
	var files []File
	if err := store.db.WithCollection(store.collection, func(c *mgo.Collection) error {
		return c.Find(bson.M{"ownerid": bson.ObjectIdHex(ownerId)}).All(&files)
	}); err != nil {
		return []File{}, err
	}
	return files, nil
}
func (store *fileStore) GetByAllowedId(userId string) ([]File, error) {
	if !bson.IsObjectIdHex(userId) {
		return []File{}, fmt.Errorf("User id '%s' is invalid", userId)
	}
	var files []File

	if err := store.db.WithCollection(store.collection, func(c *mgo.Collection) error {
		return c.Find(bson.M{(fmt.Sprintf("allowedids.%s", userId)): bson.M{"$exists": true}}).All(&files)
	}); err != nil {
		return []File{}, err
	}
	return files, nil
}
