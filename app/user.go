package main

import (
	"fmt"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Email     string        `bson:"email,omitempty" json:"email,omitempty"`
	Hash      string        `bson:"hash,omitempty" json:"hash,omitempty"`
	Roles     []string      `bson:"roles,omitempty" json:"roles,omitempty"`
	CreatedOn time.Time     `bson:"created,omitempty" json:"created,omitempty"`
	UpdatedOn time.Time     `bson:"updated,omitempty" json:"updated,omitempty"`
	Activated bool          `bson:"activated,omitempty" json:"activated,omitempty"`
	OldHashes []string      `bson:"oldhashes,omitempty" json:"oldhashes,omitempty"`
}

type UserStore interface {
	Create(email, password string) (User, error)
	ReadById(id string) (User, error)
	ReadByEmail(email string) (User, error)
	Update(id string, updates UserUpdates) (User, error)
	Delete(id string) error
	CheckPassword(password, hash string) error
}

type mongoUserStore struct {
	db         *Database
	hasher     Hasher
	collection string
}

func NewUserStore(db *Database) (UserStore, error) {
	var collection = "users"
	if err := db.EnsureUnique(collection, []string{"email"}); err != nil {
		return nil, err
	}
	return &mongoUserStore{
		db:         db,
		hasher:     NewBCryptHasher(),
		collection: collection,
	}, nil
}

func (store *mongoUserStore) Create(email, password string) (User, error) {
	if email == "" {
		return User{}, fmt.Errorf("Email missing")
	}
	if password == "" {
		return User{}, fmt.Errorf("Password missing")
	}

	hash, err := store.hasher.Hash(password)
	if err != nil {
		return User{}, err
	}
	newUser := User{
		ID:        bson.NewObjectId(),
		Email:     email,
		Hash:      hash,
		UpdatedOn: time.Now(),
		CreatedOn: time.Now(),
	}

	if err := store.db.WithCollection(store.collection, func(c *mgo.Collection) error {
		return c.Insert(newUser)
	}); err != nil {
		return User{}, err
	}

	return newUser, nil
}

func (store *mongoUserStore) ReadById(id string) (User, error) {
	if !bson.IsObjectIdHex(id) {
		return User{}, fmt.Errorf("Id '%s' is invalid", id)
	}

	var user User
	if err := store.db.WithCollection(store.collection, func(c *mgo.Collection) error {
		return c.FindId(bson.ObjectIdHex(id)).One(&user)
	}); err != nil {
		return User{}, err
	}

	return user, nil
}

func (store *mongoUserStore) ReadByEmail(email string) (User, error) {
	if email == "" {
		return User{}, fmt.Errorf("Email '%s' is invalid", email)
	}

	var user User
	if err := store.db.WithCollection(store.collection, func(c *mgo.Collection) error {
		return c.Find(bson.M{"email": email}).One(&user)
	}); err != nil {
		return User{}, err
	}

	return user, nil
}

func (store *mongoUserStore) Update(id string, updates UserUpdates) (User, error) {
	if !bson.IsObjectIdHex(id) {
		return User{}, fmt.Errorf("Id '%s' is invalid", id)
	}

	if v, ok := updates["password"].(string); ok {
		old, err := store.ReadById(id)
		if err != nil {
			return User{}, err
		}

		hash, err := store.hasher.Hash(v)
		if err != nil {
			return User{}, err
		}
		delete(updates, "password")
		updates["oldhashes"] = append(old.OldHashes, hash)
		updates["hash"] = hash
	}

	var user User
	if err := store.db.WithCollection(store.collection, func(c *mgo.Collection) error {
		if err := c.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": updates}); err != nil {
			return err
		}
		return c.FindId(bson.ObjectIdHex(id)).One(&user)
	}); err != nil {
		return User{}, err
	}

	return user, nil
}

func (store *mongoUserStore) Delete(id string) error {
	if !bson.IsObjectIdHex(id) {
		return fmt.Errorf("Id '%s' is invalid", id)
	}
	return store.db.WithCollection(store.collection, func(c *mgo.Collection) error {
		return c.RemoveId(bson.ObjectIdHex(id))
	})
}

func (store *mongoUserStore) CheckPassword(id, password string) error {
	user, err := store.ReadById(id)
	if err != nil {
		return err
	}
	return store.hasher.Check(password, user.Hash)
}

func NewUpdates() UserUpdates {
	return UserUpdates{}
}

type UserUpdates bson.M

func (uu UserUpdates) Email(email string) UserUpdates {
	uu["email"] = email

	return uu
}

func (uu UserUpdates) Password(password string) UserUpdates {
	uu["password"] = password

	return uu
}

func (uu UserUpdates) Roles(roles []string) UserUpdates {
	uu["roles"] = roles

	return uu
}

func (uu UserUpdates) Activated(activated bool) UserUpdates {
	uu["activated"] = activated

	return uu
}
