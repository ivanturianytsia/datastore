package main

import (
	"fmt"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type PasswordlessRequestStore interface {
	Add(userId string) (PasswordlessRequest, error)
	Verify(code, userId string) error
}

type PasswordlessRequest struct {
	UserId    bson.ObjectId `bson:"userid,omitempty" json:"userid,omitempty"`
	Code      string        `bson:"code,omitempty" json:"code,omitempty"`
	CreatedOn time.Time     `bson:"createdon,omitempty" json:"createdon,omitempty"`
}

type passwordlessRequestStore struct {
	db         *Database
	collection string
	expiration time.Duration
	hasher     Hasher
}

func NewPasswordlessRequestStore(db *Database) (PasswordlessRequestStore, error) {
	var collection = "emailrequests"
	if err := db.EnsureUnique(collection, []string{"userid"}); err != nil {
		return nil, err
	}
	return &passwordlessRequestStore{
		db:         db,
		hasher:     NewBCryptHasher(),
		collection: collection,
		expiration: (time.Minute * 5),
	}, nil
}

func (store *passwordlessRequestStore) Add(userId string) (PasswordlessRequest, error) {
	if !bson.IsObjectIdHex(userId) {
		return PasswordlessRequest{}, fmt.Errorf("User id '%s' is not valid.", userId)
	}

	if err := store.RemoveByUser(userId); err != nil {
		return PasswordlessRequest{}, err
	}

	request := PasswordlessRequest{
		UserId:    bson.ObjectIdHex(userId),
		Code:      genCode(4),
		CreatedOn: time.Now(),
	}

	hashed, err := store.hasher.Hash(request.Code)
	if err != nil {
		return PasswordlessRequest{}, err
	}

	if err := store.db.WithCollection(store.collection, func(c *mgo.Collection) error {
		return c.Insert(PasswordlessRequest{
			UserId:    request.UserId,
			Code:      hashed,
			CreatedOn: request.CreatedOn,
		})
	}); err != nil {
		return PasswordlessRequest{}, err
	}

	return request, nil
}

func (store *passwordlessRequestStore) Verify(code, userId string) error {
	if !bson.IsObjectIdHex(userId) {
		return fmt.Errorf("User id '%s' is not valid.", userId)
	}
	if code == "" {
		return fmt.Errorf("Code '%s' is not valid.", code)
	}
	var result PasswordlessRequest
	if err := store.db.WithCollection(store.collection, func(c *mgo.Collection) error {
		return c.Find(bson.M{"userid": bson.ObjectIdHex(userId)}).One(&result)
	}); err != nil {
		return err
	}
	if result.CreatedOn.Add(store.expiration).Before(time.Now()) {
		return fmt.Errorf("Code '%s' expired.", code)
	}
	if err := store.hasher.Check(code, result.Code); err != nil {
		return err
	}
	return store.RemoveByUser(userId)
}

func (store *passwordlessRequestStore) RemoveByUser(userId string) error {
	if !bson.IsObjectIdHex(userId) {
		return fmt.Errorf("User id '%s' is not valid.", userId)
	}
	return store.db.WithCollection(store.collection, func(c *mgo.Collection) error {
		_, err := c.RemoveAll(bson.M{"userid": bson.ObjectIdHex(userId)})
		return err
	})
}
