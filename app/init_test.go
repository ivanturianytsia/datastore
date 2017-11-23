package main

import (
	"os"
)

var defaultTestDB = "mongodb://ivan:ivan@cluster0-shard-00-00-p1nri.mongodb.net:27017,cluster0-shard-00-01-p1nri.mongodb.net:27017,cluster0-shard-00-02-p1nri.mongodb.net:27017/agh-datastore-test?ssl=true&replicaSet=Cluster0-shard-0&authSource=admin"

func init() {
	testDB := os.Getenv("TEST_DB")
	if testDB == "" {
		testDB = defaultTestDB
	}

	db := NewDatabase("agh-datastore-test")
	if err := db.Connect(testDB); err != nil {
		panic(err)
		return
	}

	if err := db.Drop(); err != nil {
		panic(err)
		return
	}
	var err error
	if tUserData.store, err = NewUserStore(db); err != nil {
		panic(err)
		return
	}

	tAuthData.service = NewAuthService(tUserData.store)

	if tPasswordlessRequestStore, err = NewPasswordlessRequestStore(db); err != nil {
		panic(err)
		return
	}
}
