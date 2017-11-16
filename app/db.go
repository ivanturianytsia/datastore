package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"regexp"
	"strings"

	mgo "gopkg.in/mgo.v2"
)

type Database struct {
	Session *mgo.Session
	Name    string
}

func NewDatabase(name string) *Database {
	return &Database{
		Name: name,
	}
}

func (d *Database) interpretConnectionString(url string) (bool, string) {
	ssl := false
	parts := strings.SplitN(url, "?", 2)
	if len(parts) == 2 {
		r := regexp.MustCompile(`ssl=true\&?`)
		replace := r.ReplaceAllString(parts[1], "")
		if replace != parts[1] {
			// SSL detected
			ssl = true
			url = strings.Join([]string{parts[0], replace}, "?")
		}
	}
	return ssl, url
}

func (d *Database) Drop() error {
	session := d.Session.Copy()
	defer session.Close()
	return session.DB(d.Name).DropDatabase()
}

func (d *Database) EnsureUnique(collectionName string, keys []string) error {
	session := d.Session.Copy()
	defer session.Close()

	if err := session.DB(d.Name).C(collectionName).EnsureIndex(mgo.Index{
		Key:    keys,
		Unique: true,
	}); err != nil {
		return err
	}
	return nil
}

func (d *Database) Connect(url string) error {
	ssl, url := d.interpretConnectionString(url)

	if !ssl {
		// If no-ssl connection
		s, err := mgo.Dial(url)
		if err != nil {
			return err
		}
		d.Session = s
	} else {
		// If ssl connection
		tlsConfig := &tls.Config{}
		dialInfo, err := mgo.ParseURL(url)
		if err != nil {
			return err
		}
		dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
			return conn, err
		}
		s, err := mgo.DialWithInfo(dialInfo)
		if err != nil {
			return err
		}
		d.Session = s
	}
	return nil
}

func (d *Database) WithCollection(collectionName string, f func(c *mgo.Collection) error) error {
	if d.Session == nil {
		return fmt.Errorf("Database is not connected, use:\n err := db.Connect(url)")
	}

	session := d.Session.Copy()
	defer session.Close()
	c := session.DB(d.Name).C(collectionName)

	return f(c)
}
